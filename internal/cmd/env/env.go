package env

import (
    "os/user"
    "wio/pkg/log"
    "wio/pkg/util/sys"

    "github.com/joho/godotenv"
    "github.com/urfave/cli"
)

var constantEnv = map[string]bool{
    "OS":      true,
    "WIOROOT": true,
}

type Env struct {
    Context *cli.Context
    Reset   bool
}

// get context for the command
func (env Env) GetContext() *cli.Context {
    return env.Context
}

// Runs the build command when cli build option is provided
func (env Env) Execute() error {
    currUser, err := user.Current()
    if err != nil {
        log.WriteFailure()
        return err
    }
    envFilePath := sys.Path(currUser.HomeDir, ".wio-usr", "wio.env")

    switch env.Reset {
    case true:
        log.Write(log.Cyan, "resetting wio environment... ")

        if err := CreateEnv(envFilePath); err != nil {
            log.WriteFailure()
            return err
        }
        log.WriteSuccess()
    case false:
        if err := handleEnvArguments(env.Context, envFilePath); err != nil {
            return err
        }
        break
    }

    return nil
}

func handleEnvArguments(cli *cli.Context, envFilePath string) error {
    if cli.NArg() == 1 {
        return showEnvironment(envFilePath, cli.Args().Get(0), false)
    } else if cli.NArg() == 2 {
        return editEnvironment(envFilePath, cli.Args().Get(0), cli.Args().Get(1))
    } else {
        return showEnvironment(envFilePath, "", true)
    }

    return nil
}

func showEnvironment(envFilePath string, key string, showAll bool) error {
    envData, err := godotenv.Read(envFilePath)
    if err != nil {
        return err
    }

    if showAll {
        for key, val := range envData {
            log.Write(log.Cyan, "%s=", key)
            log.Writeln("%s", val)
        }
    } else {
        if val, ok := envData[key]; ok {
            log.Write(log.Cyan, "%s=", key)
            log.Writeln("%s", val)
        } else {
            log.Write("%s => ", key)
            log.Writeln(log.Cyan, "no such environment key found!")
        }
    }

    return nil
}

func editEnvironment(envFilePath string, key string, value string) error {
    envData, err := godotenv.Read(envFilePath)
    if err != nil {
        return err
    }

    if val, ok := envData[key]; ok {
        log.Write(log.Cyan, "found environment key: ")
        log.Writeln("%s=%s", key, val)

        if _, ok := constantEnv[key]; ok {
            log.Write("%s => ", key)
            log.Writeln(log.Cyan, "env cannot be edited and is read only!")
        } else {
            log.Write(log.Cyan, "updated the environment key: ")
            log.Writeln("%s=%s", key, value)

            envData[key] = value

            // update wio.env file
            if err := godotenv.Write(envData, envFilePath); err != nil {
                return err
            }
        }
    } else {
        log.Write("%s => ", key)
        log.Writeln(log.Cyan, "no such environment key found!")
    }

    return nil
}

// Creates environment and overrides if there is an old environment
func CreateEnv(envFilePath string) error {
    wioRoot, err := sys.NormalIO.GetRoot()
    if err != nil {
        return err
    }

    // create wio.env file if it does not exist
    if err := godotenv.Write(map[string]string{
        "WIOROOT": wioRoot, "OS": sys.GetOS()}, envFilePath); err != nil {
        return err
    }

    return nil
}

func LoadEnv(envFilePath string) error {
    return godotenv.Load(envFilePath)
}
