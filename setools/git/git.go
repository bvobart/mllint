package git

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/utils"
	"github.com/bvobart/mllint/utils/exec"
)

// Detect detects whether this directory is inside a Git repository
func Detect(dir string) bool {
	_, err := exec.CommandOutput(dir, "git", "rev-parse", "--git-dir")
	return err == nil
}

func MakeGitInfo(dir string) api.GitInfo {
	if !Detect(dir) {
		return api.GitInfo{}
	}

	remote, _ := GetRemoteURL(dir)
	commit, _ := GetCurrentCommit(dir)
	branch, _ := GetCurrentBranch(dir)
	dirty := IsDirty(dir)
	return api.GitInfo{RemoteURL: remote, Commit: commit, Branch: branch, Dirty: dirty}
}

func GetRemoteURL(dir string) (string, error) {
	output, err := exec.CommandOutput(dir, "git", "remote", "get-url", "origin")
	return strings.TrimSpace(string(output)), err
}

func GetCurrentCommit(dir string) (string, error) {
	output, err := exec.CommandOutput(dir, "git", "rev-parse", "HEAD")
	return strings.TrimSpace(string(output)), err
}

func GetCurrentBranch(dir string) (string, error) {
	output, err := exec.CommandOutput(dir, "git", "branch", "--show-current")
	return strings.TrimSpace(string(output)), err
}

func IsDirty(dir string) bool {
	_, err := exec.CommandOutput(dir, "git", "diff", "--no-ext-diff", "--quiet")
	return err != nil
}

// IsTracking checks whether the Git repository in the given folder is tracking the files specified
// by the given pattern. This can be a literal folder or file name, but can also be a pattern
// containing wildcards, e.g. 'foo.*'
func IsTracking(dir string, pattern string) bool {
	_, err := exec.CommandOutput(dir, "git", "ls-files", "--error-unmatch", pattern)
	return err == nil
}

// FileSize is the return type for FindLargeFiles. Contains the path to the file and its filesize,
// and, if specified, the commit hash on which the given file was created.
type FileSize struct {
	Path       string
	CommitHash string
	Size       uint64
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

func FindLargeFilesInHistory(dir string, threshold uint64) ([]FileSize, error) {
	output, err := exec.PipelineOutput(dir, [][]string{
		{"git", "rev-list", "--objects", "--all"},
		{"git", "cat-file", "--batch-check=%(objecttype) %(objectname) %(objectsize) %(rest)"},
	}...)
	if err != nil {
		return nil, fmt.Errorf("failed to read Git files: %w", utils.WrapExitError(err))
	}

	files := []FileSize{}
	for _, line := range strings.Split(string(output), "\n") {
		if !strings.HasPrefix(line, "blob") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			return nil, fmt.Errorf("expecting 4 fields in this message but it has %d: '%s'", len(fields), line)
		}

		sizeStr := fields[2]
		size, err := strconv.ParseUint(sizeStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse filesize from '%s': %w", sizeStr, err)
		}

		if size > threshold {
			file := FileSize{Path: fields[3], CommitHash: fields[1], Size: size}
			files = append(files, file)
		}
	}

	// sort files by filesize in descending order
	sort.Slice(files, func(i, j int) bool {
		return files[i].Size > files[j].Size
	})

	return files, err
}
