package resolve

import (
    "strings"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/toolchain/npm"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
)

func findVersion(name string, ver string, dir string) (*npm.Version, error) {
    config, err := tryFindConfig(name, ver, dir)
    if err != nil {
        return nil, err
    }
    if config == nil {
        return nil, nil
    }
    return configToVersion(config), nil
}

// Only Name, Version, and Dependencies are needed for dependency resolution
func configToVersion(config *types.PkgConfig) *npm.Version {
    return &npm.Version{
        Name:         config.Name(),
        Version:      config.Version(),
        Dependencies: config.Dependencies(),
    }
}

// This function searches remote packages location (node_modules) for the `wio.yml` of the
// desired package and version.
func tryFindConfig(name string, ver string, dir string) (*types.PkgConfig, error) {
    path := io.Path(dir, io.Folder, io.Modules, name) + "__" + ver

    if !utils.PathExists(path) {
        return nil, nil
    }

    tryConfig, err := tryGetConfig(path)
    if err != nil {
        return nil, err
    }

    if tryConfig.Name() != name {
        // this happens when wio.yml is of unsupported type
        if strings.Trim(tryConfig.GetMainTag().GetConfigurations().WioVersion, " ") == "" {
            return nil, errors.Stringf("config [%s] => unsupported wio version: < 0.3.0", path)
        } else {
            return nil, errors.Stringf("config [%s] has wrong name: %s", path, tryConfig.Name())
        }
    }

    if tryConfig.Version() != ver {
        return nil, errors.Stringf("config [%s] has wrong version: %s", path, tryConfig.Version())
    }
    return tryConfig, nil
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
