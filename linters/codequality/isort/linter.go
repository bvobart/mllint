package isort

import (
	"fmt"
	"strconv"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/cqlinters"
	"github.com/bvobart/mllint/utils/markdowngen"
)

func NewLinter() api.Linter {
	return &ISortLinter{}
}

type ISortLinter struct{}

func (l *ISortLinter) Name() string {
	return "`isort`"
}

func (l *ISortLinter) Rules() []*api.Rule {
	return []*api.Rule{&RuleNoIssues, &RuleIsConfigured}
}

func (l *ISortLinter) LintProject(project api.Project) (api.Report, error) {
	report := api.NewReport()
	linter := cqlinters.ByType[cqlinters.TypeISort]

	if !RuleIsConfigured.Disabled {
		if linter.IsProperlyConfigured(project) {
			report.Scores[RuleIsConfigured] = 100
		} else {
			report.Scores[RuleIsConfigured] = 0
			report.Details[RuleIsConfigured] = DetailsNotProperlyConfigured
		}
	}

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
		return report, fmt.Errorf("isort failed to run: %w", err)
	}

	if len(results) == 0 {
		report.Scores[RuleNoIssues] = 100
		report.Details[RuleNoIssues] = "Congratulations, `isort` is happy with your project!"
	} else {
		report.Scores[RuleNoIssues] = 0
		report.Details[RuleNoIssues] = "isort reported **" + strconv.Itoa(len(results)) + "** files in your project that it would fix:\n\n" + markdowngen.List(asInterfaceList(results)) +
			"\nisort can fix these issues automatically when you run `isort .` in your project."
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
