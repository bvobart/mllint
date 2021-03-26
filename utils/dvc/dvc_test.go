package dvc_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/bvobart/mllint/utils/dvc"
)

func TestIsInstalled(t *testing.T) {
	dvc.ExecLookupPath = func(file string) (string, error) {
		require.Equal(t, "dvc", file)
		return file, nil
	}
	require.True(t, dvc.IsInstalled())

	dvc.ExecLookupPath = func(file string) (string, error) {
		require.Equal(t, "dvc", file)
		return file, fmt.Errorf("not on path or something")
	}
	require.False(t, dvc.IsInstalled())
}

func TestRemotes(t *testing.T) {
	t.Run("with remotes", func(t *testing.T) {
		dvc.ExecCommandOutput = func(dir, name string, args ...string) ([]byte, error) {
			require.Equal(t, ".", dir)
			require.Equal(t, "dvc", name)
			require.True(t, arrayEqual(args, []string{"remote", "list"}))
			return []byte(
				`
tmpfolder1       /tmp/dvcstore

tmpfolder2       /tmp/dvcstore
`,
			), nil
		}

		actualRemotes := dvc.Remotes(".")
		expectedRemotes := []string{"tmpfolder1", "tmpfolder2"}
		require.Equal(t, expectedRemotes, actualRemotes)
	})

	t.Run("no remotes", func(t *testing.T) {
		dvc.ExecCommandOutput = func(dir, name string, args ...string) ([]byte, error) {
			require.Equal(t, ".", dir)
			require.Equal(t, "dvc", name)
			require.True(t, arrayEqual(args, []string{"remote", "list"}))
			return []byte("\n"), nil
		}
		require.Equal(t, []string{}, dvc.Remotes("."))
	})

	t.Run("error", func(t *testing.T) {
		dvc.ExecCommandOutput = func(dir, name string, args ...string) ([]byte, error) {
			require.Equal(t, ".", dir)
			require.Equal(t, "dvc", name)
			require.True(t, arrayEqual(args, []string{"remote", "list"}))
			return []byte("\n"), fmt.Errorf("dvc not on path or something")
		}
		require.Nil(t, dvc.Remotes("."))
	})
}

func TestFiles(t *testing.T) {
	t.Run("with files", func(t *testing.T) {
		dvc.ExecCommandOutput = func(dir, name string, args ...string) ([]byte, error) {
			require.Equal(t, ".", dir)
			require.Equal(t, "dvc", name)
			require.True(t, arrayEqual(args, []string{"list", ".", "-R", "--dvc-only"}))
			return []byte(
				`
data/data.xml                                                                                                                                                                                                                                 
data/prepared/test.tsv
data/prepared/train.tsv
`,
			), nil
		}

		actualFiles := dvc.Files(".")
		expectedFiles := []string{"data/data.xml", "data/prepared/test.tsv", "data/prepared/train.tsv"}
		require.Equal(t, expectedFiles, actualFiles)
	})

	t.Run("no files", func(t *testing.T) {
		dvc.ExecCommandOutput = func(dir, name string, args ...string) ([]byte, error) {
			require.Equal(t, ".", dir)
			require.Equal(t, "dvc", name)
			require.True(t, arrayEqual(args, []string{"list", ".", "-R", "--dvc-only"}))
			return []byte("\n"), nil
		}
		require.Equal(t, []string{}, dvc.Files("."))
	})

	t.Run("error", func(t *testing.T) {
		dvc.ExecCommandOutput = func(dir, name string, args ...string) ([]byte, error) {
			require.Equal(t, ".", dir)
			require.Equal(t, "dvc", name)
			require.True(t, arrayEqual(args, []string{"list", ".", "-R", "--dvc-only"}))
			return []byte("\n"), fmt.Errorf("dvc not on path or something")
		}
		require.Nil(t, dvc.Files("."))
	})
}

func arrayEqual(one, other []string) bool {
	if len(one) != len(other) {
		return false
	}

	for i, valueOne := range one {
		if valueOne != other[i] {
			return false
		}
	}

	return true
}
