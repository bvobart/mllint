package cqlinters

import (
	"encoding/json"
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

func (p Pylint) DependencyName() string {
	return "pylint"
}

func (p Pylint) IsInstalled() bool {
	_, err := exec.LookPath("pylint")
	return err == nil
}

func (p Pylint) IsConfigured(project api.Project) bool {
	return utils.FileExists(path.Join(project.Dir, "pylintrc")) || utils.FileExists(path.Join(project.Dir, ".pylintrc"))
}

func (p Pylint) Run(project api.Project) ([]api.CQLinterResult, error) {
	if len(project.PythonFiles) == 0 {
		return []api.CQLinterResult{}, nil
	}

	pylintArgs := []string{"-f", "json"}
	pylintArgs = append(pylintArgs, project.PythonFiles...)
	output, _ := exec.CommandCombinedOutput(project.Dir, "pylint", pylintArgs...)
	// Pylint always exits with an error when there are messages, so we ignore the error.

	var messages pylintMessageList
	if err := json.Unmarshal(output, &messages); err != nil {
		return nil, fmt.Errorf("error parsing Pylint output '%s': %w", output, err)
	}

	return messages.Decode(), nil
}

type pylintMessageList []PylintMessage

func (messages pylintMessageList) Decode() []api.CQLinterResult {
	results := make([]api.CQLinterResult, len(messages))
	for i, msg := range messages {
		results[i] = msg
	}
	return results
}
