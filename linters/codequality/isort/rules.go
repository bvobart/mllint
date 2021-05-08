package isort

import (
	"fmt"

	"github.com/bvobart/mllint/api"
)

var RuleNoIssues = api.Rule{
	Slug: "isort/no-issues",
	Name: "isort reports no issues with this project",
	Details: fmt.Sprintf(`> [%s](https://github.com/PyCQA/isort) is a Python utility / library to sort imports alphabetically, and automatically separated into sections and by type. It provides a command line utility, Python library and plugins for various editors to quickly sort all your imports.

This rule checks whether %s finds any files it would fix in your project.`, "`isort`", "`isort`"),
	Weight: 1,
}

var RuleIsConfigured = api.Rule{
	Slug: "isort/is-configured",
	Name: "isort is properly configured",
	Details: fmt.Sprintf(`[%s](https://github.com/PyCQA/isort) can be configured using several configuration files,
of which `+"`.isort.cfg` and `pyproject.toml` are preferred, according to `isort`'s documentation."+`
These are both recognised by mllint, although we recommend centralising tool configurations in your project's `+"`pyproject.toml`"+`

Since mllint also recommends using [Black](https://github.com/psf/black), you should configure `+"`isort`"+` to be compatible with Black.
This is done by putting the following in your `+"`pyproject.toml`"+`

`+"```toml"+`
[tool.isort]
profile = "black"
`+"```"+`

Links to `+"`isort`s"+` documentation:
- [Supported Config Files](https://pycqa.github.io/isort/docs/configuration/config_files/)
- [Black Compatibility](https://pycqa.github.io/isort/docs/configuration/black_compatibility/)
	`, "`isort`"),
}

var DetailsNotProperlyConfigured = `isort is not properly configured.
In order to be compatible with [Black](https://github.com/psf/black), which mllint also recommends using,
you should configure ` + "`isort`" + ` to use the ` + "`black`" + ` profile.
Furthermore, we recommend centralising your configuration in your ` + "`pyproject.toml`" + `

Thus, ensure that your ` + "`pyproject.toml`" + ` contains at least the following section:

` + "```toml" + `
[tool.isort]
profile = "black"
` + "```" + `
`
