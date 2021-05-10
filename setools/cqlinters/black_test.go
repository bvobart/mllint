package cqlinters_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/cqlinters"
	"github.com/bvobart/mllint/utils"
	"github.com/bvobart/mllint/utils/exec"
	"github.com/bvobart/mllint/utils/exec/mockexec"
)

func TestBlack(t *testing.T) {
	l := cqlinters.Black{}
	require.Equal(t, cqlinters.TypeBlack, l.Type())
	require.Equal(t, "Black", l.String())
	require.Equal(t, "black", l.DependencyName())

	exec.LookPath = mockexec.ExpectLookPath(t, "black").ToBeError()
	require.False(t, l.IsInstalled())
	exec.LookPath = mockexec.ExpectLookPath(t, "black").ToBeFound()
	require.True(t, l.IsInstalled())
	exec.LookPath = exec.DefaultLookPath

	project := api.Project{Dir: "."}
	require.False(t, l.IsConfigured(project))
	project.Dir = "test-resources"
	require.True(t, l.IsConfigured(project))
}

const testBlackOutput = `would reformat utils/test-resources/python-files/some_other_script.py
would reformat utils/test-resources/python-files/subfolder/yet_another_script.py
would reformat utils/test-resources/python-files/some_script.py
Oh no! 💥 💔 💥
3 files would be reformatted, 4 files would be left unchanged.
`

var expectedBlackOutput = [3]string{
	"`utils/test-resources/python-files/some_other_script.py`",
	"`utils/test-resources/python-files/subfolder/yet_another_script.py`",
	"`utils/test-resources/python-files/some_script.py`",
}

const testBlackSuccessOutput = `All done! ✨ 🍰 ✨
1 file would be left unchanged.
`

func TestBlackRun(t *testing.T) {
	l := cqlinters.Black{}
	t.Run("EmptyProject", func(t *testing.T) {
		results, err := l.Run(api.Project{})
		require.NoError(t, err)
		require.Equal(t, results, []api.CQLinterResult{})
	})

	t.Run("NormalProject+String", func(t *testing.T) {
		project := api.Project{
			Dir:         ".",
			PythonFiles: utils.Filenames{"file1", "file2", "file3"},
		}

		exec.CommandCombinedOutput = mockexec.ExpectCommand(t).Dir(project.Dir).
			CommandName("black").CommandArgs("--check", project.Dir).
			ToOutput([]byte(testBlackOutput), errors.New("black always exits with an error when there are messages"))

		results, err := l.Run(project)
		require.NoError(t, err)
		require.Len(t, results, 3)
		for i, result := range results {
			require.Equal(t, expectedBlackOutput[i], result.String())
		}
	})

	t.Run("NoMessages", func(t *testing.T) {
		project := api.Project{
			Dir:         ".",
			PythonFiles: utils.Filenames{"file1", "file2", "file3"},
		}

		exec.CommandCombinedOutput = mockexec.ExpectCommand(t).Dir(project.Dir).
			CommandName("black").CommandArgs("--check", project.Dir).
			ToOutput([]byte(testBlackSuccessOutput), nil)

		results, err := l.Run(project)
		require.NoError(t, err)
		require.Len(t, results, 0)
	})
}
