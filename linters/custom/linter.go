package custom

import (
	"fmt"
	"strings"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/commands/mllint"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/utils/exec"
	"gopkg.in/yaml.v3"

	"github.com/google/shlex"
	"github.com/hashicorp/go-multierror"
)

// Returns an api.Linter or an api.ConfigurableLinter if also implementing api.Configurable
func NewLinter() mllint.ConfigurableLinterWithRunner {
	return &CustomLinter{}
}

// Your linter object. Give this a nice name.
type CustomLinter struct {
	customRules map[config.CustomRule]*api.Rule
	rules       []*api.Rule
	runner      mllint.Runner
}

func (l *CustomLinter) Name() string {
	return "Custom Rules"
}

func (l *CustomLinter) Configure(conf *config.Config) error {
	l.rules = make([]*api.Rule, 0, len(conf.Rules.Custom))
	l.customRules = make(map[config.CustomRule]*api.Rule, len(conf.Rules.Custom))

	for _, customRule := range conf.Rules.Custom {
		rule := api.NewCustomRule(customRule)
		l.rules = append(l.rules, &rule)
		l.customRules[customRule] = &rule
	}

	return nil
}

func (l *CustomLinter) Rules() []*api.Rule {
	return l.rules
}

func (l *CustomLinter) SetRunner(runner mllint.Runner) {
	l.runner = runner
}

func (l *CustomLinter) LintProject(project api.Project) (api.Report, error) {
	report := api.NewReport()

	// create linters from each of the rules and schedule each of them for execution on the mllint.Runner
	tasks := []*mllint.RunnerTask{}
	for customRule, rule := range l.customRules {
		customLinter := customRuleLinter{customRule, *rule}
		task := l.runner.RunLinter(rule.Slug, &customLinter, project)
		tasks = append(tasks, task)
	}

	// collect all the results
	var multiErr *multierror.Error
	mllint.ForEachTask(l.runner.CollectTasks(tasks...), func(task *mllint.RunnerTask, result mllint.LinterResult) {
		report = api.MergeReports(report, result.Report)
		multiErr = multierror.Append(multiErr, result.Err)
	})

	return report, multiErr.ErrorOrNil()
}

//---------------------------------------------------------------------------------------

// describes the expected structure of what the execution of a custom rule should result in
type customRuleResult struct {
	Score   float64 `json:"score" yaml:"score"`
	Details string  `json:"details" yaml:"details"`
}

//---------------------------------------------------------------------------------------

type customRuleLinter struct {
	customRule config.CustomRule
	rule       api.Rule
}

func (l *customRuleLinter) Name() string {
	return "Custom Rule - " + l.rule.Name
}

// otherwise unused, but it's here to ensure customRuleLinter implements api.Linter
func (l *customRuleLinter) Rules() []*api.Rule { return []*api.Rule{&l.rule} }

// runs the custom rule definition's `run` command in the project's root directory and parses the result as YAML ()
func (l *customRuleLinter) LintProject(project api.Project) (api.Report, error) {
	report := api.NewReport()

	cmdparts, err := shlex.Split(l.customRule.Run)
	if err != nil {
		return report, fmt.Errorf("custom rule `%s` has invalid run command `%s`: %w", l.customRule.Slug, l.customRule.Run, err)
	}

	output, err := exec.CommandCombinedOutput(project.Dir, cmdparts[0], cmdparts[1:]...)
	if err != nil {
		return report, fmt.Errorf("custom rule `%s` was run, but exited with an error: %w.%s", l.customRule.Slug, err, formatOutput(output))
	}

	var result customRuleResult
	if err := yaml.Unmarshal(output, &result); err != nil {
		return report, fmt.Errorf("custom rule `%s` executed successfully, but the output was not a valid YAML or JSON object: %w.%s", l.customRule.Slug, err, formatOutput(output))
	}

	report.Scores[l.rule] = result.Score
	report.Details[l.rule] = result.Details
	return report, nil
}

// formats the output of an executed command such that it can be appended to an error.
// Returns an empty string if the output is empty, returns a single line with inline code block if the trimmed output is a single line,
// returns with a multiline code block if the trimmed output is multi-line.
func formatOutput(output []byte) string {
	if len(output) == 0 {
		return ""
	}

	trimmed := strings.TrimSpace(string(output))
	if strings.Contains(trimmed, "\n") {
		return fmt.Sprintf("\nOutput:\n\n```\n%s\n```", trimmed)
	}
	return fmt.Sprintf(" Output: `%s`", trimmed)
}
