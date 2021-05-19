package codequality

import (
	"fmt"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/utils/markdowngen"
)

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

This rule will be satisfied, iff for each of these linters:
- **Either** there is a configuration file for this linter in the project
- **Or** the linter is a dependency of the project`,
	Weight: 1,
}

func detailsUseLinters(detectedLinters []api.CQLinter, missingLinters []api.CQLinter) string {
	if len(missingLinters) == 0 {
		return fmt.Sprintf("Hooray, all linters detected:\n\n%s", markdowngen.List(asInterfaceList(detectedLinters)))
	}

	extraDetails := fmt.Sprintf(`This rule will be satisfied, iff for each of these linters:
- **Either** there is a configuration file for this linter in the project
- **Or** the linter is a dependency of the project

Specifically, we recommend adding each linter to the development dependencies of your dependency manager,
e.g. using `+"`poetry add --dev %s` or `pipenv install --dev %s`\n", missingLinters[0].DependencyName(), missingLinters[0].DependencyName())

	if len(detectedLinters) == 0 {
		return fmt.Sprintf(`Oh my, no code quality linters were detected in your project!

We recommend that you start using the following linters in your project to help you measure and maintain the quality of your code:

%s

%s`, markdowngen.List(asInterfaceList(missingLinters)), extraDetails)
	}

	return fmt.Sprintf(`Linters detected:

%s

However, these linters were **missing** from your project:

%s

We recommend that you start using these linters in your project to help you measure and maintain the quality of your code.

%s`, markdowngen.List(asInterfaceList(detectedLinters)), markdowngen.List(asInterfaceList(missingLinters)), extraDetails)
}

var RuleLintersInstalled = api.Rule{
	Name: "All code quality linters should be installed in the current environment",
	Slug: "linters-installed",
	Details: fmt.Sprintf(`In order for mllint to be able to run the recommended code quality linters, they must be installed in the current environment,
i.e. they must be on PATH. 

This can be done in a variety of ways, such as installing them globally and / or appending to PATH,
but a more recommended way is to install them into a virtualenv, then activating this virtual environment and running mllint within it.
Poetry and Pipenv do this automatically, simply install them as development dependencies (%s) and run e.g. %s to open a shell in which to run mllint.`, "`--dev`", "`poetry shell`"),
	Weight: 1,
}
