package api

import "fmt"

type CQLinterType string

type CQLinter interface {
	// String should return the human-text-friendly name of this linter
	fmt.Stringer

	Type() CQLinterType

	// DependencyName returns the name of the PyPI package that implements this linter
	DependencyName() string

	// IsConfigured returns true if there is a configuration for this linter in the given project.
	IsConfigured(project Project) bool

	// IsInstalled returns true if the linter is installed (e.g. its executable is on PATH),
	// such that Run() can be called without errorring.
	IsInstalled() bool

	// Run runs the linter, collects the issues that it reports and returns them, or an error if that failed.
	Run(projectdir string) ([]CQLinterResult, error)
}

type CQLinterResult interface {
	fmt.Stringer
}
