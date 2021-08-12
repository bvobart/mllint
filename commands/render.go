package commands

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/bvobart/mllint/utils"
	"github.com/bvobart/mllint/utils/markdown"
	"github.com/spf13/cobra"
)

func NewRenderCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "render FILE",
		Short:         "Render an " + formatInlineCode("mllint") + " report to your terminal.",
		Long:          fmt.Sprintf(`Renders an %s report to your terminal in a pretty way.`, formatInlineCode("mllint")),
		RunE:          render,
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	return cmd
}

func render(cmd *cobra.Command, args []string) error {
	filename := args[0]
	if !utils.FileExists(filename) {
		return fmt.Errorf("cannot find file: %s", filename)
	}

	if path.Ext(filename) != ".md" {
		return fmt.Errorf("mllint can only render Markdown files, but the provided filename does not end with '.md': %s", filename)
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cannot read from the given file: %w", err)
	}

	fmt.Println(markdown.Render(string(contents)))
	return nil
}
