package common

import (
	"fmt"
	"strings"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/config"
)

func NewCompositeLinter(name string, linters ...api.Linter) api.ConfigurableLinter {
	rules, lintersByRule := collectRules(linters...)
	return &CompositeLinter{name, linters, rules, lintersByRule}
}

type CompositeLinter struct {
	name    string
	linters []api.Linter
	rules   []*api.Rule
	// maps the rules that the CompositeLinter returns to the underlying linter that checks it.
	lintersByRule map[*api.Rule]api.Linter
}

func (l *CompositeLinter) Name() string {
	return l.name
}

func (l *CompositeLinter) Rules() []*api.Rule {
	return l.rules
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

func (l *CompositeLinter) LintProject(projectdir string) (api.Report, error) {
	finalReport := api.Report{Scores: make(map[api.Rule]float64), Details: make(map[api.Rule]string)}

	for _, linter := range l.linters {
		report, err := linter.LintProject(projectdir)
		if err != nil {
			return api.Report{}, fmt.Errorf("linting error in linter '%s': %w", linter.Name(), err)
		}

		for rule, score := range report.Scores {
			finalReport.Scores[compositeRule(rule, linter.Name())] = score
		}
		for rule, details := range report.Details {
			finalReport.Details[compositeRule(rule, linter.Name())] = details
		}
	}

	return finalReport, nil
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
		for _, rule := range linter.Rules() {
			r := compositeRule(*rule, linter.Name())

			rules = append(rules, &r)
			lintersByRule[&r] = linter
		}
	}

	return rules, lintersByRule
}
