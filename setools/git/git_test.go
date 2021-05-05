package git_test

import (
	"io/ioutil"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/setools/git"
)

func TestDetect(t *testing.T) {
	dir := "."
	require.True(t, git.Detect(dir))

	dir = ".."
	require.True(t, git.Detect(dir))

	dir = os.TempDir()
	require.False(t, git.Detect(dir))
}

func TestIsTracking(t *testing.T) {
	dir := "."
	require.True(t, git.IsTracking(dir, "git_test.go"))
	require.True(t, git.IsTracking(dir, "git*.go"))
	require.False(t, git.IsTracking(dir, "non-existant-file"))

	file, err := ioutil.TempFile(dir, "git.is-tracking.test-resource.*.txt")
	require.NoError(t, err)
	require.False(t, git.IsTracking(dir, file.Name()))

	require.NoError(t, os.Remove(file.Name())) // cleanup
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
