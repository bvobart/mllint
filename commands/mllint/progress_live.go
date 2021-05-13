package mllint

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
)

// LiveRunnerProgress is used by an `mllint.Runner` to keep track of and pretty-print the progress of the runner in running its tasks.
//
// Unless you're changing the implementation of `mllint.Runner`, you probably don't need to interact with this.
type LiveRunnerProgress struct {
	w *uilive.Writer

	running chan *RunnerTask
	done    chan *RunnerTask

	stopped chan struct{}

	// The following properties may only be edited inside p.printWorker()

	// list of all tasks that have been run / scheduled
	tasks []taskStatus
	// maps the tasks' IDs to their index in `tasks`, since iterating over a map in Go does not happen in order of insertion and I don't want to do an O(n) search through `tasks` when updating a task's status.
	taskIndexes map[string]int
}

func NewLiveRunnerProgress() RunnerProgress {
	writer := uilive.New()
	writer.RefreshInterval = time.Hour
	return &LiveRunnerProgress{
		w:           writer,
		running:     make(chan *RunnerTask, queueSize),
		done:        make(chan *RunnerTask, queueSize),
		stopped:     make(chan struct{}),
		tasks:       []taskStatus{},
		taskIndexes: make(map[string]int),
	}
}

// Start starts the printWorker process on a new go-routine.
func (p *LiveRunnerProgress) Start() {
	p.w.Start()
	go p.printWorker()
}

// RunningTask is the way for the `mllint.Runner` to signal that it has started running a task.
func (p *LiveRunnerProgress) RunningTask(task *RunnerTask) {
	p.running <- task
}

// CompletedTask is the way for the `mllint.Runner` to signal that it has completed running a task.
func (p *LiveRunnerProgress) CompletedTask(task *RunnerTask) {
	p.done <- task
}

// AllTasksDone is the way for the `mllint.Runner` to signal that it has finished running all tasks,
// and that it won't call p.CompletedTask anymore (if it does, it panics because `p.done` is closed).
// This method will wait until the printWorker has finished printing and has shutdown.
func (p *LiveRunnerProgress) AllTasksDone() {
	close(p.done)
	<-p.stopped
}

func (p *LiveRunnerProgress) printWorker() {
	for {
		select {
		case task, open := <-p.running:
			if !open {
				break
			}

			p.onTaskRunning(task)
		case task, open := <-p.done:
			if !open {
				p.w.Stop()
				close(p.stopped)
				return
			}

			p.onTaskDone(task)
		}
	}
}

func (p *LiveRunnerProgress) onTaskRunning(task *RunnerTask) {
	p.tasks = append(p.tasks, taskStatus{task, statusRunning})
	p.taskIndexes[task.Id] = len(p.tasks) - 1

	p.printTasks()
}

func (p *LiveRunnerProgress) onTaskDone(task *RunnerTask) {
	index, found := p.taskIndexes[task.Id]
	if !found {
		return
	}

	status := p.tasks[index]
	status.Status = statusDone
	p.tasks[index] = status

	p.printTasks()
}

func (p *LiveRunnerProgress) printTasks() {
	if len(p.tasks) == 0 {
		color.New(color.Bold).Fprintln(p.w, "Waiting for tasks...")
		p.w.Flush()
		return
	}

	allDone := true
	for _, task := range p.tasks {
		if task.Status != statusDone {
			allDone = false
		}

		writer := p.w.Newline()
		task.PrintStatus(writer)
	}

	if allDone {
		fmt.Fprintln(p.w.Newline())
		color.New(color.Bold, color.FgGreen).Fprintln(p.w.Newline(), "✔️ All done!")
		fmt.Fprint(p.w.Newline(), "\n---\n\n")
	}

	p.w.Flush()
}
