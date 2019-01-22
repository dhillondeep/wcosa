package downloader

import (
    "fmt"
    "strings"
    "wio-utils/cmd/io"
    "wio/internal/config/root"
    "wio/pkg/util"
    "wio/pkg/util/sys"
    "wio/pkg/util/template"
)

var SupportedToolchains = map[string]string{
    "cosa":    "wio-framework-avr-cosa",
    "arduino": "wio-framework-avr-arduino",
}

const WioConfigSet = "set({{VAR_NAME}} \"{{VAR_VALUE}}\")"

type ModuleData struct {
    Name        string `json:"name"`
    Description string `json:"description"`

    Author  interface{} `json:"author"`
    License interface{} `json:"license"`

    Dependencies  map[string]string `json:"dependencies"`
    ToolchainFile string            `json:"toolchain-file"`
}

type Downloader interface {
    DownloadModule(path, url, reference string, retool bool) (string, error)
}

func DownloadToolchain(toolchainLink string, retool bool) (string, error) {
    // link must be plain without these accessors
    if strings.Contains(toolchainLink, "https://") || strings.Contains(toolchainLink, "http://") {
        return "", util.Error("toolchain link provided must be without https or http, ex: github.com/foo")
    }

    toolchainName, toolchainRef := func() (string, string) {
        split := strings.Split(toolchainLink, ":")

        if len(split) <= 1 {
            return split[0], ""
        } else {
            return split[0], split[1]
        }
    }()

    // use alias
    toolchainName = func() string {
        if val, exists := SupportedToolchains[toolchainName]; exists {
            return val
        } else {
            return toolchainName
        }
    }()

    var d Downloader

    // this is not a valid name so it must be a url
    if strings.Contains(toolchainName, ".") && strings.Contains(toolchainName, "/") {
        d = GitDownloader{}
    } else {
        d = NpmDownloader{}
    }

    if path, err := d.DownloadModule(root.GetToolchainPath(), toolchainName, toolchainRef, retool); err != nil {
        return "", err
    } else {
        _, err := createDepAttributes(path, "")
        if err != nil {
            return "", err
        }

        return path, nil
    }
}

func createDepAttributes(file string, configData string) (string, error) {
    wioConfigPath := sys.Path(file, "WioConfig.cmake")

    moduleData := &ModuleData{}

    if err := io.NormalIO.ParseJson(sys.Path(file, "package.json"), moduleData); err != nil {
        return "", err
    }

    for name, version := range moduleData.Dependencies {
        depPath := sys.Path(root.GetToolchainPath(), name+"__"+version)

        configData += template.Replace(WioConfigSet, map[string]string{
            "VAR_NAME":  fmt.Sprintf("WIO_DEP_%s_PATH", strings.ToUpper(name)),
            "VAR_VALUE": depPath,
        }) + "\n"
        configData += template.Replace(WioConfigSet, map[string]string{
            "VAR_NAME":  fmt.Sprintf("WIO_DEP_%s_VERSION", strings.ToUpper(name)),
            "VAR_VALUE": version,
        }) + "\n"

        returnedData, err := createDepAttributes(depPath, "")
        configData += returnedData
        if err != nil {
            return configData, err
        }
    }

    writeConfigData := configData
    writeConfigData += template.Replace(WioConfigSet, map[string]string{
        "VAR_NAME":  "WIO_TOOLCHAIN_PATH",
        "VAR_VALUE": root.GetToolchainPath(),
    }) + "\n"
    writeConfigData += template.Replace(WioConfigSet, map[string]string{
        "VAR_NAME":  "WIO_CURRENT_TOOLCHAIN_PATH",
        "VAR_VALUE": file,
    }) + "\n"

    if err := io.NormalIO.WriteFile(wioConfigPath, []byte(writeConfigData)); err != nil {
        return configData, err
    }

    return configData, nil
}
