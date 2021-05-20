package mypy

import (
	"fmt"

	"github.com/bvobart/mllint/api"
)

var RuleNoIssues = api.Rule{
	Slug: "code-quality/mypy/no-issues",
	Name: "Mypy reports no issues with this project",
	Details: fmt.Sprintf(`> [Mypy](http://mypy-lang.org/) is an optional static type checker for Python that aims to combine the benefits of dynamic (or "duck") typing and static typing. Mypy combines the expressive power and convenience of Python with a powerful type system and compile-time type checking.

This rule checks whether Mypy finds any type issues when running it on all Python files in this project.

Per default, mllint is configured to make Mypy enforce static typing.

The score for this rule is determined as a function of the number of messages Mypy returns and the lines of Python code that your project has.
In the ideal case, Mypy does not recognise any code smells in your project, in which case the score is 100%%.
When there is one Mypy issue for every 20 lines of code, then the score is 50%%.
When there is one Mypy issue for every 10 lines of code, then the score is 0%%.

More specifically, in pseudocode, %s. Note that the measured amount of lines of code includes any non-hidden Python files in the repository, including those that are ignored by Mypy.

To learn more about how type-checking works and how to use it in Python, see:
- https://realpython.com/python-type-checking/
`,
		"`score = 100 - 100 * min(1, 10 * number of msgs / lines of code)`"),
	Weight: 1,
}

// var RuleIsConfigured = api.Rule{
// 	Slug: "pylint/is-configured",
// 	Name: "Mypy is configured for this project",
// 	Details: `[Mypy](http://mypy-lang.org/) has a good default configuration,
// though there are likely to be rules that you may want to enable, disable or customise for your project.
// For example, you may want to configure your indentation width, your maximum line length, or configure which files to ignore while linting.

// Additionally, some IDEs have their own default configuration for these linters, which may only enable a subset of Mypy's rules.
// For example, VS Code [is known to do this](https://code.visualstudio.com/docs/python/linting#_default-pylint-rules).
// However, those IDEs generally _do_ pick up your project's own Mypy configuration.

// Having a Mypy configuration in the project also ensures that you, each of your colleages, as well as the CI,
// use the same linting configuration.`,
// 	Weight: 1,
// }
