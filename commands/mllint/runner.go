package mllint

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/bvobart/mllint/api"
)

const queueSize = 20 // arbitrarily chosen

// NewRunner initialises an *mllint.Runner
func NewRunner() *Runner {
	return &Runner{
		queue:    make(chan *RunnerTask, queueSize),
		done:     make(chan *RunnerTask, runtime.NumCPU()),
		progress: NewRunnerProgress(),
		running:  sync.Map{},
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
	progress *RunnerProgress
	running  sync.Map
	nRunning int32
}

// RunnerTask represents a task to run a linter on a project that was created by a call to runner.RunLinter(...)
type RunnerTask struct {
	Id      string
	Linter  api.Linter
	Project api.Project
	Result  chan LinterResult
}

// LinterResult represents the two-valued return type of a Linter, containing a report and an error.
type LinterResult struct {
	api.Report
	Err error
}

// Start starts the runner by running a queue worker go-routine in the background that will await tasks and run them as they come in.
// After calling `runner.Start()`, make sure to also `defer runner.Close()`
func (r *Runner) Start() {
	r.progress.Start()
	go r.queueWorker()
}

// Close stops the runner by closing its queue channel. Calls to `runner.RunLinter()` after `Close()` will panic.
// While running, yet uncompleted tasks may still complete, note that neither this method nor the queue worker
// will wait for them to be done. Parked tasks will not be run.
func (r *Runner) Close() {
	r.progress.Stop()
	close(r.queue)
}

// RunLinter creates a task to run an api.Linter on a project.This method does not block.
// The task will be executed in parallel with other tasks by the runner.
// If the runner is `nil`, then the linter will be executed directly on the current thread, returning its result the usual way.
//
// Once the task completes, the linter's report and error will be sent to the task's `Result` channel,
// i.e. use `<-task.Result` to await the linter's result.
func (r *Runner) RunLinter(id string, linter api.Linter, project api.Project) *RunnerTask {
	result := make(chan LinterResult, 1)
	task := RunnerTask{id, linter, project, result}

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

	for {
		select {
		case task, open := <-r.queue:
			if !open {
				fmt.Println("Closing")
				return
			}

			if r.nRunning >= int32(runtime.NumCPU()) {
				parked = append(parked, task)
				break
			}

			r.runTask(task)

		case task := <-r.done:
			r.nRunning--
			r.running.Delete(task.Id)
			r.progress.Done(task)
			r.progress.Update(&r.running, r.nRunning)

			if len(parked) == 0 {
				break
			}

			var next *RunnerTask
			next, parked = parked[0], parked[1:]
			r.runTask(next)
		}
	}
}

func (r *Runner) runTask(task *RunnerTask) {
	r.running.Store(task.Id, task)
	r.nRunning++
	r.progress.Update(&r.running, r.nRunning)

	go func() {
		if l, ok := task.Linter.(WithRunner); ok {
			l.SetRunner(r)
		}

		report, err := task.Linter.LintProject(task.Project)
		task.Result <- LinterResult{Report: report, Err: err}
		r.done <- task
	}()
}
