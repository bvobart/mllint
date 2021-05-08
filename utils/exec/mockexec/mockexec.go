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

func ExpectCommand(t *testing.T) *mockCommand {
	return &mockCommand{t: t}
}

type mockCommand struct {
	t    *testing.T
	dir  *string
	name *string
	args *[]string
}

func (mc *mockCommand) Dir(dir string) *mockCommand {
	mc.dir = &dir
	return mc
}

func (mc *mockCommand) CommandName(name string) *mockCommand {
	mc.name = &name
	return mc
}

func (mc *mockCommand) CommandArgs(args ...string) *mockCommand {
	mc.args = &args
	return mc
}

func (mc mockCommand) ToOutput(output []byte, err error) func(string, string, ...string) ([]byte, error) {
	return func(dir, name string, args ...string) ([]byte, error) {
		if mc.dir != nil {
			require.Equal(mc.t, *mc.dir, dir)
		}
		if mc.name != nil {
			require.Equal(mc.t, *mc.name, name)
		}
		if mc.args != nil {
			require.Equal(mc.t, *mc.args, args)
		}
		return output, err
	}
}
