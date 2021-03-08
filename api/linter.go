package api

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

	// LintProject is the main method that runs this linter. It will receive the full path to the
	// directory in which the project is located. It is then expected to perform its respective analysis
	// and return a list of issues (or nil if there were no issues) or an error if there was one.
	LintProject(projectdir string) ([]Issue, error)
}
