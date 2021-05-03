package pylint

import "github.com/bvobart/mllint/api"

var RuleNoIssues = api.Rule{
	Slug: "pylint/no-issues",
	Name: "Pylint reports no issues with this project",
	Details: `[Pylint](https://pypi.org/project/pylint/) is a static analysis tool for finding generic programming errors.
This rule asserts that it does not return any errors when running it on all Python files in this project.`,
}
