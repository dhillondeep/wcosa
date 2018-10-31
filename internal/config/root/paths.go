package root

type configPaths struct {
    WioUserPath   string
    ToolchainPath string
    UpdatePath    string
    EnvFilePath   string
}

var wioInternalConfigPaths = configPaths{}
