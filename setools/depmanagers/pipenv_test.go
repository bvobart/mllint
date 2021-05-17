package depmanagers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/depmanagers"
)

func TestPipenv(t *testing.T) {
	require.Equal(t, "Pipenv", depmanagers.TypePipenv.String())

	project := api.Project{Dir: "test-resources"}
	manager, err := depmanagers.TypePipenv.Detect(project)
	require.NoError(t, err)
	require.Equal(t, depmanagers.TypePipenv, manager.Type())

	require.True(t, manager.HasDependency("flask"))
	require.True(t, manager.HasDependency("numpy"))
	require.True(t, manager.HasDependency("pytest"))
	require.True(t, manager.HasDependency("pylint"))
	require.False(t, manager.HasDependency("mllint"))
}

func TestPipenvDependencies(t *testing.T) {
	project := api.Project{Dir: "test-resources"}
	manager, err := depmanagers.TypePipenv.Detect(project)
	require.NoError(t, err)

	expectedDeps := []string{"flask", "numpy", "requests", "pytest", "pylint"}
	deps := manager.Dependencies()
	require.ElementsMatch(t, expectedDeps, deps)
	require.ElementsMatch(t, deps, expectedDeps)
}
