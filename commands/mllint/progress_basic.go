package mllint

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
)

type BasicRunnerProgress struct {
	Out io.Writer
}

func NewBasicRunnerProgress() RunnerProgress {
	return &BasicRunnerProgress{os.Stdout}
}

func (p *BasicRunnerProgress) Start() {}

func (p *BasicRunnerProgress) RunningTask(task *RunnerTask) {
	taskStatus{task, statusRunning, time.Since(task.startTime)}.PrintStatus(p.Out)
}

func (p *BasicRunnerProgress) TaskAwaiting(task *RunnerTask) {
	taskStatus{task, statusAwaiting, time.Since(task.startTime)}.PrintStatus(p.Out)
}

func (p *BasicRunnerProgress) TaskResuming(task *RunnerTask) {
	taskStatus{task, statusRunning, time.Since(task.startTime)}.PrintStatus(p.Out)
}

func (p *BasicRunnerProgress) CompletedTask(task *RunnerTask) {
	taskStatus{task, statusDone, time.Since(task.startTime)}.PrintStatus(p.Out)
}

func (p *BasicRunnerProgress) AllTasksDone() {
	fmt.Fprintln(p.Out)
	color.New(color.Bold, color.FgGreen).Fprintln(p.Out, "✔️ All done!")
	fmt.Fprint(p.Out, "\n---\n\n")
}
