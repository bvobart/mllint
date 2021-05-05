package pylint

import (
	"fmt"
	"math"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/cqlinters"
	"github.com/bvobart/mllint/utils/markdowngen"
)

// Maximum number of lines of code per Pylint message reported.
// Increasing this means that users are expected to have less code smells per line of code.
const maxLoCperMsg = 10

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

	// check whether Pylint is installed so we can actually run it
	if !linter.IsInstalled() {
		report.Scores[RuleNoIssues] = 0
		report.Details[RuleNoIssues] = fmt.Sprint("Error: ", linter, " is not installed, so it could not be run.")
		return report, nil
	}

	// check if there are Python files to run Pylint on
	loc := project.PythonFiles.CountLoC()
	if loc == 0 {
		report.Scores[RuleNoIssues] = 100
		report.Details[RuleNoIssues] = "No Python code was found in the project's repository."
		return report, nil
	}

	// actually run Pylint
	results, err := linter.Run(project.Dir)
	if err != nil {
		return report, fmt.Errorf("Pylint failed to run: %w", err)
	}

	// calculate score. No Pylint messages = 100%, 1 Pylint message per 20 lines of code = 50%, 1 Pylint message per 10 lines of code = 0%
	report.Scores[RuleNoIssues] = 100 - 100*math.Min(1, float64(len(results)*maxLoCperMsg)/float64(loc))
	if len(results) == 0 {
		report.Details[RuleNoIssues] = "Congratulations, Pylint is happy with your project!"
	} else {
		report.Details[RuleNoIssues] = fmt.Sprintf("Pylint reported %d issues with your project:\n\n", len(results)) + markdowngen.List(asInterfaceList(results))
	}

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
