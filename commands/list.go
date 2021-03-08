package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gitlab.com/bvobart/mllint/projectlinters"
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
		Use:          "all",
		Short:        "Lists all available linting rules",
		Long:         "Lists all available linting rules",
		RunE:         listAll,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
	}
}

func NewListEnabledCommand() *cobra.Command {
	return &cobra.Command{
		Use:          "enabled",
		Short:        "Lists all enabled linting rules",
		Long:         "Lists all enabled linting rules",
		RunE:         listEnabled,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
	}
}

func listAll(_ *cobra.Command, args []string) error {
	linters := projectlinters.GetAllLinters()
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
	return nil
}

func listEnabled(_ *cobra.Command, args []string) error {
	return nil
}
