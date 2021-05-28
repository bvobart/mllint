package ci_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/linters/ci"
	"github.com/bvobart/mllint/setools/ciproviders"
	"github.com/bvobart/mllint/utils/exec"
)

const tmpkey = "mllint-tests-ci"

func createTestGitDir(t *testing.T, name string) string {
	dir, err := ioutil.TempDir(os.TempDir(), tmpkey+"-"+name)
	require.NoError(t, err)
	_, err = exec.CommandOutput(dir, "git", "init")
	require.NoError(t, err)
	return dir
}

func TestCILinter(t *testing.T) {
	linter := ci.NewLinter()
	require.Equal(t, "Continuous Integration (CI)", linter.Name())
	require.Equal(t, []*api.Rule{&ci.RuleUseCI}, linter.Rules())

	t.Run("None", func(t *testing.T) {
		dir := createTestGitDir(t, "")
		project := api.Project{Dir: dir}
		defer os.RemoveAll(dir)

		report, err := linter.LintProject(project)
		require.NoError(t, err)

		score, found := report.Scores[ci.RuleUseCI]
		require.True(t, found)
		require.EqualValues(t, 0, score)
	})

	t.Run("Azure", func(t *testing.T) {
		dir := createTestGitDir(t, "azure")
		project := api.Project{Dir: dir}
		defer os.RemoveAll(dir)

		require.NoError(t, ioutil.WriteFile(ciproviders.Azure{}.ConfigFile(dir), []byte("\n"), 0644))
		_, err := exec.CommandOutput(dir, "git", "add", ".")
		require.NoError(t, err)

		report, err := linter.LintProject(project)
		require.NoError(t, err)
		require.EqualValues(t, 100, report.Scores[ci.RuleUseCI])
	})

	t.Run("GitHubActions", func(t *testing.T) {
		dir := createTestGitDir(t, "ghactions")
		project := api.Project{Dir: dir}
		defer os.RemoveAll(dir)

		require.NoError(t, os.MkdirAll(ciproviders.GHActions{}.ConfigFile(dir), 0755))
		require.NoError(t, ioutil.WriteFile(path.Join(ciproviders.GHActions{}.ConfigFile(dir), "workflow.yml"), []byte("\n"), 0644))
		_, err := exec.CommandOutput(dir, "git", "add", ".")
		require.NoError(t, err)

		report, err := linter.LintProject(project)
		require.NoError(t, err)
		require.EqualValues(t, 100, report.Scores[ci.RuleUseCI])
	})

	t.Run("GitlabCI", func(t *testing.T) {
		dir := createTestGitDir(t, "gitlab")
		project := api.Project{Dir: dir}
		defer os.RemoveAll(dir)

		require.NoError(t, ioutil.WriteFile(ciproviders.Gitlab{}.ConfigFile(dir), []byte("\n"), 0644))
		_, err := exec.CommandOutput(dir, "git", "add", ".")
		require.NoError(t, err)

		report, err := linter.LintProject(project)
		require.NoError(t, err)
		require.EqualValues(t, 100, report.Scores[ci.RuleUseCI])
	})

	t.Run("Travis", func(t *testing.T) {
		dir := createTestGitDir(t, "travis")
		project := api.Project{Dir: dir}
		defer os.RemoveAll(dir)

		require.NoError(t, ioutil.WriteFile(ciproviders.Travis{}.ConfigFile(dir), []byte("\n"), 0644))
		_, err := exec.CommandOutput(dir, "git", "add", ".")
		require.NoError(t, err)

		report, err := linter.LintProject(project)
		require.NoError(t, err)
		require.EqualValues(t, 100, report.Scores[ci.RuleUseCI])
	})

	t.Run("GitlabUntracked", func(t *testing.T) {
		dir := createTestGitDir(t, "gitlab-untracked")
		project := api.Project{Dir: dir}
		defer os.RemoveAll(dir)

		require.NoError(t, ioutil.WriteFile(ciproviders.Gitlab{}.ConfigFile(dir), []byte("\n"), 0644))

		report, err := linter.LintProject(project)
		require.NoError(t, err)
		require.EqualValues(t, 25, report.Scores[ci.RuleUseCI])
	})
}
