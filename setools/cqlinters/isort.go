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

	if utils.FileExists(path.Join(project.Dir, "pyproject.toml")) {
		poetry := depmanagers.TypePoetry.For(project).(depmanagers.Poetry)
		return poetry.Config().Has("tool.isort")
	}

	return false
}

func (p ISort) IsProperlyConfigured(project api.Project) bool {
	if utils.FileExists(path.Join(project.Dir, "pyproject.toml")) {
		poetry := depmanagers.TypePoetry.For(project).(depmanagers.Poetry)
		return poetry.Config().Get("tool.isort.profile") == "black"
	}

	return false
}

func (p ISort) Run(project api.Project) ([]api.CQLinterResult, error) {
	if len(project.PythonFiles) == 0 {
		return []api.CQLinterResult{}, nil
	}

	output, err := exec.CommandCombinedOutput(project.Dir, "isort", "-c", project.Dir)
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
