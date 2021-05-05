package depmanagers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/depmanagers"
)

func TestRequirementsTxt(t *testing.T) {
	require.Equal(t, "requirements.txt", depmanagers.TypeRequirementsTxt.String())

	project := api.Project{Dir: "test-resources"}
	require.True(t, depmanagers.TypeRequirementsTxt.Detect(project))

	manager := depmanagers.TypeRequirementsTxt.For(project)
	require.Equal(t, depmanagers.TypeRequirementsTxt, manager.Type())

	require.True(t, manager.HasDependency("flask"))
	require.True(t, manager.HasDependency("numpy"))
	require.True(t, manager.HasDependency("pytest"))
	require.True(t, manager.HasDependency("pylint"))
	require.False(t, manager.HasDependency("mllint"))
}

func TestSetupPy(t *testing.T) {
	require.Equal(t, "setup.py", depmanagers.TypeSetupPy.String())

	project := api.Project{Dir: "test-resources"}
	require.False(t, depmanagers.TypeSetupPy.Detect(project))

	manager := depmanagers.TypeSetupPy.For(project)
	require.Equal(t, depmanagers.TypeSetupPy, manager.Type())

	require.False(t, manager.HasDependency("")) // actually always returns false
}
