package sys

import (
	"github.com/pkg/errors"
	"io"
	"os"
	"testing"

	"github.com/dhillondeep/afero"
	"github.com/stretchr/testify/require"
)

func TestSetFileSystem(t *testing.T) {
	SetFileSystem(fs)
	require.Equal(t, fs, afero.NewOsFs())

	fsToUse := afero.NewMemMapFs()
	SetFileSystem(fsToUse)
	require.Equal(t, fs, fsToUse)
}

func TestGetFileSystem(t *testing.T) {
	require.Equal(t, fs, GetFileSystem())
}

func TestCopyFile(t *testing.T) {
	t.Run("Happy path - file copied no override needed", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		fs := GetFileSystem()
		file, err := fs.Create("/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file.Close())
		}()

		err = CopyFile(file.Name(), "/file2.txt", false)
		require.NoError(t, err)

		file, err = fs.Open("/file2.txt")
		require.NoError(t, err)
		require.NotEqual(t, nil, file)
	})

	t.Run("Happy path - file copied override false", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		fs := GetFileSystem()
		file1, err := fs.Create("/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		_, err = file1.Write([]byte(file1.Name()))
		require.NoError(t, err)

		file2, err := fs.Create("/file2.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file2.Close())
		}()

		file2Data := []byte(file2.Name())
		_, err = file2.Write(file2Data)
		require.NoError(t, err)

		err = CopyFile(file1.Name(), file2.Name(), false)
		require.NoError(t, err)

		file2New, err := fs.Open(file2.Name())
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file2New.Close())
		}()

		file2NewData := make([]byte, len(file2.Name()))
		require.NoError(t, err)
		_, err = file2New.Read(file2NewData)

		require.Equal(t, file2Data, file2NewData)
	})

	t.Run("Happy path - file copied override true", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		fs := GetFileSystem()
		file1, err := fs.Create("/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		file1Data := []byte(file1.Name())
		_, err = file1.Write(file1Data)
		require.NoError(t, err)

		file2, err := fs.Create("/file2.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file2.Close())
		}()

		file2Data := []byte(file2.Name())
		_, err = file2.Write(file2Data)
		require.NoError(t, err)

		err = CopyFile(file1.Name(), file2.Name(), true)
		require.NoError(t, err)

		file2New, err := fs.Open(file2.Name())
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file2New.Close())
		}()

		file2NewData := make([]byte, len(file2.Name()))
		require.NoError(t, err)
		_, err = file2New.Read(file2NewData)

		require.Equal(t, file1Data, file2NewData)
	})

	t.Run("Error - Source file does not exist", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		err := CopyFile("randomFile", "moreRandom", true)
		require.Error(t, err)
	})

	t.Run("Error - dest file create", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		oldFsCreate := fsCreate
		defer func() {
			fsCreate = oldFsCreate
		}()
		fsCreate = func(name string) (file afero.File, e error) {
			return nil, errors.New("some Error")
		}

		fs := GetFileSystem()
		file1, err := fs.Create("/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		err = CopyFile(file1.Name(), "/moreRandom.txt", true)
		require.Error(t, err)
	})

	t.Run("Error - io copy", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		oldIoCopy := ioCopy
		defer func() {
			ioCopy = oldIoCopy
		}()
		ioCopy = func(dest io.Writer, src io.Reader) (i int64, e error) {
			return 0, errors.New("some error")
		}

		fs := GetFileSystem()
		file1, err := fs.Create("/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		err = CopyFile(file1.Name(), "/moreRandom.txt", true)
		require.Error(t, err)
	})

	t.Run("Error - file sync", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		oldFileSync := fileSync
		defer func() {
			fileSync = oldFileSync
		}()
		fileSync = func(file afero.File) error {
			return errors.New("some error")
		}

		fs := GetFileSystem()
		file1, err := fs.Create("/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		err = CopyFile(file1.Name(), "/moreRandom.txt", true)
		require.Error(t, err)
	})
}

