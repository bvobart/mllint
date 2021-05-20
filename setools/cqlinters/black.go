package cqlinters

import (
	"strings"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/depmanagers"
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
	pyprojectToml, err := depmanagers.ReadPyProjectTOML(project.Dir)
	if err != nil {
		return false
	}

	return pyprojectToml.Tool.Black != nil
}

func (p Black) IsProperlyConfigured(project api.Project) bool {
	return true // Black doesn't really need configuration
}

func (p Black) Run(project api.Project) ([]api.CQLinterResult, error) {
	if len(project.PythonFiles) == 0 {
		return []api.CQLinterResult{}, nil
	}

	output, err := exec.CommandCombinedOutput(project.Dir, "black", "--check", "--extend-exclude", ".venv", project.Dir)
	if err == nil {
		return []api.CQLinterResult{}, nil
	}
	return decodeBlackOutput(string(output), project.Dir), nil
}

func decodeBlackOutput(output string, projectdir string) []api.CQLinterResult {
	results := []api.CQLinterResult{}

	prefix := "would reformat "
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			results = append(results, formatFilename(strings.TrimPrefix(line, prefix), projectdir))
		}
	}

	return results
}

func formatFilename(filename string, projectdir string) stringer {
	return stringer("`" + trimProjectDir(filename, projectdir) + "`")
}
