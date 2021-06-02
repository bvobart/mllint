package template

import "github.com/bvobart/mllint/api"

var RuleHasTests = api.Rule{
	Name:    "Project has automated tests",
	Slug:    "testing/has-tests",
	Details: "TODO",
	Weight:  1,
}

var RuleTestsPass = api.Rule{
	Name:    "Project passes all of its automated tests",
	Slug:    "testing/pass",
	Details: "TODO",
	Weight:  1,
}

var RuleTestCoverage = api.Rule{
	Name:    "Project provides a test coverage report",
	Slug:    "testing/coverage",
	Details: "TODO: fill this in with details about what the rule checks, what the reasoning behind it is, where people can get more information on how to implement what this rule checks.",
	Weight:  1,
}

var RuleTestsFolder = api.Rule{
	Name:    "Tests should be placed in the tests folder",
	Slug:    "testing/tests-folder",
	Details: "TODO",
	Weight:  1,
}
