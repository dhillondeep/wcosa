package config

import (
    "os"
    "os/user"
    "wio/internal/cmd/env"
    "wio/internal/constants"
    "wio/pkg/util/sys"
)

type config struct {
    WioUserPath    string
    FrameworksPath string
    EnvFilePath    string
}

var WioInternalConfig = config{}

func CreateAndSetupWioUsr() error {
    currUser, err := user.Current()
    if err != nil {
        return err
    }

    WioInternalConfig.WioUserPath = sys.Path(currUser.HomeDir, constants.WioUsr)

    // create .wio folder if it does not exist
    if !sys.Exists(WioInternalConfig.WioUserPath) {
        if err := os.Mkdir(WioInternalConfig.WioUserPath, os.ModePerm); err != nil {
            return err
        }
    }

    // create frameworks directory if it does not exist
    WioInternalConfig.FrameworksPath = sys.Path(WioInternalConfig.WioUserPath, constants.FrameworkName)
    if !sys.Exists(WioInternalConfig.FrameworksPath) {
        if err := os.Mkdir(WioInternalConfig.FrameworksPath, os.ModePerm); err != nil {
            return err
        }
    }

    // create environment file if it does not exist
    WioInternalConfig.EnvFilePath = sys.Path(WioInternalConfig.WioUserPath, constants.EnvFileName)
    if !sys.Exists(WioInternalConfig.EnvFilePath) {
        if err := env.CreateEnv(WioInternalConfig.EnvFilePath); err != nil {
            return err
        }
    }

    // load the environment path
    if err := env.LoadEnv(WioInternalConfig.EnvFilePath); err != nil {
        return err
    }

    return nil
}

func GetFrameworksPath() string {
    return WioInternalConfig.FrameworksPath
}
