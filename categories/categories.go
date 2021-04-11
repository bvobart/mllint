package categories

import (
	"fmt"

	"github.com/bvobart/mllint/api"
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
	Description: `This category assesses your project's code quality.

It is not implemented yet. Examples of rules you might expect to see here in the future:
- Project uses linters (Pylint, Black, isort, mypy, Bandit, dslinter, mllint, etc). Will be configurable
- Project configures linters (make sure you have a configuration for the linters you employ. Default configuration should also be fine)
- rules about mllint's config. This will make it so there is little incentive to just disable all of mllint's rules, get a high score and just call it a day.
- will probably also include actually running the aforementioned linters to see what kind of issues they produce.
`,
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
	Description: `This category deals with the way your project is being tested.

It is not implemented yet. The idea is that this will contain some rules to check whether you have tests, what your latest test results were, how good your test coverage is,
and probably also something on whether you're actually testing your ML code or not.`,
}

var ContinuousIntegration = api.Category{
	Name: "Continuous Integration",
	Slug: "ci",
	Description: `This category evaluates your Continuous Integration (CI) setup.

It is not implemented yet, but will contain rules on whether you are using CI (GH Actions, Gitlab CI, Travis CI, Azure Pipelines).
Later, it may be expanded with rules on the structure of your pipelines (e.g. has stages build, test, deploy, that actually build, test and deploy the project.`,
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
