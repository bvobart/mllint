package depmanagers

import (
	"path"

	"github.com/fatih/color"
	"github.com/pelletier/go-toml"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/utils"
)

//---------------------------------------------------------------------------------------

type typePoetry string

func (p typePoetry) String() string {
	return string(p)
}

func (p typePoetry) Detect(project api.Project) bool {
	poetryFile := path.Join(project.Dir, "pyproject.toml")
	if !utils.FileExists(poetryFile) {
		return false
	}

	contents, err := toml.LoadFile(poetryFile)
	if err != nil {
		color.Red("Error: Poetry.Detect - failed to read %s: %s", poetryFile, err.Error())
		return false
	}

	backend := contents.Get("build-system.build-backend").(string)
	return backend == "poetry.core.masonry.api"
}

func (p typePoetry) For(project api.Project) api.DependencyManager {
	poetryFile := path.Join(project.Dir, "pyproject.toml")
	tomlConf, err := toml.LoadFile(poetryFile)
	if err != nil {
		panic(err)
	}
	return Poetry{Project: project, config: tomlConf}
}

//---------------------------------------------------------------------------------------

type Poetry struct {
	Project api.Project
	config  *toml.Tree
}

func (p Poetry) Type() api.DependencyManagerType {
	return TypePoetry
}

func (p Poetry) HasDependency(dependency string) bool {
	return p.config.Has("tool.poetry.dependencies."+dependency) || p.HasDevDependency(dependency)
}

func (p Poetry) HasDevDependency(dependency string) bool {
	return p.config.Has("tool.poetry.dev-dependencies." + dependency)
}

func (p Poetry) HasConfigSection(section string) bool {
	return p.config.Has(section)
}

//---------------------------------------------------------------------------------------
