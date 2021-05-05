package template

import (
	"errors"

	"github.com/bvobart/mllint/api"
)

// Returns an api.Linter or an api.ConfigurableLinter if also implementing api.Configurable
func NewLinter() api.Linter {
	return &Linter{}
}

// Your linter object. Give this a nice name.
type Linter struct{}

func (l *Linter) Name() string {
	return "Linter Template"
}

func (l *Linter) Rules() []*api.Rule {
	return []*api.Rule{&RuleSomething} // add all the rules that your linter may report on here.
}

func (l *Linter) LintProject(project api.Project) (api.Report, error) {
	report := api.NewReport()

	// Implement me by doing your checks and appending scores and details to the report,
	// for example:
	report.Scores[RuleSomething] = 80
	report.Details[RuleSomething] = "Something is not quite right, here are more details / this is how to fix it"

	return report, errors.New("not implemented") // replace with nil when implemented
}
