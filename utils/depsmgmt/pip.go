package depsmgmt

import (
	"path"

	"gitlab.com/bvobart/mllint/utils"
)

type Pip struct{}

func (p Pip) Detect(projectdir string) bool {
	return utils.FileExists(path.Join(projectdir, "requirements.txt"))
}

func (p Pip) Type() DependencyManagerType {
	return TypePip
}
