package root

type configPaths struct {
    WioUserPath   string
    ToolchainPath string
    UpdatePath    string
    EnvFilePath   string
}

var wioInternalConfigPaths = configPaths{}

func GetWioUserPath() string {
    return wioInternalConfigPaths.WioUserPath
}

func GetToolchainPath() string {
    return wioInternalConfigPaths.ToolchainPath
}

func GetUpdatePath() string {
    return wioInternalConfigPaths.UpdatePath
}

func GetEnvFilePath() string {
    return wioInternalConfigPaths.EnvFilePath
}
