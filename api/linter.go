package api

import (
	"fmt"

	"github.com/bvobart/mllint/config"
)

// Linter is the main interface for a struct that defines a linter over a certain category.
// There must be one Linter per Category, which may be a linters.CompositeLinter that employs
// several other Linters to subdivide checking all the rules within that category.
// It is recommended to implement this interface as a struct with methods that have pointer receivers.
type Linter interface {
	Name() string

	// Rules returns all the rules that this linter can check while linting a project
	Rules() []Rule

	// Configure will be called before LintProject() is called to analyse a project.
	// Implement this method such that the linter configures itself to use the settings in the config object,
	// and to disable any (computationally expensive) linting rules that are disabled in the config anyways.
	//
	// Return a non-nil error if there is a problem with how mllint is configured. E.g. if a configuration option
	// for this linter has a value that is outside of valid ranges.
	Configure(conf *config.Config) error

	// LintProject is the main method that runs this linter. It will receive the full path to the
	// directory in which the project is located. It is then expected to perform its respective analysis
	// and return a Report or an error if there was one.
	//
	// The returned Report should contain a mapping of each checked Rule to a percentual score between 0 and 100.
	// A linter may also add additional details to a report related to a specific rule, which is especially
	// recommended if the rule scored less than 100.
	LintProject(projectdir string) (Report, error)
}

// Report is the type of object returned by a Linter after linting a project.
type Report struct {
	// Scores maps each evaluated rule to a score
	Scores map[Rule]float64

	// Details contains any additional details to accompany a Rule's evaluation.
	// Typically, when a Linter detects that a project does not conform to a Rule,
	// it will want to provide some form of reasoning about it, pointers to which
	// parts of the project repo the Rule violation occcurs in, and what the user can
	// do to fix the issue.
	//
	// The mapped string may be formatted using Markdown.
	Details map[Rule]string
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
