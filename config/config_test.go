package config_test

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/bvobart/mllint/config"
)

const yamlRulesDisabled = `
rules:
  disabled:
    - use-git
    - another-rule
`

const yamlLinters = `
code-quality:
  linters:
    - pylint
    - mypy
    - black
`

const yamlTesting = `
testing:
    report: junit-report.xml
    targets:
        minimum: 2
        ratio:
            tests: 2
            other: 8
    coverage:
        report: coverage.xml
        targets:
            line: 50 # percent line coverage.
`

const yamlCustomRule = `
rules:
  custom:
    - name: Custom Test Rule
      slug: custom/test-rule
      details: Tests whether parsing a custom rule from a YAML config works
      weight: 420
      run: python ./scripts/mllint-test-rule.py
    - name: Custom Test Rule
      slug: custom/test-rule
      details: Tests whether parsing a custom rule from a YAML config works
      weight: 420
      run: python ./scripts/mllint-test-rule.py
`

const yamlInvalid = `
rules:
  disabled: nothing
`

const tomlRulesDisabled = `
[tool.mllint]
  [tool.mllint.rules]
  disabled = ["use-git", "another-toml-rule"]
  
	[tool.mllint.git]
  maxFileSize = 1337
`

const tomlLinters = `
[tool.mllint.code-quality]
linters = ["pylint", "mypy"]
`

const tomlTesting = `
[tool.mllint.testing]
report = "tests-report.xml"
targets = { minimum = 2, ratio = { tests = 2, other = 8 }}
coverage = { report = "coverage.xml", targets = { line = 100.0 }}
`
const tomlCustomRule = `
[tool.mllint.rules]

[[tool.mllint.rules.custom]]
name = "Custom Test Rule"
slug = "custom/test-rule"
details = "Tests whether parsing a custom rule from a YAML config works"
weight = 420.0
run = "python ./scripts/mllint-test-rule.py"

[[tool.mllint.rules.custom]]
name = "Custom Test Rule"
slug = "custom/test-rule"
details = "Tests whether parsing a custom rule from a YAML config works"
weight = 420.0
run = "python ./scripts/mllint-test-rule.py"
`

const tomlInvalid = `
[tool.mllint.rules]
disabled = "nothing"
`

func TestParseYAML(t *testing.T) {
	finishedReader := strings.NewReader("")
	_, err := ioutil.ReadAll(finishedReader)
	require.NoError(t, err)

	tests := []struct {
		Name     string
		File     io.Reader
		Expected *config.Config
		Err      error
	}{
		{
			Name:     "FinishedReader",
			File:     finishedReader,
			Expected: config.Default(),
			Err:      nil,
		},
		{
			Name:     "EmptyFile",
			File:     strings.NewReader(""),
			Expected: config.Default(),
			Err:      nil,
		},
		{
			Name: "RulesDisabled",
			File: strings.NewReader(yamlRulesDisabled),
			Expected: func() *config.Config {
				c := config.Default()
				c.Rules.Disabled = []string{"use-git", "another-rule"}
				return c
			}(),
			Err: nil,
		},
		{
			Name: "YamlLinters",
			File: strings.NewReader(yamlLinters),
			Expected: func() *config.Config {
				c := config.Default()
				c.CodeQuality.Linters = []string{"pylint", "mypy", "black"}
				return c
			}(),
			Err: nil,
		},
		{
			Name: "YamlTesting",
			File: strings.NewReader(yamlTesting),
			Expected: func() *config.Config {
				c := config.Default()
				c.Testing.Report = "junit-report.xml"
				c.Testing.Targets.Minimum = 2
				c.Testing.Targets.Ratio.Tests = 2
				c.Testing.Targets.Ratio.Other = 8
				c.Testing.Coverage.Report = "coverage.xml"
				c.Testing.Coverage.Targets.Line = 50
				return c
			}(),
			Err: nil,
		},
		{
			Name: "YamlCustomRule",
			File: strings.NewReader(yamlCustomRule),
			Expected: func() *config.Config {
				c := config.Default()
				customRule := config.CustomRule{
					Name:    "Custom Test Rule",
					Slug:    "custom/test-rule",
					Details: "Tests whether parsing a custom rule from a YAML config works",
					Weight:  420,
					Run:     "python ./scripts/mllint-test-rule.py",
				}
				c.Rules.Custom = []config.CustomRule{customRule, customRule}
				return c
			}(),
			Err: nil,
		},
		{
			Name:     "YamlError",
			File:     strings.NewReader(yamlInvalid),
			Expected: nil,
			Err:      &yaml.TypeError{Errors: []string{"line 3: cannot unmarshal !!str `nothing` into []string"}},
		},
	}

	t.Parallel()
	for _, test := range tests {
		tt := test
		t.Run(tt.Name, func(t *testing.T) {
			conf, err := config.ParseYAML(tt.File)

			if tt.Err == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.Err.Error())
			}

			require.Equal(t, tt.Expected, conf)
		})
	}
}

