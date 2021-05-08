package mockexec

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func ExpectLookPath(t *testing.T, expectedFile string) mockLookPath {
	return mockLookPath{t, expectedFile}
}

type mockLookPath struct {
	t            *testing.T
	expectedFile string
}

func (mlp mockLookPath) ToBeFound() func(string) (string, error) {
	return func(file string) (string, error) {
		require.Equal(mlp.t, mlp.expectedFile, file)
		return "", nil
	}
}

func (mlp mockLookPath) ToBeError() func(string) (string, error) {
	return func(file string) (string, error) {
		require.Equal(mlp.t, mlp.expectedFile, file)
		return "", errors.New("not found")
	}
}
