package mllint

import (
	"io"

	"github.com/fatih/color"
)

type RunnerProgress interface {
	Start()
	RunningTask(task *RunnerTask)
	CompletedTask(task *RunnerTask)
	AllTasksDone()
}

type status string

const (
	statusRunning status = "⏳ Running -"
	statusDone    status = "✔️ Done -"
)

type taskStatus struct {
	*RunnerTask
	Status status
}

func (s taskStatus) PrintStatus(writer io.Writer) {
	if s.Status == statusRunning {
		color.New(color.FgYellow).Fprintln(writer, s.Status, s.displayName)
	}
	if s.Status == statusDone {
		color.New(color.FgGreen).Fprintln(writer, s.Status, s.displayName)
	}
}
