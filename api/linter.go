package api

import (
	"github.com/bvobart/mllint/config"
)

// Linter is the main interface for a struct that defines a linter over a certain category.
// There must be one Linter per Category, which may be a linters.CompositeLinter that employs
// several other Linters to subdivide checking all the rules within that category.
// It is recommended to implement this interface as a struct with methods that have pointer receivers.
type Linter interface {
	// Name of the linter
	Name() string

	// Rules returns all the rules that this linter can check while linting a project
	Rules() []*Rule

	// LintProject is the main method that runs this linter. It will receive the full path to the
	// directory in which the project is located. It is then expected to perform its respective analysis
	// and return a Report or an error if there was one.
	//
	// The returned Report should contain a mapping of each checked Rule to a percentual score between 0 and 100.
	// A linter may also add additional details to a report related to a specific rule, which is especially
	// recommended if the rule scored less than 100.
	LintProject(project Project) (Report, error)
}

// Configure should be implemented such that the struct that implements it configures itself to use the settings
// from the config object. If implemented on a Linter, this will be called before LintProject is called.
type Configurable interface {
	// Configure the implementing struct with the config. Return a non-nil error if there is a problem with how
	// mllint is configured. E.g. if a configuration option for this linter has a value that is outside of valid ranges.
	Configure(conf *config.Config) error
}

// ConfigurableLinter is simply a Linter that also implements Configurable.
type ConfigurableLinter interface {
	Linter
	Configurable
}
