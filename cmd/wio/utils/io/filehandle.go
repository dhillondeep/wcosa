package io

import (
    "os"
    "io"
    "io/ioutil"
    "errors"
)

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
