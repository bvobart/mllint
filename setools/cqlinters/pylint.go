package cqlinters

import (
	"fmt"
	"path"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/utils"
	"github.com/bvobart/mllint/utils/exec"
)

type Pylint struct{}

func (p Pylint) Type() api.CQLinterType {
	return TypePylint
}

func (p Pylint) String() string {
	return "Pylint"
}

func (p Pylint) IsInstalled() bool {
	_, err := exec.LookPath("pylint")
	return err == nil
}

func (p Pylint) Detect(project api.Project) bool {
	if len(project.DepManagers) > 0 && project.DepManagers.Main().HasDependency("pylint") {
		return true
	}

	if utils.FileExists(path.Join(project.Dir, "pylintrc")) || utils.FileExists(path.Join(project.Dir, ".pylintrc")) {
		return true
	}

	return false
}

func (p Pylint) Run(projectdir string) error {
	// TODO: copy from python-ml-analysis
	return fmt.Errorf("not implemented")
}
