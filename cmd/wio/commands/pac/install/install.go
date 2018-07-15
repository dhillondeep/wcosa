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
}

func (cmd Cmd) GetContext() *cli.Context {
    return cmd.Context
}

func (cmd Cmd) Execute() error {
    dir, err := commands.GetDirectory(cmd)
	if err != nil {
		return err
	}
    info := resolve.NewInfo(dir)
    name, ver, err := cmd.getArgs(info)
    if err != nil {
        return err
    }
    config, err := utils.ReadWioConfig(dir)
    if err != nil {
        return err
    }
    config.GetDependencies()[name] = &types.DependencyTag{
        Version:        ver,
        Vendor:         false,
        LinkVisibility: "PRIVATE",
    }
    log.Infoln(log.Cyan, "Adding dependency %s@%s", name, ver)
    if err := utils.WriteWioConfig(dir, config); err != nil {
        return err
    }

	return nil
}
