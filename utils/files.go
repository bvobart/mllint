package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hhatto/gocloc"
)

// FileExists checks if a file exists and is not a directory
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	return !os.IsNotExist(err) && info != nil && !info.IsDir()
}

// FolderExists checks if a folder exists
func FolderExists(filename string) bool {
	info, err := os.Stat(filename)
	return !os.IsNotExist(err) && info.IsDir()
}

// FolderIsEmpty checks if a folder is empty
func FolderIsEmpty(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	_, err = file.Readdirnames(1)
	if errors.Is(err, io.EOF) {
		return true, nil
	}
	return false, err
}

// OpenFile looks inside of the given folder for a file matching the given pattern.
// Will return a non-nil error when either no or more than one files match.
// Returns the opened file otherwise.
func OpenFile(folder string, pattern string) (*os.File, error) {
	matches, err := filepath.Glob(path.Join(folder, pattern))
	if err != nil {
		return nil, err
	}

	if len(matches) == 1 {
		return os.Open(matches[0])
	}

	folderpath := AbsolutePath(folder)
	if len(matches) == 0 {
		return nil, fmt.Errorf("did not find a file matching %s in folder %s: %w", pattern, folderpath, os.ErrNotExist)
	}

	return nil, fmt.Errorf("pattern %s in folder %s matches multiple files: %+v", pattern, folderpath, matches)
}

// AbsolutePath returns the absolute path to the given file.
// If the filename is not already an absolute path, then we assume that it is a file path relative to the current working directory
func AbsolutePath(filename string) string {
	if path.IsAbs(filename) {
		return filename
	}

	cwd, _ := os.Getwd()
	return path.Join(cwd, filename)
}

// FindPythonFilesIn finds all Python (*.py) files in the given directory and subdirectories
// Returns their filepaths, relative to the given directory
// Ignores hidden folders (folders whose names start with a '.'), but not hidden files.
func FindPythonFilesIn(dir string) (Filenames, error) {
	return FindFilesByExtInDir(dir, ".py")
}

// FindIPynbFilesIn finds all Jupyter Notebook (*.ipynb) files in the given directory and
// subdirectories. Returns their filepaths, relative to the given directory.
// Ignores hidden folders (folders whose names start with a '.'), but not hidden files.
func FindIPynbFilesIn(dir string) (Filenames, error) {
	return FindFilesByExtInDir(dir, ".ipynb")
}

// FindFilesByExtInDir finds all files in the given directory and subdirectories that have
// a certain file extension. File extension must start with a '.', e.g. ".py" or ".ipynb"
// Returns filepaths relative to the given directory.
// Ignores hidden folders (folders whose names start with a '.'), but not hidden files.
// Also explicitly ignores `venv`, `env`. `venv.bak` and `env.bak` folders
func FindFilesByExtInDir(dir string, extension string) (Filenames, error) {
	files := Filenames{}
	err := filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if file.IsDir() && strings.HasPrefix(file.Name(), ".") {
			return filepath.SkipDir
		}

		// ignore virtualenv directories from official Python Gitignore: https://github.com/github/gitignore/blob/991e760c1c6d50fdda246e0178b9c58b06770b90/Python.gitignore#L107
		if file.IsDir() && (file.Name() == "venv" || file.Name() == "env" || file.Name() == "ENV" || file.Name() == "env.bak" || file.Name() == "venv.bak") {
			return filepath.SkipDir
		}

		if !file.IsDir() && filepath.Ext(path) == extension {
			relpath, _ := filepath.Rel(dir, path)
			files = append(files, relpath)
		}

		return nil
	})

	return files, err
}

// Filenames is simply an alias for []string, but allows me to add some methods.
type Filenames []string

func (names Filenames) Concat(extra Filenames) Filenames {
	return append(names, extra...)
}

func (names Filenames) Filter(shouldInclude func(filename string) bool) Filenames {
	result := Filenames{}
	for _, filename := range names {
		if shouldInclude(filename) {
			result = append(result, filename)
		}
	}
	return result
}

// Prefix prefixes each of the filenames with a directory name.
// i.e. Filenames{"name.py"}.Prefix("something") becomes Filenames{"something/name.py"}
func (names Filenames) Prefix(dir string) Filenames {
	for i, name := range names {
		names[i] = filepath.Join(dir, name)
	}
	return names
}

var langPython = gocloc.NewLanguage("Python", []string{"#"}, [][]string{{"\"\"\"", "\"\"\""}})

func (names Filenames) CountLoC() int32 {
	total := int32(0)
	opts := gocloc.NewClocOptions()
	for _, name := range names {
		analysed := gocloc.AnalyzeFile(name, langPython, opts)
		total += analysed.Code
	}
	return total
}
