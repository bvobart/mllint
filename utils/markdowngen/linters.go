package markdowngen

import (
	"fmt"
	"strings"

	"github.com/bvobart/mllint/api"
)

func LinterRules(linter api.Linter) string {
	builder := strings.Builder{}
	for _, rule := range linter.Rules() {
		builder.WriteString(fmt.Sprintf("- `%s` — %s\n", rule.Slug, rule.Name))
	}
	return builder.String()
}

func LintersOverview(linters map[api.Category]api.Linter) string {
	builder := strings.Builder{}
	for cat, linter := range linters {
		builder.WriteString(fmt.Sprintf("## %s (`%s`)\n\n", cat.Name, cat.Slug))
		builder.WriteString(LinterRules(linter))
		builder.WriteString("\n")
	}
	return builder.String()
}

func CategoryDetails(cat api.Category, linter api.Linter) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("## Category — %s (`%s`)\n\n", cat.Name, cat.Slug))
	builder.WriteString(cat.Description)
	builder.WriteString("\n\n")
	builder.WriteString(fmt.Sprintf("### Rules — %s (`%s`)\n\n", cat.Name, cat.Slug))
	builder.WriteString(LinterRules(linter))
	return builder.String()
}

func RuleDetails(rule api.Rule) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("## Rule — %s (`%s`)\n\n", rule.Name, rule.Slug))
	builder.WriteString(fmt.Sprintf("> Weight: %.0f\n\n", rule.Weight))
	builder.WriteString(rule.Details)
	builder.WriteString("\n")
	return builder.String()
}
