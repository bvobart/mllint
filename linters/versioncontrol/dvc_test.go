package versioncontrol_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/linters/versioncontrol"
	"github.com/bvobart/mllint/utils/exec"
)

func TestDVCName(t *testing.T) {
	linter := &versioncontrol.DVCLinter{}
	require.Equal(t, "Data", linter.Name())
}

func TestDVCRules(t *testing.T) {
	linter := &versioncontrol.DVCLinter{}
	rules := linter.Rules()
	require.Equal(t, []*api.Rule{
		&versioncontrol.RuleDVC,
		&versioncontrol.RuleDVCIsInstalled,
		&versioncontrol.RuleCommitDVCFolder,
		&versioncontrol.RuleDVCHasRemote,
		&versioncontrol.RuleDVCHasFiles,
		&versioncontrol.RuleCommitDVCLock,
	}, rules)
}

func TestDVCLinter(t *testing.T) {
	linter := &versioncontrol.DVCLinter{}

	t.Run("no dvc", func(t *testing.T) {
		dir := "test-resources/dvc"
		report, err := linter.LintProject(dir)
		require.NoError(t, err)

		require.EqualValues(t, 0, report.Scores[versioncontrol.RuleDVC])
		require.EqualValues(t, 0, report.Scores[versioncontrol.RuleDVCIsInstalled])
		require.EqualValues(t, 0, report.Scores[versioncontrol.RuleCommitDVCFolder])
		require.EqualValues(t, 0, report.Scores[versioncontrol.RuleDVCHasRemote])
		require.EqualValues(t, 0, report.Scores[versioncontrol.RuleDVCHasFiles])
	})

	t.Run("commit DVC folder", func(t *testing.T) {
		// Given: a temporary folder with contents after 'dvc init' in it. temp folder will not be tracked by Git.
		dir, err := ioutil.TempDir("test-resources/dvc", "commit-dvc-folder-")
		require.NoError(t, err)
		require.NoError(t, os.MkdirAll(path.Join(dir, ".dvc"), 0755))
		require.NoError(t, ioutil.WriteFile(path.Join(dir, ".dvc", "config"), []byte("\n"), 0644))
		// clean up temporary folder
		defer os.RemoveAll(dir)

		// ensure DVC will seem installed
		exec.LookPath = func(file string) (string, error) { return "", nil }
		exec.CommandOutput = exec.DefaultCommandOutput

		report, err := linter.LintProject(dir)
		require.NoError(t, err)

		// Then:
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVC])
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVCIsInstalled])
		require.EqualValues(t, 0, report.Scores[versioncontrol.RuleCommitDVCFolder])
		require.EqualValues(t, 0, report.Scores[versioncontrol.RuleDVCHasRemote])
		require.EqualValues(t, 0, report.Scores[versioncontrol.RuleDVCHasFiles])
	})

	t.Run("commit DVC folder and not installed", func(t *testing.T) {
		// Given: a temporary folder with contents after 'dvc init' in it. temp folder will not be tracked by Git.
		dir, err := ioutil.TempDir("test-resources/dvc", "commit-dvc-folder-")
		require.NoError(t, err)
		require.NoError(t, os.MkdirAll(path.Join(dir, ".dvc"), 0755))
		require.NoError(t, ioutil.WriteFile(path.Join(dir, ".dvc", "config"), []byte("\n"), 0644))
		// clean up temporary folder
		defer os.RemoveAll(dir)

		// ensure DVC will seem *not* installed
		exec.LookPath = func(file string) (string, error) { return "", fmt.Errorf("dvc not on path or something") }
		exec.CommandOutput = exec.DefaultCommandOutput

		report, err := linter.LintProject(dir)
		require.NoError(t, err)

		// Then:
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVC])
		require.EqualValues(t, 0, report.Scores[versioncontrol.RuleDVCIsInstalled])
		require.EqualValues(t, 0, report.Scores[versioncontrol.RuleCommitDVCFolder])
		require.EqualValues(t, 0, report.Scores[versioncontrol.RuleDVCHasRemote])
		require.EqualValues(t, 0, report.Scores[versioncontrol.RuleDVCHasFiles])
	})

	t.Run("DVC folder committed but DVC is not installed", func(t *testing.T) {
		dir := "test-resources/dvc/dvc-init"
		exec.LookPath = func(file string) (string, error) { return "", fmt.Errorf("dvc not on path or something") }
		exec.CommandOutput = exec.DefaultCommandOutput

		report, err := linter.LintProject(dir)
		require.NoError(t, err)

		// Then:
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVC])
		require.EqualValues(t, 0, report.Scores[versioncontrol.RuleDVCIsInstalled])
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleCommitDVCFolder])
		require.EqualValues(t, 0, report.Scores[versioncontrol.RuleDVCHasRemote])
		require.EqualValues(t, 0, report.Scores[versioncontrol.RuleDVCHasFiles])
	})

	t.Run("DVC folder committed and DVC installed, but no remotes and no files", func(t *testing.T) {
		dir := "test-resources/dvc/dvc-init"
		exec.LookPath = func(file string) (string, error) { return "", nil }
		exec.CommandOutput = func(cmddir, name string, args ...string) ([]byte, error) {
			if name == "git" {
				return exec.DefaultCommandOutput(cmddir, name, args...)
			}
			require.Equal(t, dir, cmddir)
			require.Equal(t, name, "dvc")
			return []byte("\n"), nil
		}

		report, err := linter.LintProject(dir)
		require.NoError(t, err)

		// Then:
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVC])
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVCIsInstalled])
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleCommitDVCFolder])
		require.EqualValues(t, 0, report.Scores[versioncontrol.RuleDVCHasRemote])
		require.EqualValues(t, 0, report.Scores[versioncontrol.RuleDVCHasFiles])
	})

	t.Run("DVC with remotes configured but no files", func(t *testing.T) {
		dir := "test-resources/dvc/dvc-init"
		exec.LookPath = func(file string) (string, error) { return "", nil }
		exec.CommandOutput = func(cmddir, name string, args ...string) ([]byte, error) {
			if name == "git" {
				return exec.DefaultCommandOutput(cmddir, name, args...)
			}

			require.Equal(t, dir, cmddir)
			require.Equal(t, name, "dvc")

			if args[0] == "remote" && args[1] == "list" {
				return []byte("testremote			/path/to/remote\n"), nil
			}
			return []byte("\n"), nil
		}

		report, err := linter.LintProject(dir)
		require.NoError(t, err)

		// Then:
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVC])
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVCIsInstalled])
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleCommitDVCFolder])
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVCHasRemote])
		require.EqualValues(t, 0, report.Scores[versioncontrol.RuleDVCHasFiles])
	})

	t.Run("DVC with remotes and files configured", func(t *testing.T) {
		dir := "test-resources/dvc/dvc-init"
		exec.LookPath = func(file string) (string, error) { return "", nil }
		exec.CommandOutput = func(cmddir, name string, args ...string) ([]byte, error) {
			if name == "git" {
				return exec.DefaultCommandOutput(cmddir, name, args...)
			}

			require.Equal(t, dir, cmddir)
			require.Equal(t, name, "dvc")

			if args[0] == "remote" && args[1] == "list" {
				return []byte("testremote			/path/to/remote\n"), nil
			}
			if args[0] == "list" && args[1] == "." {
				return []byte("dvcfile1\ndvcfile2\n"), nil
			}
			return []byte("\n"), nil
		}

		report, err := linter.LintProject(dir)
		require.NoError(t, err)

		// Then:
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVC])
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVCIsInstalled])
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleCommitDVCFolder])
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVCHasRemote])
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVCHasFiles])
	})

	t.Run("DVC with remotes and files configured and uncommitted dvc.lock", func(t *testing.T) {
		dir := "test-resources/dvc/dvc-init"
		lockfilePath := path.Join(dir, "dvc.lock")
		require.NoError(t, ioutil.WriteFile(lockfilePath, []byte("\n"), 0644))
		defer os.Remove(lockfilePath)

		exec.LookPath = func(file string) (string, error) { return "", nil }
		exec.CommandOutput = func(cmddir, name string, args ...string) ([]byte, error) {
			if name == "git" {
				return exec.DefaultCommandOutput(cmddir, name, args...)
			}

			require.Equal(t, dir, cmddir)
			require.Equal(t, name, "dvc")

			if args[0] == "remote" && args[1] == "list" {
				return []byte("testremote			/path/to/remote\n"), nil
			}
			if args[0] == "list" && args[1] == "." {
				return []byte("dvcfile1\ndvcfile2\n"), nil
			}
			return []byte("\n"), nil
		}

		report, err := linter.LintProject(dir)
		require.NoError(t, err)

		// Then:
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVC])
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVCIsInstalled])
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleCommitDVCFolder])
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVCHasRemote])
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVCHasFiles])
		require.EqualValues(t, 0, report.Scores[versioncontrol.RuleCommitDVCLock])
	})

	t.Run("DVC Perfect Score", func(t *testing.T) {
		dir := "test-resources/dvc/dvc-lock-committed"
		exec.LookPath = func(file string) (string, error) { return "", nil }
		exec.CommandOutput = func(cmddir, name string, args ...string) ([]byte, error) {
			if name == "git" {
				return exec.DefaultCommandOutput(cmddir, name, args...)
			}

			require.Equal(t, dir, cmddir)
			require.Equal(t, name, "dvc")

			if args[0] == "remote" && args[1] == "list" {
				return []byte("testremote			/path/to/remote\n"), nil
			}
			if args[0] == "list" && args[1] == "." {
				return []byte("dvcfile1\ndvcfile2\n"), nil
			}
			return []byte("\n"), nil
		}

		report, err := linter.LintProject(dir)
		require.NoError(t, err)

		// Then:
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVC])
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVCIsInstalled])
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleCommitDVCFolder])
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVCHasRemote])
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleDVCHasFiles])
		require.EqualValues(t, 100, report.Scores[versioncontrol.RuleCommitDVCLock])
	})
}
