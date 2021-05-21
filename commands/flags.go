package commands

import (
	"fmt"
	"os"

	"github.com/bvobart/mllint/utils"
	"github.com/spf13/cobra"
)

var (
	quiet         bool
	outputFile    string
	force         bool
	progressPlain bool
)

func SetQuietFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Set this to true to only print to the bare minimum.")
}

func SetOutputFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", fmt.Sprintf(`Export the generated Markdown output to a file at the given location. Note that %s does not add the %s file extension to this filename.
Set this to %s (a single dash) in order to print the raw Markdown directly to the console.`, formatInlineCode("mllint"), formatInlineCode(".md"), formatInlineCode("-")))
}

func outputToStdout() bool {
	return outputFile == "-"
}

func outputToFile() bool {
	return outputFile != "" && !outputToStdout()
}

func checkOutputFlag() error {
	if outputToFile() && utils.FileExists(outputFile) {
		if !force {
			return fmt.Errorf("%w: %s", ErrOutputFileAlreadyExists, formatInlineCode(utils.AbsolutePath(outputFile)))
		}
		if err := os.Remove(outputFile); err != nil {
			return fmt.Errorf("tried to remove %s, but got error: %w", utils.AbsolutePath(outputFile), err)
		}
	}

	if outputToStdout() {
		quiet = true
	}

	return nil
}

func SetForceFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Use this flag to remove the output file provided with "+formatInlineCode("--output")+" in case that already exists.")
}

func SetProgressPlainFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&progressPlain, "progress-plain", false, "Use this flag to print linting progress plainly, without rewriting terminal output. Enabled automatically in non-interactive terminals.")
}

// only execute f when quiet is nil or false.
func shush(f func()) {
	if !quiet {
		f()
	}
}
