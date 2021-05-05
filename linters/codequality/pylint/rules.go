package pylint

import (
	"fmt"

	"github.com/bvobart/mllint/api"
)

var RuleNoIssues = api.Rule{
	Slug: "pylint/no-issues",
	Name: "Pylint reports no issues with this project",
	Details: fmt.Sprintf(`[Pylint](https://pypi.org/project/pylint/) is a static analysis tool for finding generic programming errors.
This rule checks whether Pylint returns any errors when running it on all Python files in this project.

The score for this rule is determined as a function of the number of messages Pylint returns and the lines of Python code that your project has.
In the ideal case, Pylint does not recognise any code smells in your project, in which case the score is 100%%.
When there is one Pylint issue for every 20 lines of code, then the score is 50%%.
When there is one Pylint issue for every 10 lines of code, then the score is 0%%.

More specifically, in pseudocode, %s.

Note that the measured amount of lines of code includes any non-hidden Python files in the repository, including those that are ignored by Pylint.`,
		"`score = 100 - 100 * min(1, 10 * number of msgs / lines of code)`"),
	Weight: 1,
}
