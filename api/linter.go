package api

import (
	"fmt"

	"gitlab.com/bvobart/mllint/config"
)

// Linter is the main interface for any linter usable by `mllint`.
// It is recommended to implement this interface as a struct with methods that have pointer receivers.
type Linter interface {
	// Name returns the name of this linter. Linter names should be lowercased, use dashes for spaces and contain no special characters.
	// e.g. 'use-dependency-manager'
	Name() string

	// Rules returns a list of the rules that this linter enforces.
	// The name of a linting rule should also be lowercased, use dashes for spaces, contain no special characters and should not contain the linter's name, so,
	// e.g. 'no-requirements-txt'.
	//
	// These linting rule names will be combined with the linter's name to form the full name of the rule,
	// which will be used to enable and disable these specific rules.
	// e.g. 'use-dependency-manager/no-requirements-txt'
	//
	// If your linter only enforces one rule, or if the linter has one main rule that should have the name of the linter,
	// then include an empty string in this method's response, e.g. `[]string{""}`.
	Rules() []string

	// Configure will be called before LintProject() is called to analyse a project.
	// Implement this method such that the linter configures itself to use the settings in the config object,
	// and to disable any (computationally expensive) linting rules that are disabled in the config anyways.
	//
	// Return a non-nil error if there is a problem with how mllint is configured. E.g. if a configuration option
	// for this linter has a value that is outside of valid ranges.
	Configure(conf *config.Config) error

	// LintProject is the main method that runs this linter. It will receive the full path to the
	// directory in which the project is located. It is then expected to perform its respective analysis
	// and return a list of issues (or nil if there were no issues) or an error if there was one.
	//
	// Any issues of rules that are not enabled will be filtered out before the issues are displayed to the user.
	LintProject(projectdir string) ([]Issue, error)
}

type LinterList []Linter

func (list LinterList) FilterEnabled(conf config.RuleConfig) LinterList {
	linters := map[string]Linter{}
	for _, linter := range list {
		linters[linter.Name()] = linter
	}

	for _, rule := range conf.Disabled {
		if _, isEnabled := linters[rule]; isEnabled {
			linters[rule] = nil
		}
	}

	enabledLinters := LinterList{}
	for _, linter := range list {
		if linters[linter.Name()] != nil {
			enabledLinters = append(enabledLinters, linter)
		}
	}

	return enabledLinters
}

func (list LinterList) Configure(conf *config.Config) (LinterList, error) {
	for _, linter := range list {
		if err := linter.Configure(conf); err != nil {
			return nil, fmt.Errorf("%s configuration error: %w", linter.Name(), err)
		}
	}
	return list, nil
}
