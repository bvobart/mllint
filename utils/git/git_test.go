package git_test

import (
	"math"
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

func TestFindLargeFiles(t *testing.T) {
	dir := "."

	threshold := uint64(1)
	largeFiles, err := git.FindLargeFiles(dir, threshold)
	require.NoError(t, err)
	require.Len(t, largeFiles, 2)

	// test that largeFiles is sorted by filesize in descending order (i.e. largest files first)
	prevSize := uint64(math.MaxUint64)
	for _, file := range largeFiles {
		require.Truef(t, file.Size < prevSize, "Should be sorted by filesize in descending order: %+v", largeFiles)
		prevSize = file.Size
	}

	threshold = uint64(1000000000)
	largeFiles, err = git.FindLargeFiles(dir, threshold)
	require.NoError(t, err)
	require.Len(t, largeFiles, 0)
}
