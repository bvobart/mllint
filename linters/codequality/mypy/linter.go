package mypy

import (
	"fmt"
	"math"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/cqlinters"
	"github.com/bvobart/mllint/utils/markdowngen"
)

// Maximum number of lines of code per Mypy message reported.
// Increasing this means that users are expected to have less code smells per line of code.
const maxLoCperMsg = 10

func NewLinter() api.Linter {
	return &MypyLinter{}
}

type MypyLinter struct{}

func (l *MypyLinter) Name() string {
	return "Mypy"
}

func (l *MypyLinter) Rules() []*api.Rule {
	return []*api.Rule{&RuleNoIssues}
}

func (l *MypyLinter) LintProject(project api.Project) (api.Report, error) {
	report := api.NewReport()
	linter := cqlinters.ByType[cqlinters.TypeMypy]

	if RuleNoIssues.Disabled {
		return report, nil
	}

	if !linter.IsInstalled() {
		report.Scores[RuleNoIssues] = 0
		report.Details[RuleNoIssues] = fmt.Sprint("Error: ", linter, " is not installed, so it could not be run.")
		return report, nil
	}

	loc := project.PythonFiles.CountLoC()
	if loc == 0 {
		report.Scores[RuleNoIssues] = 100
		report.Details[RuleNoIssues] = "No Python code was found in the project's repository"
	}

	results, err := linter.Run(project)
	if err != nil {
		return report, fmt.Errorf("Mypy failed to run: %w", err)
	}

	// calculate score. No Mypy messages = 100%, 1 Mypy message per 20 lines of code = 50%, 1 Mypy message per 10 lines of code = 0%
	report.Scores[RuleNoIssues] = 100 - 100*math.Min(1, float64(len(results)*maxLoCperMsg)/float64(loc))
	if len(results) == 0 {
		report.Details[RuleNoIssues] = "Congratulations, Mypy is happy with your project!"
	} else {
		report.Details[RuleNoIssues] = fmt.Sprintf("Mypy reported **%d** issues with your project:\n\n", len(results)) + markdowngen.List(asInterfaceList(results))
	}

	return report, nil
}

func asInterfaceList(list []api.CQLinterResult) []interface{} {
	res := make([]interface{}, len(list))
	for i, item := range list {
		res[i] = item
	}
	return res
}
