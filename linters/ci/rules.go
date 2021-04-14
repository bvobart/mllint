package ci

import "github.com/bvobart/mllint/api"

var RuleUseCI = api.Rule{
	Slug: "use",
	Name: "Project uses Continuous Integration (CI)",
	Details: `blablabla CI is important because ...

Something something CI providers Gitlab CI, GH Actions, Azure DevOps, Travis CI

Links to documentation`,
	Weight: 1,
}
