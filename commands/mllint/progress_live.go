package mllint

import (
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
)

// LiveRunnerProgress is used by an `mllint.Runner` to keep track of and pretty-print the progress of the runner in running its tasks.
//
// Unless you're changing the implementation of `mllint.Runner`, you probably don't need to interact with this.
type LiveRunnerProgress struct {
	w *uilive.Writer

	print   chan struct{}
	stopped chan struct{}

	// The following properties are protected by this mutex
	mu sync.RWMutex

	// list of all tasks that have been run / scheduled with their status (running or done)
	tasks []taskStatus
	// maps the tasks' IDs to their index in `tasks`, since iterating over a map in Go does not happen in order of insertion and I don't want to do an O(n) search through `tasks` when updating a task's status.
	taskIndexes map[string]int
}

func NewLiveRunnerProgress() RunnerProgress {
	writer := uilive.New()
	writer.RefreshInterval = time.Hour
	return &LiveRunnerProgress{
		w:           writer,
		print:       make(chan struct{}, queueSize),
		stopped:     make(chan struct{}),
		mu:          sync.RWMutex{},
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
	p.mu.Lock()

	p.taskIndexes[task.Id] = len(p.tasks)
	p.tasks = append(p.tasks, taskStatus{task, statusRunning, 0})

	p.mu.Unlock()

	p.print <- struct{}{}
}

func (p *LiveRunnerProgress) TaskAwaiting(task *RunnerTask) {
	p.updateTaskStatus(task, statusAwaiting, time.Since(task.startTime))
	p.print <- struct{}{}
}

func (p *LiveRunnerProgress) TaskResuming(task *RunnerTask) {
	p.updateTaskStatus(task, statusRunning, time.Since(task.startTime))
	p.print <- struct{}{}
}

// CompletedTask is the way for the `mllint.Runner` to signal that it has completed running a task.
func (p *LiveRunnerProgress) CompletedTask(task *RunnerTask) {
	p.updateTaskStatus(task, statusDone, time.Since(task.startTime))
	p.print <- struct{}{}
}

func (p *LiveRunnerProgress) updateTaskStatus(task *RunnerTask, newStatus status, timeRunning time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	index, found := p.taskIndexes[task.Id]
	if !found {
		return
	}

	taskStatus := p.tasks[index]
	taskStatus.Status = newStatus
	taskStatus.TimeRunning = timeRunning
	p.tasks[index] = taskStatus
}

// AllTasksDone is the way for the `mllint.Runner` to signal that it has finished running all tasks,
// and that it won't call p.CompletedTask anymore (if it does, it panics because `p.done` is closed).
// This method will wait until the printWorker has finished printing and has shutdown.
func (p *LiveRunnerProgress) AllTasksDone() {
	close(p.print)
	<-p.stopped
}

// waits for signals on p.print to print the current list of tasks.
func (p *LiveRunnerProgress) printWorker() {
	for {
		_, open := <-p.print
		if !open {
			p.w.Stop()
			close(p.stopped)
			return
		}

		p.printTasks()
	}
}

func (p *LiveRunnerProgress) printTasks() {
	firstLine := p.w.Newline()

	p.mu.RLock() // lock

	nDone := 0
	nTotal := len(p.tasks)
	for _, task := range p.tasks {
		if task.Status == statusDone {
			nDone++
		}

		writer := p.w.Newline()
		task.PrintStatus(writer)
	}

	p.mu.RUnlock() // unlock

	if nDone == nTotal {
		fmt.Fprintln(p.w.Newline())
		color.New(color.Bold, color.FgGreen).Fprintln(p.w.Newline(), "✔️ All done!")
		fmt.Fprint(p.w.Newline(), "\n---\n\n")
	} else {
		fmt.Fprintln(p.w.Newline())
		color.New(color.Bold).Fprintf(firstLine, "Analysing project... (%d/%d)\n", nDone, nTotal)
		fmt.Fprintln(p.w.Newline())
	}

	p.w.Flush()
}
