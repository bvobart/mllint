package mllint

import (
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/fatih/color"

	"github.com/bvobart/mllint/api"
)

const queueSize = 20 // arbitrarily chosen

func NewRunner() *Runner {
	return &Runner{
		queue:   make(chan *RunnerTask, queueSize),
		running: sync.Map{},
		done:    make(chan *RunnerTask, runtime.NumCPU()),
	}
}

type Runner struct {
	queue    chan *RunnerTask
	running  sync.Map
	nRunning int32
	done     chan *RunnerTask
}

type RunnerTask struct {
	Id      string
	Linter  api.Linter
	Project api.Project
	Result  chan LinterResult
}

type LinterResult struct {
	api.Report
	Err error
}

func (r *Runner) Start() {
	go r.queueWorker()
}

func (r *Runner) Close() {
	close(r.queue)
}

func (r *Runner) RunLinter(id string, linter api.Linter, project api.Project) *RunnerTask {
	result := make(chan LinterResult, 1)
	task := RunnerTask{id, linter, project, result}
	r.queue <- &task
	return &task
}

func (r *Runner) queueWorker() {
	parked := []*RunnerTask{}

	for {
		select {
		case task, open := <-r.queue:
			if !open {
				fmt.Println("Closing")
				r.PrintRunning()
				return
			}

			if r.nRunning >= int32(runtime.NumCPU()) {
				parked = append(parked, task)
				color.Blue("Scheduled: %s, Running: %d", task.Linter.Name(), r.nRunning)
				break
			}

			r.runTask(task)

		case task := <-r.done:
			r.nRunning--
			r.running.Delete(task.Id)

			color.Green("Done: %s, Running: %d", task.Linter.Name(), r.nRunning)
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
	color.Yellow("Running: %s, Running: %d", task.Linter.Name(), r.nRunning)
	r.running.Store(task.Id, task)
	r.nRunning++

	go func() {
		report, err := task.Linter.LintProject(task.Project)
		task.Result <- LinterResult{Report: report, Err: err}
		r.done <- task
	}()
}

func (r *Runner) PrintRunning() {
	fmt.Print(r.PrintRunningString())
}

func (r *Runner) PrintRunningString() string {
	builder := strings.Builder{}

	r.running.Range(func(key, value interface{}) bool {
		id := key.(string)
		task := value.(*RunnerTask)
		builder.WriteString(fmt.Sprintf("%s - %s\n", id, task.Linter.Name()))
		return true
	})

	return builder.String()
}
