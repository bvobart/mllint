package depmanagers_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/depmanagers"
)

func TestRequirementsTxt(t *testing.T) {
	require.Equal(t, "requirements.txt", depmanagers.TypeRequirementsTxt.String())

	project := api.Project{Dir: "test-resources"}
	manager, err := depmanagers.TypeRequirementsTxt.Detect(project)
	require.NoError(t, err)
	require.Equal(t, depmanagers.TypeRequirementsTxt, manager.Type())

	require.True(t, manager.HasDependency("flask"))
	require.True(t, manager.HasDependency("numpy"))
	require.True(t, manager.HasDependency("pytest"))
	require.True(t, manager.HasDependency("pylint"))
	require.False(t, manager.HasDependency("mllint"))
}

func TestSetupPy(t *testing.T) {
	require.Equal(t, "setup.py", depmanagers.TypeSetupPy.String())

	project := api.Project{Dir: "."}
	manager, err := depmanagers.TypeSetupPy.Detect(project)
	require.ErrorIs(t, err, os.ErrNotExist)
	require.Nil(t, manager)

	project = api.Project{Dir: "test-resources"}
	manager, err = depmanagers.TypeSetupPy.Detect(project)
	require.NoError(t, err)
	require.NotNil(t, manager)
	require.Equal(t, depmanagers.TypeSetupPy, manager.Type())

	require.False(t, manager.HasDependency("")) // actually always returns false
}
