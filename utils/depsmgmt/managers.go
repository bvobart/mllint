package depsmgmt

const (
	TypePoetry DependencyManagerType = "poetry"
	TypePipenv DependencyManagerType = "pipenv"
	TypePip    DependencyManagerType = "pip"
)

var (
	poetry DependencyManager = Poetry{}
	pipenv DependencyManager = Pipenv{}
	pip    DependencyManager = Pip{}
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

	if pip.Detect(projectdir) {
		managers = append(managers, pip)
	}

	return managers
}
