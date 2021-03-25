package commands

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/projectlinters"
	"github.com/bvobart/mllint/utils"
)

var (
	ErrNotAFolder  = errors.New("not a folder")
	ErrIssuesFound = errors.New(color.RedString("issues found:"))
)

func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "run [dir]",
		Short:         "Run mllint on the project",
		Long:          "Run mllint on the project in the given directory, or the current directory if none was given.",
		RunE:          lint,
		Args:          cobra.MaximumNArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	return cmd
}

func lint(cmd *cobra.Command, args []string) error {
	projectdir, err := parseProjectDir(args)
	if err != nil {
		return fmt.Errorf("invalid argument: %w", err)
	}

	shush(func() { color.Green("Linting project at  %s", color.HiWhiteString(projectdir)) })
	conf, _, err := getConfig(projectdir)
	if err != nil {
		return err
	}
	shush(func() { fmt.Print("---\n\n") })

	allIssues := api.IssueList{}
	linters, err := projectlinters.GetAllLinters().FilterEnabled(conf.Rules).Configure(conf)
	if err != nil {
		return err
	}

	for _, linter := range linters {
		issues, err := linter.LintProject(projectdir)
		if err != nil {
			return fmt.Errorf("%s failed to lint project: %w", linter.Name(), err)
		}

		allIssues = append(allIssues, issues...)
	}

	enabledIssues := allIssues.FilterEnabled(conf.Rules)
	prettyPrintIssues(enabledIssues)

	if len(enabledIssues) > 0 {
		return fmt.Errorf("%s %w %s", color.RedString("❌"), ErrIssuesFound, color.HiWhiteString("%d", len(allIssues)))
	}
	color.Green("✔️ Passed!")
	fmt.Println()
	return nil
}

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

func prettyPrintIssues(issues []api.Issue) {
	for i, issue := range issues {
		fmt.Printf("%d:  %s\n\n", i+1, issue.String())
	}

	if len(issues) > 0 {
		fmt.Println()
	}
}
