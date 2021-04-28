package depmanagers

import "github.com/bvobart/mllint/api"

const (
	TypePoetry          api.DependencyManagerType = "Poetry"
	TypePipenv          api.DependencyManagerType = "Pipenv"
	TypeRequirementsTxt api.DependencyManagerType = "requirements.txt"
	TypeSetupPy         api.DependencyManagerType = "setup.py"
)

var (
	poetry          api.DependencyManager = Poetry{}
	pipenv          api.DependencyManager = Pipenv{}
	requirementstxt api.DependencyManager = RequirementsTxt{}
	setuppy         api.DependencyManager = SetupPy{}
)

// Detect checks the files in the project directory to detect which dependency manager(s) is / are
// being used in this project.
func Detect(project api.Project) []api.DependencyManager {
	managers := []api.DependencyManager{}
	if poetry.Detect(project) {
		managers = append(managers, poetry)
	}

	if pipenv.Detect(project) {
		managers = append(managers, pipenv)
	}

	if requirementstxt.Detect(project) {
		managers = append(managers, requirementstxt)
	}

	if setuppy.Detect(project) {
		managers = append(managers, setuppy)
	}

	return managers
}
