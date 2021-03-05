package depsmgmt

import (
	"path"

	"github.com/fatih/color"
	"github.com/pelletier/go-toml"

	"gitlab.com/bvobart/mllint/utils"
)

type Poetry struct{}

func (p Poetry) Detect(projectdir string) bool {
	poetryFile := path.Join(projectdir, "pyproject.toml")
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

func (p Poetry) Type() DependencyManagerType {
	return TypePoetry
}
