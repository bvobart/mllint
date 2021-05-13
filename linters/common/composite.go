package common

import (
	"fmt"
	"strings"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/commands/mllint"
	"github.com/bvobart/mllint/config"
	"github.com/hashicorp/go-multierror"
)

func NewCompositeLinter(name string, linters ...api.Linter) mllint.ConfigurableLinterWithRunner {
	rules, lintersByRule := collectRules(linters...)
	return &CompositeLinter{name, linters, rules, nil, lintersByRule}
}

type CompositeLinter struct {
	name    string
	linters []api.Linter
	rules   []*api.Rule
	runner  *mllint.Runner
	// maps the rules that the CompositeLinter returns to the underlying linter that checks it.
	lintersByRule map[*api.Rule]api.Linter
}

func (l *CompositeLinter) Name() string {
	return l.name
}

func (l *CompositeLinter) Rules() []*api.Rule {
	return l.rules
}

func (l *CompositeLinter) SetRunner(r *mllint.Runner) {
	l.runner = r
}

func (l *CompositeLinter) DisableRule(rule *api.Rule) {
	linter, ok := l.lintersByRule[rule]
	if !ok {
		return
	}

	rule.Disable()

	slashIndex := strings.Index(rule.Slug, "/")
	originalSlug := rule.Slug[slashIndex+1:]

	for _, originalRule := range linter.Rules() {
		if originalSlug == originalRule.Slug {
			if compLinter, ok := linter.(*CompositeLinter); ok {
				compLinter.DisableRule(originalRule)
				return
			}

			originalRule.Disable()
			return
		}
	}
}

func (l *CompositeLinter) Configure(conf *config.Config) error {
	for _, linter := range l.linters {
		configurable, ok := linter.(api.Configurable)
		if ok {
			if err := configurable.Configure(conf); err != nil {
				return fmt.Errorf("configuration error in linter '%s': %w", linter.Name(), err)
			}
		}
	}
	return nil
}

func (l *CompositeLinter) LintProject(project api.Project) (api.Report, error) {
	finalReport := api.NewReport()

	tasks := make([]*mllint.RunnerTask, len(l.linters))
	for i, linter := range l.linters {
		displayName := l.name + " - " + linter.Name()
		tasks[i] = l.runner.RunLinter(fmt.Sprint(i), linter, project, mllint.DisplayName(displayName))
	}

	var err *multierror.Error
	mllint.ForEachTask(mllint.CollectTasks(tasks...), func(task *mllint.RunnerTask, result mllint.LinterResult) {
		if result.Err != nil {
			err = multierror.Append(err, fmt.Errorf("linting error in linter '%s': %w", task.Linter.Name(), result.Err))
		}

		for rule, score := range result.Report.Scores {
			finalReport.Scores[compositeRule(rule, task.Linter.Name())] = score
		}
		for rule, details := range result.Report.Details {
			finalReport.Details[compositeRule(rule, task.Linter.Name())] = details
		}
	})

	return finalReport, err.ErrorOrNil()
}

func compositeRule(rule api.Rule, linterName string) api.Rule {
	rule.Name = linterName + ": " + rule.Name
	rule.Slug = api.Slug(linterName) + "/" + rule.Slug
	return rule
}

func collectRules(linters ...api.Linter) (rules []*api.Rule, lintersByRule map[*api.Rule]api.Linter) {
	rules = []*api.Rule{}
	lintersByRule = map[*api.Rule]api.Linter{}

	for _, linter := range linters {
		linterName := linter.Name()
		for _, rule := range linter.Rules() {
			r := compositeRule(*rule, linterName)

			rules = append(rules, &r)
			lintersByRule[&r] = linter
		}
	}

	return rules, lintersByRule
}