func TestMustCopyFile(t *testing.T) {
	t.Run("Happy path - file copied no override needed", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		fs := GetFileSystem()
		file, err := fs.Create("/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file.Close())
		}()

		err = MustCopyFile(file.Name(), "/file2.txt")
		require.NoError(t, err)

		file, err = fs.Open("/file2.txt")
		require.NoError(t, err)
		require.NotEqual(t, nil, file)
	})

	t.Run("Happy path - file copied override needed", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		fs := GetFileSystem()
		file1, err := fs.Create("/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		file1Data := []byte(file1.Name())
		_, err = file1.Write(file1Data)
		require.NoError(t, err)

		file2, err := fs.Create("/file2.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file2.Close())
		}()

		file2Data := []byte(file2.Name())
		_, err = file2.Write(file2Data)
		require.NoError(t, err)

		err = MustCopyFile(file1.Name(), file2.Name())
		require.NoError(t, err)

		file2New, err := fs.Open(file2.Name())
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file2New.Close())
		}()

		file2NewData := make([]byte, len(file2.Name()))
		require.NoError(t, err)
		_, err = file2New.Read(file2NewData)

		require.Equal(t, file1Data, file2NewData)
	})

	t.Run("Error - Source file does not exist", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		err := MustCopyFile("randomFile", "moreRandom")
		require.Error(t, err)
	})

	t.Run("Error - dest file create", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		oldFsCreate := fsCreate
		defer func() {
			fsCreate = oldFsCreate
		}()
		fsCreate = func(name string) (file afero.File, e error) {
			return nil, errors.New("some Error")
		}

		fs := GetFileSystem()
		file1, err := fs.Create("/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		err = MustCopyFile(file1.Name(), "/moreRandom.txt")
		require.Error(t, err)
	})

	t.Run("Error - io copy", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		oldIoCopy := ioCopy
		defer func() {
			ioCopy = oldIoCopy
		}()
		ioCopy = func(dest io.Writer, src io.Reader) (i int64, e error) {
			return 0, errors.New("some error")
		}

		fs := GetFileSystem()
		file1, err := fs.Create("/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		err = MustCopyFile(file1.Name(), "/moreRandom.txt")
		require.Error(t, err)
	})

	t.Run("Error - file sync", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		oldFileSync := fileSync
		defer func() {
			fileSync = oldFileSync
		}()
		fileSync = func(file afero.File) error {
			return errors.New("some error")
		}

		fs := GetFileSystem()
		file1, err := fs.Create("/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		err = MustCopyFile(file1.Name(), "/moreRandom.txt")
		require.Error(t, err)
	})
}

