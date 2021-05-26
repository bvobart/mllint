package mllint

import (
	"io"
	"time"

	"github.com/bvobart/mllint/api"
)

const queueSize = 20 // arbitrarily chosen

// NewMLLintRunner initialises an *mllint.MLLintRunner
func NewMLLintRunner(progress RunnerProgress) *MLLintRunner {
	if progress == nil {
		progress = &BasicRunnerProgress{Out: io.Discard}
	}
	return &MLLintRunner{
		queue:    make(chan *RunnerTask, queueSize),
		awaiting: make(chan *RunnerTask, queueSize),
		resuming: make(chan *RunnerTask, queueSize),
		done:     make(chan *RunnerTask, queueSize),
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
type MLLintRunner struct {
	queue    chan *RunnerTask
	awaiting chan *RunnerTask
	resuming chan *RunnerTask
	done     chan *RunnerTask

	progress RunnerProgress
	closed   chan struct{}
	nRunning int32
}

// Start starts the runner by running a queue worker go-routine in the background that will await tasks and run them as they come in.
// After calling `runner.Start()`, make sure to also `defer runner.Close()`
func (r *MLLintRunner) Start() {
	r.progress.Start()
	go r.queueWorker()
}

// Close stops the runner by closing its queue channel. Calls to `runner.RunLinter()` after Close will panic.
// Close blocks until all tasks added to its queue by `runner.RunLinter()` have completed,
// and all progress output has finished printing to the terminal.
func (r *MLLintRunner) Close() {
	close(r.queue)
	<-r.closed
}

// RunLinter creates a task to run an api.Linter on a project.This method does not block.
// The task will be executed in parallel with other tasks by the runner.
// If the runner is `nil`, then the linter will be executed directly on the current thread, returning its result the usual way.
//
// Once the task completes, the linter's report and error will be sent to the task's `Result` channel,
// i.e. use `<-task.Result` to await the linter's result.
func (r *MLLintRunner) RunLinter(id string, linter api.Linter, project api.Project, options ...TaskOption) *RunnerTask {
	result := make(chan LinterResult, 1)
	task := RunnerTask{id, linter, project, result, linter.Name(), time.Now()}
	for _, optionFunc := range options {
		optionFunc(&task)
	}

	// a nil runner simply runs the task on the current thread.
	if r == nil {
		report, err := linter.LintProject(project)
		task.Result <- LinterResult{report, err}
		return &task
	}

	r.queue <- &task
	return &task
}

func (r *MLLintRunner) CollectTasks(tasks ...*RunnerTask) chan *RunnerTask {
	return collectTasks(func() {}, tasks...)
}

type childRunner struct {
	// the parent runner on which all tasks will actually be scheduled
	parent *MLLintRunner
	// the task for which this child runner was created
	task *RunnerTask
}

func (r *childRunner) RunLinter(id string, linter api.Linter, project api.Project, options ...TaskOption) *RunnerTask {
	return r.parent.RunLinter(id, linter, project, options...)
}

func (r *childRunner) CollectTasks(tasks ...*RunnerTask) chan *RunnerTask {
	if len(tasks) == 0 {
		funnel := make(chan *RunnerTask)
		close(funnel)
		return funnel
	}

	r.parent.awaiting <- r.task
	return collectTasks(func() { r.parent.resuming <- r.task }, tasks...)
}
