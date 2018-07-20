package resolve

import (
    "wio/cmd/wio/errors"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
)

func tryFindConfig(name, ver, path string, strict bool) (*types.PkgConfig, error) {
	config, err := tryGetConfig(path)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}
	if config.Name() != name {
		return nil, errors.Stringf("config %s has wrong name", path)
	}
	if config.Version() != ver {
		if strict {
			return nil, errors.Stringf("config %s has wrong version", path)
		} else {
			return nil, nil
		}
	}
	return config, nil
}

func tryGetConfig(path string) (*types.PkgConfig, error) {
    wioPath := io.Path(path, io.Config)
    if !io.Exists(wioPath) {
        return nil, nil
    }
    isApp, err := utils.IsAppType(wioPath)
    if err != nil {
        return nil, err
    }
    if isApp {
        return nil, errors.Stringf("config %s is supposed to be package")
    }
    config := &types.PkgConfig{}
    err = io.NormalIO.ParseYml(wioPath, config)
    return config, err
}
