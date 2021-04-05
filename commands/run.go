package commands

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/linters"
	"github.com/bvobart/mllint/utils"
)

var (
	ErrNotAFolder        = errors.New("not a folder")
	ErrRulesUnsuccessful = errors.New(color.RedString("rules unsuccessful:"))
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

	linters.DisableAll(conf.Rules.Disabled)
	if err = linters.ConfigureAll(conf); err != nil {
		return err
	}

	reports := map[api.Category]api.Report{}
	for cat, linter := range linters.ByCategory {
		report, err := linter.LintProject(projectdir)
		if err != nil {
			return fmt.Errorf("linter %s failed to lint project: %w", linter.Name(), err)
		}

		reports[cat] = report
	}

	rulesFailed := prettyPrintReports(reports)
	if rulesFailed > 0 {
		return fmt.Errorf("%s %w %s", color.RedString("❌"), ErrRulesUnsuccessful, color.HiWhiteString("%d", rulesFailed))
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

func prettyPrintReports(reports map[api.Category]api.Report) int {
	rulesFailed := 0
	for cat, report := range reports {
		color.Set(color.Bold).Println(cat)
		color.Unset()

		for rule, score := range report.Scores {
			if !rule.Disabled {
				if score < 100 {
					rulesFailed++
				}

				fmt.Println(fmt.Sprintf("%s: %.2f", rule.Name, score) + "%")
			}
		}

		fmt.Println()
	}

	return rulesFailed
}
