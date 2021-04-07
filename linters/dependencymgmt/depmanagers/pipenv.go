package depmanagers

import (
	"path"

	"github.com/bvobart/mllint/utils"
)

type Pipenv struct{}

func (p Pipenv) Detect(projectdir string) bool {
	return utils.FileExists(path.Join(projectdir, "Pipfile"))
}

func (p Pipenv) Type() DependencyManagerType {
	return TypePipenv
}
