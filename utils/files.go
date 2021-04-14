package utils

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

// FileExists checks if a file exists and is not a directory
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// FolderExists checks if a folder exists
func FolderExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// FolderIsEmpty checks if a folder is empty
func FolderIsEmpty(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	_, err = file.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

// OpenFile looks inside of the given folder for a file matching the given pattern.
// Will return a non-nil error when either no or more than one files match.
// Returns the opened file otherwise.
func OpenFile(folder string, pattern string) (*os.File, error) {
	matches, err := filepath.Glob(folder + "/" + pattern)
	if err != nil {
		return nil, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("did not find a file matching %s in folder %s/%s", pattern, cwd, folder)
	} else if len(matches) > 1 {
		return nil, fmt.Errorf("pattern %s in folder %s/%s matches multiple files: %+v", pattern, cwd, folder, matches)
	} else {
		return os.Open(matches[0])
	}
}

// AbsolutePath returns the absolute path to the given file,
// assuming that it is a file path relative to the current working directory
func AbsolutePath(filename string) string {
	cwd, _ := os.Getwd()
	return path.Join(cwd, filename)
}
