package depmanagers

import (
	"path"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/utils"
)

type Pipenv struct{}

func (p Pipenv) Detect(project api.Project) bool {
	return utils.FileExists(path.Join(project.Dir, "Pipfile"))
}

func (p Pipenv) Type() api.DependencyManagerType {
	return TypePipenv
}
