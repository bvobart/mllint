package dependencymgmt

// This array holds the names of all (known) dependencies that should always be dev dependencies, not normal dependencies.
// Of course, feel free to extend this if necessary!
var ShouldBeDevDependencies = []string{
	"mllint",
	"pylint",
	"mypy",
	"black",
	"isort",
	"bandit",
	"pylama",
	"flake8",
	"pyflakes",
	"mccabe",
	"tox",
	"dvc",
	"pytest",
}
