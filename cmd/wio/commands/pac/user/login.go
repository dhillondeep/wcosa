package user

import (
    "errors"
    "wio/cmd/wio/commands"
    "wio/cmd/wio/log"
    "wio/cmd/wio/toolchain/npm/login"
)

type loginArgs struct {
    dir   string
    name  string
    pass  string
    email string
}

func (cmd Login) getArgs() (*loginArgs, error) {
    ctx := cmd.Context
    dir, err := commands.GetDirectory(cmd)
    if err != nil {
        return nil, err
    }
    args := ctx.Args()
    if len(args) < 3 {
        return nil, errors.New("wio login [username] [password] [email]")
    }
    return &loginArgs{
        dir:   dir,
        name:  args[0],
        pass:  args[1],
        email: args[2],
    }, nil
}

func (cmd Login) Execute() error {
    args, err := cmd.getArgs()
    if err != nil {
        return err
    }
    log.Info(log.Cyan, "Sending login info ... ")
    token, err := login.GetToken(args.name, args.pass, args.email)
    if err != nil {
        log.WriteFailure()
        return err
    }
    log.WriteSuccess()
    log.Info(log.Cyan, "Saving login token ... ")
    if err := token.Save(args.dir); err != nil {
        log.WriteFailure()
        return err
    }
    log.WriteSuccess()
    return nil
}
