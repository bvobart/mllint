package depmanagers

import (
	"path"

	"github.com/pelletier/go-toml"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/utils"
)

type typePipenv string

func (p typePipenv) String() string {
	return string(p)
}

func (p typePipenv) Detect(project api.Project) bool {
	return utils.FileExists(path.Join(project.Dir, "Pipfile"))
}

func (p typePipenv) For(project api.Project) api.DependencyManager {
	pipfilePath := path.Join(project.Dir, "Pipfile")
	pipfile, err := toml.LoadFile(pipfilePath)
	if err != nil {
		panic(err)
	}
	return Pipenv{Project: project, pipfile: pipfile}
}

//---------------------------------------------------------------------------------------

type Pipenv struct {
	Project api.Project
	pipfile *toml.Tree
}

func (p Pipenv) Type() api.DependencyManagerType {
	return TypePipenv
}

func (p Pipenv) HasDependency(dependency string) bool {
	return p.pipfile.Has("packages."+dependency) || p.pipfile.Has("dev-packages."+dependency)
}
