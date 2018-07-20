package resolve

import (
    "wio/cmd/wio/errors"
    "wio/cmd/wio/toolchain/npm"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
)

func findVersion(node *Node, dir string) (*npm.Version, error) {
    config, err := tryFindConfig(node, dir)
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

// This function searched local filesystem for the `wio.yml` of the
// desired package and version. The function looks in the places
// -- $BASE_DIR/vendor/[name]
// -- $BASE_DIR/vendor/[name]__[version]
// -- $BASE_DIR/.wio/node_modules/[name]__[version]
//
// Function returns nil error and nil result if not found.
// Vendor is preferred to allow overrides.
func tryFindConfig(node *Node, dir string) (*types.PkgConfig, error) {
    name := node.Name
    ver := node.ResolvedVersion.Str()

    paths := []string{
        io.Path(dir, io.Vendor, name),
        io.Path(dir, io.Vendor, name+"__"+ver),
        io.Path(dir, io.Folder, io.Modules, name+"__"+ver),
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
            return nil, errors.Stringf("config %s has wrong version", paths[i])
        }
        config = tryConfig

        if i == 0 || i == 1 {
            // package is found in vendor directory
            node.Vendor = true
        } else {
            node.Vendor = false
        }

        node.Path = paths[i]
        node.Config = config
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
