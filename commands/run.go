package commands

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/linters"
	"github.com/bvobart/mllint/utils"
	"github.com/bvobart/mllint/utils/markdown"
)

var ErrNotAFolder = errors.New("not a folder")
var ErrOutputFileAlreadyExists = errors.New("output file already exists")

var outputFile string

func NewRunCommand() *cobra.Command {
	runner := runCommand{}
	cmd := &cobra.Command{
		Use:           "run [dir]",
		Short:         "Run mllint on the project",
		Long:          "Run mllint on the project in the given directory, or the current directory if none was given.",
		RunE:          runner.RunLint,
		Args:          cobra.MaximumNArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", `Export the report generated for your project to a Markdown file at the given location.
Set this to '-' (a single dash) in order to print the raw Markdown directly to the console.`)
	return cmd
}

type runCommand struct {
	Project api.Project
	Config  *config.Config
}

func outputToStdout() bool {
	return outputFile == "-"
}

func outputToFile() bool {
	return outputFile != "" && !outputToStdout()
}

func (rc *runCommand) RunLint(cmd *cobra.Command, args []string) error {
	if outputToFile() && utils.FileExists(outputFile) {
		return fmt.Errorf("%w: %s", ErrOutputFileAlreadyExists, utils.AbsolutePath(outputFile))
	}
	if outputToStdout() {
		quiet = true
	}

	var err error
	rc.Project = api.Project{}
	rc.Project.Dir, err = parseProjectDir(args)
	if err != nil {
		return fmt.Errorf("invalid project path: %w", err)
	}

	shush(func() { color.Green("Linting project at  %s", color.HiWhiteString(rc.Project.Dir)) })
	rc.Config, rc.Project.ConfigType, err = getConfig(rc.Project.Dir)
	if err != nil {
		return err
	}
	shush(func() { fmt.Print("---\n\n") })

	linters.DisableAll(rc.Config.Rules.Disabled)
	if err = linters.ConfigureAll(rc.Config); err != nil {
		return err
	}

	rc.Project.Reports = map[api.Category]api.Report{}
	for cat, linter := range linters.ByCategory {
		report, err := linter.LintProject(rc.Project.Dir)
		if err != nil {
			return fmt.Errorf("linter %s failed to lint project: %w", linter.Name(), err)
		}

		rc.Project.Reports[cat] = report
	}

	output := markdown.FromProject(rc.Project)

	if outputToStdout() {
		fmt.Println(output)
		return nil
	}

	if outputToFile() {
		if err := ioutil.WriteFile(outputFile, []byte(output), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %s", err)
		}
		shush(func() { fmt.Println("Your report is complete, see", utils.AbsolutePath(outputFile)+"\n") })
	} else {
		fmt.Println(markdown.Render(output))
	}
	shush(func() { fmt.Println("---") })

	rulesFailed := rc.countRulesFailed()
	if rulesFailed == 0 {
		printPassed()
	} else {
		printFailed(rulesFailed)
	}

	return nil
}

func (rc *runCommand) countRulesFailed() int {
	rulesFailed := 0
	for _, report := range rc.Project.Reports {
		for _, score := range report.Scores {
			if score < 100 {
				rulesFailed++
			}
		}
	}
	return rulesFailed
}

func printPassed() {
	shush(func() { color.Green("✔️ Passed!") })
	shush(func() { color.Green("Well done, great job!") })
	shush(func() { fmt.Println() })
}

func printFailed(rulesFailed int) {
	shush(func() { color.Red("❌ rules unsuccessful: %s", color.HiWhiteString("%d", rulesFailed)) })

	if rulesFailed <= 2 {
		msg := "You're almost there! There's still a few improvements to be done to get your project up to quality."
		shush(func() { color.Yellow(msg) })
		msg = "Use %s " + color.YellowString("with each rule's slug to learn more about what you can do to get the rules to pass and improve the quality of your ML project.")
		shush(func() { color.Yellow(msg, color.HiWhiteString("mllint describe")) })
		shush(func() { fmt.Println() })
		return
	}

	shush(func() { color.Red("Your project is still lacking in quality and could do with some improvements.") })
	msg := "Use %s " + color.RedString("with each rule's slug to learn more about what you can do to get the rules to pass and improve the quality of your ML project.")
	shush(func() { color.Red(msg, color.YellowString("mllint describe")) })
	shush(func() { fmt.Println() })
}
