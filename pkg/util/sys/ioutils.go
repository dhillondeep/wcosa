// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package io contains helper functions related to io
// This file contains all the utilities available to be used from copying files to reading JSON
package sys

import (
    "bytes"
    "errors"
    "io"
    "io/ioutil"
    "os"
    "path"
    "path/filepath"
    "runtime"
)

const (
    Folder   = ".wio"
    Temp     = ".tmp"
    Config   = "wio.yml"
    Modules  = "node_modules"
    Vendor   = "vendor"
    Download = "cache"
)

const (
    WINDOWS = "windows"
    DARWIN  = "darwin"
    LINUX   = "linux"
)

var NormalIO NormalHandler = 0
var AssetIO AssetHandler = 1
var Sep = string(filepath.Separator) // separator based on the OS

// Returns the root path to the files in terms of this executable
func (normalHandler NormalHandler) GetRoot() (string, error) {
    ex, err := os.Executable()
    if err != nil {
        return "", err
    }

    fileInfo, err := os.Lstat(ex)
    if err != nil {
        return "", err
    }

    if fileInfo.Mode()&os.ModeSymlink != 0 {
        newPath, err := os.Readlink(ex)
        if err != nil {
            return "", nil
        }

        newPath = filepath.Dir(newPath)

        // check if the path is relative
        if !filepath.IsAbs(newPath) {
            oldPath := filepath.Dir(ex)
            ex, err = filepath.Abs(path.Join(oldPath, newPath))
            if err != nil {
                return "", err
            }
        } else {
            ex = newPath
        }
    } else {
        ex = filepath.Dir(ex)
    }

    return ex, nil
}

// Returns the root path to the asset files in terms of assets folder
func (assetHandler AssetHandler) GetRoot() (string, error) {
    return "assets", nil
}

// Returns operating system from three types (windows, darwin, and linux)
func GetOS() string {
    goos := runtime.GOOS

    if goos == "windows" {
        return WINDOWS
    } else if goos == "darwin" {
        return DARWIN
    } else {
        return LINUX
    }
}

func Exists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}

func Path(values ...string) string {
    var buffer bytes.Buffer
    for _, value := range values {
        buffer.WriteString(value)
        buffer.WriteString(Sep)
    }
    pth := buffer.String()
    return filepath.Clean(pth[:len(pth)-1])
}

func CopyFile(src string, dst string) error {
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()
    dstFile, err := os.Create(dst) // creates if file doesn't exist
    if err != nil {
        return err
    }
    defer dstFile.Close()
    if _, err := io.Copy(dstFile, srcFile); err != nil {
        return err
    }
    return dstFile.Sync()
}

func Copy(src string, dst string) error {
    src = filepath.Clean(src)
    dst = filepath.Clean(dst)
    if !Exists(src) {
        return errors.New("source path [" + src + "] does not exist")
    }
    if err := os.RemoveAll(dst); err != nil {
        return err
    }
    si, err := os.Stat(src)
    if err != nil {
        return err
    }
    if !si.IsDir() {
        return CopyFile(src, dst)
    }
    if _, err := os.Stat(dst); err != nil && !os.IsNotExist(err) {
        return err
    }
    if err := os.MkdirAll(dst, si.Mode()); err != nil {
        return err
    }
    entries, err := ioutil.ReadDir(src)
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

// Checks if the give path is a director and based on the returns
// true or false. If path does not exist, it throws an error
func IsDir(path string) (bool, error) {
    fi, err := os.Stat(path)
    if err != nil {
        return false, err
    }

    return fi.IsDir(), nil
}
