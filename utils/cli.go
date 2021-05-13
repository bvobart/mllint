package utils

import (
	"os"

	"github.com/mattn/go-isatty"
)

func IsInteractive() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}
