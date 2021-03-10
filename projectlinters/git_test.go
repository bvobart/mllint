package projectlinters_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/bvobart/mllint/api"
	"gitlab.com/bvobart/mllint/projectlinters"
)

func TestUseGit(t *testing.T) {
	linter := projectlinters.UseGit{}
	dir := "test-resources"
	issues, err := linter.LintProject(dir)
	require.NoError(t, err)
	require.Nil(t, issues)

	dir = os.TempDir()
	issues, err = linter.LintProject(dir)
	require.NoError(t, err)
	require.Equal(t, []api.Issue{api.NewIssue(linter.Name(), "", api.SeverityError, projectlinters.MsgUseGit)}, issues)
}

func TestNoBigFiles(t *testing.T) {
	linter := projectlinters.GitNoBigFiles{
		Threshold: 10_000_000, // 10 MB
	}

	dir := "."
	issues, err := linter.LintProject(dir)
	require.NoError(t, err)
	require.Nil(t, issues)

	// TODO: add test for when there are files larger than threshold.
}
