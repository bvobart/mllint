package main

import (
	"os"

	"github.com/bvobart/mllint/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
