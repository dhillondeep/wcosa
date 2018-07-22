package vendor

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

func (cmd Cmd) GetContext() *cli.Context {
    return cmd.Context
}

func (cmd Cmd) Execute() error {
    dir, err := commands.GetDirectory(cmd)
    if err != nil {
        return err
    }
    if len(cmd.Context.Args()) <= 0 {
        return errors.String("missing vendor package name")
    }
    info := &Info{Dir: dir, Name: cmd.Context.Args()[0]}
    switch cmd.Op {
    case Add:
        return info.AddVendorPackage()
    case Remove:
        return info.RemoveVendorPackage()
    default:
        return nil
    }
}

func (info *Info) AddVendorPackage() error {
    config, err := utils.ReadWioConfig(info.Dir)
    if err != nil {
        return err
    }
    pkgDir := io.Path(info.Dir, io.Vendor, info.Name)
    exists := io.Exists(pkgDir)
    if !exists {
        return errors.Stringf("failed to find vendor/%s", info.Name)
    }
    vendorConfig, err := utils.ReadWioConfig(pkgDir)
    if err != nil {
        return err
    }
    if vendorConfig.GetType() != constants.Pkg {
        return errors.Stringf("project %s is not a package", info.Name)
    }
    if vendorConfig.GetName() != info.Name {
        log.Warnln("package name %s does not match folder name %s", vendorConfig.GetName(), info.Name)
    }
    tag := &types.DependencyImpl{
        Version: vendorConfig.GetVersion(),
        Vendor:  true,
    }
    config.AddDependency(vendorConfig.GetName(), tag)
    if err := utils.WriteWioConfig(info.Dir, config); err != nil {
        return err
    }
    log.Info(log.Cyan, "Added vendor dependency: ")
    log.Infoln(log.Green, "%s", info.Name)
    return nil
}

func (info *Info) RemoveVendorPackage() error {
    config, err := utils.ReadWioConfig(info.Dir)
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
    if err := utils.WriteWioConfig(info.Dir, config); err != nil {
        return err
    }
    log.Info(log.Cyan, "Removed vendor dependency: ")
    log.Infoln(log.Green, "%s", info.Name)
    return nil

NoRemove:
    log.Warnln("Vendor dependency %s not found", info.Name)
    return nil
}
