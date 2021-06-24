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
	Testing     TestingConfig     `yaml:"testing" toml:"testing"`
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

// CodeQualityConfig contains the configuration for the CQ linters used in the Code Quality category
type CodeQualityConfig struct {
	// Defines all code linters to use in the Code Quality category
	Linters []string `yaml:"linters" toml:"linters"`
}

// TestingConfig contains the configuration for the rules in the Testing category.
type TestingConfig struct {
	// Filename of the project's test execution report, either absolute or relative to the project's root.
	// Expects a JUnit XML file, which when using `pytest` can be generated with `pytest --junitxml=tests-report.xml`
	Report string `yaml:"report" toml:"report"`

	// Settings about how many tests there should be in a project.
	Targets TestingTargets `yaml:"targets" toml:"targets"`

	// Settings about the rules for checking project test coverage.
	Coverage TestCoverage `yaml:"coverage" toml:"coverage"`
}

type TestingTargets struct {
	// Minimum amount of test files to have in a project. Absolute number. Defaults to 1.
	Minimum uint64 `yaml:"minimum" toml:"minimum"`

	// Ratio of test files to have in a project, i.e. number of test files per other Python file.
	// Defaults to 1 part tests to 4 parts non-tests
	Ratio TestingTargetsRatio `yaml:"ratio" toml:"ratio"`
}

type TestingTargetsRatio struct {
	// Number of parts of test files.
	Tests uint64 `yaml:"tests" toml:"tests"`
	// Number of parts of other Python files.
	Other uint64 `yaml:"other" toml:"other"`
}

type TestCoverage struct {
	// Filename of the project's test coverage report, either absolute or relative to the project's root.
	// Expects a Cobertura-compatible XML file, which can be generated after `coverage run -m pytest --junitxml=tests-report.xml`
	// with `coverage xml -o tests-coverage.xml`, or using the `pytest-cov` plugin.
	Report string `yaml:"report" toml:"report"`

	// Specifies the target amount of line / branch / whatever coverage that the user wants want to have in the project
	// Only line coverage is implemented so far.
	Targets TestCoverageTargets `yaml:"targets" toml:"targets"`
}

type TestCoverageTargets struct {
	// Target amount of overall line coverage to achieve in tests.
	Line float64 `yaml:"line" toml:"line"`
}

//---------------------------------------------------------------------------------------

func Default() *Config {
	return &Config{
		Rules:       RuleConfig{Disabled: []string{}},
		Git:         GitConfig{MaxFileSize: 10_000_000}, // 10 MB
		CodeQuality: CodeQualityConfig{Linters: []string{"pylint", "mypy", "black", "isort", "bandit"}},
		Testing: TestingConfig{
			Targets: TestingTargets{
				Minimum: 1,
				Ratio: TestingTargetsRatio{
					Tests: 1,
					Other: 4,
				},
			},
			Coverage: TestCoverage{
				Targets: TestCoverageTargets{
					Line: 80,
				},
			},
		},
	}
}

//---------------------------------------------------------------------------------------

type FileType string

const (
	TypeDefault FileType = "default"
	TypeYAML    FileType = ".mllint.yml"
	TypeTOML    FileType = "pyproject.toml"
)

func (t FileType) String() string {
	if t == "" {
		return "unknown"
	}
	return string(t)
}

//---------------------------------------------------------------------------------------

// ParseFromDir parses the mllint config from the given project directory.
// If an `.mllint.yml` file is present, then this will be used,
// otherwise, if a `pyproject.toml` file is present, then this will be used,
// otherwise, the default config is returned.
// The returned FileType will be either config.TypeYAML, config.TypeTOML, or config.TypeDefault.
func ParseFromDir(projectdir string) (*Config, FileType, error) {
	yamlFile := path.Join(projectdir, string(TypeYAML))
	if utils.FileExists(yamlFile) {
		file, err := os.Open(yamlFile)
		if err != nil {
			return nil, "", fmt.Errorf("failed to open YAML config file '%s': %w", yamlFile, err)
		}
		conf, err := ParseYAML(file)
		if err != nil {
			return nil, "", fmt.Errorf("failed to parse YAML config file '%s': %w", yamlFile, err)
		}
		return conf, TypeYAML, nil
	}

	tomlFile := path.Join(projectdir, string(TypeTOML))
	if utils.FileExists(tomlFile) {
		file, err := os.Open(tomlFile)
		if err != nil {
			return nil, "", fmt.Errorf("failed to open TOML config from '%s': %w", tomlFile, err)
		}
		conf, err := ParseTOML(file)
		if err != nil {
			return nil, "", fmt.Errorf("failed to parse TOML config from '%s': %w", tomlFile, err)
		}
		return conf, TypeTOML, nil
	}

	return Default(), TypeDefault, nil
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

//---------------------------------------------------------------------------------------

func (conf *Config) YAML() ([]byte, error) {
	return yaml.Marshal(conf)
}

func (conf *Config) TOML() ([]byte, error) {
	pyproject := pyprojectTOML{}
	pyproject.Tool.Mllint = conf
	return toml.Marshal(pyproject)
}
