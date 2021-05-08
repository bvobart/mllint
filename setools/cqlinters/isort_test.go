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

func TestISort(t *testing.T) {
	l := cqlinters.ISort{}
	require.Equal(t, cqlinters.TypeISort, l.Type())
	require.Equal(t, "isort", l.String())
	require.Equal(t, "isort", l.DependencyName())

	exec.LookPath = func(file string) (string, error) { return "", errors.New("nope") }
	require.False(t, l.IsInstalled())
	exec.LookPath = func(file string) (string, error) { return "", nil }
	require.True(t, l.IsInstalled())
	exec.LookPath = exec.DefaultLookPath

	project := api.Project{Dir: "."}
	require.False(t, l.IsConfigured(project))
	project.Dir = "test-resources"
	require.True(t, l.IsConfigured(project))
}

const testISortOutput = `ERROR: src/evaluate.py Imports are incorrectly sorted and/or formatted.
ERROR: src/train.py Imports are incorrectly sorted and/or formatted.
ERROR: src/train_old.py Imports are incorrectly sorted and/or formatted.
ERROR: src/prepare.py Imports are incorrectly sorted and/or formatted.
ERROR: src/featurization.py Imports are incorrectly sorted and/or formatted.
Skipped 3 files
`

var expectedISortOutput = [5]string{
	"`src/evaluate.py` - Imports are incorrectly sorted and/or formatted.",
	"`src/train.py` - Imports are incorrectly sorted and/or formatted.",
	"`src/train_old.py` - Imports are incorrectly sorted and/or formatted.",
	"`src/prepare.py` - Imports are incorrectly sorted and/or formatted.",
	"`src/featurization.py` - Imports are incorrectly sorted and/or formatted.",
}

const testISortSuccessSkippedOutput = "Skipped 4 files\n"
const testISortSuccessEmptyOutput = "\n"

var exiterr = errors.New("isort always exits with an error when there are messages")

func TestISortRun(t *testing.T) {
	l := cqlinters.ISort{}

	// helper mocking function to assert that isort is called correctly
	isortCommandCombinedOutput := func(project api.Project, output string, err error) func(dir, name string, args ...string) ([]byte, error) {
		return func(dir, name string, args ...string) ([]byte, error) {
			require.Equal(t, project.Dir, dir)
			require.Equal(t, "isort", name)
			require.Equal(t, []string{"-c", project.Dir}, args)
			return []byte(output), err
		}
	}

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

		exec.CommandCombinedOutput = isortCommandCombinedOutput(project, testISortOutput, exiterr)

		results, err := l.Run(project)
		require.NoError(t, err)
		require.Len(t, results, 5)
		for i, result := range results {
			require.Equal(t, expectedISortOutput[i], result.String())
			require.IsType(t, cqlinters.ISortProblem{}, result)
		}
	})

	t.Run("SkippedFiles", func(t *testing.T) {
		project := api.Project{
			Dir:         "test",
			PythonFiles: utils.Filenames{"file1", "file2", "file3"},
		}

		exec.CommandCombinedOutput = isortCommandCombinedOutput(project, testISortSuccessSkippedOutput, nil)

		results, err := l.Run(project)
		require.NoError(t, err)
		require.Len(t, results, 0)
	})

	t.Run("EmptyOutput", func(t *testing.T) {
		project := api.Project{
			Dir:         "test",
			PythonFiles: utils.Filenames{"file1", "file2", "file3"},
		}

		exec.CommandCombinedOutput = isortCommandCombinedOutput(project, testISortSuccessEmptyOutput, nil)

		results, err := l.Run(project)
		require.NoError(t, err)
		require.Len(t, results, 0)
	})
}
