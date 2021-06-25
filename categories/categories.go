package categories

import (
	"fmt"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/cqlinters"
	"github.com/bvobart/mllint/utils/markdowngen"
)

var VersionControl = api.Category{
	Name: "Version Control",
	Slug: "version-control",
	Description: fmt.Sprintf(`This category contains rules relating to version controlling the code and data.
Version control software allows you to track changes to your project and helps to work collaboratively with other people within the same project.
It also allows you to easily return to an earlier version of your project or merge two versions together.

Git is the ubiquitously used tool for version controlling *code*, but Git is not very efficient at handling
large or binary files. It is therefore *not* directly possible to use Git to version control *data*.
Since *data* plays just as important of a role in ML as the *code* does, %s will also check how
a project manages its data. This can be done with a tool like [Data Version Control (DVC)](https://dvc.org).`,
		"`mllint`"),
}

var FileStructure = api.Category{
	Name: "File Structure",
	Slug: "file-structure",
	Description: `This category deals with the file and folder structure of your ML project.

It is not implemented yet. Examples of rules you might see here in the future:
- Project keeps its data in the './data' folder
- Project maintains documentation in a './docs' folder.
- Project's source code is kept in a './src' folder, or a folder with the same name as the project / package.`,
}

var DependencyMgmt = api.Category{
	Name: "Dependency Management",
	Slug: "dependency-management",
	Description: fmt.Sprintf(`This category deals with how your project manages its dependencies:
the Python packages that your project uses to make it work, such as %s.

Proper dependency management, i.e., properly specifying which packages your project uses and which exact versions of those packages are being used, 
is important for being able to recreate the environment that your project was developed in.
This allows other developers, automated deployment systems, or even just yourself, to install exactly those Python packages that you had installed while developing your project.
Therefore, there is no risk that they are either not able to run your code due to missing dependencies, or run into unexpected bugs caused by secretly updated dependencies.
In engineering terms, this relates to the concept of reproducibility: given the same project at the same revision with the same inputs, the same output should be produced.

Additionally, proper dependency management helps with the maintainability of your project.
In this case, that means how easy it will be later on to update the packages that your project uses,
but also to add new packages, remove old ones. This is especially useful for indirect dependencies,
as no-one has or likes to take the time to go through the changelogs of every sub-package you are using
to see if it is compatible with all other (sub-)packages.`,
		"`scikit-learn`, `pandas`, `tensorflow` and `pytorch`"),
}

var CodeQuality = api.Category{
	Name: "Code Quality",
	Slug: "code-quality",
	Description: `This category assesses your project's code quality by running several static analysis tools on your project.
Static analysis tools analyse your code without actually running it, in an attempt to find potential bugs, refactoring opportunities and/or coding style violations.

The linter for this category will check whether your project is using the configured set of code quality linters.
` + "`mllint`" + ` supports (and by default requires) the following linters:
` + markdowngen.List(asInterfaceList(cqlinters.AllTypes)) + `

For your project to be considered to be using a linter...
- **Either** there is a configuration file for this linter in the project
- **Or** the linter is a dependency of the project (preferably a dev dependency)

You can configure which linters ` + "`mllint`" + ` requires your project to use, using the following snippet of YAML in a ` + "`.mllint.yml`" + ` configuration file:
` + "```yaml" + `
code-quality:
	linters:
		- pylint
		- mypy
		- black
		- isort
		- bandit
` + "```" + `

or TOML:
` + "```toml" + `
[code-quality]
  linters = ["pylint", "mypy", "black", "isort", "bandit"]
` + "```" + `

We recommend that you configure each of these linters as you see fit using their respective configuration options.
Those will then automatically be picked up as ` + "`mllint`" + ` runs them.`,
}

func asInterfaceList(list []api.CQLinterType) []interface{} {
	res := make([]interface{}, len(list))
	for i := range list {
		res[i] = list[i]
	}
	return res
}

var DataQuality = api.Category{
	Name: "Data Quality",
	Slug: "data-quality",
	Description: `This category assesses your project's data quality.

It is not implemented yet. The idea is that this will contain rules on whether you have proper cleaning scripts
and may also include dynamic checks on the data that is currently in the repository
(e.g. is it complete (not missing values), are types of each value consistent, that sorta stuff. Perhaps with data-linter or tensorflow-data-validation)`,
}

