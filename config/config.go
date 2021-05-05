package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"

	"github.com/bvobart/mllint/utils"
)

// Config describes the structure of an `.mllint.yml` file
type Config struct {
	Rules       RuleConfig        `yaml:"rules" toml:"rules"`
	Git         GitConfig         `yaml:"git" toml:"git"`
	CodeQuality CodeQualityConfig `yaml:"code-quality" toml:"code-quality"`
}

// RuleConfig contains info about which rules are enabled / disabled.
type RuleConfig struct {
	Disabled []string `yaml:"disabled" toml:"disabled"`
}

// GitConfig contains the configuration for the Git linters.
type GitConfig struct {
	// Maximum size of files in bytes tolerated by the 'git-no-big-files' linter
	// Default is 10 MB
	MaxFileSize uint64 `yaml:"maxFileSize" toml:"maxFileSize"`
}

type CodeQualityConfig struct {
	// Defines all code linters to use in the Code Quality category
	Linters []string `yaml:"linters" toml:"linters"`
}

func Default() *Config {
	return &Config{
		Rules:       RuleConfig{Disabled: []string{}},
		Git:         GitConfig{MaxFileSize: 10_000_000},             // 10 MB
		CodeQuality: CodeQualityConfig{Linters: []string{"pylint"}}, // TODO: add more
	}
}

type FileType string

const (
	YAMLFile FileType = ".mllint.yml"
	TOMLFile FileType = "pyproject.toml"
)

func (t FileType) String() string {
	if t == "" {
		return "default"
	}
	return string(t)
}

// ParseFromDir parses the mllint config from the given project directory.
// If an `.mllint.yml` file is present, then this will be used,
// otherwise, if a `pyproject.toml` file is present, then this will be used,
// otherwise, the default config is returned.
// The returned FileType will be either config.YAMLFile, config.TOMLFile, or "".
func ParseFromDir(projectdir string) (*Config, FileType, error) {
	yamlFile := path.Join(projectdir, string(YAMLFile))
	if utils.FileExists(yamlFile) {
		file, err := os.Open(yamlFile)
		if err != nil {
			return nil, "", fmt.Errorf("failed to open YAML config file '%s': %w", yamlFile, err)
		}
		conf, err := ParseYAML(file)
		if err != nil {
			return nil, "", fmt.Errorf("failed to parse YAML config file '%s': %w", yamlFile, err)
		}
		return conf, YAMLFile, nil
	}

	tomlFile := path.Join(projectdir, string(TOMLFile))
	if utils.FileExists(tomlFile) {
		file, err := os.Open(tomlFile)
		if err != nil {
			return nil, "", fmt.Errorf("failed to open TOML config from '%s': %w", tomlFile, err)
		}
		conf, err := ParseTOML(file)
		if err != nil {
			return nil, "", fmt.Errorf("failed to parse TOML config from '%s': %w", tomlFile, err)
		}
		return conf, TOMLFile, nil
	}

	return Default(), "", nil
}

// ParseYAML parses the YAML config from the given reader (tip: *os.File implements io.Reader)
func ParseYAML(reader io.Reader) (*Config, error) {
	contents, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	config := Default()
	if err := yaml.Unmarshal(contents, config); err != nil {
		return nil, err
	}

	return config, nil
}

type pyprojectTOML struct {
	Tool struct {
		Mllint *Config `toml:"mllint"`
	} `toml:"tool"`
}

// ParseYAML parses the TOML config from the given reader (tip: *os.File implements io.Reader)
func ParseTOML(reader io.Reader) (*Config, error) {
	contents, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	tomlFile := pyprojectTOML{}
	tomlFile.Tool.Mllint = Default()
	if err := toml.Unmarshal(contents, &tomlFile); err != nil {
		return nil, err
	}

	return tomlFile.Tool.Mllint, nil
}
