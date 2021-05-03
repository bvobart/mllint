package pylint

import (
	"fmt"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/cqlinters"
)

func NewLinter() api.Linter {
	return &PylintLinter{}
}

type PylintLinter struct{}

func (l *PylintLinter) Name() string {
	return "Pylint"
}

func (l *PylintLinter) Rules() []*api.Rule {
	return []*api.Rule{}
}

func (l *PylintLinter) LintProject(project api.Project) (api.Report, error) {
	report := api.NewReport()
	linter := cqlinters.ByType[cqlinters.TypePylint]

	if !linter.IsInstalled() {
		report.Scores[RuleNoIssues] = 0
		report.Details[RuleNoIssues] = fmt.Sprint("Error: ", linter, " is not installed, so it could not be run.")
	}

	err := linter.Run(project.Dir)
	if err != nil {
		return report, fmt.Errorf("Pylint failed to run: %w", err)
	}

	// TODO: gather linting issues
	// TODO: find a way to score the specific rule for this CQlinter

	return report, nil
}
