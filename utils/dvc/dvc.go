package dvc

import (
	"os/exec"
	"strings"
)

var (
	// ExecLookupPath is a function that performs `exec.LookPath`.
	// The sole purpose of this variable is to be able to mock calls to the exec module during tests.
	ExecLookupPath = DefaultExecLookupPath
	// ExecCommandOutput is a function that performs `exec.Command` in a certain dir, returning the command's Output().
	// The sole purpose of this variable is to be able to mock calls to the exec module during tests.
	ExecCommandOutput = DefaultExecCommandOutput
)

// DefaultExecLookupPath simply calls os/exec.LookPath
func DefaultExecLookupPath(file string) (string, error) { return exec.LookPath(file) }

// DefaultExecLookupPath simply calls os/exec.Command, sets the directory and returns the command's Output().
func DefaultExecCommandOutput(dir string, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.Output()
}

//---------------------------------------------------------------------------------------

// Checks if DVC is installed (i.e. can be found on PATH)
func IsInstalled() bool {
	_, err := ExecLookupPath("dvc")
	return err == nil
}

// Remotes runs 'dvc remote list' to figure out what remotes are configured for DVC in the current project.
func Remotes(dir string) []string {
	output, err := ExecCommandOutput(dir, "dvc", "remote", "list")
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
	output, err := ExecCommandOutput(dir, "dvc", "list", ".", "-R", "--dvc-only")
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
