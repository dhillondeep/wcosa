package env

import (
    "wio/internal/config/root"
    "wio/pkg/log"

    "github.com/joho/godotenv"

    "github.com/urfave/cli"
)

const (
    RESET = 0
    UNSET = 1
    SET   = 2
    VIEW  = 3
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
    case VIEW:
        return env.viewCommand()
    }

    return nil
}

// Display environment
func (env Env) viewCommand() error {
    envData, err := godotenv.Read(root.GetEnvFilePath())
    if err != nil {
        return err
    }

    isReadOnly := func(envName string) string {
        readOnlyText := "readonly"
        otherText := "        "
        if envValue, exists := envMeta[envName]; !exists {
            return otherText
        } else if envValue {
            return readOnlyText
        } else {
            return otherText
        }
    }

    printEnvs := func(envName, envValue string) {
        log.Write(log.Green, isReadOnly(envName)+"  ")
        log.Write(log.Cyan, envName+"=")
        log.Writeln(envValue)
    }

    numArgs := env.Context.NArg()

    if numArgs == 0 {
        for envName, envValue := range envData {
            printEnvs(envName, envValue)
        }
    } else {
        for i := 0; i < numArgs; i++ {
            if envValue, exists := envData[env.Context.Args().Get(i)]; exists {
                printEnvs(env.Context.Args().Get(i), envValue)
            }
        }
    }

    return nil
}
