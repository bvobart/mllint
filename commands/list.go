package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/projectlinters"
)

func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list all|enabled",
		Short: "Lists all available or all enabled linting rules",
		Long:  "Lists all available or all enabled linting rules that will be used while analysing a project. All rules are enabled by default, but you can configure rules to enable or disable in a .mllint.yml file in the root of your project folder.",
	}
	cmd.AddCommand(NewListAllCommand())
	cmd.AddCommand(NewListEnabledCommand())
	return cmd
}

func NewListAllCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "Lists all available linting rules",
		Long:  "Lists all available linting rules",
		RunE:  listAll,
		Args:  cobra.NoArgs,
	}
}

func NewListEnabledCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "enabled [dir]",
		Short: "Lists all enabled linting rules in the current project",
		Long:  "Lists all enabled linting rules in the project in the given directory, or the current directory if none was given.",
		RunE:  listEnabled,
		Args:  cobra.MaximumNArgs(1),
	}
}

func listAll(_ *cobra.Command, args []string) error {
	linters := projectlinters.GetAllLinters()
	prettyPrintRules(linters)
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

	linters, err := projectlinters.GetAllLinters().FilterEnabled(conf.Rules).Configure(conf)
	if err != nil {
		return err
	}
	prettyPrintRules(linters)

	shush(func() { fmt.Println("\n---") })
	return nil
}

func prettyPrintRules(linters []api.Linter) {
	for _, linter := range linters {
		rules := linter.Rules()
		if len(rules) == 1 {
			fmt.Println("-", color.BlueString(linter.Name()))
			continue
		}

		faintSlash := color.Set(color.Faint).Sprint("/")
		color.Unset()
		for _, rule := range rules {
			fmt.Println("-", color.BlueString(linter.Name())+faintSlash+rule)
		}
	}
}
