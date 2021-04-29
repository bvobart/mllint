package cqlinters_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/cqlinters"
	"github.com/bvobart/mllint/setools/depmanagers"
	"github.com/bvobart/mllint/utils/exec"
)

func TestPylint(t *testing.T) {
	l := cqlinters.Pylint{}
	require.Equal(t, cqlinters.TypePylint, l.Type())
	require.Equal(t, "Pylint", l.String())

	exec.LookPath = func(file string) (string, error) { return "", errors.New("nope") }
	require.False(t, l.IsInstalled())
	exec.LookPath = func(file string) (string, error) { return "", nil }
	require.True(t, l.IsInstalled())
	exec.LookPath = exec.DefaultLookPath
}

func TestPylintDetect(t *testing.T) {
	l := cqlinters.Pylint{}

	project := api.Project{Dir: "."}
	require.False(t, l.Detect(project))

	project = api.Project{Dir: "test-resources"}
	require.True(t, depmanagers.TypePoetry.Detect(project))
	project.DepManagers = []api.DependencyManager{depmanagers.TypePoetry.For(project)}

	require.True(t, l.Detect(project))
}
