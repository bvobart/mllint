package depsmgmt

import (
	"path"

	"github.com/bvobart/mllint/utils"
)

type RequirementsTxt struct{}

func (p RequirementsTxt) Detect(projectdir string) bool {
	return utils.FileExists(path.Join(projectdir, "requirements.txt"))
}

func (p RequirementsTxt) Type() DependencyManagerType {
	return TypeRequirementsTxt
}

type SetupPy struct{}

func (p SetupPy) Detect(projectdir string) bool {
	return utils.FileExists(path.Join(projectdir, "setup.py"))
}

func (p SetupPy) Type() DependencyManagerType {
	return TypeSetupPy
}
