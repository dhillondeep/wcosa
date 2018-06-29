package utils

import (
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "regexp"

    goerr "errors"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/log"
    "wio/cmd/wio/types"
    wio "wio/cmd/wio/utils/io"
    "strings"
)

// Checks if path exists and returns true and false based on that
func PathExists(path string) bool {
    if _, err := os.Stat(path); err != nil {
        return false
    }
    return true
}

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
    if err == io.EOF {
        return true, nil
    }
    return false, err // Either not empty or error, suits both cases
}

// It takes in a slice and an element and then ut appends that element to the slice only
// if that element in not already in the slice
func AppendIfMissingElem(slice []string, i string) []string {
    for _, ele := range slice {
        if ele == i {
            return slice
        }
    }
    return append(slice, i)
}

// It takes two slices and appends the second one onto the first one. It does
// not allow duplicates
func AppendIfMissing(slice []string, slice2 []string) []string {
    newSlice := make([]string, 0)

    for _, ele1 := range slice {
        newSlice = AppendIfMissingElem(newSlice, ele1)
    }

    for _, ele2 := range slice2 {
        newSlice = AppendIfMissingElem(newSlice, ele2)
    }

    return newSlice
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
    if !PathExists(src) {
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

    _, err = io.Copy(out, in)
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

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
    if !PathExists(src) {
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

// Checks if the config file is of App type or Pkg type
func IsAppType(wioPath string) (bool, error) {
    // read wio.yml file to see which project type we are building
    data, err := wio.NormalIO.ReadFile(wioPath)
    if err != nil {
        return false, err
    }

    // regex expression to check for app type
    pat := regexp.MustCompile(`(^app:)|((\s| |^\w)app:(\s+|))`)
    s := pat.FindString(string(data))

    return s != "", nil
}

//  Eeturns elements in a that aren't in b
func Difference(a, b []string) []string {
    mb := map[string]bool{}
    for _, x := range b {
        mb[x] = true
    }
    var ab []string
    for _, x := range a {
        if _, ok := mb[x]; !ok {
            ab = append(ab, x)
        }
    }
    return ab
}

// Read config file and return config object
func ReadWioConfig(path string) (*types.Config, error) {
    defer func() {
        if r := recover(); r != nil {
            configError := errors.ConfigParsingError{
                Err: goerr.New("fatal error occurred while parsing wio.yml file"),
            }

            log.WriteErrorlnExit(configError)
        }
    }()

    isApp, err := IsAppType(path)
    if err != nil {
        return nil, err
    }

    if !isApp {
        pkgConfig := &types.PkgConfig{}

        // try to read pkg config file
        if err := wio.NormalIO.ParseYml(path, pkgConfig); err != nil {
            configError := errors.ConfigParsingError{
                Err: goerr.New("wio.yml file could not be parsed for project of type: project"),
            }

            return nil, configError
        }

        return &types.Config{Config: pkgConfig, Type: types.Pkg}, nil
    } else {
        appConfig := &types.AppConfig{}

        // try to read pkg config file
        if err := wio.NormalIO.ParseYml(path, appConfig); err != nil {
            configError := errors.ConfigParsingError{
                Err: goerr.New("wio.yml file could not be parsed for project of type: application"),
            }

            return nil, configError
        }

        return &types.Config{Config: appConfig, Type: types.App}, nil
    }
}

func Contains(slice []string, value string) bool {
    for _, element := range slice {
        if element == value {
            return true
        }
    }
    return false
}

func ContainsNoCase(slice []string, value string) bool {
    for _, element := range slice {
        if strings.ToLower(element) == strings.ToLower(value) {
            return true
        }
    }
    return false
}