func TestParseTOML(t *testing.T) {
	finishedReader := strings.NewReader("")
	_, err := ioutil.ReadAll(finishedReader)
	require.NoError(t, err)

	tests := []struct {
		Name     string
		File     io.Reader
		Expected *config.Config
		Err      error
	}{
		{
			Name:     "FinishedReader",
			File:     finishedReader,
			Expected: config.Default(),
			Err:      nil,
		},
		{
			Name:     "EmptyFile",
			File:     strings.NewReader(""),
			Expected: config.Default(),
			Err:      nil,
		},
		{
			Name: "RulesDisabled",
			File: strings.NewReader(tomlRulesDisabled),
			Expected: func() *config.Config {
				c := config.Default()
				c.Rules.Disabled = []string{"use-git", "another-toml-rule"}
				c.Git.MaxFileSize = 1337
				return c
			}(),
			Err: nil,
		},
		{
			Name: "TomlLinters",
			File: strings.NewReader(tomlLinters),
			Expected: func() *config.Config {
				c := config.Default()
				c.CodeQuality.Linters = []string{"pylint", "mypy"}
				return c
			}(),
			Err: nil,
		},
		{
			Name: "TomlTesting",
			File: strings.NewReader(tomlTesting),
			Expected: func() *config.Config {
				c := config.Default()
				c.Testing.Report = "tests-report.xml"
				c.Testing.Targets.Minimum = 2
				c.Testing.Targets.Ratio.Tests = 2
				c.Testing.Targets.Ratio.Other = 8
				c.Testing.Coverage.Report = "coverage.xml"
				c.Testing.Coverage.Targets.Line = 100
				return c
			}(),
			Err: nil,
		},
		{
			Name: "TomlCustomRule",
			File: strings.NewReader(tomlCustomRule),
			Expected: func() *config.Config {
				c := config.Default()
				customRule := config.CustomRule{
					Name:    "Custom Test Rule",
					Slug:    "custom/test-rule",
					Details: "Tests whether parsing a custom rule from a YAML config works",
					Weight:  420,
					Run:     "python ./scripts/mllint-test-rule.py",
				}
				c.Rules.Custom = []config.CustomRule{customRule, customRule}
				return c
			}(),
			Err: nil,
		},
		{
			Name:     "TomlError",
			File:     strings.NewReader(tomlInvalid),
			Expected: nil,
			Err:      errors.New("(3, 1): Can't convert nothing(string) to []string(slice)"),
		},
	}

	t.Parallel()
	for _, test := range tests {
		tt := test
		t.Run(tt.Name, func(t *testing.T) {
			conf, err := config.ParseTOML(tt.File)

			if tt.Err == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.Err.Error())
			}

			require.Equal(t, tt.Expected, conf)
		})
	}
}

func TestParseFromDir(t *testing.T) {
	tests := []struct {
		Name     string // test name
		Dir      string // dir to parse config from
		Expected *config.Config
		Type     config.FileType
	}{
		{
			Name: ".mllint.yml",
			Dir:  "test-resources/yaml",
			Type: config.TypeYAML,
			Expected: func() *config.Config {
				// preconditions, check whether .mllint.yml file is present in test dir.
				filename := path.Join("test-resources/yaml", ".mllint.yml")
				require.FileExists(t, filename)
				configfile, err := os.Open(filename)
				require.NoError(t, err)

				expectedConfig, err := config.ParseYAML(configfile)
				require.NoError(t, err)
				return expectedConfig
			}(),
		},
		{
			Name: "pyproject.toml",
			Dir:  "test-resources/toml",
			Type: config.TypeTOML,
			Expected: func() *config.Config {
				// preconditions, check whether pyproject.toml is present in test dir.
				filename := path.Join("test-resources/toml", "pyproject.toml")
				require.FileExists(t, filename)
				configfile, err := os.Open(filename)
				require.NoError(t, err)

				expectedConfig, err := config.ParseTOML(configfile)
				require.NoError(t, err)
				return expectedConfig
			}(),
		},
		{
			Name: "precedence", // tests that a config specified in a .mllint.yml has precendence over specifying them in the pyproject.toml.
			Dir:  "test-resources/precedence",
			Type: config.TypeYAML,
			Expected: func() *config.Config {
				// preconditions, check whether .mllint.yml and pyproject.toml are present in test dir.
				yamlFile := path.Join("test-resources/precedence", ".mllint.yml")
				tomlFile := path.Join("test-resources/precedence", "pyproject.toml")
				require.FileExists(t, yamlFile)
				require.FileExists(t, tomlFile)

				configfile, err := os.Open(yamlFile)
				require.NoError(t, err)

				expectedConfig, err := config.ParseYAML(configfile)
				require.NoError(t, err)
				return expectedConfig
			}(),
		},
		{
			Name:     "EmptyDir",
			Dir:      "test-resources",
			Expected: config.Default(),
			Type:     config.TypeDefault,
		},
	}

	t.Parallel()
	for _, test := range tests {
		tt := test
		t.Run(tt.Name, func(t *testing.T) {
			conf, typee, err := config.ParseFromDir(tt.Dir)
			require.NoError(t, err)
			require.Equal(t, tt.Expected, conf)
			require.Equal(t, tt.Type, typee)
		})
	}
}

func TestConfigTypeString(t *testing.T) {
	require.Equal(t, "unknown", config.FileType("").String())
	require.Equal(t, string(config.TypeDefault), config.TypeDefault.String())
	require.Equal(t, string(config.TypeYAML), config.TypeYAML.String())
	require.Equal(t, string(config.TypeTOML), config.TypeTOML.String())
}
