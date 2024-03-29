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
	conf := config.CodeQualityConfig{Linters: []string{"pylint", "mypy", "black", "isort", "bandit"}}
	expected := []api.CQLinter{cqlinters.Pylint{}, cqlinters.Mypy{}, cqlinters.Black{}, cqlinters.ISort{}, cqlinters.Bandit{}}

	linters, err := cqlinters.FromConfig(conf)
	require.NoError(t, err)
	require.Equal(t, expected, linters)
}

func TestFromConfigNotFound(t *testing.T) {
	conf := config.CodeQualityConfig{Linters: []string{"eslint", "pylint"}}
	expected := []api.CQLinter{cqlinters.Pylint{}}

	linters, err := cqlinters.FromConfig(conf)
	require.EqualError(t, err, "unknown code quality linters in mllint's config: [eslint]")
	require.Equal(t, expected, linters)
}

func TestDetect(t *testing.T) {
	project := api.Project{Dir: "test-resources"}
	poetry, err := depmanagers.TypePoetry.Detect(project)
	require.NoError(t, err)
	project.DepManagers = []api.DependencyManager{poetry}

	linters := cqlinters.Detect(project)
	require.Len(t, linters, 5)
	require.Subset(t, linters, []api.CQLinter{cqlinters.Pylint{}, cqlinters.Mypy{}, cqlinters.Black{}, cqlinters.ISort{}, cqlinters.Bandit{}})
}
