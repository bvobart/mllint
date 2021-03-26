package commands

import (
	"errors"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	quiet *bool
)

func Execute() error {
	startTime := time.Now()
	err := NewRootCommand().Execute()
	if err != nil && errors.Is(err, ErrIssuesFound) {
		color.HiWhite("%s", err)
	} else if err != nil {
		color.Red("Fatal: %s", err)
	}
	shush(func() { fmt.Println("took:", time.Since(startTime)) })
	return err
}

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
	quiet = cmd.PersistentFlags().BoolP("quiet", "q", false, "Set this to true to minimise printing to the bare minimum.")
	cmd.AddCommand(NewRunCommand())
	cmd.AddCommand(NewListCommand())
	cmd.AddCommand(NewConfigCommand())
	cmd.AddCommand(NewVersionCommand())
	return cmd
}

func runRoot(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("accepts at most 1 arg, received %d", len(args))
	}
	return NewRunCommand().RunE(cmd, args)
}

// only execute f when quiet is nil or false.
func shush(f func()) {
	if quiet == nil || !*quiet {
		f()
	}
}
