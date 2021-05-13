package commands

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var quiet bool

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
		Use:           "mllint [dir]",
		Short:         "Machine Learning project linter",
		Long:          "mllint is a command-line utility to evaluate the quality of Machine Learning (ML) projects by means of static analysis of the project's repository.",
		RunE:          runRoot,
		Args:          cobra.ArbitraryArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Set this to true to minimise printing to the bare minimum.")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", `Export the report generated for your project to a Markdown file at the given location.
Set this to '-' (a single dash) in order to print the raw Markdown directly to the console.`)
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Use this flag to remove the output file provided with '--output' in case that already exists.")
	cmd.Flags().BoolVar(&progressPlain, "progress-plain", false, "Use this flag to print linting progress plainly, without rewriting terminal output. Enabled automatically in non-interactive terminals.")

	cmd.AddCommand(NewRunCommand())
	cmd.AddCommand(NewListCommand())
	cmd.AddCommand(NewConfigCommand())
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

// only execute f when quiet is nil or false.
func shush(f func()) {
	if !quiet {
		f()
	}
}
