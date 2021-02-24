package api

type Linter interface {
	Name() string
	LintProject(projectdir string) ([]Issue, error)
}
