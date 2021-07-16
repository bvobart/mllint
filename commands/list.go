package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/linters"
	"github.com/bvobart/mllint/utils/markdowngen"
)

func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list all|enabled",
		Short: "Lists all available or all enabled linting rules for each category",
		Long: fmt.Sprintf(`For each category of evaluation, this command lists all available or all enabled linting rules that will be used to analyse a project.

All rules are enabled by default, but if desired, it is possible to disable rules in your project's %s configuration.
See %s's ReadMe or run %s to learn more about configuring %s.`, formatInlineCode("mllint"), formatInlineCode("mllint"), formatInlineCode("mllint help config"), formatInlineCode("mllint")),
	}
	SetOutputFlag(cmd)
	SetForceFlag(cmd)
	cmd.AddCommand(NewListAllCommand())
	cmd.AddCommand(NewListEnabledCommand())
	return cmd
}

func NewListAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all",
		Short: "Lists all available linting rules for each category",
		Long:  "Lists all available linting rules for each category",
		RunE:  listAll,
		Args:  cobra.NoArgs,
	}
	SetOutputFlag(cmd)
	SetForceFlag(cmd)
	return cmd
}

func NewListEnabledCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enabled [dir]",
		Short: "Lists all enabled linting rules in the current project for each category",
		Long: fmt.Sprintf(`For each category, lists all enabled linting rules in the project in the given directory, or the current directory if none was given.

All rules are enabled by default, but if desired, it is possible to disable rules in your project's %s configuration.
See %s's ReadMe or run %s to learn more about configuring %s.`, formatInlineCode("mllint"), formatInlineCode("mllint"), formatInlineCode("mllint help config"), formatInlineCode("mllint")),
		RunE: listEnabled,
		Args: cobra.MaximumNArgs(1),
	}
	SetOutputFlag(cmd)
	SetForceFlag(cmd)
	return cmd
}

func listAll(_ *cobra.Command, args []string) error {
	if err := checkOutputFlag(); err != nil {
		return err
	}

	conf, conftype, err := config.ParseFromDir(".")
	if err == nil && len(conf.Rules.Custom) > 0 {
		shush(func() {
			color.Green("Including %d custom rules from the %s file in the current directory\n\n", len(conf.Rules.Custom), conftype.String())
		})

		if err := linters.ConfigureAll(conf); err != nil {
			color.HiYellow("Warning! The mllint configuration file %s contains an error: %s\n\n", conftype.String(), err.Error())
		}
	}

	return listLinters(linters.ByCategory)
}

func listEnabled(_ *cobra.Command, args []string) error {
	if err := checkOutputFlag(); err != nil {
		return err
	}

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

	if err := linters.ConfigureAll(conf); err != nil {
		return err
	}
	linters.DisableAll(conf.Rules.Disabled)

	if err := listLinters(linters.ByCategory); err != nil {
		return err
	}

	shush(func() { fmt.Println("---") })
	return nil
}

func listLinters(linters map[api.Category]api.Linter) error {
	if outputToFile() || outputToStdout() {
		md := markdowngen.LintersOverview(linters)

		if outputToFile() {
			return writeToOutputFile(md)
		}

		if outputToStdout() {
			fmt.Println(md)
			return nil
		}
	}

	prettyPrintLinters(linters)
	return nil
}
