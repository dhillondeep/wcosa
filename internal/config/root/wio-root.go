package root

import (
    "os"
    "os/user"
    "wio/internal/constants"
    "wio/pkg/util/sys"

    "github.com/joho/godotenv"
)

func CreateWioRoot() error {
    currUser, err := user.Current()
    if err != nil {
        return err
    }

    wioInternalConfigPaths.WioUserPath = sys.Path(currUser.HomeDir, constants.WioRoot)

    // create root folder if it does not exist
    if !sys.Exists(wioInternalConfigPaths.WioUserPath) {
        if err := os.Mkdir(wioInternalConfigPaths.WioUserPath, os.ModePerm); err != nil {
            return err
        }
    }

    // create toolchain directory if it does not exist
    wioInternalConfigPaths.ToolchainPath = sys.Path(wioInternalConfigPaths.WioUserPath, constants.RootToolchain)
    if !sys.Exists(wioInternalConfigPaths.ToolchainPath) {
        if err := os.Mkdir(wioInternalConfigPaths.ToolchainPath, os.ModePerm); err != nil {
            return err
        }
    }

    // create update directory if it does not exist
    wioInternalConfigPaths.UpdatePath = sys.Path(wioInternalConfigPaths.WioUserPath, constants.RootUpdate)
    if !sys.Exists(wioInternalConfigPaths.UpdatePath) {
        if err := os.Mkdir(wioInternalConfigPaths.UpdatePath, os.ModePerm); err != nil {
            return err
        }
    }

    // create environment file if it does not exist
    wioInternalConfigPaths.EnvFilePath = sys.Path(wioInternalConfigPaths.WioUserPath, constants.RootEnv)
    if !sys.Exists(wioInternalConfigPaths.EnvFilePath) {
        if err := CreateEnv(); err != nil {
            return err
        }
    }

    return nil
}

// Creates environment and overrides if there is an old environment
func CreateEnv() error {
    wioRoot, err := sys.NormalIO.GetRoot()
    if err != nil {
        return err
    }

    envs := map[string]string{
        "WIOROOT": wioRoot,
        "WIOOS":   sys.GetOS(),
    }

    // create wio.env file
    if err := godotenv.Write(envs, GetEnvFilePath()); err != nil {
        return err
    }

    return nil
}

// Loads environment
func LoadEnv() error {
    return godotenv.Load(GetEnvFilePath())
}