func TestCopy(t *testing.T) {
	t.Run("Happy path - file copied no override needed", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		fs := GetFileSystem()
		file, err := fs.Create("/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file.Close())
		}()

		err = Copy(file.Name(), "/file2.txt")
		require.NoError(t, err)

		newFile, err := fs.Open("/file2.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, newFile.Close())
		}()
		require.NotEqual(t, nil, file)
	})

	t.Run("Happy path - folder copied no override needed", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		fs := GetFileSystem()

		require.NoError(t, fs.MkdirAll("/folder", os.ModePerm))

		file1, err := fs.Create("/folder/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		file2, err := fs.Create("/folder/file2.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file2.Close())
		}()

		err = Copy("/folder", "/folder2")
		require.NoError(t, err)

		_, err = fs.Stat("/folder2")
		require.Equal(t, nil, err)

		newFile1, err := fs.Open("/folder2/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, newFile1.Close())
		}()
		require.NotEqual(t, nil, file1)

		newFile2, err := fs.Open("/folder2/file2.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, newFile2.Close())
		}()
		require.NotEqual(t, nil, file2)
	})

	t.Run("Happy path - file copied override needed", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		fs := GetFileSystem()
		file1, err := fs.Create("/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		file1Data := []byte(file1.Name())
		_, err = file1.Write(file1Data)
		require.NoError(t, err)

		file2, err := fs.Create("/file2.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file2.Close())
		}()

		file2Data := []byte(file2.Name())
		_, err = file2.Write(file2Data)
		require.NoError(t, err)

		err = MustCopyFile(file1.Name(), file2.Name())
		require.NoError(t, err)

		file2New, err := fs.Open(file2.Name())
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file2New.Close())
		}()

		file2NewData := make([]byte, len(file2.Name()))
		require.NoError(t, err)
		_, err = file2New.Read(file2NewData)

		require.Equal(t, file1Data, file2NewData)
	})

	t.Run("Happy path - folder copied folder override needed", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		fs := GetFileSystem()

		require.NoError(t, fs.MkdirAll("/folder", os.ModePerm))

		file1, err := fs.Create("/folder/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		file2, err := fs.Create("/folder/file2.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file2.Close())
		}()

		require.NoError(t, fs.MkdirAll("/folder2", os.ModePerm))

		folder2File6, err := fs.Create("/folder2/file6.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, folder2File6.Close())
		}()

		folder2File7, err := fs.Create("/folder2/file7.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, folder2File7.Close())
		}()

		err = Copy("/folder", "/folder2")
		require.NoError(t, err)

		_, err = fs.Stat("/folder2")
		require.Equal(t, nil, err)

		_, err = fs.Stat("/folder2/file6.txt")
		require.NotEqual(t, nil, err)

		_, err = fs.Stat("/folder2/file7.txt")
		require.NotEqual(t, nil, err)

		newFile1, err := fs.Open("/folder2/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, newFile1.Close())
		}()
		require.NotEqual(t, nil, file1)

		newFile2, err := fs.Open("/folder2/file2.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, newFile2.Close())
		}()
		require.NotEqual(t, nil, file2)
	})

	t.Run("Error - Source file does not exist", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		err := Copy("randomFile", "moreRandom")
		require.Error(t, err)
	})

	t.Run("Error - dest file create", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		oldFsCreate := fsCreate
		defer func() {
			fsCreate = oldFsCreate
		}()
		fsCreate = func(name string) (file afero.File, e error) {
			return nil, errors.New("some Error")
		}

		fs := GetFileSystem()
		file1, err := fs.Create("/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		err = Copy(file1.Name(), "/moreRandom.txt")
		require.Error(t, err)
	})

	t.Run("Error - io copy", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		oldIoCopy := ioCopy
		defer func() {
			ioCopy = oldIoCopy
		}()
		ioCopy = func(dest io.Writer, src io.Reader) (i int64, e error) {
			return 0, errors.New("some error")
		}

		fs := GetFileSystem()
		file1, err := fs.Create("/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		err = Copy(file1.Name(), "/moreRandom.txt")
		require.Error(t, err)
	})

	t.Run("Error - file sync", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		oldFileSync := fileSync
		defer func() {
			fileSync = oldFileSync
		}()
		fileSync = func(file afero.File) error {
			return errors.New("some error")
		}

		fs := GetFileSystem()
		file1, err := fs.Create("/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		err = Copy(file1.Name(), "/moreRandom.txt")
		require.Error(t, err)
	})

	t.Run("Error - fs mkdirAll", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		oldFsMkdirAll := fsMkdirAll
		defer func() {
			fsMkdirAll = oldFsMkdirAll
		}()
		fsMkdirAll = func(path string, perm os.FileMode) error {
			return errors.New("some error")
		}

		fs := GetFileSystem()

		require.NoError(t, fs.MkdirAll("/folder", os.ModePerm))

		file1, err := fs.Create("/folder/file.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		err = Copy("/folder", "/brother")
		require.Error(t, err)
	})

	t.Run("Error - fs remove all", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		oldFsRemoveAll := fsRemoveAll
		defer func() {
			fsRemoveAll = oldFsRemoveAll
		}()
		fsRemoveAll = func(dst string) error {
			return errors.New("some error")
		}

		fs := GetFileSystem()

		require.NoError(t, fs.MkdirAll("/folder", os.ModePerm))

		file1, err := fs.Create("/folder/file.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		err = Copy("/folder", ".")
		require.Error(t, err)
	})

	t.Run("Error - read dir", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		oldReadDir := aferoReadDir
		defer func() {
			aferoReadDir = oldReadDir
		}()
		aferoReadDir = func(fs afero.Fs, dirname string) (infos []os.FileInfo, e error) {
			return nil, errors.New("some error")
		}

		fs := GetFileSystem()

		require.NoError(t, fs.MkdirAll("/folder", os.ModePerm))

		file1, err := fs.Create("/folder/file.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		err = Copy("/folder", "/someDir")
		require.Error(t, err)
	})

	t.Run("Error - recursive error", func(t *testing.T) {
		SetFileSystem(afero.NewMemMapFs())

		firstTime := true
		oldReadDir := aferoReadDir
		defer func() {
			aferoReadDir = oldReadDir
		}()
		aferoReadDir = func(fs afero.Fs, dirname string) (infos []os.FileInfo, e error) {
			if firstTime {
				infos, e = afero.ReadDir(fs, dirname)
				firstTime = false
			} else {
				infos, e = nil, errors.New("some error")
			}

			return
		}

		fs := GetFileSystem()

		require.NoError(t, fs.MkdirAll("/folder", os.ModePerm))
		require.NoError(t, fs.MkdirAll("/folder/folder2", os.ModePerm))

		file1, err := fs.Create("/folder1/folder2/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		err = Copy("/folder", "/someDir")
		require.Error(t, err)
	})

	t.Run("Happy path - skip symlink", func(t *testing.T) {
		SetFileSystem(afero.NewOsFs())

		fs := GetFileSystem().(*afero.OsFs)

		dir, err := os.Getwd()
		require.NoError(t, err)

		require.NoError(t, fs.MkdirAll(dir+"/testFolder", os.ModePerm))
		defer func() {
			require.NoError(t, fs.RemoveAll(dir+"/testFolder"))
		}()

		file1, err := fs.Create(dir + "/testFolder/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		handled, err := fs.SymlinkIfPossible(file1.Name(), dir+"/testFolder/file1Symlink.txt")
		require.Equal(t, true, handled)
		require.NoError(t, err)

		err = Copy(dir+"/testFolder", dir+"/testFolderCopied")
		defer func() {
			require.NoError(t, fs.RemoveAll(dir+"/testFolderCopied"))
		}()
		require.NoError(t, err)

		srcFiles, err := afero.ReadDir(fs, dir+"/testFolder")
		require.NoError(t, err)
		copiedFiles, err := afero.ReadDir(fs, dir+"/testFolderCopied")
		require.NoError(t, err)
		require.NotEqual(t, len(srcFiles), len(copiedFiles))
		require.Equal(t, len(srcFiles)-1, len(copiedFiles))
	})
}

func TestWriteFile(t *testing.T) {
	SetFileSystem(afero.NewMemMapFs())

	fs := GetFileSystem()

	t.Run("Happy path - file does not exist so it get's created", func(t *testing.T) {
		data := []byte("Hello World")

		err := WriteFile("newFile.txt", data)
		require.NoError(t, err)

		newFile, err := fs.Open("newFile.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, newFile.Close())
		}()

		newFileData := make([]byte, len(data))
		require.NoError(t, err)
		_, err = newFile.Read(newFileData)

		require.Equal(t, data, newFileData)
	})

	t.Run("Happy path - file exists so it get's overiden", func(t *testing.T) {
		data := []byte("Hello World Updated")

		err := WriteFile("newFile.txt", data)
		require.NoError(t, err)

		newFileOverride, err := fs.Open("newFile.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, newFileOverride.Close())
		}()

		newFileOverrideData := make([]byte, len(data))
		require.NoError(t, err)
		_, err = newFileOverride.Read(newFileOverrideData)

		require.Equal(t, data, newFileOverrideData)
	})

	t.Run("Error - fs create error", func(t *testing.T) {
		oldFsCreate := fsCreate
		defer func() {
			fsCreate = oldFsCreate
		}()
		fsCreate = func(name string) (file afero.File, e error) {
			return nil, errors.New("some Error")
		}

		err := WriteFile("newFile.txt", []byte("yay"))
		require.Error(t, err)
	})
}

func TestReadFile(t *testing.T) {
	SetFileSystem(afero.NewMemMapFs())

	fs := GetFileSystem()

	data := []byte("Hello World")

	t.Run("Happy path - file is read and data matches", func(t *testing.T) {
		file1, err := fs.Create("/file1.txt")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file1.Close())
		}()

		_, err = file1.Write(data)
		require.NoError(t, err)

		read, err := ReadFile("/file1.txt")
		require.NoError(t, err)

		require.Equal(t, data, read)
	})

	t.Run("Error - file does not exist", func(t *testing.T) {
		_, err := ReadFile("newFile.txt")
		require.Error(t, err)
	})

}
