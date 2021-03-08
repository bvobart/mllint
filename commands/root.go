package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "mllint [dir]",
		Short:         "Machine Learning project linter",
		Long:          "mllint is a command-line utility to evaluate the quality of Machine Learning (ML) projects by means of static analysis of the project's repository.",
		RunE:          runRoot,
		Args:          cobra.ArbitraryArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.AddCommand(NewRunCommand())
	cmd.AddCommand(NewListCommand())
	return cmd
}

func runRoot(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("accepts at most 1 arg, received %d", len(args))
	}
	return NewRunCommand().RunE(cmd, args)
}
