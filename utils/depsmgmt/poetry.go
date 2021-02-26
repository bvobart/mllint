package depsmgmt

import (
	"path"

	"gitlab.com/bvobart/mllint/utils"
)

type Poetry struct{}

func (p Poetry) Detect(projectdir string) bool {
	return utils.FileExists(path.Join(projectdir, "pyproject.toml"))
}

func (p Poetry) Type() DependencyManagerType {
	return TypePoetry
}
