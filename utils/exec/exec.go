package exec

import "os/exec"

var (
	// LookPath is a function that performs `exec.LookPath`.
	// The sole purpose of this variable is to be able to mock calls to the exec module during tests.
	//
	// LookPath searches for an executable named file in the directories named by the PATH environment variable.
	// If file contains a slash, it is tried directly and the PATH is not consulted.
	// The result may be an absolute path or a path relative to the current directory.
	// (from exec.LookPath docstring)
	LookPath = DefaultLookPath

	// CommandOutput is a function that performs `exec.Command` in a certain dir, returning the command's Output().
	// The sole purpose of this variable is to be able to mock calls to the exec module during tests.
	CommandOutput = DefaultCommandOutput

	// CommandCombinedOutput is a function that performs `exec.Command` in a certain dir, returning the command's CombinedOutput().
	// The sole purpose of this variable is to be able to mock calls to the exec module during tests.
	CommandCombinedOutput = DefaultCommandCombinedOutput
)

// DefaultLookPath simply calls os/exec.LookPath
func DefaultLookPath(file string) (string, error) { return exec.LookPath(file) }

// DefaultCommandOutput simply calls os/exec.Command, sets the directory and returns the command's Output().
func DefaultCommandOutput(dir string, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.Output()
}

// DefaultCommandOutput simply calls os/exec.Command, sets the directory and returns the command's Output().
func DefaultCommandCombinedOutput(dir string, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}
