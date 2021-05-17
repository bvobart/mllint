package depmanagers

import (
	"os"
	"path"

	"github.com/pelletier/go-toml"
)

type PyProjectTOML struct {
	Tool struct {
		Black  *toml.Tree    `toml:"black,omitempty"`
		ISort  *toml.Tree    `toml:"isort,omitempty"`
		Poetry *PoetryConfig `toml:"poetry,omitempty"`
	} `toml:"tool"`
	BuildSystem struct {
		BuildBackend string `toml:"build-backend"`
	} `toml:"build-system"`
}

type PoetryConfig struct {
	Dependencies    *toml.Tree `toml:"dependencies"`
	DevDependencies *toml.Tree `toml:"dev-dependencies"`
}

func ReadPyProjectTOML(dir string) (*PyProjectTOML, error) {
	filepath := path.Join(dir, "pyproject.toml")
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	pyprojectToml := PyProjectTOML{}
	if err := toml.Unmarshal(contents, &pyprojectToml); err != nil {
		return nil, err
	}

	return &pyprojectToml, nil
}
