package utils

import (
    "wio/cmd/wio/types"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/utils/io"
)

func ReadWioConfig(dir string) (types.IConfig, error) {
    wioPath := io.Path(dir, io.Config)
    exists, err := io.Exists(wioPath)
    if err != nil {
        return nil, err
    }
    if !exists {
        return nil, errors.Stringf("Path does not contain a wio.yml: %s", dir)
    }
    isApp, err := IsAppType(wioPath)
    if err != nil {
        return nil, err
    }
    var config types.IConfig
    if isApp {
        config = &types.AppConfig{}
    } else {
        config = &types.PkgConfig{}
    }
    err = io.NormalIO.ParseYml(wioPath, config)
    return config, err
}

func WriteWioConfig(dir string, config types.IConfig) error {
    wioPath := io.Path(dir, io.Config)
    exists, err := io.Exists(wioPath)
    if err != nil {
        return err
    }
    if !exists {
        return errors.Stringf("Path does not contain a wio.yml: %s", dir)
    }
    return types.PrettyPrint(config, wioPath)
}
