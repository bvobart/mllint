package cqlinters_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/setools/cqlinters"
	"github.com/bvobart/mllint/utils/exec"
)

func TestPylint(t *testing.T) {
	l := cqlinters.Pylint{}
	require.Equal(t, cqlinters.TypePylint, l.Type())
	require.Equal(t, "Pylint", l.String())
	require.Equal(t, "pylint", l.DependencyName())

	exec.LookPath = func(file string) (string, error) { return "", errors.New("nope") }
	require.False(t, l.IsInstalled())
	exec.LookPath = func(file string) (string, error) { return "", nil }
	require.True(t, l.IsInstalled())
	exec.LookPath = exec.DefaultLookPath
}
