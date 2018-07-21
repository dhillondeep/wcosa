package install

import (
    "wio/cmd/wio/commands"
    "wio/cmd/wio/log"
    "wio/cmd/wio/toolchain/npm/resolve"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"

    "github.com/urfave/cli"
)

type Cmd struct {
    Context *cli.Context

    dir    string
    info   *resolve.Info
    config types.Config
}

func (cmd Cmd) GetContext() *cli.Context {
    return cmd.Context
}

func (cmd Cmd) Execute() error {
    var err error
    cmd.dir, err = commands.GetDirectory(cmd)
    if err != nil {
        return err
    }
    cmd.config, err = utils.ReadWioConfig(cmd.dir)
    if err != nil {
        return err
    }
    cmd.info = resolve.NewInfo(cmd.dir)

    if len(cmd.Context.Args()) > 0 {
        if err := cmd.AddDependency(); err != nil {
            return err
        }
    }

    if err := cmd.info.ResolveRemote(cmd.config); err != nil {
        return err
    }
    return cmd.info.InstallResolved()
}

func (cmd Cmd) AddDependency() error {
    name, ver, err := cmd.getArgs(cmd.info)
    if err != nil {
        return err
    }
    log.Info(log.Cyan, "Adding dependency: ")
    log.Infoln(log.Green, "%s@%s", name, ver)
    deps := cmd.config.GetDependencies()
    if prev, exists := deps[name]; exists && prev.GetVersion() != ver {
        log.Warnln("Replacing previous version %s", prev.GetVersion())
    } else if exists {
        log.Warnln("Same version already exists")
    }
    deps[name] = &types.DependencyImpl{
        Version: ver,
        Vendor:  false,
    }
    return utils.WriteWioConfig(cmd.dir, cmd.config)
}
