package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/utils"
)

func NewConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config [dir]",
		Short: "Prints the current mllint configuration.",
		Long: fmt.Sprintf(`Prints the mllint configuration as parsed from a configuration file in the root of the given (or current) directory. 

This can be either:
  - %s  Uses the YAML syntax as output by this command.
  - %s  Uses the TOML syntax configuration in the [tool.mllint] section. Has the same structure as the YAML
  - the default configuration if none of the files above was found.

Specifying %s or %s will cause this command to purely print the current or default config, allowing for e.g. %s`,
			formatInlineCode(string(config.TypeYAML)), formatInlineCode(string(config.TypeTOML)), formatInlineCode("--quiet"), formatInlineCode("-q"), formatInlineCode("mllint config -q > .mllint.yml")),
		RunE: runConfig,
		Args: cobra.MaximumNArgs(1),
	}
	return cmd
}

func runConfig(_ *cobra.Command, args []string) error {
	// catch `mllint config default`
	if len(args) == 1 && args[0] == "default" && !utils.FolderExists("default") {
		return runConfigDefault()
	}

	projectdir, err := parseProjectDir(args)
	if err != nil {
		return err
	}
	shush(func() { color.Green("Using project at  %s", color.HiWhiteString(projectdir)) })

	conf, _, err := getConfig(projectdir)
	if err != nil {
		return err
	}
	shush(func() { fmt.Print("---\n\n") })

	// print the config
	output, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	fmt.Println(string(output))

	shush(func() { fmt.Println("---") })
	return nil
}

func runConfigDefault() error {
	shush(func() { color.Green("Using default configuration\n\n") })

	output, err := yaml.Marshal(config.Default())
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}

// Parses the config from the project dir and prints a nice message about where it came from.
func getConfig(projectdir string) (*config.Config, config.FileType, error) {
	conf, typee, err := config.ParseFromDir(projectdir)
	if err != nil {
		return conf, typee, err
	}

	isDefault := cmp.Equal(conf, config.Default())
	if typee == config.TypeYAML || typee == config.TypeTOML {
		shush(func() { color.Green("Using configuration from %s (default: %v)\n", typee.String(), isDefault) })
	} else {
		shush(func() {
			color.Yellow("No .mllint.yml or pyproject.toml found in project folder, using default configuration\n")
		})
	}

	return conf, typee, err
}
