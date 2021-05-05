package cqlinters_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/setools/cqlinters"
	"github.com/bvobart/mllint/setools/depmanagers"
)

func TestFromConfig(t *testing.T) {
	conf := config.CodeQualityConfig{Linters: []string{"pylint"}}
	expected := []api.CQLinter{cqlinters.Pylint{}}

	linters, err := cqlinters.FromConfig(conf)
	require.NoError(t, err)
	require.Equal(t, expected, linters)
}

func TestFromConfigNotFound(t *testing.T) {
	conf := config.CodeQualityConfig{Linters: []string{"eslint", "pylint"}}
	expected := []api.CQLinter{cqlinters.Pylint{}}

	linters, err := cqlinters.FromConfig(conf)
	require.EqualError(t, err, "could not parse these code quality linters from mllint's config: [eslint]")
	require.Equal(t, expected, linters)
}

func TestDetect(t *testing.T) {
	project := api.Project{Dir: "test-resources"}
	require.True(t, depmanagers.TypePoetry.Detect(project))
	project.DepManagers = []api.DependencyManager{depmanagers.TypePoetry.For(project)}

	linters := cqlinters.Detect(project)
	require.Equal(t, []api.CQLinter{cqlinters.Pylint{}}, linters)
}
