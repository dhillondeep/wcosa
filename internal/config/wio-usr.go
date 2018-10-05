package config

import (
    "os"
    "os/user"
    "wio/internal/cmd/env"
    "wio/internal/constants"
    "wio/pkg/util/sys"
)

func CreateAndSetupWioUsr() error {
    currUser, err := user.Current()
    if err != nil {
        return err
    }

    wioUserPath := sys.Path(currUser.HomeDir, constants.WioUsr)

    // create .wio folder if it does not exist
    if !sys.Exists(wioUserPath) {
        if err := os.Mkdir(wioUserPath, os.ModePerm); err != nil {
            return err
        }
    }

    // create frameworks directory
    wioUserFramework := sys.Path(wioUserPath, constants.FrameworkName)
    if !sys.Exists(wioUserFramework) {
        if err := os.Mkdir(wioUserFramework, os.ModePerm); err != nil {
            return err
        }
    }

    envFilePath := sys.Path(wioUserPath, constants.EnvFileName)
    if !sys.Exists(envFilePath) {
        if err := env.CreateEnv(wioUserPath); err != nil {
            return err
        }
    }

    if err := env.LoadEnv(envFilePath); err != nil {
        return err
    }

    return nil
}
