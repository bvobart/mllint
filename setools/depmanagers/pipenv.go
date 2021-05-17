package depmanagers

import (
	"os"
	"path"

	"github.com/pelletier/go-toml"

	"github.com/bvobart/mllint/api"
)

type typePipenv string

func (p typePipenv) String() string {
	return string(p)
}

func (p typePipenv) Detect(project api.Project) (api.DependencyManager, error) {
	pipfilePath := path.Join(project.Dir, "Pipfile")
	pipfileContents, err := os.ReadFile(pipfilePath)
	if err != nil {
		return nil, err
	}

	pipfile := Pipfile{}
	if err := toml.Unmarshal(pipfileContents, &pipfile); err != nil {
		return nil, err
	}

	return Pipenv{Project: project, Pipfile: pipfile}, nil
}

//---------------------------------------------------------------------------------------

type Pipenv struct {
	Project api.Project
	Pipfile
}

type Pipfile struct {
	Packages    *toml.Tree `toml:"packages"`
	DevPackages *toml.Tree `toml:"dev-packages"`
}

func (p Pipenv) Type() api.DependencyManagerType {
	return TypePipenv
}

func (p Pipenv) HasDependency(dependency string) bool {
	return p.Pipfile.Packages != nil && p.Pipfile.Packages.Has(dependency) || p.HasDevDependency(dependency)
}

func (p Pipenv) HasDevDependency(dependency string) bool {
	return p.Pipfile.DevPackages != nil && p.Pipfile.DevPackages.Has(dependency)
}

func (p Pipenv) Dependencies() []string {
	deps := []string{}
	if p.Pipfile.Packages != nil {
		deps = append(deps, p.Pipfile.Packages.Keys()...)
	}
	if p.Pipfile.DevPackages != nil {
		deps = append(deps, p.Pipfile.DevPackages.Keys()...)
	}
	return deps
}
