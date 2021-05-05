package git

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/bvobart/mllint/utils"
	"github.com/bvobart/mllint/utils/exec"
)

// Detect detects whether this directory is inside a Git repository
func Detect(dir string) bool {
	_, err := exec.CommandOutput(dir, "git", "rev-parse", "--git-dir")
	return err == nil
}

// IsTracking checks whether the Git repository in the given folder is tracking the files specified
// by the given pattern. This can be a literal folder or file name, but can also be a pattern
// containing wildcards, e.g. 'foo.*'
func IsTracking(dir string, pattern string) bool {
	_, err := exec.CommandOutput(dir, "git", "ls-files", "--error-unmatch", pattern)
	return err == nil
}

// FileSize is the return type for FindLargeFiles. Contains the path to the file and its filesize.
type FileSize struct {
	Path string
	Size uint64
}

// FindLargeFiles looks for any files being tracked in the current Git repository that have a
// filesize larger than the given threshold, measured in bytes.
func FindLargeFiles(dir string, threshold uint64) ([]FileSize, error) {
	output, err := exec.CommandOutput(dir, "git", "ls-tree", "-r", "-t", "-l", "--full-name", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("failed to read Git files: %w", utils.WrapExitError(err))
	}

	files := []FileSize{}
	if string(output) == "" {
		return files, nil
	}

	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			return nil, fmt.Errorf("unexpected output from git ls-tree: %s", string(output))
		}

		sizeStr := fields[3]
		if sizeStr == "-" {
			continue
		}

		size, err := strconv.ParseUint(sizeStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse filesize from '%s': %w", sizeStr, err)
		}

		if size > threshold {
			files = append(files, FileSize{Path: fields[4], Size: size})
		}
	}

	// sort files by filesize in descending order
	sort.Slice(files, func(i, j int) bool {
		return files[i].Size > files[j].Size
	})

	return files, nil
}
