package utils

import (
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"
)

func ReadWioConfig(directory string) (types.IConfig, error) {
    wioPath := directory + io.Sep + "wio.yml"
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
