package utils

import (
	"errors"
	"fmt"
	"os/exec"
)

// WrapExitError checks whether `err` is an *exec.ExitError and wraps it in a nice way
// Includes the Stderr output in the error message.
func WrapExitError(err error) error {
	var exiterr *exec.ExitError
	if errors.As(err, &exiterr) {
		return fmt.Errorf("exit code %d - %s", exiterr.ExitCode(), string(exiterr.Stderr))
	}
	return err
}

// WrapExitErrorf checks whether `err` is an *exec.ExitError and wraps it in a nice way
// Instead of the stderr output, this function allows you to specify your own message
// and possible formatting arguments.
func WrapExitErrorf(err error, message string, args ...interface{}) error {
	var exiterr *exec.ExitError
	if errors.As(err, &exiterr) {
		msg := fmt.Sprintf(message, args...)
		return fmt.Errorf("exit code %d - %s", exiterr.ExitCode(), msg)
	}
	return err
}
