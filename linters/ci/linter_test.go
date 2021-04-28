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
)

func TestCILinter(t *testing.T) {
	linter := ci.NewLinter()
	require.Equal(t, "Continuous Integration (CI)", linter.Name())
	require.Equal(t, []*api.Rule{&ci.RuleUseCI}, linter.Rules())

	t.Run("None", func(t *testing.T) {
		project := api.Project{Dir: "test-resources/none"}
		report, err := linter.LintProject(project)
		require.NoError(t, err)

		score, found := report.Scores[ci.RuleUseCI]
		require.True(t, found)
		require.EqualValues(t, 0, score)
	})

	t.Run("Azure", func(t *testing.T) {
		project := api.Project{Dir: "test-resources/azure"}
		report, err := linter.LintProject(project)
		require.NoError(t, err)
		require.EqualValues(t, 100, report.Scores[ci.RuleUseCI])
	})

	t.Run("GitHubActions", func(t *testing.T) {
		project := api.Project{Dir: "test-resources/ghactions"}
		report, err := linter.LintProject(project)
		require.NoError(t, err)
		require.EqualValues(t, 100, report.Scores[ci.RuleUseCI])
	})

	t.Run("GitlabCI", func(t *testing.T) {
		project := api.Project{Dir: "test-resources/gitlab"}
		report, err := linter.LintProject(project)
		require.NoError(t, err)
		require.EqualValues(t, 100, report.Scores[ci.RuleUseCI])
	})

	t.Run("Travis", func(t *testing.T) {
		project := api.Project{Dir: "test-resources/travis"}
		report, err := linter.LintProject(project)
		require.NoError(t, err)
		require.EqualValues(t, 100, report.Scores[ci.RuleUseCI])
	})

	t.Run("GitlabUntracked", func(t *testing.T) {
		dir, err := ioutil.TempDir("test-resources", "gitlab-untracked")
		require.NoError(t, err)
		require.NoError(t, ioutil.WriteFile(path.Join(dir, ciproviders.Gitlab{}.ConfigFile()), []byte("\n"), 0644))
		defer os.RemoveAll(dir)

		project := api.Project{Dir: dir}
		report, err := linter.LintProject(project)
		require.NoError(t, err)
		require.EqualValues(t, 25, report.Scores[ci.RuleUseCI])
	})
}
