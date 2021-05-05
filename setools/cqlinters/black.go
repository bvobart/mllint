package cqlinters

import (
	"path"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/depmanagers"
	"github.com/bvobart/mllint/utils"
	"github.com/bvobart/mllint/utils/exec"
)

type Black struct{}

func (p Black) Type() api.CQLinterType {
	return TypeBlack
}

func (p Black) String() string {
	return "Black"
}

func (p Black) DependencyName() string {
	return "black"
}

func (p Black) IsInstalled() bool {
	_, err := exec.LookPath("black")
	return err == nil
}

func (p Black) IsConfigured(project api.Project) bool {
	if !utils.FileExists(path.Join(project.Dir, "pyproject.toml")) {
		return false
	}

	poetry := depmanagers.TypePoetry.For(project).(depmanagers.Poetry)
	return poetry.HasConfigSection("tool.black")
}

func (p Black) Run(project api.Project) ([]api.CQLinterResult, error) {
	if len(project.PythonFiles) == 0 {
		return []api.CQLinterResult{}, nil
	}

	output, err := exec.CommandCombinedOutput(project.Dir, "black", "--check", project.Dir)
	if err == nil {
		return []api.CQLinterResult{}, nil
	}
	return []api.CQLinterResult{stringer(string(output))}, nil
}
