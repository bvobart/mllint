package commands

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/linters"
	"github.com/bvobart/mllint/utils/markdown"
	"github.com/bvobart/mllint/utils/markdowngen"
)

func NewDescribeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe RULE...",
		Short: "Describe an " + formatInlineCode("mllint") + "category or rule by its slug.",
		Long: fmt.Sprintf(`Describe an %s category or rule by its slug.
The slug is the lowercased, dashed reference string that every category and rule have. %s often displays these together.

To list all rules and their respective slugs, use %s`, formatInlineCode("mllint"), formatInlineCode("mllint"), formatInlineCode("mllint list all")),
		RunE:          describe,
		Args:          cobra.MinimumNArgs(1),
		ValidArgs:     collectAllSlugs(),
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	SetOutputFlag(cmd)
	SetForceFlag(cmd)
	return cmd
}

func describe(cmd *cobra.Command, args []string) error {
	if err := checkOutputFlag(); err != nil {
		return err
	}

	output := strings.Builder{}
	for i, slug := range args {
		if i > 0 {
			output.WriteString("\n")
		}

		if cat, ok := categories.BySlug[slug]; ok {
			output.WriteString(describeCategory(cat))
		} else if rules := linters.FindRules(slug); len(rules) > 0 {
			output.WriteString(describeRules(rules))
		} else {
			output.WriteString(color.RedString("No rule or category found that matched: %s\n", color.Set(color.Reset).Sprint(slug)))
		}

		if outputToFile() || outputToStdout() {
			if i < len(args)-1 {
				output.WriteString("\n---\n")
			}
		} else {
			color.Green(markdown.Render("---"))
		}
	}

	if outputToFile() {
		return writeToOutputFile(output.String())
	}

	if outputToStdout() {
		fmt.Println(output.String())
		return nil
	}

	return nil
}

func describeCategory(cat api.Category) string {
	if outputToFile() || outputToStdout() {
		linter := linters.ByCategory[cat]
		return markdowngen.CategoryDetails(cat, linter)
	}

	prettyPrintCategory(cat)
	color.New(color.Faint).Println("Category")
	fmt.Println()
	fmt.Println(markdown.Render(cat.Description))

	color.New(color.Bold).Println("Rules")
	linter := linters.ByCategory[cat]
	prettyPrintLinter(linter)

	fmt.Println()
	return ""
}

func describeRules(rules []*api.Rule) string {
	output := strings.Builder{}
	for i, rule := range rules {
		output.WriteString(describeRule(*rule))

		if outputToFile() || outputToStdout() {
			output.WriteString("\n---\n\n")
		} else if i < len(rules)-1 {
			color.Green(markdown.Render("---"))
		}
	}
	return output.String()
}

func describeRule(rule api.Rule) string {
	if outputToFile() || outputToStdout() {
		return markdowngen.RuleDetails(rule)
	}
	prettyPrintRule(rule)
	return ""
}

func collectAllSlugs() []string {
	slugs := []string{}
	for _, linter := range linters.ByCategory {
		rules := linter.Rules()
		for _, rule := range rules {
			slugs = append(slugs, rule.Slug)
		}
	}
	return slugs
}
