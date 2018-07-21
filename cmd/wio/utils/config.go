package utils

import (
    "wio/cmd/wio/errors"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"
)

func ReadWioConfig(dir string) (types.Config, error) {
    path := io.Path(dir, io.Config)
    if !io.Exists(path) {
        return nil, errors.Stringf("path does not contain a wio.yml: %s", dir)
    }
    ret := &types.ConfigImpl{}
    err := io.NormalIO.ParseYml(path, ret)
    return ret, err
}

func WriteWioConfig(dir string, config types.Config) error {
    return io.NormalIO.WriteYml(io.Path(dir, io.Config), config)
}
