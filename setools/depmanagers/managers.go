package depmanagers

import "github.com/bvobart/mllint/api"

var (
	TypePoetry          api.DependencyManagerType = typePoetry("Poetry")
	TypePipenv          api.DependencyManagerType = typePipenv("Pipenv")
	TypeRequirementsTxt api.DependencyManagerType = typeRequirementsTxt("requirements.txt")
	TypeSetupPy         api.DependencyManagerType = typeSetupPy("setup.py")
)

// all is ordered by how recommended each manager is, i.e. the first one in this list is the type of dependency manager we want to recommend the most.
var all = []api.DependencyManagerType{
	TypePoetry,
	TypePipenv,
	TypeRequirementsTxt,
	TypeSetupPy,
}

//---------------------------------------------------------------------------------------

// Detect checks the files in the project directory to detect which dependency manager(s) is / are
// being used in this project.
func Detect(project api.Project) api.DependencyManagerList {
	managers := api.DependencyManagerList{}

	for _, managerType := range all {
		if managerType.Detect(project) {
			managers = append(managers, managerType.For(project))
		}
	}

	return managers
}
