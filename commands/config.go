package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gitlab.com/bvobart/mllint/config"
	"gitlab.com/bvobart/mllint/utils"
	"gopkg.in/yaml.v3"
)

func NewConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config [dir]",
		Short: "Prints the current mllint configuration.",
		Long: `Prints the mllint configuration as parsed from the '.mllint.yml' file in the root of the given (or current) directory, or the default configuration if none was found.
Specifying --quiet or -q will cause this command to purely print the current or default config, allowing for e.g. 'mllint config -q > .mllint.yml'`,
		RunE: runConfig,
		Args: cobra.MaximumNArgs(1),
	}
	return cmd
}

func runConfig(_ *cobra.Command, args []string) error {
	// catch `mllint config default`
	if len(args) == 1 && args[0] == "default" && !utils.FolderExists("default") {
		shush(func() { color.Green("Using default configuration\n\n") })
		output, err := yaml.Marshal(config.Default())
		if err != nil {
			return err
		}
		fmt.Println(string(output))
		return nil
	}

	projectdir, err := parseProjectDir(args)
	if err != nil {
		return err
	}
	shush(func() { color.Green("Using project at  %s", color.HiWhiteString(projectdir)) })

	conf, err := getConfig(projectdir)
	if err != nil {
		return err
	}

	output, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func getConfig(projectdir string) (*config.Config, error) {
	conf, err := config.ParseFromDir(projectdir)
	if err == nil {
		shush(func() { color.Green("Using configuration from project\n") })
		return conf, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		shush(func() { color.Yellow("No .mllint.yml found in project folder, using default configuration\n\n") })
		return config.Default(), nil
	}

	return nil, err
}
