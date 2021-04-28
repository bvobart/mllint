package depmanagers

import (
	"path"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/utils"
)

type RequirementsTxt struct{}

func (p RequirementsTxt) Detect(project api.Project) bool {
	return utils.FileExists(path.Join(project.Dir, "requirements.txt"))
}

func (p RequirementsTxt) Type() api.DependencyManagerType {
	return TypeRequirementsTxt
}

type SetupPy struct{}

func (p SetupPy) Detect(project api.Project) bool {
	return utils.FileExists(path.Join(project.Dir, "setup.py"))
}

func (p SetupPy) Type() api.DependencyManagerType {
	return TypeSetupPy
}
