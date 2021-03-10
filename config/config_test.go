package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"gitlab.com/bvobart/mllint/config"
)

func TestParse(t *testing.T) {
	tests := []struct {
		Name     string
		Filename string
		Expected *config.Config
		ErrIs    error
		Err      error
	}{
		{Name: "EmptyFilename", Filename: "", Expected: nil, ErrIs: os.ErrNotExist},
		{Name: "FileDoesNotExist", Filename: "does-not-exist.yml", Expected: nil, ErrIs: os.ErrNotExist},
		{
			Name:     "EmptyFile",
			Filename: "test-resources/empty.config.yml",
			Expected: config.Default(),
			Err:      nil,
		},
		{
			Name:     "RulesDisabled",
			Filename: "test-resources/rules.config.yml",
			Expected: func() *config.Config {
				c := config.Default()
				c.Rules.Disabled = []string{"use-git", "another-rule"}
				return c
			}(),
			Err: nil,
		},
		{
			Name:     "YamlError",
			Filename: "test-resources/yamlerr.config.yml",
			Expected: nil,
			Err:      fmt.Errorf("YAML error in config file 'test-resources/yamlerr.config.yml': %w", &yaml.TypeError{Errors: []string{"line 2: cannot unmarshal !!str `nothing` into []string"}}),
		},
	}

	t.Parallel()
	for _, test := range tests {
		tt := test
		t.Run(tt.Name, func(t *testing.T) {
			conf, err := config.Parse(tt.Filename)

			if tt.ErrIs == nil && tt.Err == nil {
				require.NoError(t, err)
			}
			if tt.ErrIs != nil {
				require.ErrorIs(t, err, tt.ErrIs)
			}
			if tt.Err != nil {
				require.EqualError(t, err, tt.Err.Error())
			}

			require.Equal(t, tt.Expected, conf)
		})
	}
}

func TestParseFromDir(t *testing.T) {
	conf, err := config.ParseFromDir("test-resources")
	require.NoError(t, err)
	expectedConfig, err := config.Parse("test-resources/.mllint.yml")
	require.NoError(t, err)

	require.Equal(t, expectedConfig, conf)
}
