package commands

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"gitlab.com/bvobart/mllint/api"
	"gitlab.com/bvobart/mllint/projectlinters"
	"gitlab.com/bvobart/mllint/utils"
)

var (
	ErrNotAFolder  = errors.New("not a folder")
	ErrIssuesFound = errors.New("issues found")
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "mllint [dir]",
		Short:         "Machine Learning project linter",
		Long:          "mllint is a command-line utility to evaluate the quality of Machine Learning (ML) projects by means of static analysis of the project's repository.",
		RunE:          lint,
		Args:          cobra.MaximumNArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	return cmd
}

func lint(cmd *cobra.Command, args []string) error {
	projectdir, err := parseArgs(args)
	if err != nil {
		return fmt.Errorf("invalid argument: %w", err)
	}

	color.Green("Linting project at  %s", color.HiWhiteString(projectdir))

	allIssues := []api.Issue{}
	linters := projectlinters.GetAllLinters()
	for _, linter := range linters {
		issues, err := linter.LintProject(projectdir)
		if err != nil {
			return fmt.Errorf("%s failed to lint project: %w", linter.Name(), err)
		}

		allIssues = append(allIssues, issues...)
	}

	prettyPrintIssues(allIssues)
	if len(allIssues) > 0 {
		return fmt.Errorf("%w: %s", ErrIssuesFound, color.HiWhiteString("%d", len(allIssues)))
	}
	color.Green("Passed!")
	fmt.Println()
	return nil
}

func parseArgs(args []string) (string, error) {
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

func prettyPrintIssues(issues []api.Issue) {
	fmt.Println()
	for i, issue := range issues {
		rule := color.Set(color.Faint).Sprint(issue.Rule)
		msg := prettyPrintMessage(issue.Message)
		fmt.Printf("%d:  %s  %s  %s\n", i+1, issue.Severity.String(), rule, msg)
	}

	if len(issues) > 0 {
		fmt.Println()
	}
}

func prettyPrintMessage(msg string) string {
	return strings.ReplaceAll(msg, ">", color.HiYellowString(">"))
}
