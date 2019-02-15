package root

type configPaths struct {
	WioUserPath    string
	ToolchainPath  string
	SecurityPath   string
	UpdatePath     string
	EnvFilePath    string
	ConfigFilePath string
}

var wioInternalConfigPaths = configPaths{}

func GetWioUserPath() string {
	return wioInternalConfigPaths.WioUserPath
}

func GetToolchainPath() string {
	return wioInternalConfigPaths.ToolchainPath
}

func GetSecurityPath() string {
	return wioInternalConfigPaths.SecurityPath
}

func GetUpdatePath() string {
	return wioInternalConfigPaths.UpdatePath
}

func GetEnvFilePath() string {
	return wioInternalConfigPaths.EnvFilePath
}

func GetConfigFilePath() string {
	return wioInternalConfigPaths.ConfigFilePath
}
