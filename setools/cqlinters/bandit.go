package cqlinters

import (
	"fmt"
	"path"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/utils"
	"github.com/bvobart/mllint/utils/exec"
	"gopkg.in/yaml.v3"
)

type Bandit struct{}

func (p Bandit) Type() api.CQLinterType {
	return TypeBandit
}

func (p Bandit) String() string {
	return "Bandit"
}

func (p Bandit) DependencyName() string {
	return "bandit"
}

func (p Bandit) IsInstalled() bool {
	_, err := exec.LookPath("bandit")
	return err == nil
}

func (p Bandit) IsConfigured(project api.Project) bool {
	return utils.FileExists(path.Join(project.Dir, ".bandit"))
}

func (p Bandit) IsProperlyConfigured(project api.Project) bool {
	return true // Bandit doesn't necessarily need to be configured.
}

func (p Bandit) Run(project api.Project) ([]api.CQLinterResult, error) {
	if len(project.PythonFiles) == 0 {
		return []api.CQLinterResult{}, nil
	}

	// We need to explicitly ignore the project's venv and .venv folders since Bandit doesn't do that by default
	// These folders also have to be referenced using their full path, see https://github.com/PyCQA/bandit/issues/488
	excludeDirs := path.Join(project.Dir, ".venv") + "," + path.Join(project.Dir, "venv")
	output, err := exec.CommandCombinedOutput(project.Dir, "bandit", "-f", "yaml", "-x", excludeDirs, "-r", project.Dir)
	if err == nil {
		return []api.CQLinterResult{}, nil
	}
	return decodeBanditOutput(output)
}

func decodeBanditOutput(output []byte) ([]api.CQLinterResult, error) {
	parsedOutput := banditYamlOutput{}
	if err := yaml.Unmarshal(output, &parsedOutput); err != nil {
		return nil, fmt.Errorf("failed to parse Bandit's YAML output: %w", err)
	}

	if len(parsedOutput.Errors) > 0 {
		return nil, fmt.Errorf("Bandit had errors: %v", parsedOutput.Errors)
	}

	results := make([]api.CQLinterResult, len(parsedOutput.Results))
	for i, result := range parsedOutput.Results {
		results[i] = result
	}
	return results, nil
}

type banditYamlOutput struct {
	Errors  []string        `yaml:"errors"`
	Results []BanditMessage `yaml:"results"`
}
