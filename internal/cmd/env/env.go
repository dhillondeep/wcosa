package env

import (
    "os"
    "strings"
    "wio/internal/config/root"
    "wio/internal/constants"
    "wio/pkg/log"
    "wio/pkg/util"
    "wio/pkg/util/sys"

    "github.com/joho/godotenv"

    "github.com/urfave/cli"
)

const (
    RESET = 0
    UNSET = 1
    SET   = 2
    VIEW  = 3
)

const (
    ENVBOOLVALUE = "$BOOL___VAL$"
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
    envPath := root.GetEnvFilePath()
    localFlagDefined := env.Context.Bool("local")

    var err error
    if localFlagDefined {
        envPath, err = getLocalEnvPath()
        if err != nil {
            return err
        }

        if !sys.Exists(envPath) {
            if err := root.CreateLocalEnv(envPath); err != nil {
                return err
            }
        }
    }

    switch env.Command {
    case RESET:
        log.Write(log.Cyan, "resetting wio environment... ")

        if localFlagDefined {
            if err := root.CreateLocalEnv(envPath); err != nil {
                log.WriteFailure()
                return err
            }
        } else if err := root.CreateEnv(); err != nil {
            log.WriteFailure()
            return err
        }

        log.WriteSuccess()
        break
    case UNSET:
        env.unsetCommand(env.Context.Args(), env.Context.NArg(), envPath)
        break
    case SET:
        env.setCommand(env.Context.Args(), env.Context.NArg(), envPath)
        break
    case VIEW:
        return env.viewCommand()
    }

    return nil
}

// unset env variables
func (env Env) unsetCommand(args cli.Args, numArgs int, envPath string) error {
    envData, err := godotenv.Read(envPath)
    if err != nil {
        return err
    }

    unsetEnv := func(envName string) bool {
        if isReadOnly(envName) {
            log.WriteFailure()
            log.Errln("%s is readonly and cannot be changed", envName)
            return false
        } else {
            if _, exists := envData[envName]; exists {
                delete(envData, envName)
                log.WriteSuccess()
                return true
            } else {
                log.WriteFailure()
                log.Errln("%s does not exist as an environment variable", envName)
                return false
            }
        }
    }

    envDataChanged := false
    if numArgs == 0 {
        log.Warnln("no environment variable to unset")
        return nil
    } else {
        for i := 0; i < numArgs; i++ {
            providedEnv := args.Get(i)

            log.Info("un-setting %s... ", providedEnv)

            if len(providedEnv) == 1 {
                envDataChanged = unsetEnv(providedEnv)
            } else {
                envDataChanged = unsetEnv(providedEnv)
            }
        }
    }

    if envDataChanged {
        // update wio.env file
        if err := godotenv.Write(envData, envPath); err != nil {
            return err
        }
    }

    return nil
}

// set env variables
func (env Env) setCommand(args cli.Args, numArgs int, envPath string) error {
    envData, err := godotenv.Read(envPath)
    if err != nil {
        return err
    }

    setEnv := func(envName string, envValueProvided string) bool {
        isValueBoolType := func() bool {
            if envValueProvided == ENVBOOLVALUE {
                return true
            } else {
                return false
            }
        }

        if isReadOnly(envName) {
            log.WriteFailure()
            log.Errln("%s is readonly and cannot be changed", envName)
            return false
        } else {
            if envValue, exists := envData[envName]; exists {
                isBoolType := isValueBoolType()

                if (isBoolType && envValue != ENVBOOLVALUE) || (!isBoolType && envValue == ENVBOOLVALUE) {
                    log.WriteFailure()
                    log.Errln("%s is current assigned to a value and there is type mismatch", envName)
                    return false
                }
            }

            envData[envName] = envValueProvided
            log.WriteSuccess()
            return true
        }
    }

    envDataChanged := false
    if numArgs == 0 {
        log.Warnln("no environment variable to set")
        return nil
    } else {
        for i := 0; i < numArgs; i++ {
            providedEnv := args.Get(i)

            // parse args
            envSplit := strings.Split(providedEnv, "=")

            log.Info("setting %s... ", envSplit[0])

            if len(envSplit) == 1 {
                envDataChanged = setEnv(envSplit[0], ENVBOOLVALUE)
            } else {
                envDataChanged = setEnv(envSplit[0], envSplit[1])
            }
        }
    }

    if envDataChanged {
        // update wio.env file
        if err := godotenv.Write(envData, envPath); err != nil {
            return err
        }
    }

    return nil
}

// Display environment
func (env Env) viewCommand() error {
    envPaths := []string{root.GetEnvFilePath()}

    localEnvPath, err := getLocalEnvPath()
    if err == nil && sys.Exists(localEnvPath) {
        envPaths = append(envPaths, localEnvPath)
    }

    envData, err := godotenv.Read(envPaths...)
    if err != nil {
        return err
    }

    readOnlyTextMapper := func(envName string) string {
        readOnlyText := "readonly"
        otherText := "        "

        if isReadOnly(envName) {
            return readOnlyText
        } else {
            return otherText
        }
    }

    printEnvs := func(envName, envValue string) {
        log.Write(log.Green, readOnlyTextMapper(envName)+"  ")

        if envValue == ENVBOOLVALUE {
            log.Write(log.Green, "BOOL  ")
            log.Writeln(log.Cyan, envName)
        } else {
            log.Write(log.Green, "STR   ")
            log.Write(log.Cyan, envName+"=")
            log.Writeln(envValue)
        }
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

func isReadOnly(envName string) bool {
    if envValue, exists := envMeta[envName]; !exists {
        return false
    } else if envValue {
        return true
    } else {
        return false
    }
}

func getLocalEnvPath() (string, error) {
    directory, err := os.Getwd()
    envPath := sys.Path(directory, sys.WioFolder, constants.RootEnv)

    if err != nil {
        return "", err
    }

    if !sys.Exists(sys.Path(directory, sys.Config)) {
        return "", util.Error("not a valid wio project: path is missing wio.yml: %s", directory)
    }

    return envPath, nil
}
