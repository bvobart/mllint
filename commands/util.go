package commands

import (
	"fmt"
	"os"
	"path"

	"github.com/bvobart/mllint/utils"
	"github.com/fatih/color"
)

// returns the current dir if args is empty, or the absolute path to the folder pointed to by args[0]
func parseProjectDir(args []string) (string, error) {
	currentdir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if len(args) == 0 {
		return currentdir, nil
	}

	projectdir := path.Join(currentdir, args[0])
	if !utils.FolderExists(projectdir) {
		return "", fmt.Errorf("%w: %s", ErrNotAFolder, projectdir)
	}

	return projectdir, nil
}

func formatInlineCode(text string) string {
	return color.New(color.Reset, color.Italic, color.FgYellow).Sprint(text)
}
