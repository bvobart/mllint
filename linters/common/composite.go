package common

import (
	"fmt"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/config"
)

func NewCompositeLinter(name string, linters ...api.Linter) api.Linter {
	return &CompositeLinter{name, linters}
}

type CompositeLinter struct {
	name    string
	linters []api.Linter
}

func (l *CompositeLinter) Name() string {
	return l.name
}

func (l *CompositeLinter) Rules() []api.Rule {
	rules := []api.Rule{}
	for _, linter := range l.linters {
		for _, rule := range linter.Rules() {
			rule.Name = l.prefixRule(rule.Name)
			rules = append(rules, rule)
		}
	}
	return rules
}

func (l *CompositeLinter) Configure(conf *config.Config) error {
	for _, linter := range l.linters {
		if err := linter.Configure(conf); err != nil {
			return fmt.Errorf("configuration error in linter '%s': %w", linter.Name(), err)
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
			rule.Name = l.prefixRule(rule.Name)
			finalReport.Scores[rule] = score
		}
		for rule, details := range report.Details {
			rule.Name = l.prefixRule(rule.Name)
			finalReport.Details[rule] = details
		}
	}

	return finalReport, nil
}

// prefixes the name of a rule with the name of this linter and a colon.
func (l *CompositeLinter) prefixRule(name string) string {
	return l.name + ": " + name
}
