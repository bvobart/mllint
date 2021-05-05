package black

import (
	"fmt"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/cqlinters"
	"github.com/bvobart/mllint/utils/markdowngen"
)

func NewLinter() api.Linter {
	return &BlackLinter{}
}

type BlackLinter struct{}

func (l *BlackLinter) Name() string {
	return "Black"
}

func (l *BlackLinter) Rules() []*api.Rule {
	return []*api.Rule{&RuleNoIssues}
}

func (l *BlackLinter) LintProject(project api.Project) (api.Report, error) {
	report := api.NewReport()
	linter := cqlinters.ByType[cqlinters.TypeBlack]

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
		return report, fmt.Errorf("Black failed to run: %w", err)
	}

	if len(results) == 0 {
		report.Scores[RuleNoIssues] = 100
		report.Details[RuleNoIssues] = "Congratulations, Black is happy with your project!"
	} else {
		report.Scores[RuleNoIssues] = 0
		report.Details[RuleNoIssues] = "Black reported about your project that it ...\n\n" + markdowngen.CodeBlock(results[0].String()) +
			"\nBlack can fix these issues automatically when you run `black .` in your project."
	}

	return report, nil
}
