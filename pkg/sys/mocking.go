package sys

import (
	"github.com/dhillondeep/afero"
	"io"
	"os"
	"path/filepath"
)

////////////// Dependency injection for testing ////////////////
var fsCreate = func(name string) (afero.File, error) {
	return fs.Create(name)
}

var ioCopy = func(dest io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dest, src)
}

var fileSync = func(file afero.File) error {
	return file.Sync()
}

var aferoReadDir = func(fs afero.Fs, dirname string) ([]os.FileInfo, error) {
	return afero.ReadDir(fs, dirname)
}

var fsMkdirAll = func(path string, perm os.FileMode) error {
	return fs.MkdirAll(path, perm)
}

var fsRemoveAll = func(dst string) error {
	return fs.RemoveAll(dst)
}

var osLstat = func(name string) (os.FileInfo, error) {
	return os.Lstat(name)
}

var osExecutable = func() (string, error) {
	return os.Executable()
}

var osReadlink = func(name string) (string, error) {
	return os.Readlink(name)
}

var filepathAbs = func(path string) (string, error) {
	return filepath.Abs(path)
}
