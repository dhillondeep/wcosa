package sys

import (
    "encoding/json"

    yaml "gopkg.in/yaml.v2"
)

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
