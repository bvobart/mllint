package mllint

import (
	"fmt"
	"io"
	"time"

	"github.com/fatih/color"
)

type RunnerProgress interface {
	Start()
	RunningTask(task *RunnerTask)
	TaskAwaiting(task *RunnerTask)
	TaskResuming(task *RunnerTask)
	CompletedTask(task *RunnerTask)
	AllTasksDone()
}

type status string

const (
	statusRunning  status = "🏃 Running -"
	statusAwaiting status = "⏳ Waiting -"
	statusDone     status = "✔️ Done -"
)

type taskStatus struct {
	*RunnerTask
	Status      status
	TimeRunning time.Duration
}

func (s taskStatus) PrintStatus(writer io.Writer) {
	msg := fmt.Sprintf("%s %s", s.Status, s.displayName)
	timeRunning := color.New(color.Faint, color.Italic).Sprint("(", formatDuration(s.TimeRunning), ")")

	switch s.Status {
	case statusRunning:
		fmt.Fprintln(writer, color.New(color.Italic, color.FgYellow).Sprint(msg), timeRunning)
	case statusAwaiting:
		fmt.Fprintln(writer, color.New(color.FgYellow).Sprint(msg), timeRunning)
	case statusDone:
		fmt.Fprintln(writer, color.New(color.FgGreen).Sprint(msg), timeRunning)
	}
}

func formatDuration(duration time.Duration) string {
	if duration.Milliseconds() == 0 {
		return fmt.Sprintf("%.2f ms", float64(duration.Microseconds())/1000)
	}
	if duration.Milliseconds() > 1000 {
		return fmt.Sprintf("%.2f s", duration.Seconds())
	}
	return fmt.Sprint(duration.Milliseconds(), " ms")
}
