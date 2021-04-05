package exec

import "os/exec"

var (
	// LookPath is a function that performs `exec.LookPath`.
	// The sole purpose of this variable is to be able to mock calls to the exec module during tests.
	LookPath = DefaultLookPath

	// CommandOutput is a function that performs `exec.Command` in a certain dir, returning the command's Output().
	// The sole purpose of this variable is to be able to mock calls to the exec module during tests.
	CommandOutput = DefaultCommandOutput
)

// DefaultLookPath simply calls os/exec.LookPath
func DefaultLookPath(file string) (string, error) { return exec.LookPath(file) }

// DefaultLookupPath simply calls os/exec.Command, sets the directory and returns the command's Output().
func DefaultCommandOutput(dir string, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.Output()
}
