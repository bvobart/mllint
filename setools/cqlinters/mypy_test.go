package cqlinters_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/cqlinters"
	"github.com/bvobart/mllint/utils"
	"github.com/bvobart/mllint/utils/exec"
)

func TestMypy(t *testing.T) {
	l := cqlinters.Mypy{}
	require.Equal(t, cqlinters.TypeMypy, l.Type())
	require.Equal(t, "Mypy", l.String())
	require.Equal(t, "mypy", l.DependencyName())

	exec.LookPath = func(file string) (string, error) { return "", errors.New("nope") }
	require.False(t, l.IsInstalled())
	exec.LookPath = func(file string) (string, error) { return "", nil }
	require.True(t, l.IsInstalled())
	exec.LookPath = exec.DefaultLookPath
}

const testMypyOutput = `src/evaluate.py:6: error: Cannot find implementation or library stub for module named 'sklearn.metrics'
src/evaluate.py:6: note: See https://mypy.readthedocs.io/en/latest/running_mypy.html#missing-imports
src/evaluate.py:6: error: Cannot find implementation or library stub for module named 'sklearn'
src/evaluate.py:37: error: Incompatible types in assignment (expression has type "TextIO", variable has type "BinaryIO")
Found 4 errors in 1 file (checked 2 source files)
`

const testMypySuccessOutput = "Success: no issues found in 1 source file"

var expectedMypyMessageStrings = [4]string{
	"src/evaluate.py:6: error: Cannot find implementation or library stub for module named 'sklearn.metrics'",
	"src/evaluate.py:6: note: See https://mypy.readthedocs.io/en/latest/running_mypy.html#missing-imports",
	"src/evaluate.py:6: error: Cannot find implementation or library stub for module named 'sklearn'",
	`src/evaluate.py:37: error: Incompatible types in assignment (expression has type "TextIO", variable has type "BinaryIO")`,
}

func TestMypyRun(t *testing.T) {
	l := cqlinters.Mypy{}
	t.Run("EmptyProject", func(t *testing.T) {
		results, err := l.Run(api.Project{})
		require.NoError(t, err)
		require.Equal(t, results, []api.CQLinterResult{})
	})

	t.Run("NormalProject+String", func(t *testing.T) {
		project := api.Project{
			Dir:         "test",
			PythonFiles: utils.Filenames{"file1", "file2", "file3"},
		}

		exec.CommandOutput = func(dir, name string, args ...string) ([]byte, error) {
			require.Equal(t, project.Dir, dir)
			require.Equal(t, "mypy", name)
			require.Equal(t, []string{project.Dir}, args)
			return []byte(testMypyOutput), errors.New("mypy always exits with an error when there are messages")
		}

		results, err := l.Run(project)
		require.NoError(t, err)
		require.Len(t, results, 4)
		for i, result := range results {
			require.Equal(t, expectedMypyMessageStrings[i], result.String())
		}
	})

	t.Run("NoMessages", func(t *testing.T) {
		project := api.Project{
			Dir:         "test",
			PythonFiles: utils.Filenames{"file1", "file2", "file3"},
		}

		exec.CommandOutput = func(dir, name string, args ...string) ([]byte, error) {
			require.Equal(t, project.Dir, dir)
			require.Equal(t, "mypy", name)
			require.Equal(t, []string{project.Dir}, args)
			return []byte(testMypySuccessOutput), nil
		}

		results, err := l.Run(project)
		require.NoError(t, err)
		require.Len(t, results, 0)
	})
}