var Testing = api.Category{
	Name: "Testing",
	Slug: "testing",
	Description: `Testing in the context of Software Engineering refers to the practice of writing automated checks to ensure that something works as intended.
Testing ML systems is, however, different from testing traditional software systems.
In traditional software systems, **humans write** all the logic that processes whatever data the system handles,
whereas in ML systems, **humans provide examples** (training data) of what we want the desired behaviour to be and the **machine learns** the logic required to produce this behaviour.

Properly testing ML systems is not only limited to testing the output behaviour of the system, but also entails, e.g.:
- ensuring that data preparation is done correctly and consistently
- ensuring that data featurisation is done correctly and consistent
- ensuring that the data is fed into the learning process correctly, e.g. testing helper functions
- ensuring that the learned logic consistently and accurately produces the desired behaviour

This category contains several rules relating to whether and to what degree you are testing the code of your ML project.
Per default, ` + "`mllint`" + ` expects **at least one test file** to be implemented in your project ` + "(i.e. a Python file starting with `test_` or ending with `_test.py`)" + `
and recommends that you have **at least 1 test file** for **every 4 non-test files**, though both these targets are configurable.
See the default configuration and the description of rule ` + "`testing/has-tests`" + ` for more information on how to configure this.

For ` + "`mllint`" + ` to be able to assess whether your project's tests pass and what coverage these tests achieve,
we will **not** actually run your tests. Instead, we expect you to run your project's tests yourself and provide 
the filenames to a JUnit-compatible XML test report and a Cobertura-compatible XML coverage report in your project's ` + "`mllint`" + ` configuration.
See the description of rule ` + "`testing/pass` and `testing/coverage`" + ` for more information on how to generate and configure these.

---

Here are some links to interesting blogs that give more in-depth information about different techniques for testing ML systems:
- [MadeWithML - Testing ML Systems: Code, Data and Models](https://madewithml.com/courses/mlops/testing/)
- [Jeremy Jordan - Effective testing for machine learning systems](https://www.jeremyjordan.me/testing-ml/)

> *"When writing tests for machine learning systems, one must not only test the student (the ML model), but also the teacher (the code that produces the ML model)." — Bart van Oort (bvobart)*`,
}

var ContinuousIntegration = api.Category{
	Name: "Continuous Integration",
	Slug: "ci",
	Description: fmt.Sprintf(`This category evaluates checks whether your project uses Continuous Integration (CI) and how you are using it.

Continuous Integration is the practice of automating the integration (merging) of all changes that multiple developers make to a software project.
This is done by running an automated process for every commit to your project's Git repository.
This process then downloads your project's source code at that commit, builds it, runs the linters configured for the project—we hope you include %s—and runs the project's tests against the system.

The core idea is that the CI server should be the unbiased arbiter of whether the project's code **works** after a certain set of changes, 
while providing a standardised environment to your whole team for verifying that the project truly works as intended.
No more 'but it worked on my machine' excuses.

Explore these sources to learn more about what CI entails:
- [WikiPedia - Continuous Integration](https://en.wikipedia.org/wiki/Continuous_integration)
- [ThoughtWorks](https://www.thoughtworks.com/continuous-integration), though a sales pitch, succinctly describes CI and several best practices relating to it, as well as its primary advantages.
- [SE4ML Best Practice - Use Continuous Integration](https://se-ml.github.io/best_practices/03-cont-int/)

To learn how to implement CI, see also the description of rule %s

Note: that this category is not fully implemented yet. It may later be expanded with rules on the structure of your CI pipelines (e.g. has stages build, test, deploy, that actually build, test and deploy the project.`,
		"`mllint`", "`ci/use`"),
}

var Deployment = api.Category{
	Name: "Deployment",
	Slug: "deployment",
	Description: `This category evaluates your project's ability to be deployed in the real world.

It is not yet implemented, but may contain rules about Dockerfiles and configurability, among others.`,
}

var All = []api.Category{
	VersionControl,
	FileStructure,
	DependencyMgmt,
	CodeQuality,
	DataQuality,
	Testing,
	ContinuousIntegration,
	Deployment,
}

var BySlug = makeSlugMap()

func makeSlugMap() map[string]api.Category {
	res := map[string]api.Category{}
	for _, cat := range All {
		res[cat.Slug] = cat
	}
	return res
}
