package bandit

import (
	"fmt"
	"strconv"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/cqlinters"
	"github.com/bvobart/mllint/utils/markdowngen"
)

func NewLinter() api.Linter {
	return &BanditLinter{}
}

type BanditLinter struct{}

func (l *BanditLinter) Name() string {
	return "Bandit"
}

func (l *BanditLinter) Rules() []*api.Rule {
	return []*api.Rule{&RuleNoIssues}
}

func (l *BanditLinter) LintProject(project api.Project) (api.Report, error) {
	report := api.NewReport()
	linter := cqlinters.ByType[cqlinters.TypeBandit]

	if RuleNoIssues.Disabled {
		return report, nil
	}

	if !linter.IsInstalled() {
		report.Scores[RuleNoIssues] = 0
		report.Details[RuleNoIssues] = fmt.Sprint("Error: ", linter, " is not installed, so it could not be run.")
		return report, nil
	}

	if len(project.PythonFiles) == 0 {
		report.Scores[RuleNoIssues] = 100
		report.Details[RuleNoIssues] = "No Python files were found in the project's repository"
	}

	results, err := linter.Run(project)
	if err != nil {
		return report, fmt.Errorf("Bandit failed to run: %w", err)
	}

	if len(results) == 0 {
		report.Scores[RuleNoIssues] = 100
		report.Details[RuleNoIssues] = "Congratulations, Bandit is happy with your project!"
	} else {
		report.Scores[RuleNoIssues] = 0
		report.Details[RuleNoIssues] = "Bandit reported " + strconv.Itoa(len(results)) + " issues with your project:\n\n" + markdowngen.List(asInterfaceList(results))
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
