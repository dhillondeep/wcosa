// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package io contains helper functions related to io
// This file contains all the utilities available to be used from copying files to reading JSON
package sys

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var sep = string(filepath.Separator)

// GetRoot returns root folder where the executable is located
func GetRoot() (string, error) {
	ex, err := osExecutable()
	if err != nil {
		return "", err
	}

	fileInfo, err := osLstat(ex)
	if err != nil {
		return "", err
	}

	if fileInfo.Mode()&os.ModeSymlink != 0 {
		newPath, err := osReadlink(ex)
		if err != nil {
			return "", err
		}

		newPath = filepath.Dir(newPath)

		// check if the path is relative
		if !filepath.IsAbs(newPath) {
			ex, err = filepathAbs(filepath.Join(filepath.Dir(ex), newPath))
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

// GetOS returns operating system from three types (windows, darwin, and linux)
func GetOS() string {
	return runtime.GOOS
}

// GetArch returns architecture
func GetArch() string {
	return strings.ToLower(runtime.GOARCH)
}

// Exists checks if path
func Exists(path string) bool {
	_, err := GetFileSystem().Stat(path)
	return err == nil
}

// GetSep returns separator based on the OS
func GetSep() string {
	return sep
}

// JoinPaths joins paths provided using operating system specific separator
func JoinPaths(values ...string) string {
	if len(values) > 0 {
		var buffer bytes.Buffer
		for _, value := range values {
			buffer.WriteString(value)
			buffer.WriteString(GetSep())
		}
		pth := buffer.String()
		return filepath.Clean(pth[:len(pth)-1])
	} else {
		return ""
	}
}

// IsDir checks if the given path is a directory. If path does not exist, error is thrown
func IsDir(path string) (bool, error) {
	fi, err := GetFileSystem().Stat(path)
	if err != nil {
		return false, err
	}

	return fi.IsDir(), nil
}
