package projectlinters_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/projectlinters"
	"github.com/bvobart/mllint/utils/dvc"
)

func TestUseDVC(t *testing.T) {
	linter := projectlinters.UseDVC{}
	expectLintIssues := func(dir string, expected []api.Issue) {
		issues, err := linter.LintProject(dir)
		require.NoError(t, err)
		require.Equal(t, expected, issues)
	}

	t.Run("no dvc", func(t *testing.T) {
		dir := "test-resources/dvc"
		expected := []api.Issue{
			api.NewIssue(linter.Name(), "", api.SeverityError, projectlinters.MsgUseDVC),
		}

		expectLintIssues(dir, expected)
	})

	t.Run("commit DVC folder", func(t *testing.T) {
		// create temporary folder with contents after 'dvc init' in it. temp folder will not be tracked by Git.
		dir, err := ioutil.TempDir("test-resources/dvc", "commit-dvc-folder-")
		require.NoError(t, err)
		require.NoError(t, os.MkdirAll(path.Join(dir, ".dvc"), 0755))
		require.NoError(t, ioutil.WriteFile(path.Join(dir, ".dvc", "config"), []byte("\n"), 0644))
		// clean up temporary folder
		defer os.RemoveAll(dir)

		dvc.ExecLookupPath = func(file string) (string, error) { return "", nil }
		expected := []api.Issue{
			api.NewIssue(linter.Name(), projectlinters.RuleCommitDVCFolder, api.SeverityError, projectlinters.MsgCommitDVCFolder),
		}

		expectLintIssues(dir, expected)
	})

	t.Run("commit DVC folder and not installed", func(t *testing.T) {
		// create temporary folder with contents after 'dvc init' in it. temp folder will not be tracked by Git.
		dir, err := ioutil.TempDir("test-resources/dvc", "commit-dvc-folder-")
		require.NoError(t, err)
		require.NoError(t, os.MkdirAll(path.Join(dir, ".dvc"), 0755))
		require.NoError(t, ioutil.WriteFile(path.Join(dir, ".dvc", "config"), []byte("\n"), 0644))
		// clean up temporary folder
		defer os.RemoveAll(dir)

		dvc.ExecLookupPath = func(file string) (string, error) { return "", fmt.Errorf("dvc not on path or something") }
		expected := []api.Issue{
			api.NewIssue(linter.Name(), projectlinters.RuleCommitDVCFolder, api.SeverityError, projectlinters.MsgCommitDVCFolder),
			api.NewIssue(linter.Name(), projectlinters.RuleDVCIsInstalled, api.SeverityError, projectlinters.MsgUseDVCIsInstalled),
		}

		expectLintIssues(dir, expected)
	})

	t.Run("not installed", func(t *testing.T) {
		dir := "test-resources/dvc/dvc-init"
		dvc.ExecLookupPath = func(file string) (string, error) { return "", fmt.Errorf("dvc not on path or something") }

		expected := []api.Issue{
			api.NewIssue(linter.Name(), projectlinters.RuleDVCIsInstalled, api.SeverityError, projectlinters.MsgUseDVCIsInstalled),
		}

		expectLintIssues(dir, expected)
	})

	t.Run("no remotes and no files", func(t *testing.T) {
		dir := "test-resources/dvc/dvc-init"
		dvc.ExecLookupPath = func(file string) (string, error) { return "", nil }
		dvc.ExecCommandOutput = func(cmddir, name string, args ...string) ([]byte, error) {
			require.Equal(t, dir, cmddir)
			return []byte("\n"), nil
		}

		expected := []api.Issue{
			api.NewIssue(linter.Name(), projectlinters.RuleDVCAddRemote, api.SeverityError, projectlinters.MsgDVCAddRemote),
			api.NewIssue(linter.Name(), projectlinters.RuleDVCAddFiles, api.SeverityWarning, projectlinters.MsgDVCAddFiles),
		}

		expectLintIssues(dir, expected)
	})
}
