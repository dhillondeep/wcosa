package npm

import (
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/utils"
)

func getOrFetchVersion(name string, ver string, dir string) (*Version, error) {
    config, err := tryFindConfig(name, ver, dir)
    if err != nil {
        return nil, err
    }
    if config != nil {
        return configToVersion(config), nil
    }
    return fetchPackageVersion(name, ver)
}

// Only Name, Version, and Dependencies are needed for dependency resolution
func configToVersion(config *types.PkgConfig) *Version {
    return &Version{
        Name: config.Name(),
        Version: config.Version(),
        Dependencies: config.Dependencies(),
    }
}

// This function searched local filesystem for the `wio.yml` of the
// desired package and version. The function looks in the places
// -- $BASE_DIR/vendor/[name]
// -- $BASE_DIR/vendor/[name]__[version]
// -- $BASE_DIR/.wio/node_modules/[name]__[version]
//
// Function returns nil error and nil result if not found.
// Vendor is preferred to allow overrides.
func tryFindConfig(name string, ver string, baseDir string) (*types.PkgConfig, error) {
    paths := []string{
        io.Path(baseDir, io.Vendor, name),
        io.Path(baseDir, io.Vendor, name+"__"+ver),
        io.Path(baseDir, io.Folder, io.Modules, name+"__"+ver),
    }
    var config *types.PkgConfig = nil
    for i := 0; config == nil && i < len(paths); i++ {
        tryConfig, err := tryGetConfig(paths[i])
        if err != nil {
            return nil, err
        }
        if tryConfig == nil {
            continue
        }
        if tryConfig.Name() != name {
            return nil, errors.Stringf("config %s has wrong name", paths[i])
        }
        if tryConfig.Version() != ver {
            if i != 0 {
                return nil, errors.Stringf("config %s has wrong version", paths[i])
            } else {
                // version-less path
                continue
            }
        }
        config = tryConfig
    }
    return config, nil
}

func tryGetConfig(wioPath string) (*types.PkgConfig, error) {
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
