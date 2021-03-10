package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gitlab.com/bvobart/mllint/utils"
	"gopkg.in/yaml.v3"
)

// Config describes the structure of an `.mllint.yml` file
type Config struct {
	Rules RuleConfig `yaml:"rules"`
	Git   GitConfig  `yaml:"git"`
}

// RuleConfig contains info about which rules are enabled / disabled.
type RuleConfig struct {
	Disabled []string `yaml:"disabled"`
}

// GitConfig contains the configuration for the Git linters.
type GitConfig struct {
	// Maximum size of files in bytes tolerated by the 'git-no-big-files' linter
	// Default is 10 MB
	MaxFileSize uint64 `yaml:"maxFileSize"`
}

func Default() *Config {
	return &Config{
		Rules: RuleConfig{Disabled: []string{}},
		Git:   GitConfig{MaxFileSize: 10_000_000}, // 10 MB
	}
}

// ParseFromDir parses the `.mllint.yml` file in the given project directory.
func ParseFromDir(projectdir string) (*Config, error) {
	return Parse(path.Join(projectdir, ".mllint.yml"))
}

// Parse parses the config file at the given file location.
func Parse(filename string) (*Config, error) {
	if !utils.FileExists(filename) {
		return nil, fmt.Errorf("cannot parse config file '%s': %w", filename, os.ErrNotExist)
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read from config file '%s': %w", filename, err)
	}

	config := Default()
	if err := yaml.Unmarshal(contents, config); err != nil {
		return nil, fmt.Errorf("YAML error in config file '%s': %w", filename, err)
	}

	return config, nil
}
