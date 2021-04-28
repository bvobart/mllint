package depmanagers

import (
	"path"

	"github.com/fatih/color"
	"github.com/pelletier/go-toml"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/utils"
)

type Poetry struct{}

func (p Poetry) Detect(project api.Project) bool {
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

func (p Poetry) Type() api.DependencyManagerType {
	return TypePoetry
}
