package cqlinters

import (
	"path"
	"strings"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/depmanagers"
	"github.com/bvobart/mllint/utils"
	"github.com/bvobart/mllint/utils/exec"
)

type ISort struct{}

func (p ISort) Type() api.CQLinterType {
	return TypeISort
}

func (p ISort) String() string {
	return "isort"
}

func (p ISort) DependencyName() string {
	return "isort"
}

func (p ISort) IsInstalled() bool {
	_, err := exec.LookPath("isort")
	return err == nil
}

func (p ISort) IsConfigured(project api.Project) bool {
	if utils.FileExists(path.Join(project.Dir, ".isort.cfg")) {
		return true
	}

	if pyprojectToml, err := depmanagers.ReadPyProjectTOML(project.Dir); err == nil {
		return pyprojectToml.Tool.ISort != nil
	}

	return false
}

func (p ISort) IsProperlyConfigured(project api.Project) bool {
	if pyprojectToml, err := depmanagers.ReadPyProjectTOML(project.Dir); err == nil {
		return pyprojectToml.Tool.ISort != nil && pyprojectToml.Tool.ISort.Get("profile") == "black"
	}

	return false
}

func (p ISort) Run(project api.Project) ([]api.CQLinterResult, error) {
	if len(project.PythonFiles) == 0 {
		return []api.CQLinterResult{}, nil
	}

	// Enforce explicit ignoring of virtualenv folders.
	// Folders to be ignored taken from official Python Gitignore: https://github.com/github/gitignore/blob/991e760c1c6d50fdda246e0178b9c58b06770b90/Python.gitignore#L107
	excludeDirs := []string{".env", ".venv", "env", "venv", "ENV", "env.bak", "venv.bak"}
	excludeArgs := make([]string, 0, len(excludeDirs)*2)
	for _, excludeDir := range excludeDirs {
		excludeArgs = append(excludeArgs, "--extend-skip", excludeDir)
	}

	args := append([]string{"-c", project.Dir}, excludeArgs...)
	output, err := exec.CommandCombinedOutput(project.Dir, "isort", args...)
	if err == nil {
		return []api.CQLinterResult{}, nil
	}
	return decodeISortOutput(string(output), project.Dir), nil
}

func decodeISortOutput(output string, projectdir string) []api.CQLinterResult {
	results := []api.CQLinterResult{}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ERROR:") {
			parts := strings.Split(line, " ")
			problem := ISortProblem{
				Path:    trimProjectDir(parts[1], projectdir),
				Message: strings.Join(parts[2:], " "),
			}
			results = append(results, problem)
		}
	}

	return results
}
