package sys

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dhillondeep/afero"
)

var fs = afero.NewOsFs()

// SetFileSystem set filesystem used inside the application
func SetFileSystem(givenFs afero.Fs) {
	fs = givenFs
}

// GetFileSystem provides the filesystem being used
func GetFileSystem() afero.Fs {
	return fs
}

// CopyFile copies file from src to destination and if destination file exists, it overrides the file
// content based on if override is specified
func CopyFile(source, destination string, override bool) (err error) {
	if _, err := fs.Stat(destination); err == nil && !override {
		return nil
	}

	srcFile, err := fs.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := fsCreate(destination) // creates if file doesn't exist
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = ioCopy(destFile, srcFile) // check first var for number of bytes copied
	if err != nil {
		return err
	}

	err = fileSync(destFile)
	if err != nil {
		return err
	}

	return nil
}

// MustCopyFile copies file from src to destination and override if destination file exists
func MustCopyFile(source, destination string) error {
	return CopyFile(source, destination, true)
}

// Copy copies one path to another. The path could be a file or a directory
// It overrides the destination
func Copy(src string, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	if !Exists(src) {
		return errors.New("source path [" + src + "] does not exist")
	}
	if err := fsRemoveAll(dst); err != nil {
		return err
	}
	si, _ := fs.Stat(src)

	if !si.IsDir() {
		return MustCopyFile(src, dst)
	}

	if err := fsMkdirAll(dst, si.Mode()); err != nil {
		return err
	}

	entries, err := aferoReadDir(fs, src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.Mode()&os.ModeSymlink != 0 {
			continue // skip symlinks
		}
		if err := Copy(srcPath, dstPath); err != nil {
			return err
		}
	}
	return nil
}

// ReadFile reads the file and provides it's content
func ReadFile(fileName string) (data []byte, err error) {
	file, err := fs.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}

// WriteFile writes text to a file on normal filesystem
func WriteFile(fileName string, data []byte) (err error) {
	f, err := fsCreate(fileName)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(data)
	return
}
