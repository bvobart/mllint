package commands

import (
	"errors"
	"fmt"

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
	"github.com/bvobart/mllint/setools/git"
	"github.com/bvobart/mllint/utils"
	"github.com/bvobart/mllint/utils/markdown"
)

var ErrNotAFolder = errors.New("not a folder")
var ErrOutputFileAlreadyExists = errors.New("output file already exists")

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
	SetOutputFlag(cmd)
	SetForceFlag(cmd)
	SetProgressPlainFlag(cmd)
	return cmd
}

type runCommand struct {
	ProjectR api.ProjectReport
	Config   *config.Config
	Runner   *mllint.MLLintRunner
}

// Runs pre-analysis checks:
// - Retrieve some info about project's Git state
// - Detect dependency managers used in the project
// - Detect code quality linters used in the project
// - Detect the Python files in the project repository.
func (rc *runCommand) runPreAnalysisChecks() error {
	rc.ProjectR.Git = git.MakeGitInfo(rc.ProjectR.Dir)
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
	err := checkOutputFlag()
	if err != nil {
		return err
	}

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
	rc.ProjectR.Config = *rc.Config
	shush(func() { fmt.Print("---\n\n") })

	// configure all linters with config
	if err = linters.ConfigureAll(rc.Config); err != nil {
		return err
	}

	// disable any rules from config.
	// This is done after configuring each linter, such that any rules arising from the configuration (e.g. custom rules) can also be disabled.
	rulesDisabled := linters.DisableAll(rc.Config.Rules.Disabled)

	// run pre-analysis checks
	if err = rc.runPreAnalysisChecks(); err != nil {
		return fmt.Errorf("failed to run pre-analysis checks: %w", err)
	}

	// start the runner and do all linting
	progress := createRunnerProgress()
	rc.Runner = mllint.NewMLLintRunner(progress)
	rc.Runner.Start()

	tasks := scheduleLinters(rc.Runner, rc.ProjectR.Project, linters.ByCategory)
	rc.ProjectR.Reports, rc.ProjectR.Errors = collectReports(rc.Runner, tasks...)

	// convert project report to Markdown
	output := markdown.FromProject(rc.ProjectR)

	rc.Runner.Close()

	if outputToStdout() {
		fmt.Println(output)
		return nil
	}

	if outputToFile() {
		return writeToOutputFile(output)
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

	printSurveyMessage()
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

func scheduleLinters(runner mllint.Runner, project api.Project, linters map[api.Category]api.Linter) []*mllint.RunnerTask {
	tasks := make([]*mllint.RunnerTask, 0, len(linters))
	for cat, linter := range linters {
		if len(linter.Rules()) == 0 {
			continue
		}

		// use cat.Slug as ID so we can retrieve the category from categories.BySlug later, see collectReports(..)
		task := runner.RunLinter(cat.Slug, linter, project)
		tasks = append(tasks, task)
	}
	return tasks
}

func collectReports(runner mllint.Runner, tasks ...*mllint.RunnerTask) (map[api.Category]api.Report, *multierror.Error) {
	var err *multierror.Error
	reports := map[api.Category]api.Report{}

	mllint.ForEachTask(runner.CollectTasks(tasks...), func(task *mllint.RunnerTask, result mllint.LinterResult) {
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

// BvOBart: As much as I hate programs asking me to fill in a survey for me, it's a necessity for being able to evaluate mllint for my thesis.
func printSurveyMessage() {
	yellow := color.New(color.FgYellow)
	code := color.New(color.FgHiWhite, color.Italic)
	shush(func() {
		fmt.Println("---")
		yellow.Print("Thank you for using ", code.Sprint("mllint"), yellow.Sprintln(", I'm very interested to know how your experiences have been! 😊"))
		yellow.Print("It would be of great help to me, ", code.Sprint("mllint"), yellow.Sprint(" and, in particular, my ", code.Sprint("MSc thesis"), yellow.Sprintln(" if you are able to spend 15 minutes of your time filling in this feedback survey for me:")))
		yellow.Add(color.Italic).Println("➡️   https://forms.office.com/r/pXtfUKWUDA")
		fmt.Println()
	})
}
