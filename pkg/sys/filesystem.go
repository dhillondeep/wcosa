package sys

import (
	"errors"
	"github.com/spf13/afero"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var fs = afero.NewOsFs()

func SetFileSystem(givenFs afero.Fs) {
	fs = givenFs
}

func GetFileSystem() afero.Fs {
	return fs
}

// CopyFile copies file from src to destination and if destination file exists, it overrides the file
// content based on if override is specified
func CopyFile(source, destination string, override bool) error {
	if _, err := os.Stat(destination); err == nil && !override {
		return nil
	}

	srcFile, err := fs.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := fs.Create(destination) // creates if file doesn't exist
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

// Copy copies one path to another. The path could be a file or a directory
// It overrides the destination
func Copy(src string, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	if !Exists(src) {
		return errors.New("source path [" + src + "] does not exist")
	}
	if err := fs.RemoveAll(dst); err != nil {
		return err
	}
	si, err := fs.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return MustCopyFile(src, dst)
	}
	if _, err := fs.Stat(dst); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(dst, si.Mode()); err != nil {
		return err
	}
	entries, err := afero.ReadDir(fs, src)
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

// MustCopy copies file from src to destination and override if destination file exists
func MustCopyFile(source, destination string) error {
	return CopyFile(source, destination, true)
}

// CopyMultipleFiles copies multiple files from source to destination
func CopyMultipleFiles(sources []string, destinations []string, overrides []bool) error {
	if len(sources) != len(destinations) || len(destinations) != len(overrides) {
		return errors.New("length of sources, destinations, and overrides is not equal")
	}

	for i := 0; i < len(sources); i++ {
		if err := CopyFile(sources[i], destinations[i], overrides[i]); err != nil {
			return err
		}
	}

	return nil
}

// ReadFile reads the file and provides it's content
func ReadFile(fileName string) ([]byte, error) {
	file, err := fs.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}

// WriteFile writes text to a file on normal filesystem
func WriteFile(fileName string, data []byte) error {
	return ioutil.WriteFile(fileName, data, os.ModePerm)
}
