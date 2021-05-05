package black

import (
	"fmt"

	"github.com/bvobart/mllint/api"
)

var RuleNoIssues = api.Rule{
	Slug: "black/no-issues",
	Name: "Black reports no issues with this project",
	Details: fmt.Sprintf(`> [Black](https://github.com/psf/black) is the uncompromising Python code formatter. By using it, you agree to cede control over minutiae of hand-formatting. In return, Black gives you speed, determinism, and freedom from %s nagging about formatting. You will save time and mental energy for more important matters.

This rule checks whether Black finds any files it would fix in your project.`, "`pycodestyle`"),
	Weight: 1,
}
