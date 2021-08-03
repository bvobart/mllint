package markdown

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-multierror"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/linters"
)

// FromReport creates an ML project report formatted as a Markdown string
func FromProject(project api.ProjectReport) string {
	output := strings.Builder{}
	writeProjectHeader(&output, project)
	writeConfigDetails(&output, project.Config)
	writeProjectReports(&output, project.Reports)
	writeProjectErrors(&output, project.Errors)
	return output.String()
}

func humanizeBool(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

func writeProjectHeader(output *strings.Builder, project api.ProjectReport) {
	output.WriteString("# ML Project Report\n")
	output.WriteString("**Project** | **Details**\n")
	output.WriteString("--------|--------\n")
	output.WriteString(fmt.Sprintf("Date    | %s \n", time.Now().Format(time.RFC1123Z)))
	output.WriteString("Path    | `" + project.Dir + "`\n")
	output.WriteString("Config  | `" + project.ConfigType.String() + "`\n")
	configIsDefault := cmp.Equal(project.Config, *config.Default())
	output.WriteString("Default | " + humanizeBool(configIsDefault) + "\n")

	if project.Git.RemoteURL != "" {
		output.WriteString("Git: Remote URL | `" + project.Git.RemoteURL + "`\n")
		output.WriteString("Git: Commit     | `" + project.Git.Commit + "`\n")
		output.WriteString("Git: Branch     | `" + project.Git.Branch + "`\n")
		output.WriteString("Git: Dirty Workspace?  | " + humanizeBool(project.Git.Dirty) + "\n")
	}

	output.WriteString(fmt.Sprintf("Number of Python files | %d\n", len(project.PythonFiles)))
	output.WriteString(fmt.Sprintf("Lines of Python code   | %d\n", project.PythonFiles.CountLoC()))
	output.WriteString("\n---\n\n")
}

func writeConfigDetails(output *strings.Builder, config config.Config) {
	if len(config.Rules.Disabled) > 0 {
		output.WriteString("## Config\n\n")
		output.WriteString("**Note** — The following rules were disabled in `mllint`'s configuration:\n")
		for _, slug := range config.Rules.Disabled {
			output.WriteString(fmt.Sprintf("- `%s`\n", slug))
		}
		output.WriteString("\n")
	}
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
	overallScore := report.OverallScore()

	output.WriteString(fmt.Sprintf("### %s (`%s`) — **%.1f**%%\n", category.Name, category.Slug, overallScore))
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
		writeRuleScore(output, *rule, score)

		// include any details for the rule if the linter decided to report any.
		if linterDetails, ok := report.Details[*rule]; ok {
			writeRuleDetails(&details, *rule, score, linterDetails)
		}
	}

	output.WriteString(" | _Total_ | | | \n")
	output.WriteString(fmt.Sprintf("%s | **%.1f**%% | | %s | `%s`\n", getPassedEmoji(overallScore), overallScore, category.Name, category.Slug))

	output.WriteString("\n")
	output.WriteString(details.String())
}

func writeRuleScore(output *strings.Builder, rule api.Rule, score float64) {
	passed := getPassedEmoji(score)
	line := fmt.Sprintf("%s | %.1f%% | %.0f | %s | `%s`\n", passed, score, rule.Weight, rule.Name, rule.Slug)
	output.WriteString(line)
}

func writeRuleDetails(output *strings.Builder, rule api.Rule, score float64, details string) {
	passed := getPassedEmoji(score)
	output.WriteString("#### Details — " + rule.Name + " — " + passed + "\n\n")
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

func getPassedEmoji(score float64) string {
	if score < 100 {
		return "❌"
	}
	return "✅"
}
