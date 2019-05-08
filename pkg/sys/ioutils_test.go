package sys

import (
	"github.com/dhillondeep/afero"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGetRoot(t *testing.T) {
	t.Run("Happy path - os executable path found", func(t *testing.T) {
		ex, err := os.Executable()
		require.NoError(t, err)
		folder := filepath.Dir(ex)

		root, err := GetRoot()
		require.NoError(t, err)

		require.Equal(t, folder, root)
	})

	t.Run("Error - os Executable error", func(t *testing.T) {
		oldOsExecutable := osExecutable
		defer func() {
			osExecutable = oldOsExecutable
		}()
		osExecutable = func() (s string, e error) {
			return "", errors.New("some error")
		}

		_, err := GetRoot()
		require.Error(t, err)
	})

	t.Run("Error - os Lstat error", func(t *testing.T) {
		oldOsLstat := osLstat
		defer func() {
			osLstat = oldOsLstat
		}()
		osLstat = func(name string) (info os.FileInfo, e error) {
			return nil, errors.New("some error")
		}

		_, err := GetRoot()
		require.Error(t, err)
	})

	t.Run("Happy path - os executable symlink path found", func(t *testing.T) {
		ex, err := os.Executable()
		require.NoError(t, err)
		folder := filepath.Dir(ex)

		currDir, err := os.Getwd()
		require.NoError(t, err)

		symLinkPath := currDir + "/tmp"
		require.NoError(t, os.MkdirAll(symLinkPath, os.ModePerm))
		defer func() {
			os.RemoveAll(symLinkPath)
		}()

		symEx := symLinkPath + "/sysTest"
		require.NoError(t, os.Symlink(ex, symEx))
		defer func() {
			os.RemoveAll(symEx)
		}()

		oldOsExecutable := osExecutable
		defer func() {
			osExecutable = oldOsExecutable
		}()
		osExecutable = func() (s string, e error) {
			return symEx, nil
		}

		root, err := GetRoot()
		require.NoError(t, err)

		require.Equal(t, folder, root)
	})

	t.Run("Error - os executable symlink path found os readlink error", func(t *testing.T) {
		ex, err := os.Executable()
		require.NoError(t, err)

		currDir, err := os.Getwd()
		require.NoError(t, err)

		symLinkPath := currDir + "/tmp"
		require.NoError(t, os.MkdirAll(symLinkPath, os.ModePerm))
		defer func() {
			os.RemoveAll(symLinkPath)
		}()

		symEx := symLinkPath + "/sysTest"
		require.NoError(t, os.Symlink(ex, symEx))
		defer func() {
			os.RemoveAll(symEx)
		}()

		oldOsExecutable := osExecutable
		defer func() {
			osExecutable = oldOsExecutable
		}()
		osExecutable = func() (s string, e error) {
			return symEx, nil
		}

		oldOsReadlink := osReadlink
		defer func() {
			osReadlink = oldOsReadlink
		}()
		osReadlink = func(name string) (s string, e error) {
			return "", errors.New("some error")
		}

		_, err = GetRoot()
		require.Error(t, err)
	})

	t.Run("Happy path - os executable symlink path found relative", func(t *testing.T) {
		ex, err := os.Executable()
		require.NoError(t, err)
		folder := filepath.Dir(ex)

		currDir, err := os.Getwd()
		require.NoError(t, err)

		symLinkPath := currDir + "/tmp"
		require.NoError(t, os.MkdirAll(symLinkPath, os.ModePerm))
		defer func() {
			os.RemoveAll(symLinkPath)
		}()

		symEx := symLinkPath + "/sysTest"
		require.NoError(t, os.Symlink(ex, symEx))
		defer func() {
			os.RemoveAll(symEx)
		}()

		oldOsExecutable := osExecutable
		defer func() {
			osExecutable = oldOsExecutable
		}()
		osExecutable = func() (s string, e error) {
			return symEx, nil
		}

		oldOsReadlink := osReadlink
		defer func() {
			osReadlink = oldOsReadlink
		}()
		osReadlink = func(name string) (s string, e error) {
			rel, err := filepath.Rel(symLinkPath, folder)
			require.NoError(t, err)

			return rel + "/" + filepath.Base(ex), nil
		}

		root, err := GetRoot()
		require.NoError(t, err)

		require.Equal(t, folder, root)
	})

	t.Run("Error - os executable symlink path found relative but abs error", func(t *testing.T) {
		ex, err := os.Executable()
		require.NoError(t, err)
		folder := filepath.Dir(ex)

		currDir, err := os.Getwd()
		require.NoError(t, err)

		symLinkPath := currDir + "/tmp"
		require.NoError(t, os.MkdirAll(symLinkPath, os.ModePerm))
		defer func() {
			os.RemoveAll(symLinkPath)
		}()

		symEx := symLinkPath + "/sysTest"
		require.NoError(t, os.Symlink(ex, symEx))
		defer func() {
			os.RemoveAll(symEx)
		}()

		oldOsExecutable := osExecutable
		defer func() {
			osExecutable = oldOsExecutable
		}()
		osExecutable = func() (s string, e error) {
			return symEx, nil
		}

		oldOsReadlink := osReadlink
		defer func() {
			osReadlink = oldOsReadlink
		}()
		osReadlink = func(name string) (s string, e error) {
			rel, err := filepath.Rel(symLinkPath, folder)
			require.NoError(t, err)

			return rel + "/" + filepath.Base(ex), nil
		}

		oldFilepathAbs := filepathAbs
		defer func() {
			filepathAbs = oldFilepathAbs
		}()
		filepathAbs = func(path string) (s string, e error) {
			return "", errors.New("some error")
		}

		_, err = GetRoot()
		require.Error(t, err)
	})
}

func TestGetOS(t *testing.T) {
	require.Equal(t, runtime.GOOS, GetOS())
}

func TestGetArch(t *testing.T) {
	require.Equal(t, runtime.GOARCH, GetArch())
}

func TestExists(t *testing.T) {
	SetFileSystem(afero.NewMemMapFs())

	fs := GetFileSystem()

	t.Run("Happy path - file exists", func(t *testing.T) {
		_, err := fs.Create("/file.txt")
		require.NoError(t, err)

		require.True(t, Exists("/file.txt"))
	})

	t.Run("Happy path - file does not exist", func(t *testing.T) {
		require.False(t, Exists("/file2.txt"))
	})
}

func TestGetSep(t *testing.T) {
	require.Equal(t, string(filepath.Separator), GetSep())
}

func TestJoinPaths(t *testing.T) {
	require.Equal(t, "", JoinPaths())
	require.Equal(t, "path", JoinPaths("path"))
	require.Equal(t, "path/path2", JoinPaths("path", "path2"))
	require.Equal(t, "path/path2/path3", JoinPaths("path/path2", "path3"))
	require.Equal(t, "path/path2/path3", JoinPaths("path///path2//", "path3"))
}

func TestIsDir(t *testing.T) {
	SetFileSystem(afero.NewMemMapFs())

	fs := GetFileSystem()

	require.NoError(t, fs.MkdirAll("folder", os.ModePerm))

	_, err := fs.Create("/file.txt")
	require.NoError(t, err)

	dir, err := IsDir("folder")
	require.NoError(t, err)
	require.True(t, dir)

	dir, err = IsDir("/file.txt")
	require.NoError(t, err)
	require.False(t, dir)

	_, err = IsDir("folderDoesNotExist")
	require.Error(t, err)
}
