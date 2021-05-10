package cqlinters

import "path/filepath"

type stringer string

func (s stringer) String() string { return string(s) }

func trimProjectDir(path string, projectdir string) string {
	relpath, err := filepath.Rel(projectdir, path)
	if err != nil {
		return path
	}
	return relpath
}
