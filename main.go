package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"gitlab.com/bvobart/mllint/commands"
)

func main() {
	if err := Execute(); err != nil {
		os.Exit(1)
	}
}

func Execute() error {
	startTime := time.Now()
	err := commands.NewRootCommand().Execute()
	if err != nil && errors.Is(err, commands.ErrIssuesFound) {
		color.HiWhite("%s", err)
	} else if err != nil {
		color.Red("Fatal: %s", err)
	}
	fmt.Println("took:", time.Since(startTime))
	return err
}
