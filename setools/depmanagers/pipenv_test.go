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
	require.True(t, depmanagers.TypePipenv.Detect(project))

	manager := depmanagers.TypePipenv.For(project)
	require.Equal(t, depmanagers.TypePipenv, manager.Type())

	require.True(t, manager.HasDependency("flask"))
	require.True(t, manager.HasDependency("numpy"))
	require.True(t, manager.HasDependency("pytest"))
	require.True(t, manager.HasDependency("pylint"))
	require.False(t, manager.HasDependency("mllint"))
}
