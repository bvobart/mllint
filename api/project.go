package api

import (
	"github.com/hashicorp/go-multierror"

	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/utils"
)

// Project contains general information about the project that will be filled in before the linters start their analysis.
type Project struct {
	// The project's assumed root directory, absolute path.
	Dir string
	// Type of mllint configuration
	ConfigType config.FileType
	// Dependency managers that this project uses, e.g. requirements.txt, Poetry or Pipenv
	DepManagers DependencyManagerList
	// Code Quality linters that this project uses, i.e. static analysis tools that focus on analysing code, such as Pylint, Mypy and Bandit.
	CQLinters []CQLinter
	// Absolute paths to the Python files that are in this project's repository
	PythonFiles utils.Filenames
}

// ProjectReport is what you end up with after mllint finishes analysing a project.
type ProjectReport struct {
	Project
	Reports map[Category]Report
	Errors  *multierror.Error
}
