package env

import (
    "wio/internal/config/root"
    "wio/pkg/log"

    "github.com/urfave/cli"
)

const (
    RESET = 0
    UNSET = 1
    SET   = 2
)

type Env struct {
    Context *cli.Context
    Command byte
}

// get context for the command
func (env Env) GetContext() *cli.Context {
    return env.Context
}

// Runs the build command when cli build option is provided
func (env Env) Execute() error {
    switch env.Command {
    case RESET:
        log.Write(log.Cyan, "resetting wio environment... ")

        if err := root.CreateEnv(); err != nil {
            log.WriteFailure()
            return err
        }
        log.WriteSuccess()
        break
    case UNSET:
        break
    case SET:
        break
    }

    return nil
}
