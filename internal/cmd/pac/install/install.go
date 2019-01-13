package install

import (
    "wio/internal/cmd"
    "wio/internal/types"
    "wio/pkg/log"
    "wio/pkg/npm/resolve"

    "github.com/urfave/cli"
)

type Cmd struct {
    Context *cli.Context

    dir    string
    info   *resolve.Info
    config types.Config
}

func (c Cmd) GetContext() *cli.Context {
    return c.Context
}

func (c Cmd) Execute() error {
    var err error
    c.dir, err = cmd.GetDirectory(c)
    if err != nil {
        return err
    }
    c.config, err = types.ReadWioConfig(c.dir, false)
    if err != nil {
        return err
    }
    c.info = resolve.NewInfo(c.dir)

    if len(c.Context.Args()) > 0 {
        if err := c.AddDependency(); err != nil {
            return err
        }
    }

    if err := c.info.ResolveRemote(c.config); err != nil {
        return err
    }
    return c.info.InstallResolved()
}

func (c Cmd) AddDependency() error {
    name, ver, err := c.getArgs(c.info)
    if err != nil {
        return err
    }
    log.Info(log.Cyan, "Adding dependency: ")
    log.Infoln(log.Green, "%s@%s", name, ver)
    deps := c.config.GetDependencies()
    if prev, exists := deps[name]; exists && prev.GetVersion() != ver {
        log.Warnln("Replacing previous version %s", prev.GetVersion())
    } else if exists {
        log.Warnln("Same version already exists")
    }
    c.config.AddDependency(name, &types.DependencyImpl{
        Version: ver,
        Vendor:  false,
    })
    return types.WriteWioConfig(c.dir, c.config)
}
