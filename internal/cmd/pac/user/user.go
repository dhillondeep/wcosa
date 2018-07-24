package user

import "github.com/urfave/cli"

type Login struct {
    Context *cli.Context
}

type Logout struct {
    Context *cli.Context
}

func (cmd Login) GetContext() *cli.Context {
    return cmd.Context
}

func (cmd Logout) GetContext() *cli.Context {
    return cmd.Context
}
