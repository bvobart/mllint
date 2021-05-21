package bandit

import (
	"github.com/bvobart/mllint/api"
)

var RuleNoIssues = api.Rule{
	Slug: "code-quality/bandit/no-issues",
	Name: "Bandit reports no issues with this project",
	Details: `> [Bandit](https://github.com/PyCQA/bandit) is a tool designed to find common security issues in Python code.

This rule checks whether Bandit finds any security issues in your project.

For configuring Bandit's settings, such as which directories to exclude and which rules to enable / disable,
create a ` + "`.bandit`" + `file at the root of your project. See [Bandit's documentation](https://github.com/PyCQA/bandit#per-project-command-line-args) to learn more.`,
	Weight: 1,
}
