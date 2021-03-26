package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev-snapshot"
	commit  = "unknown"
	date    = "latest"
)

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints the version of mllint",
		Long:  "Prints the version of mllint",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("mllint version: %s\n", version)
			fmt.Printf("commit: %s (date: %s)\n", commit, date)
		},
		Args: cobra.ArbitraryArgs,
	}
}
