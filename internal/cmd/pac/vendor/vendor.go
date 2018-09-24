package vendor

import (
    "wio/internal/cmd"
    "wio/internal/constants"
    "wio/internal/types"
    "wio/pkg/log"
    "wio/pkg/util"
    "wio/pkg/util/sys"

    "github.com/urfave/cli"
)

type CmdOp int

const (
    Add    CmdOp = 0
    Remove CmdOp = 1
)

type Cmd struct {
    Context *cli.Context
    Op      CmdOp
}

type Info struct {
    Dir  string
    Name string
}

func (c Cmd) GetContext() *cli.Context {
    return c.Context
}

func (c Cmd) Execute() error {
    dir, err := cmd.GetDirectory(c)
    if err != nil {
        return err
    }
    if len(c.Context.Args()) <= 0 {
        return util.Error("missing vendor package name")
    }
    info := &Info{Dir: dir, Name: c.Context.Args()[0]}
    switch c.Op {
    case Add:
        return info.AddVendorPackage()
    case Remove:
        return info.RemoveVendorPackage()
    default:
        return nil
    }
}

func (info *Info) AddVendorPackage() error {
    config, err := types.ReadWioConfig(info.Dir)
    if err != nil {
        return err
    }
    pkgDir := sys.Path(info.Dir, sys.Vendor, info.Name)
    exists := sys.Exists(pkgDir)
    if !exists {
        return util.Error("failed to find vendor/%s", info.Name)
    }
    vendorConfig, err := types.ReadWioConfig(pkgDir)
    if err != nil {
        return err
    }
    if vendorConfig.GetType() != constants.Pkg {
        return util.Error("project %s is not a package", info.Name)
    }
    if vendorConfig.GetName() != info.Name {
        log.Warnln("package name %s does not match folder name %s", vendorConfig.GetName(), info.Name)
    }
    tag := &types.DependencyImpl{
        Version: vendorConfig.GetVersion(),
        Vendor:  true,
    }
    config.AddDependency(vendorConfig.GetName(), tag)
    if err := types.WriteWioConfig(info.Dir, config); err != nil {
        return err
    }
    log.Info(log.Cyan, "Added vendor dependency: ")
    log.Infoln(log.Green, "%s", info.Name)
    return nil
}

func (info *Info) RemoveVendorPackage() error {
    config, err := types.ReadWioConfig(info.Dir)
    if err != nil {
        return err
    }
    deps := config.GetDependencies()
    if deps == nil {
        goto NoRemove
    }
    if _, exists := deps[info.Name]; !exists {
        goto NoRemove
    }
    delete(deps, info.Name)
    if err := types.WriteWioConfig(info.Dir, config); err != nil {
        return err
    }
    log.Info(log.Cyan, "Removed vendor dependency: ")
    log.Infoln(log.Green, "%s", info.Name)
    return nil

NoRemove:
    log.Warnln("Vendor dependency %s not found", info.Name)
    return nil
}
