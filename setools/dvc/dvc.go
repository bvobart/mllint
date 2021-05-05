package dvc

import (
	"strings"

	"github.com/bvobart/mllint/utils/exec"
)

// Checks if DVC is installed (i.e. can be found on PATH)
func IsInstalled() bool {
	_, err := exec.LookPath("dvc")
	return err == nil
}

// Remotes runs 'dvc remote list' to figure out what remotes are configured for DVC in the current project.
func Remotes(dir string) []string {
	output, err := exec.CommandOutput(dir, "dvc", "remote", "list")
	if err != nil {
		return nil
	}

	remotes := []string{}
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		remotes = append(remotes, fields[0])
	}

	return remotes
}

// Files returns all files tracked by DVC in the current project.
func Files(dir string) []string {
	output, err := exec.CommandOutput(dir, "dvc", "list", ".", "-R", "--dvc-only")
	if err != nil {
		return nil
	}

	files := []string{}
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		filename := strings.TrimSpace(line)
		if len(filename) > 0 {
			files = append(files, filename)
		}
	}

	return files
}

// TODO: add function to check whether there are any DVC pipelines in a project.
