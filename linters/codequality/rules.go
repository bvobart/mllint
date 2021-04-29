package codequality

import "github.com/bvobart/mllint/api"

var RuleUseLinters = api.Rule{
	Name: "Project should use code quality linters",
	Slug: "use-linters",
	Details: `If you have ever seen your code get squiggly coloured underlining while you are writing it,
then you'll be familiar with linting. Linting (or 'static code analysis' as it is more formally called)
is the process of parsing and analysing source code without running it, in an attempt to find common programming issues.

Such issues include type errors and possible buggy operations (e.g. using an undefined variable, or accessing a non-existant field on an object),
but can also show opportunities to rewrite and improve (refactor) your code (e.g. use a list comprehension instead of a for-loop, or showing duplicate (copy-pasted) code).
Additionally, static code analysis can also help enforce a specific code style to keep the code better readable and understandable.
Overall, linting / static code analysis helps ensure the quality and maintainability of your source code.

Python knows many static analysis tools. We recommend you employ the following linters:

Linter | Why?
-------|-----------
**[Pylint](https://pypi.org/project/pylint/)** | General code smells: probable bugs, warnings, refactoring opportunities, basic code style, Python programming conventions.
**[Mypy](http://mypy-lang.org/)**              | Type checking
**[Black](https://github.com/psf/black)**      | Code style
**[isort](https://pypi.org/project/isort/)**   | Code style: automatically sorts and prettyprints imports.
**[Bandit](https://pypi.org/project/bandit/)** | Security

This rule will be satisfied, iff for each of these linters (customisable via config), 
there is _either_ a configuration file in the repository, _or_ the linter is a dependency of the project.`,
	Weight: 1,
}
