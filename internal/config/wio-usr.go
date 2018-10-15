package config

import (
    "os"
    "os/user"
    "wio/internal/cmd/env"
    "wio/internal/constants"
    "wio/pkg/util/sys"
)

type config struct {
    WioUserPath   string
    ToolchainPath string
    VersionsPath  string
    EnvFilePath   string
}

var wioInternalConfig = config{}

func CreateAndSetupWioUsr() error {
    currUser, err := user.Current()
    if err != nil {
        return err
    }

    wioInternalConfig.WioUserPath = sys.Path(currUser.HomeDir, constants.WioUsr)

    // create .wio folder if it does not exist
    if !sys.Exists(wioInternalConfig.WioUserPath) {
        if err := os.Mkdir(wioInternalConfig.WioUserPath, os.ModePerm); err != nil {
            return err
        }
    }

    // create toolchain directory if it does not exist
    wioInternalConfig.ToolchainPath = sys.Path(wioInternalConfig.WioUserPath, constants.ToolchainName)
    if !sys.Exists(wioInternalConfig.ToolchainPath) {
        if err := os.Mkdir(wioInternalConfig.ToolchainPath, os.ModePerm); err != nil {
            return err
        }
    }

    // create versions directory if it does not exist
    wioInternalConfig.VersionsPath = sys.Path(wioInternalConfig.WioUserPath, constants.VersionsName)
    if !sys.Exists(wioInternalConfig.VersionsPath) {
        if err := os.Mkdir(wioInternalConfig.VersionsPath, os.ModePerm); err != nil {
            return err
        }
    }

    // create environment file if it does not exist
    wioInternalConfig.EnvFilePath = sys.Path(wioInternalConfig.WioUserPath, constants.EnvFileName)
    if !sys.Exists(wioInternalConfig.EnvFilePath) {
        if err := env.CreateEnv(wioInternalConfig.EnvFilePath); err != nil {
            return err
        }
    }

    // load the environment path
    if err := env.LoadEnv(wioInternalConfig.EnvFilePath); err != nil {
        return err
    }

    return nil
}

func GetToolchainPath() string {
    return wioInternalConfig.ToolchainPath
}

func GetVersionsPath() string {
    return wioInternalConfig.VersionsPath
}
