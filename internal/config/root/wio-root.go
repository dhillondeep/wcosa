package root

import (
    "os"
    "os/user"
    "wio/internal/constants"
    "wio/pkg/util/sys"

    "github.com/joho/godotenv"
)

type WioRootConfig struct {
    Updated bool `json:"updated"`
}

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
    wioInternalConfigPaths.ToolchainPath = sys.Path(GetWioUserPath(), constants.RootToolchain)
    if !sys.Exists(wioInternalConfigPaths.ToolchainPath) {
        if err := os.Mkdir(wioInternalConfigPaths.ToolchainPath, os.ModePerm); err != nil {
            return err
        }
    }

    // create security directory if it does not exist
    wioInternalConfigPaths.SecurityPath = sys.Path(GetWioUserPath(), constants.Security)
    if !sys.Exists(wioInternalConfigPaths.SecurityPath) {
        if err := os.Mkdir(wioInternalConfigPaths.SecurityPath, os.ModePerm); err != nil {
            return err
        }
    }

    // create update directory if it does not exist
    wioInternalConfigPaths.UpdatePath = sys.Path(GetWioUserPath(), constants.RootUpdate)
    if !sys.Exists(wioInternalConfigPaths.UpdatePath) {
        if err := os.Mkdir(wioInternalConfigPaths.UpdatePath, os.ModePerm); err != nil {
            return err
        }
    }

    // create config file if doesn't exist
    wioInternalConfigPaths.ConfigFilePath = sys.Path(GetWioUserPath(), constants.RootConfig)
    config, err := CreateConfig()
    if err != nil {
        return err
    }

    // create environment file if it does not exist
    wioInternalConfigPaths.EnvFilePath = sys.Path(GetWioUserPath(), constants.RootEnv)
    if err := CreateEnv(config); err != nil {
        return err
    }

    return nil
}

// Create config file
func CreateConfig() (*WioRootConfig, error) {
    var config *WioRootConfig
    if !sys.Exists(wioInternalConfigPaths.ConfigFilePath) {
        config = &WioRootConfig{
            Updated: false,
        }
        if err := sys.NormalIO.WriteJson(wioInternalConfigPaths.ConfigFilePath, config); err != nil {
            return nil, err
        }
    } else {
        config = &WioRootConfig{}
        if err := sys.NormalIO.ParseJson(wioInternalConfigPaths.ConfigFilePath, config); err != nil {
            return nil, err
        }
    }

    return config, nil
}

// Creates environment and overrides if there is an old environment
func CreateEnv(config *WioRootConfig) error {
    var wioRoot string
    var wioPath string
    var err error

    readValues := func() error {
        wioRoot, err = sys.NormalIO.GetRoot()
        if err != nil {
            return err
        }

        wioPath, err = os.Executable()
        if err != nil {
            return err
        }

        return nil
    }

    if !sys.Exists(wioInternalConfigPaths.EnvFilePath) {
        if err := readValues(); err != nil {
            return err
        }

        envs := map[string]string{
            "WIOROOT": wioRoot,
            "WIOOS":   sys.GetOS(),
            "WIOARCH": sys.GetArch(),
            "WIOPATH": wioPath,
        }

        return godotenv.Write(envs, wioInternalConfigPaths.EnvFilePath)
    } else {
        if err := readValues(); err != nil {
            return err
        }

        if config.Updated {
            envsRead, err := godotenv.Read(wioInternalConfigPaths.EnvFilePath)
            if err != nil {
                return err
            }

            envsRead["WIOROOT"] = wioRoot
            envsRead["WIOOS"] = sys.GetOS()
            envsRead["WIOARCH"] = sys.GetArch()
            envsRead["WIOPATH"] = wioPath

            return godotenv.Write(envsRead, wioInternalConfigPaths.EnvFilePath)
        }
    }

    return nil
}

// Creates environment and overrides if there is an old environment
func CreateLocalEnv(path string) error {
    envs := map[string]string{}

    // create wio.env file
    if err := godotenv.Write(envs, path); err != nil {
        return err
    }

    return nil
}

// Loads environment
func LoadEnv() error {
    return godotenv.Load(GetEnvFilePath())
}
