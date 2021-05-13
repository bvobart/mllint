package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/commands/mllint"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/linters"
	"github.com/bvobart/mllint/setools/cqlinters"
	"github.com/bvobart/mllint/setools/depmanagers"
	"github.com/bvobart/mllint/utils"
	"github.com/bvobart/mllint/utils/markdown"
)

var ErrNotAFolder = errors.New("not a folder")
var ErrOutputFileAlreadyExists = errors.New("output file already exists")

// Flags
var (
	outputFile    string
	progressPlain bool
	force         bool
)

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
Set this to '-' (a single dash) in order to print the raw Markdown directly to the console (implies '-q').`)
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Use this flag to remove the output file provided with '--output' in case that already exists.")
	cmd.Flags().BoolVar(&progressPlain, "progress-plain", false, "Use this flag to print linting progress plainly, without rewriting terminal output. Overrides '-q'. Enabled automatically in non-interactive terminals (except when using '-q').")
	return cmd
}

type runCommand struct {
	ProjectR api.ProjectReport
	Config   *config.Config
	Runner   *mllint.Runner
}

func outputToStdout() bool {
	return outputFile == "-"
}

func outputToFile() bool {
	return outputFile != "" && !outputToStdout()
}

// Runs pre-analysis checks:
// - Detect dependency managers used in the project
// - Detect code quality linters used in the project
// - Detect the Python files in the project repository.
func (rc *runCommand) runPreAnalysisChecks() error {
	rc.ProjectR.DepManagers = depmanagers.Detect(rc.ProjectR.Project)
	rc.ProjectR.CQLinters = cqlinters.Detect(rc.ProjectR.Project)

	pyfiles, err := utils.FindPythonFilesIn(rc.ProjectR.Dir)
	if err != nil {
		return err
	}
	rc.ProjectR.PythonFiles = pyfiles.Prefix(rc.ProjectR.Dir)

	return nil
}

func (rc *runCommand) RunLint(cmd *cobra.Command, args []string) error {
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

	var err error
	rc.ProjectR = api.ProjectReport{}
	rc.ProjectR.Dir, err = parseProjectDir(args)
	if err != nil {
		return fmt.Errorf("invalid project path: %w", err)
	}

	shush(func() { color.Green("Linting project at  %s", color.HiWhiteString(rc.ProjectR.Dir)) })
	rc.Config, rc.ProjectR.ConfigType, err = getConfig(rc.ProjectR.Dir)
	if err != nil {
		return err
	}
	shush(func() { fmt.Print("---\n\n") })

	// disable any rules from config
	rulesDisabled := linters.DisableAll(rc.Config.Rules.Disabled)

	// configure all linters with config
	if err = linters.ConfigureAll(rc.Config); err != nil {
		return err
	}

	// run pre-analysis checks
	if err = rc.runPreAnalysisChecks(); err != nil {
		return fmt.Errorf("failed to run pre-analysis checks: %w", err)
	}

	// start the runner and do all linting
	progress := createRunnerProgress()
	rc.Runner = mllint.NewRunner(progress)
	rc.Runner.Start()

	tasks := scheduleLinters(rc.Runner, rc.ProjectR.Project, linters.ByCategory)
	rc.ProjectR.Reports, rc.ProjectR.Errors = collectReports(tasks...)

	// convert project report to Markdown
	output := markdown.FromProject(rc.ProjectR)

	rc.Runner.Close()

	if outputToStdout() {
		fmt.Println(output)
		return nil
	}

	if outputToFile() {
		if err := ioutil.WriteFile(outputFile, []byte(output), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		bold := color.New(color.Bold)
		shush(func() { bold.Println("Your report is complete, see", formatInlineCode(utils.AbsolutePath(outputFile))) })
		shush(func() { bold.Println() })
	} else {
		fmt.Println(markdown.Render(output))
	}

	shush(func() { fmt.Println("---") })

	rulesFailed := countRulesFailed(rc.ProjectR.Reports)
	if rulesDisabled > 0 {
		printSkipped(rulesDisabled)
	}
	if rulesFailed == 0 {
		printPassed()
	} else {
		printFailed(rulesFailed)
	}

	return nil
}

func createRunnerProgress() mllint.RunnerProgress {
	if progressPlain {
		return mllint.NewBasicRunnerProgress()
	}
	if quiet {
		return nil
	}
	if utils.IsInteractive() {
		return mllint.NewLiveRunnerProgress()
	}
	return mllint.NewBasicRunnerProgress()
}

func scheduleLinters(runner *mllint.Runner, project api.Project, linters map[api.Category]api.Linter) []*mllint.RunnerTask {
	tasks := make([]*mllint.RunnerTask, 0, len(linters))
	for cat, linter := range linters {
		// use cat.Slug as ID so we can retrieve the category from categories.BySlug later, see collectReports(..)
		task := runner.RunLinter(cat.Slug, linter, project)
		tasks = append(tasks, task)
	}
	return tasks
}

func collectReports(tasks ...*mllint.RunnerTask) (map[api.Category]api.Report, *multierror.Error) {
	var err *multierror.Error
	reports := map[api.Category]api.Report{}

	mllint.ForEachTask(mllint.CollectTasks(tasks...), func(task *mllint.RunnerTask, result mllint.LinterResult) {
		if result.Err != nil {
			err = multierror.Append(err, fmt.Errorf("**%s** - %w", task.Linter.Name(), result.Err))
		}

		// cat.Slug was used as task ID, see scheduleLinters(..)
		cat := categories.BySlug[task.Id]
		reports[cat] = result.Report
	})

	return reports, err
}

func countRulesFailed(reports map[api.Category]api.Report) int {
	rulesFailed := 0
	for _, report := range reports {
		for rule, score := range report.Scores {
			if !rule.Disabled && score < 100 {
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
	shush(func() { color.Red("❌ rules unsuccessful:\t%s", color.HiWhiteString("%d", rulesFailed)) })

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
	shush(func() { color.Red(msg, formatInlineCode("mllint describe")) })
	shush(func() { fmt.Println() })
}

func printSkipped(rulesDisabled int) {
	shush(func() { color.Red("⏭️ rules skipped: \t%s", color.HiWhiteString("%d", rulesDisabled)) })
}
