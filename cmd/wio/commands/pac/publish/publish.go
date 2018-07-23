package publish

import (
    "wio/cmd/wio/commands"
    "wio/cmd/wio/toolchain/npm/publish"
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
    cfg, err := utils.ReadWioConfig(dir)
    if err != nil {
        return err
    }
    return publish.Do(dir, cfg)
}
