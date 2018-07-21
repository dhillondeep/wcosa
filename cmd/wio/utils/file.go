package utils

import (
    "fmt"
    sysio "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/utils/io"
)

// Checks if the give path is a director and based on the returns
// true or false. If path does not exist, it throws an error
func IsDir(path string) (bool, error) {
    fi, err := os.Stat(path)
    if err != nil {
        return false, err
    }

    return fi.IsDir(), nil
}

// This checks if the directory is empty or not
func IsEmpty(name string) (bool, error) {
    f, err := os.Open(name)
    if err != nil {
        return false, err
    }
    defer f.Close()

    _, err = f.Readdirnames(1) // Or f.Readdir(1)
    if err == sysio.EOF {
        return true, nil
    }
    return false, err // Either not empty or error, suits both cases
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) error {
    if !io.Exists(src) {
        msg := fmt.Sprintf("Path [%s] does not exist", src)
        return errors.String(msg)
    }

    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer func() {
        if e := out.Close(); e != nil {
            err = e
        }
    }()

    _, err = sysio.Copy(out, in)
    if err != nil {
        return err
    }

    err = out.Sync()
    if err != nil {
        return err
    }

    si, err := os.Stat(src)
    if err != nil {
        return err
    }
    err = os.Chmod(dst, si.Mode())
    if err != nil {
        return err
    }

    return nil
}

func Copy(src string, dst string) error {
    src = filepath.Clean(src)
    dst = filepath.Clean(dst)

    si, err := os.Stat(src)
    if err != nil {
        return err
    }
    if si.IsDir() {
        return CopyDir(src, dst)
    } else {
        return CopyFile(src, dst)
    }
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
    if !io.Exists(src) {
        return
    }

    src = filepath.Clean(src)
    dst = filepath.Clean(dst)

    si, err := os.Stat(src)
    if err != nil {
        return err
    }
    if !si.IsDir() {
        return fmt.Errorf("source is not a directory")
    }

    _, err = os.Stat(dst)
    if err != nil && !os.IsNotExist(err) {
        return
    }
    if err == nil {
        return fmt.Errorf("destination already exists")
    }

    err = os.MkdirAll(dst, si.Mode())
    if err != nil {
        return
    }

    entries, err := ioutil.ReadDir(src)
    if err != nil {
        return
    }

    for _, entry := range entries {
        srcPath := filepath.Join(src, entry.Name())
        dstPath := filepath.Join(dst, entry.Name())

        if entry.IsDir() {
            err = CopyDir(srcPath, dstPath)
            if err != nil {
                return
            }
        } else {
            // Skip symlinks.
            if entry.Mode()&os.ModeSymlink != 0 {
                continue
            }

            err = CopyFile(srcPath, dstPath)
            if err != nil {
                return
            }
        }
    }

    return
}
