package pac

import (
    "wio/cmd/wio/commands"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/log"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"

    "github.com/urfave/cli"
)

type VendorOp int

type Vendor struct {
    Context *cli.Context
    Op      VendorOp
}

type vendorInfo struct {
    dir  string
    name string
}

func (cmd Vendor) GetContext() *cli.Context {
    return cmd.Context
}

func (cmd Vendor) Execute() error {
    return nil
}

func addVendorPackage(info *vendorInfo) error {
    config, err := utils.ReadWioConfig(info.dir)
    if err != nil {
        return err
    }
    pkgDir := io.Path(info.dir, info.name)
    exists, err := io.Exists(pkgDir)
    if err != nil {
        return err
    }
    if !exists {
        return errors.Stringf("failed to find vendor/%s", info.name)
    }
    vendorConfig, err := utils.ReadWioConfig(pkgDir)
    if err != nil {
        return err
    }
    if vendorConfig.GetType() != constants.PKG {
        return errors.Stringf("project %s is not a package", info.name)
    }
    pkgConfig := vendorConfig.(*types.PkgConfig)
    pkgMeta := pkgConfig.MainTag.Meta
    if pkgMeta.Name != info.name {
        log.Warnln("package name %s does not match folder name %s", pkgMeta.Name, info.name)
    }
    tag := &types.DependencyTag{
        Version:        pkgMeta.Version,
        Vendor:         true,
        LinkVisibility: "PRIVATE",
    }
    config.GetDependencies()[pkgMeta.Name] = tag
    return utils.WriteWioConfig(info.dir, config)
}
