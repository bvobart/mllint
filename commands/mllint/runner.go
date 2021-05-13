package mllint

import (
	"io"
	"runtime"

	"github.com/bvobart/mllint/api"
)

const queueSize = 20 // arbitrarily chosen

// NewRunner initialises an *mllint.Runner
func NewRunner(progress RunnerProgress) *Runner {
	if progress == nil {
		progress = &BasicRunnerProgress{Out: io.Discard}
	}
	return &Runner{
		queue:    make(chan *RunnerTask, queueSize),
		done:     make(chan *RunnerTask, runtime.NumCPU()),
		progress: progress,
		closed:   make(chan struct{}),
		nRunning: 0,
	}
}

// Runner implements a parallel linter runner for mllint.
// Use `r := NewRunner()` to create a runner, then call r.Start() on it (and defer r.Close())
// to start a queue worker process that will watch for incoming tasks to run a linter,
// as added to the queue by calls to r.RunLinter().
//
// The Runner will ensure that all tasks added to its queue are executed,
// run in parallel across all `runtime.NumCPU()` available logical cores (CPU threads, aka the amount of bars you see in `htop` :P).
// No more than this number of tasks will be running in parallel, any additional tasks are parked and executed as running tasks complete.
//
// `r.RunLinter()` returns a *RunnerTask, use `result := <-task.Result` to await the completion of the linter task and receive the result of the task.
//
// Alternatively, if you have a list of tasks to await and want to receive their results as soon as each completes,
// use `mllint.CollectTasks(tasks...)` to get a channel where each completed task will be sent.
// A task you receive through this channel will have a `LinterResult` in the `task.Result`'s buffer, so `<-task.Result` will not block.
// The channel closes once all tasks have completed.
//
// Example usage with `mllint.ForEachTask()`:
//
// ```go
// mllint.ForEachTask(mllint.CollectTasks(tasks...), func(task *RunnerTask, result LinterResult) {
//   // do something with the completed task and its result
// })
// ```
type Runner struct {
	queue    chan *RunnerTask
	done     chan *RunnerTask
	progress RunnerProgress
	closed   chan struct{}
	nRunning int32
}

// RunnerTask represents a task to run a linter on a project that was created by a call to runner.RunLinter(...)
type RunnerTask struct {
	Id          string
	Linter      api.Linter
	Project     api.Project
	Result      chan LinterResult
	displayName string
}

// LinterResult represents the two-valued return type of a Linter, containing a report and an error.
type LinterResult struct {
	api.Report
	Err error
}

// TaskOption is an option for a task created by RunLinter, e.g. setting a custom display name.
type TaskOption func(task *RunnerTask)

func DisplayName(name string) TaskOption {
	return func(task *RunnerTask) {
		task.displayName = name
	}
}

// Start starts the runner by running a queue worker go-routine in the background that will await tasks and run them as they come in.
// After calling `runner.Start()`, make sure to also `defer runner.Close()`
func (r *Runner) Start() {
	r.progress.Start()
	go r.queueWorker()
}

// Close stops the runner by closing its queue channel. Calls to `runner.RunLinter()` after Close will panic.
// Close blocks until all tasks added to its queue by `runner.RunLinter()` have completed,
// and all progress output has finished printing to the terminal.
func (r *Runner) Close() {
	close(r.queue)
	<-r.closed
}

// RunLinter creates a task to run an api.Linter on a project.This method does not block.
// The task will be executed in parallel with other tasks by the runner.
// If the runner is `nil`, then the linter will be executed directly on the current thread, returning its result the usual way.
//
// Once the task completes, the linter's report and error will be sent to the task's `Result` channel,
// i.e. use `<-task.Result` to await the linter's result.
func (r *Runner) RunLinter(id string, linter api.Linter, project api.Project, options ...TaskOption) *RunnerTask {
	result := make(chan LinterResult, 1)
	task := RunnerTask{id, linter, project, result, linter.Name()}
	for _, optionFunc := range options {
		optionFunc(&task)
	}

	if r == nil {
		report, err := linter.LintProject(project)
		task.Result <- LinterResult{report, err}
		return &task
	}

	r.queue <- &task
	return &task
}

// Watches the queue for new jobs, running or parking them as they come in / complete.
// Run in a new go-routine using `go r.queueWorker()`
func (r *Runner) queueWorker() {
	parked := []*RunnerTask{}
	closed := false

	for {
		select {
		// when new task is scheduled...
		case task, open := <-r.queue:
			// if channel just closed and no tasks are running, signal that we're finished and exit
			if !open && r.nRunning == 0 {
				r.progress.AllTasksDone()
				close(r.closed)
				return
			}
			// if channel just closed, but there are still tasks running, signal to the next case that there will be no new tasks.
			if !open {
				closed = true
				break
			}

			// if we're already running the maximum number of tasks, park it
			if r.nRunning >= int32(runtime.NumCPU()) {
				parked = append(parked, task)
				break
			}

			// otherwise just run the task
			r.runTask(task)

		// when a task completes...
		case task := <-r.done:
			r.nRunning--
			r.progress.CompletedTask(task)

			// if the queue is closed and no other tasks are running, then signal that we're finished and exit
			if closed && r.nRunning == 0 {
				r.progress.AllTasksDone()
				close(r.closed)
				return
			}

			// else, if there are parked tasks, run one of them.
			if len(parked) > 0 {
				var next *RunnerTask
				next, parked = parked[0], parked[1:]
				r.runTask(next)
			}
		}
	}
}

// actually start running the task in a new go-routine
func (r *Runner) runTask(task *RunnerTask) {
	r.nRunning++
	r.progress.RunningTask(task)

	go func() {
		if l, ok := task.Linter.(WithRunner); ok {
			l.SetRunner(r)
		}

		report, err := task.Linter.LintProject(task.Project)
		task.Result <- LinterResult{Report: report, Err: err}
		r.done <- task
	}()
}
