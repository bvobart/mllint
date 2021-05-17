package depmanagers_test

import (
	"testing"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/depmanagers"
	"github.com/stretchr/testify/require"
)

func TestDetect(t *testing.T) {
	project := api.Project{Dir: "test-resources"}
	managers := depmanagers.Detect(project)

	require.Len(t, managers, 4)
	require.True(t, managers.ContainsAllTypes(depmanagers.TypePipenv, depmanagers.TypePoetry, depmanagers.TypeRequirementsTxt, depmanagers.TypeSetupPy))
}
