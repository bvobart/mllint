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
	return "Custom"
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

	var multiErr *multierror.Error
	for customRule, rule := range l.customRules {
		var err error
		report.Scores[*rule], report.Details[*rule], err = l.runCustomRule(project, customRule)
		if err != nil {
			multiErr = multierror.Append(multiErr, err)
			continue
		}
	}

	return report, multiErr.ErrorOrNil()
}

func (l *CustomLinter) runCustomRule(project api.Project, rule config.CustomRule) (float64, string, error) {
	// TODO: run these rules in their own processes on a runner

	cmdparts, err := shlex.Split(rule.Run)
	if err != nil {
		return 0, "", fmt.Errorf("custom rule `%s` has invalid run command `%s`: %w", rule.Slug, rule.Run, err)
	}

	output, err := exec.CommandCombinedOutput(project.Dir, cmdparts[0], cmdparts[1:]...)
	if err != nil {
		return 0, "", fmt.Errorf("custom rule `%s` was run, but exited with an error: %w.%s", rule.Slug, err, formatOutput(output))
	}

	var result customRuleResult
	if err := yaml.Unmarshal(output, &result); err != nil {
		return 0, "", fmt.Errorf("custom rule `%s` executed successfully, but the output was not a valid JSON / YAML object: %w.%s", rule.Slug, err, formatOutput(output))
	}

	return result.Score, result.Details, nil
}

// describes the expected structure of what the execution of a custom rule should result in
type customRuleResult struct {
	Score   float64 `json:"score" yaml:"score"`
	Details string  `json:"details" yaml:"details"`
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
