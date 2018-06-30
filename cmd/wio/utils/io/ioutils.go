// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package io contains helper functions related to io
// This file contains all the utilities available to be used from copying files to reading JSON
package io

import (
    "os"
    "path"
    "path/filepath"
    "runtime"
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

func Exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if os.IsNotExist(err) {
        return false, nil
    }
    if err != nil {
        return false, err
    }
    return true, nil
}
