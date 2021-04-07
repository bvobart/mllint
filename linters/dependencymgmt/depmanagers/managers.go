package depmanagers

const (
	TypePoetry          DependencyManagerType = "Poetry"
	TypePipenv          DependencyManagerType = "Pipenv"
	TypeRequirementsTxt DependencyManagerType = "requirements.txt"
	TypeSetupPy         DependencyManagerType = "setup.py"
)

var (
	poetry          DependencyManager = Poetry{}
	pipenv          DependencyManager = Pipenv{}
	requirementstxt DependencyManager = RequirementsTxt{}
	setuppy         DependencyManager = SetupPy{}
)

type DependencyManagerType string
type DependencyManager interface {
	Detect(projectdir string) bool
	Type() DependencyManagerType
}

// Detect checks the files in the project directory to detect which dependency manager(s) is / are
// being used in this project.
func Detect(projectdir string) []DependencyManager {
	managers := []DependencyManager{}
	if poetry.Detect(projectdir) {
		managers = append(managers, poetry)
	}

	if pipenv.Detect(projectdir) {
		managers = append(managers, pipenv)
	}

	if requirementstxt.Detect(projectdir) {
		managers = append(managers, requirementstxt)
	}

	if setuppy.Detect(projectdir) {
		managers = append(managers, setuppy)
	}

	return managers
}
