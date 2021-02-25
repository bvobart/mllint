package git_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/bvobart/mllint/utils/git"
)

func TestDetect(t *testing.T) {
	dir := "."
	require.True(t, git.Detect(dir))

	dir = ".."
	require.True(t, git.Detect(dir))

	dir = os.TempDir()
	require.False(t, git.Detect(dir))
}
