package commands

import (
	"fmt"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/linters"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewDescribeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe RULE...",
		Short: "Describe an mllint category or rule by its slug.",
		Long: `Describe an mllint category or rule by its slug.
The slug is the lowercased, dashed reference string that every category and rule have. mllint often displays these together.`,
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
		} else if rule, ok := linters.GetRule(slug); ok {
			describeRule(rule)
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
	fmt.Println(cat.Description)
	fmt.Println()

	color.Set(color.Bold).Println("Rules")
	color.Unset()
	linter := linters.ByCategory[cat]
	prettyPrintLinterRules(cat, linter)

	fmt.Println()
}

func describeRule(rule api.Rule) {
	color.Set(color.Bold).Print(rule.Name)
	color.Unset()
	fmt.Print(" ")
	cfmt := color.Set(color.Faint)
	cfmt.Printf("(%s)\nRule\n", rule.Slug)
	cfmt.Printf("Weight: %.0f\n", rule.Weight)
	color.Unset()

	fmt.Println()
	fmt.Println(rule.Details)
	fmt.Println()
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
