package cqlinters

import (
	"path"
	"strings"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/utils"
	"github.com/bvobart/mllint/utils/exec"
)

type Mypy struct{}

func (p Mypy) Type() api.CQLinterType {
	return TypeMypy
}

func (p Mypy) String() string {
	return "Mypy"
}

func (p Mypy) DependencyName() string {
	return "mypy"
}

func (p Mypy) IsInstalled() bool {
	_, err := exec.LookPath("mypy")
	return err == nil
}

func (p Mypy) IsConfigured(project api.Project) bool {
	return utils.FileExists(path.Join(project.Dir, "mypy.ini")) || utils.FileExists(path.Join(project.Dir, ".mypy.ini"))
}

func (p Mypy) Run(project api.Project) ([]api.CQLinterResult, error) {
	if len(project.PythonFiles) == 0 {
		return []api.CQLinterResult{}, nil
	}

	output, _ := exec.CommandOutput(project.Dir, "mypy", project.Dir)
	return decodeMypyOutput(output), nil
}

func decodeMypyOutput(output []byte) []api.CQLinterResult {
	msgs := strings.Split(string(output), "\n")
	msgs = msgs[:len(msgs)-2] // the last 2 lines are just "Found x errors in y files" and a blank line

	res := make([]api.CQLinterResult, len(msgs))
	for i, msg := range msgs {
		res[i] = stringer(msg)
	}
	return res
}

type stringer string

func (s stringer) String() string { return string(s) }
