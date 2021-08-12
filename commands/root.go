package commands

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func Execute() error {
	startTime := time.Now()
	err := NewRootCommand().Execute()
	if err != nil {
		color.Red("Error: %s", err)
	}
	shush(func() { fmt.Println("took:", time.Since(startTime)) })
	return err
}

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           formatInlineCode("mllint") + " [dir]",
		Short:         "Machine Learning project linter",
		Long:          formatInlineCode("mllint") + " is a command-line utility to evaluate the quality of Machine Learning (ML) projects by means of static analysis of the project's repository.",
		RunE:          runRoot,
		Args:          cobra.ArbitraryArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	SetQuietFlag(cmd)
	SetOutputFlag(cmd)
	SetForceFlag(cmd)
	SetProgressPlainFlag(cmd)

	cmd.AddCommand(NewRunCommand())
	cmd.AddCommand(NewListCommand())
	cmd.AddCommand(NewConfigCommand())
	cmd.AddCommand(NewRenderCommand())
	cmd.AddCommand(NewVersionCommand())
	cmd.AddCommand(NewDescribeCommand())
	return cmd
}

func runRoot(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("accepts at most 1 arg, received %d", len(args))
	}

	runner := runCommand{}
	return runner.RunLint(cmd, args)
}
