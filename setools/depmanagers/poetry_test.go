package depmanagers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/depmanagers"
)

func TestPoetry(t *testing.T) {
	require.Equal(t, "Poetry", depmanagers.TypePoetry.String())

	project := api.Project{Dir: "."}
	require.False(t, depmanagers.TypePoetry.Detect(project))

	project = api.Project{Dir: "test-resources"}
	require.True(t, depmanagers.TypePoetry.Detect(project))

	manager := depmanagers.TypePoetry.For(project)
	require.Equal(t, depmanagers.TypePoetry, manager.Type())

	require.True(t, manager.HasDependency("pandas"))
	require.True(t, manager.HasDependency("scikit-learn"))
	require.True(t, manager.HasDependency("mllint"))
	require.True(t, manager.HasDependency("dvc"))
	require.False(t, manager.HasDependency("requires"))
}
