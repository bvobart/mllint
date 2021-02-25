package utils

import (
	"fmt"
	"os/exec"
)

// WrapExitError checks whether `err` is an *exec.ExitError and wraps it in a nice way
// Includes the Stderr output in the error message.
func WrapExitError(err error) error {
	if exiterr, ok := err.(*exec.ExitError); ok {
		return fmt.Errorf("exit code %d - %s", exiterr.ExitCode(), string(exiterr.Stderr))
	}
	return err
}

// WrapExitErrorf checks whether `err` is an *exec.ExitError and wraps it in a nice way
// Instead of the stderr output, this function allows you to specify your own message
// and possible formatting arguments.
func WrapExitErrorf(err error, message string, args ...interface{}) error {
	if exiterr, ok := err.(*exec.ExitError); ok {
		msg := fmt.Sprintf(message, args...)
		return fmt.Errorf("exit code %d - %s", exiterr.ExitCode(), msg)
	}
	return err
}
