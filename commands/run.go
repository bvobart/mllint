package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"runtime"
	"sync"

	"github.com/fatih/color"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/linters"
	"github.com/bvobart/mllint/setools/cqlinters"
	"github.com/bvobart/mllint/setools/depmanagers"
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
	ProjectR api.ProjectReport
	Config   *config.Config
	Runner   *MLLintRunner
}

func outputToStdout() bool {
	return outputFile == "-"
}

func outputToFile() bool {
	return outputFile != "" && !outputToStdout()
}

type MLLintRunner struct {
	queue chan lintJob
	wg    sync.WaitGroup
}

type lintJob struct {
	id      string
	linter  api.Linter
	project api.Project
	result  chan LinterResult
}

type LinterResult struct {
	api.Report
	Err error
}

func NewRunner() *MLLintRunner {
	return &MLLintRunner{
		queue: make(chan lintJob, 20),
		wg:    sync.WaitGroup{},
	}
}

func (r *MLLintRunner) Start() {
	go r.queueWorker()
}

func (r *MLLintRunner) AwaitAll() {
	r.wg.Wait()
}

func (r *MLLintRunner) queueWorker() {
	running := 0
	parked := []lintJob{}
	done := make(chan lintJob, runtime.NumCPU())

	for {
		select {
		case job, open := <-r.queue:
			if !open {
				return
			}

			if running >= runtime.NumCPU() {
				parked = append(parked, job)
				color.Blue("Scheduled: %s", job.linter.Name())
				break
			}

			color.Yellow("Running: %s", job.linter.Name())
			running++
			go r.runJob(job, done)

		case job := <-done:
			running--

			if len(parked) > 0 {
				var nextJob lintJob
				nextJob, parked = parked[0], parked[1:]

				color.Yellow("Running: %s", job.linter.Name())
				running++
				go r.runJob(nextJob, done)
			}

			color.Green("Done: %s", job.linter.Name())
		}
	}
}

func (r *MLLintRunner) runJob(job lintJob, done chan lintJob) {
	report, err := job.linter.LintProject(job.project)
	job.result <- LinterResult{Report: report, Err: err}

	done <- job
	r.wg.Done()
}

func (r *MLLintRunner) RunLinter(id string, linter api.Linter, project api.Project) lintJob {
	result := make(chan LinterResult, 1)
	job := lintJob{id, linter, project, result}
	r.wg.Add(1)
	r.queue <- job
	return job
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
		return fmt.Errorf("%w: %s", ErrOutputFileAlreadyExists, utils.AbsolutePath(outputFile))
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

	// do all linting
	rc.Runner = NewRunner()
	rc.Runner.Start()

	jobs := make([]lintJob, len(linters.ByCategory))
	for cat, linter := range linters.ByCategory {
		job := rc.Runner.RunLinter(cat.Slug, linter, rc.ProjectR.Project)
		jobs = append(jobs, job)
	}

	// TODO: get this shit working

	rc.Runner.AwaitAll()

	rc.ProjectR.Reports = map[api.Category]api.Report{}
	for _, job := range jobs {
		color.Yellow("Awaiting: %s", job.id)
		result := <-job.result
		if result.Err != nil {
			rc.ProjectR.Errors = multierror.Append(rc.ProjectR.Errors, fmt.Errorf("**%s** - %w", job.linter.Name(), err))
		}

		cat := categories.BySlug[job.id]
		rc.ProjectR.Reports[cat] = result.Report
	}

	// convert project report to Markdown
	output := markdown.FromProject(rc.ProjectR)

	if outputToStdout() {
		fmt.Println(output)
		return nil
	}

	if outputToFile() {
		if err := ioutil.WriteFile(outputFile, []byte(output), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		shush(func() { fmt.Println("Your report is complete, see", utils.AbsolutePath(outputFile)+"\n") })
	} else {
		fmt.Println(markdown.Render(output))
	}
	shush(func() { fmt.Println("---") })

	rulesFailed := rc.countRulesFailed()
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

func (rc *runCommand) countRulesFailed() int {
	rulesFailed := 0
	for _, report := range rc.ProjectR.Reports {
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
	shush(func() { color.Red(msg, color.YellowString("mllint describe")) })
	shush(func() { fmt.Println() })
}

func printSkipped(rulesDisabled int) {
	shush(func() { color.Red("⏭️ rules skipped: \t%s", color.HiWhiteString("%d", rulesDisabled)) })
}
