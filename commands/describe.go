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

		fmt.Println("---")
	}
	return nil
}

func describeCategory(cat api.Category) {
	prettyPrintCategory(cat)
	color.Set(color.Faint).Println("Category")
	color.Unset()

	fmt.Println()
	fmt.Println(markdown.Render(cat.Description))

	color.Set(color.Bold).Println("Rules")
	color.Unset()
	linter := linters.ByCategory[cat]
	prettyPrintLinterRules(cat, linter)

	fmt.Println()
}

func describeRules(rules []*api.Rule) {
	for i, rule := range rules {
		if i > 0 {
			fmt.Println()
		}

		describeRule(*rule)

		if i < len(rules)-1 {
			fmt.Println("---")
		}
	}
}

func describeRule(rule api.Rule) {
	// TODO: just create Markdown output and pretty-print that.
	bold := color.New(color.Bold)
	faint := color.New(color.Faint)

	bold.Print(rule.Name, " ")
	faint.Println(rule.Slug)
	faint.Println("Rule")
	faint.Printf("Weight: %.0f\n", rule.Weight)

	fmt.Println()
	fmt.Println(markdown.Render(rule.Details))
}

func collectAllSlugs() []string {
	slugs := []string{}
	for cat, linter := range linters.ByCategory {
		rules := linter.Rules()
		for _, rule := range rules {
			slugs = append(slugs, cat.Slug+"/"+rule.Slug)
		}
	}
	return slugs
}
