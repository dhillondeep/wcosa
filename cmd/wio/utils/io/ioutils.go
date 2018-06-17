// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package io contains helper functions related to io
// This file contains all the utilities available to be used from copying files to reading JSON
package io

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
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

// Copies file from src to destination and if destination file exists, it overrides the file
// content based on if override is specified. Copies file from OS filesystem
func (normalHandler NormalHandler) CopyFile(source string, destination string, override bool) error {
	if _, err := os.Stat(destination); err == nil && !override {
		return nil
	}

	srcFile, err := os.Open(source)

	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destination) // creates if file doesn't exist
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

// Copies file from src to destination and if destination file exists, it overrides the file
// content based on if override is specified. Copies file from binary assets
func (assetHandler AssetHandler) CopyFile(source string, destination string, override bool) error {
	if _, err := os.Stat(destination); err == nil && !override {
		return nil
	}

	rootPath, err := assetHandler.GetRoot()
	if err != nil {
		return err
	}

	dest, err := os.Create(destination) // creates if file doesn't exist
	if err != nil {
		return err
	}
	defer dest.Close()

	srcData, err := Asset(rootPath + Sep + source)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destination, srcData, os.ModePerm)
	if err != nil {
		return err
	}

	err = dest.Sync()
	return err
}

// Copies multiple files from source to destination. Source files are from filesystem
func (normalHandler NormalHandler) CopyMultipleFiles(sources []string, destinations []string, overrides []bool) error {
	if len(sources) != len(destinations) || len(destinations) != len(overrides) {
		return errors.New("length of sources, destinations and overrides is not equal")
	}

	for i := 0; i < len(sources); i++ {
		if err := normalHandler.CopyFile(sources[i], destinations[i], overrides[i]); err != nil {
			return err
		}
	}

	return nil
}

// Copies multiple files from source to destination. Source files are from binary assets
func (assetHandler AssetHandler) CopyMultipleFiles(sources []string, destinations []string, overrides []bool) error {
	if len(sources) != len(destinations) || len(destinations) != len(overrides) {
		return errors.New("length of sources, destinations and overrides is not equal")
	}

	for i := 0; i < len(sources); i++ {
		if err := assetHandler.CopyFile(sources[i], destinations[i], overrides[i]); err != nil {
			return err
		}
	}

	return nil
}

// Reads the file and provides it's content as a string. From normal filesystem
func (normalHandler NormalHandler) ReadFile(fileName string) ([]byte, error) {
	buff, err := ioutil.ReadFile(fileName)
	return buff, err
}

// Reads the file and provides it's content as a string. From binary assets
func (assetHandler AssetHandler) ReadFile(fileName string) ([]byte, error) {
	rootPath, err := assetHandler.GetRoot()
	if err != nil {
		return nil, err
	}

	return Asset(rootPath + Sep + fileName)
}

// Writes text to a file on normal filesystem
func (normalHandler NormalHandler) WriteFile(fileName string, data []byte) error {
	return ioutil.WriteFile(fileName, data, os.ModePerm)
}

// Writes text to binary assets (invalid to do)
func (assetHandler AssetHandler) WriteFile(fileName string, data []byte) error {
	return errors.New("assets are readonly and cannot be modified")
}

// Parses JSON from the file on filesystem
func (normalHandler NormalHandler) ParseJson(fileName string, out interface{}) (err error) {
	text, err := normalHandler.ReadFile(fileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(text), out)
	return err
}

// Parses JSON from the data in assets
func (assetHandler AssetHandler) ParseJson(fileName string, out interface{}) (err error) {
	text, err := assetHandler.ReadFile(fileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(text), out)
	return err
}

// Parses YML from the file on filesystem
func (normalHandler NormalHandler) ParseYml(fileName string, out interface{}) error {
	text, err := normalHandler.ReadFile(fileName)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(text, out)
}

// Parses YML from the data in assets
func (assetHandler AssetHandler) ParseYml(fileName string, out interface{}) error {
	text, err := assetHandler.ReadFile(fileName)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(text, out)
}

// Writes JSON data to a file on filesystem
func (normalHandler NormalHandler) WriteJson(fileName string, in interface{}) error {
	data, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return err
	}

	return normalHandler.WriteFile(fileName, data)
}

// Writes JSON data to a binary asset (not valid)
func (assetHandler AssetHandler) WriteJson(fileName string, in interface{}) error {
	return assetHandler.WriteFile(fileName, nil)
}

// Writes YML data to a file on filesystem
func (normalHandler NormalHandler) WriteYml(fileName string, in interface{}) error {
	data, err := yaml.Marshal(in)
	if err != nil {
		return err
	}

	return normalHandler.WriteFile(fileName, data)
}

// Writes YML data to a binary asset (not valid)
func (assetHandler AssetHandler) WriteYml(fileName string, in interface{}) error {
	return assetHandler.WriteFile(fileName, nil)
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
