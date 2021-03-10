package main

import (
	"os"

	"gitlab.com/bvobart/mllint/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
