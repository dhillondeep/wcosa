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

const (
    Add VendorOp = 0
    Rm  VendorOp = 1
)

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
    dir, err := commands.GetDirectory(cmd)
    if err != nil {
        return err
    }
    if len(cmd.Context.Args()) <= 0 {
        return errors.String("missing vendor package name")
    }
    info := &vendorInfo{dir: dir, name: cmd.Context.Args()[0]}
    switch cmd.Op {
    case Add:
        return addVendorPackage(info)
    case Rm:
    }
    return errors.Stringf("invalid VendorOp %d", cmd.Op)
}

func addVendorPackage(info *vendorInfo) error {
    config, err := utils.ReadWioConfig(info.dir)
    if err != nil {
        return err
    }
    pkgDir := io.Path(info.dir, io.Vendor, info.name)
    exists := io.Exists(pkgDir)
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
    if config.GetDependencies() == nil {
        config.SetDependencies(make(types.DependenciesTag))
    }
    config.GetDependencies()[pkgMeta.Name] = tag
    return utils.WriteWioConfig(info.dir, config)
}
