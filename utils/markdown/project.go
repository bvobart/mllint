package markdown

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/linters"
)

// FromReport creates an ML project report formatted as a Markdown string
func FromProject(project api.ProjectReport) string {
	output := strings.Builder{}
	writeProjectHeader(&output, project)
	writeProjectReports(&output, project.Reports)
	writeProjectErrors(&output, project.Errors)
	return output.String()
}

func writeProjectHeader(output *strings.Builder, project api.ProjectReport) {
	output.WriteString("# ML Project Report\n")
	output.WriteString("Project | Details\n")
	output.WriteString("--------|--------\n")
	output.WriteString("Path    | `" + project.Dir + "`\n")
	output.WriteString("Config  | `" + project.ConfigType.String() + "`\n")
	output.WriteString(fmt.Sprintf("Date    | %s \n", time.Now().Format(time.RFC1123Z)))
	output.WriteString(fmt.Sprintf("Number of Python files | %d\n", len(project.PythonFiles)))
	output.WriteString(fmt.Sprintf("Lines of Python code | %d\n", project.PythonFiles.CountLoC()))
	output.WriteString("\n---\n\n")
}

func writeProjectReports(output *strings.Builder, reports map[api.Category]api.Report) {
	output.WriteString("## Reports\n\n")
	for _, category := range categories.All {
		// check that a linter is implemented for this category
		linter, ok := linters.ByCategory[category]
		if !ok {
			continue
		}

		// check that a report was produced for this category
		report, ok := reports[category]
		if !ok {
			continue
		}

		// if so, write the category's report to the output
		writeCategoryReport(output, category, linter, report)
	}
}

func writeCategoryReport(output *strings.Builder, category api.Category, linter api.Linter, report api.Report) {
	output.WriteString(fmt.Sprintln("###", category.Name, "(`"+category.Slug+"`)"))
	output.WriteString("\n")
	output.WriteString("Passed | Score | Weight | Rule | Slug\n")
	output.WriteString(":-----:|------:|-------:|------|-----\n")

	details := strings.Builder{}
	for _, rule := range linter.Rules() {
		// check if the rule was scored and wasn't disabled
		score, ok := report.Scores[*rule]
		if !ok || rule.Disabled {
			continue
		}
		writeRuleScore(output, category, *rule, score)

		// include any details for the rule if the linter decided to report any.
		if linterDetails, ok := report.Details[*rule]; ok {
			writeRuleDetails(&details, *rule, linterDetails)
		}
	}

	output.WriteString("\n")
	output.WriteString(details.String())
}

func writeRuleScore(output *strings.Builder, category api.Category, rule api.Rule, score float64) {
	passed := "✅"
	if score < 100 {
		passed = "❌"
	}

	line := fmt.Sprintf("%s | %.1f%% | %.0f | %s | %s\n", passed, score, rule.Weight, rule.Name, rule.FullSlug(category))
	output.WriteString(line)
}

func writeRuleDetails(output *strings.Builder, rule api.Rule, details string) {
	output.WriteString("#### " + rule.Name + "\n\n")
	output.WriteString(details)
	output.WriteString("\n\n")
}

func writeProjectErrors(output *strings.Builder, multiErr *multierror.Error) {
	if multiErr == nil {
		return
	}

	output.WriteString("## Errors\n\n")

	multiErr.ErrorFormat = func(errors []error) string {
		b := strings.Builder{}
		b.WriteString(fmt.Sprintln(len(errors), "error(s) occurred while analysing your project:"))
		for _, err := range errors {
			b.WriteString(fmt.Sprintln("- ❌", err))
		}
		return b.String()
	}

	output.WriteString(fmt.Sprint(multiErr))
}
