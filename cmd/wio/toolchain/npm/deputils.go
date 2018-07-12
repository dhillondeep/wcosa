package npm

import (
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/utils"
)

func getOrFetchVersion(node *depTreeNode, info *depTreeInfo) (*packageVersion, error) {
    config, err := tryFindConfig(node, info)
    if err != nil {
        return nil, err
    }
    if config != nil {
        return configToVersion(config), nil
    }
    return fetchPackageVersion(node.name, node.version)
}

// Only Name, Version, and Dependencies are needed for dependency resolution
func configToVersion(config *types.PkgConfig) *packageVersion {
    return &packageVersion{
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
func tryFindConfig(node *depTreeNode, info *depTreeInfo) (*types.PkgConfig, error) {
    paths := []string{
        io.Path(info.baseDir, io.Vendor, node.name),
        io.Path(info.baseDir, io.Vendor, node.name+"__"+node.version),
        io.Path(info.baseDir, io.Folder, io.Modules, node.name+"__"+node.version),
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
        if tryConfig.Name() != node.name {
            return nil, errors.Stringf("config %s has wrong name", paths[i])
        }
        if tryConfig.Version() != node.version {
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
