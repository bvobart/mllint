package versioncontrol_test

import (
	"os"
	"testing"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/linters/versioncontrol"
	"github.com/stretchr/testify/require"
)

func TestGitName(t *testing.T) {
	linter := &versioncontrol.GitLinter{}
	require.Equal(t, "Code", linter.Name())
}

func TestGitRules(t *testing.T) {
	linter := &versioncontrol.GitLinter{}
	require.Equal(t, []*api.Rule{&versioncontrol.RuleGit, &versioncontrol.RuleGitNoBigFiles}, linter.Rules())
}

func TestGitConfigure(t *testing.T) {
	linter := &versioncontrol.GitLinter{}
	conf := config.Default()
	conf.Git.MaxFileSize = 1337
	require.NoError(t, linter.Configure(conf))
	require.EqualValues(t, 1337, linter.MaxFileSize)
}

func TestProjectUsesGit(t *testing.T) {
	linter := &versioncontrol.GitLinter{}
	dir := "."
	report, err := linter.LintProject(dir)
	require.NoError(t, err)
	require.EqualValues(t, 100, report.Scores[versioncontrol.RuleGit])

	dir = os.TempDir()
	report, err = linter.LintProject(dir)
	require.NoError(t, err)
	require.EqualValues(t, 0, report.Scores[versioncontrol.RuleGit])
}

func TestGitNoBigFiles(t *testing.T) {
	linter := &versioncontrol.GitLinter{
		MaxFileSize: 10_000_000, // 10 MB
	}

	dir := "."
	report, err := linter.LintProject(dir)
	require.NoError(t, err)
	require.EqualValues(t, 100, report.Scores[versioncontrol.RuleGit])

	// TODO: add test for when there are files larger than threshold.
}
