package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/linters"
)

func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list all|enabled",
		Short: "Lists all available or all enabled linting rules for each category",
		Long: `For each category of evaluation, this command lists all available or all enabled linting rules that will be used to analyse a project.

All rules are enabled by default, but if desired, it is possible to disable rules in your project's mllint configuration.
See mllint's ReadMe or run 'mllint help config' to learn more about configuring mllint.`,
	}
	cmd.AddCommand(NewListAllCommand())
	cmd.AddCommand(NewListEnabledCommand())
	return cmd
}

func NewListAllCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "Lists all available linting rules for each category",
		Long:  "Lists all available linting rules for each category",
		RunE:  listAll,
		Args:  cobra.NoArgs,
	}
}

func NewListEnabledCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "enabled [dir]",
		Short: "Lists all enabled linting rules in the current project for each category",
		Long: `For each category, lists all enabled linting rules in the project in the given directory, or the current directory if none was given.

All rules are enabled by default, but if desired, it is possible to disable rules in your project's mllint configuration.
See mllint's ReadMe or run 'mllint help config' to learn more about configuring mllint.`,
		RunE: listEnabled,
		Args: cobra.MaximumNArgs(1),
	}
}

func listAll(_ *cobra.Command, args []string) error {
	prettyPrintLinters(linters.ByCategory)
	return nil
}

func listEnabled(_ *cobra.Command, args []string) error {
	projectdir, err := parseProjectDir(args)
	if err != nil {
		return fmt.Errorf("invalid argument: %w", err)
	}
	shush(func() { color.Green("Using project at  %s", color.HiWhiteString(projectdir)) })

	conf, _, err := getConfig(projectdir)
	if err != nil {
		return err
	}
	shush(func() { fmt.Print("---\n\n") })

	linters.DisableAll(conf.Rules.Disabled)
	if err := linters.ConfigureAll(conf); err != nil {
		return err
	}
	prettyPrintLinters(linters.ByCategory)

	shush(func() { fmt.Println("---") })
	return nil
}

func prettyPrintLinters(linters map[api.Category]api.Linter) {
	if len(linters) == 0 {
		color.Red("Oh no! Your mllint configuration has disabled ALL rules!")
		fmt.Println()
	}

	for cat, linter := range linters {
		color.Set(color.Bold).Print(cat.Name)
		color.Unset()
		fmt.Print(" ")
		color.Set(color.Faint).Printf("(%s)\n", cat.Slug)
		color.Unset()

		for _, rule := range linter.Rules() {
			if !rule.Disabled {
				coloredSlug := color.Set(color.Faint).Sprintf("(%s/%s)", cat.Slug, rule.Slug)
				color.Unset()
				fmt.Println("-", color.BlueString(rule.Name), coloredSlug)
			}
		}

		fmt.Println()
	}
}
