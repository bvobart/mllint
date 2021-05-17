package depmanagers_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/depmanagers"
)

func TestPoetryDetect(t *testing.T) {
	require.Equal(t, "Poetry", depmanagers.TypePoetry.String())

	project := api.Project{Dir: "."} // no pyproject.toml
	manager, err := depmanagers.TypePoetry.Detect(project)
	require.ErrorIs(t, err, os.ErrNotExist)
	require.Nil(t, manager)

	fakeProjectFile, err := os.Create("pyproject.toml") // empty pyproject.toml
	require.NoError(t, err)
	defer os.Remove(fakeProjectFile.Name())
	manager, err = depmanagers.TypePoetry.Detect(project)
	require.EqualError(t, err, `expecting build-system.build-backend to be 'poetry.core.masonry.api', but was: ''`)
	require.Nil(t, manager)

	project = api.Project{Dir: "test-resources"} // correct pyproject.toml
	manager, err = depmanagers.TypePoetry.Detect(project)
	require.NoError(t, err)
	require.Equal(t, depmanagers.TypePoetry, manager.Type())
}

func TestPoetryDependencies(t *testing.T) {
	project := api.Project{Dir: "test-resources"}
	manager, err := depmanagers.TypePoetry.Detect(project)
	require.NoError(t, err)

	expectedDeps := []string{"python", "pandas", "pyaml", "scikit-learn", "scipy", "mllint", "dvc", "pylint"}
	deps := manager.Dependencies()
	require.ElementsMatch(t, expectedDeps, deps)
	require.ElementsMatch(t, deps, expectedDeps)
}

func TestPoetryHasDependency(t *testing.T) {
	project := api.Project{Dir: "test-resources"}
	manager, err := depmanagers.TypePoetry.Detect(project)
	require.NoError(t, err)

	require.True(t, manager.HasDependency("pandas"))
	require.True(t, manager.HasDependency("scikit-learn"))
	require.True(t, manager.HasDependency("mllint"))
	require.True(t, manager.HasDependency("dvc"))
	require.False(t, manager.HasDependency("requires"))
}
