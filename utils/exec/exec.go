package exec

import (
	"bytes"
	"fmt"
	"os/exec"
)

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

	// PipelineOutput is a function that allows executing a pipeline of commands,
	// i.e. commands like `ls -l | grep exec | wc -l`
	PipelineOutput = DefaultPipelineOutput
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

func DefaultPipelineOutput(dir string, commands ...[]string) ([]byte, error) {
	if len(commands) == 0 {
		return []byte{}, nil
	}

	// create all command objects
	output := bytes.Buffer{}
	cmds := make([]*exec.Cmd, len(commands))
	for i, command := range commands {
		cmd := exec.Command(command[0], command[1:]...)
		cmd.Dir = dir
		cmds[i] = cmd

		// connect stdin of this command to the stdout of the previous command
		if i > 0 {
			var err error
			if cmd.Stdin, err = cmds[i-1].StdoutPipe(); err != nil {
				return nil, fmt.Errorf("failed to create output pipe: %w", err)
			}
		}

		// save the output of the last command to a buffer
		if i == len(commands)-1 {
			cmd.Stdout = &output
		}
	}

	// start all commands
	for i, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return nil, fmt.Errorf("failed to start command '%s': %w", commands[i], err)
		}
	}
	// then wait for each command to exit
	for i, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return nil, fmt.Errorf("command failed: '%s': %w", commands[i], err)
		}
	}
	return output.Bytes(), nil
}
