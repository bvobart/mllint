package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/linters"
	"github.com/bvobart/mllint/utils/markdown"
)

func NewDescribeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe RULE...",
		Short: "Describe an mllint category or rule by its slug.",
		Long: `Describe an mllint category or rule by its slug.
The slug is the lowercased, dashed reference string that every category and rule have. mllint often displays these together.
To list all rules and their slugs, use 'mllint list all'`,
		RunE:          describe,
		Args:          cobra.MinimumNArgs(1),
		ValidArgs:     collectAllSlugs(),
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	return cmd
}

func describe(cmd *cobra.Command, args []string) error {
	for i, slug := range args {
		if i > 0 {
			fmt.Println()
		}

		if cat, ok := categories.BySlug[slug]; ok {
			describeCategory(cat)
		} else if rules := linters.FindRules(slug); len(rules) > 0 {
			describeRules(rules)
		} else {
			color.Red("No rule or category found that matched: %s", color.Set(color.Reset).Sprint(slug))
		}

		color.Green(markdown.Render("---"))
	}
	return nil
}

func describeCategory(cat api.Category) {
	prettyPrintCategory(cat)
	color.New(color.Faint).Println("Category")
	fmt.Println()
	fmt.Println(markdown.Render(cat.Description))

	color.New(color.Bold).Println("Rules")
	linter := linters.ByCategory[cat]
	prettyPrintLinter(linter)

	fmt.Println()
}

func describeRules(rules []*api.Rule) {
	for i, rule := range rules {
		describeRule(*rule)

		if i < len(rules)-1 {
			color.Green(markdown.Render("---"))
		}
	}
}

func describeRule(rule api.Rule) {
	// TODO: be able to print Markdown output when using -o
	prettyPrintRule(rule)
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
