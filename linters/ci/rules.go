package ci

import (
	"fmt"

	"github.com/bvobart/mllint/api"
)

var RuleUseCI = api.Rule{
	Slug: "ci/use",
	Name: "Project uses Continuous Integration (CI)",
	Details: fmt.Sprintf(`This rule checks if your project is using Continuous Integration (CI).
To learn more about what CI is, does and entails, see the description of category %s

Implementing CI requires picking a CI provider that will run the automated builds and tests. 
There are many CI providers available and you will have to make your own decision on which fits you best,
but %s currently recognises four CI providers, namely:
- [Azure DevOps Pipelines](https://docs.microsoft.com/en-us/azure/devops/pipelines/?view=azure-devops)
- [GitHub Actions](https://docs.github.com/en/actions)
- [Gitlab CI](https://docs.gitlab.com/ee/ci/)
- [Travis CI](https://docs.travis-ci.com/)

Follow your CI provider's respective 'Getting Started' guide and set your project up with a pipeline to build, test and lint your project.`,
		"`ci`", "`mllint`"), // TODO: add a ref to a to be implemented rule on how the structure of your pipelines should be.
	Weight: 1,
}
