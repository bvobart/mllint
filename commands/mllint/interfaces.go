package mllint

import (
	"time"

	"github.com/bvobart/mllint/api"
)

type Runner interface {
	RunLinter(id string, linter api.Linter, project api.Project, options ...TaskOption) *RunnerTask
	CollectTasks(tasks ...*RunnerTask) chan *RunnerTask
}

// TaskOption is an option for a task created by RunLinter, e.g. setting a custom display name.
type TaskOption func(task *RunnerTask)

func DisplayName(name string) TaskOption {
	return func(task *RunnerTask) {
		task.displayName = name
	}
}

//---------------------------------------------------------------------------------------

// RunnerTask represents a task to run a linter on a project that was created by a call to runner.RunLinter(...)
type RunnerTask struct {
	Id          string
	Linter      api.Linter
	Project     api.Project
	Result      chan LinterResult
	displayName string
	startTime   time.Time
}

// LinterResult represents the two-valued return type of a Linter, containing a report and an error.
type LinterResult struct {
	api.Report
	Err error
}

//---------------------------------------------------------------------------------------

type WithRunner interface {
	SetRunner(runner Runner)
}

type LinterWithRunner interface {
	api.Linter
	WithRunner
}

type ConfigurableLinterWithRunner interface {
	api.ConfigurableLinter
	WithRunner
}

//---------------------------------------------------------------------------------------
