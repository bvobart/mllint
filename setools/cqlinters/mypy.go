package cqlinters

import (
	"fmt"
	"path"
	"strconv"
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

func (p Mypy) IsProperlyConfigured(project api.Project) bool {
	return p.IsConfigured(project)
}

func (p Mypy) Run(project api.Project) ([]api.CQLinterResult, error) {
	if len(project.PythonFiles) == 0 {
		return []api.CQLinterResult{}, nil
	}

	// Enforce explicit ignoring of virtualenv folders.
	// Folders to be ignored taken from official Python Gitignore: https://github.com/github/gitignore/blob/991e760c1c6d50fdda246e0178b9c58b06770b90/Python.gitignore#L107
	excludeArg := `/(\.env|\.venv|env|venv|ENV|env\.bak|venv\.bak)/`
	output, _ := exec.CommandOutput(project.Dir, "mypy", project.Dir, "--exclude", excludeArg, "--strict", "--no-pretty", "--no-error-summary", "--no-color-output", "--hide-error-context", "--show-error-codes", "--show-column-numbers")
	return decodeMypyOutput(output)
}

func decodeMypyOutput(output []byte) ([]api.CQLinterResult, error) {
	msgs := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(msgs) < 2 {
		return []api.CQLinterResult{}, nil
	}

	var err error
	res := make([]api.CQLinterResult, len(msgs))
	for i, msg := range msgs {
		res[i], err = parseMypyMessage(msg)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Mypy message '%s': %w", msg, err)
		}
	}
	return res, nil
}

func parseMypyMessage(text string) (*MypyMessage, error) {
	parts := strings.Split(text, ":")

	if len(parts) >= 3 && len(parts) < 5 {
		return &MypyMessage{
			Filename: parts[0],
			Severity: strings.TrimSpace(parts[1]),
			Message:  strings.TrimSpace(strings.Join(parts[2:], ":")),
		}, nil
	}

	if len(parts) >= 5 {
		line, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("error parsing '%s' as line number: %w", parts[1], err)
		}
		column, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("error parsing '%s' as column number: %w", parts[2], err)
		}

		return &MypyMessage{
			Filename: parts[0],
			Line:     line,
			Column:   column,
			Severity: strings.TrimSpace(parts[3]),
			Message:  strings.TrimSpace(strings.Join(parts[4:], ":")),
		}, nil
	}

	return nil, fmt.Errorf("malformed Mypy message: expecting at least 3 or 5 parts separated by colons, but found %d", len(parts))
}
