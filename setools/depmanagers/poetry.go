package depmanagers

import (
	"fmt"

	"github.com/bvobart/mllint/api"
)

//---------------------------------------------------------------------------------------

const poetryBuildBackend = "poetry.core.masonry.api"

type typePoetry string

func (p typePoetry) String() string {
	return string(p)
}

func (p typePoetry) Detect(project api.Project) (api.DependencyManager, error) {
	pyprojectToml, err := ReadPyProjectTOML(project.Dir)
	if err != nil {
		return nil, err
	}

	if pyprojectToml.BuildSystem.BuildBackend != poetryBuildBackend {
		return nil, fmt.Errorf("expecting build-system.build-backend to be '%s', but was: '%s'", poetryBuildBackend, pyprojectToml.BuildSystem.BuildBackend)
	}

	return Poetry{Config: pyprojectToml.Tool.Poetry}, nil
}

//---------------------------------------------------------------------------------------

type Poetry struct {
	Config *PoetryConfig
}

func (p Poetry) Type() api.DependencyManagerType {
	return TypePoetry
}

func (p Poetry) HasDependency(dependency string) bool {
	return p.Config != nil && p.Config.Dependencies != nil && p.Config.Dependencies.Has(dependency) || p.HasDevDependency(dependency)
}

func (p Poetry) HasDevDependency(dependency string) bool {
	return p.Config != nil && p.Config.Dependencies != nil && p.Config.DevDependencies.Has(dependency)
}

func (p Poetry) Dependencies() []string {
	deps := []string{}
	if p.Config != nil && p.Config.Dependencies != nil {
		deps = append(deps, p.Config.Dependencies.Keys()...)
	}
	if p.Config != nil && p.Config.DevDependencies != nil {
		deps = append(deps, p.Config.DevDependencies.Keys()...)
	}
	return deps
}

//---------------------------------------------------------------------------------------
