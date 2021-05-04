package pylint

import (
	"fmt"
	"math"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/cqlinters"
	"github.com/bvobart/mllint/utils/markdowngen"
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

	results, err := linter.Run(project.Dir)
	if err != nil {
		return report, fmt.Errorf("Pylint failed to run: %w", err)
	}

	// TODO: find a way to score the specific rule for this CQlinter
	// TODO: scale this score by the number of total lines of Python code.
	report.Scores[RuleNoIssues] = 100 - math.Min(100, float64(4*len(results))) // max 25 messages.

	if len(results) == 0 {
		report.Details[RuleNoIssues] = "Congratulations, Pylint is happy with your project!"
	} else {
		report.Details[RuleNoIssues] = "Pylint reported some issues with your project:\n\n" + markdowngen.List(asInterfaceList(results))
	}

	// TODO: sort all messages by severity before displaying them.

	// TODO: find a solution for multi-line error messages such as those for duplicate code.

	// for _, result := range results {
	// 	message := result.(cqlinters.PylintMessage)
	// 	// TODO: do some specific analysis of these results
	// e.g. create rules about: duplicate code, import management, other stuff from paper.
	// }

	return report, nil
}

func asInterfaceList(list []api.CQLinterResult) []interface{} {
	res := make([]interface{}, len(list))
	for i, item := range list {
		res[i] = item
	}
	return res
}
