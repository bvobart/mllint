package utils_test

import (
	"os"
	"path"
	"testing"

	"github.com/bvobart/mllint/utils"
	"github.com/stretchr/testify/require"
)

func TestFileFolderExists(t *testing.T) {
	dir := "test-resources"
	require.True(t, utils.FolderExists(dir))
	require.False(t, utils.FileExists(dir))

	dir = path.Join(dir, "python-files")
	require.True(t, utils.FolderExists(dir))

	file := path.Join(dir, "some_script.py")
	require.True(t, utils.FileExists(file))
	require.False(t, utils.FolderExists(file))

	dir = "non-existant"
	require.False(t, utils.FileExists(dir))
	require.False(t, utils.FolderExists(dir))
}

func TestFolderIsEmpty(t *testing.T) {
	dir := "test-resources"
	isEmpty, err := utils.FolderIsEmpty(dir)
	require.NoError(t, err)
	require.False(t, isEmpty)

	dir = "non-existant"
	isEmpty, err = utils.FolderIsEmpty(dir)
	require.True(t, os.IsNotExist(err))
	require.False(t, isEmpty)

	dir, err = os.MkdirTemp("", "")
	require.NoError(t, err)
	defer os.Remove(dir)
	isEmpty, err = utils.FolderIsEmpty(dir)
	require.NoError(t, err)
	require.True(t, isEmpty)
}

func TestOpenFile(t *testing.T) {
	dir := "test-resources/python-files"
	expectedFile, err := os.Open(path.Join(dir, "some_script.py"))
	require.NoError(t, err)

	file, err := utils.OpenFile(dir, "some_script*")
	require.NoError(t, err)
	require.Equal(t, expectedFile.Name(), file.Name())

	_, err = utils.OpenFile(dir, "some_*")
	require.Error(t, err)
	require.Contains(t, err.Error(), "some_*")
	require.Contains(t, err.Error(), dir)
	require.Contains(t, err.Error(), "matches multiple files")

	_, err = utils.OpenFile(dir, "non-existant")
	require.Error(t, err)
	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestFindPythonFiles(t *testing.T) {
	dir := "test-resources/python-files"
	files, err := utils.FindPythonFilesIn(dir)
	require.NoError(t, err)
	require.Equal(t, utils.Filenames{"some_other_script.py", "some_script.py", "subfolder/yet_another_script.py"}, files)
}

func TestCountLoC(t *testing.T) {
	files := utils.Filenames{"some_script.py", "some_other_script.py", "subfolder/yet_another_script.py"}
	files = files.Prefix(path.Join("test-resources", "python-files"))
	loc := files.CountLoC()
	require.EqualValues(t, 14, loc)
}

func TestPrefix(t *testing.T) {
	filenames := utils.Filenames{"file1.py", "folder/file2.py"}
	prefixed := filenames.Prefix("something/test-dir")

	expected := utils.Filenames{"something/test-dir/file1.py", "something/test-dir/folder/file2.py"}
	require.Equal(t, len(filenames), len(prefixed))
	require.Equal(t, expected, prefixed)
}

func TestAbsolutePath(t *testing.T) {
	require.Equal(t, "/dev/null", utils.AbsolutePath("/dev/null"))

	cwd, err := os.Getwd()
	require.NoError(t, err)
	require.Equal(t, path.Join(cwd, "test-resources"), utils.AbsolutePath("test-resources"))
}
