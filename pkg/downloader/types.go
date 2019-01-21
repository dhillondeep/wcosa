package downloader

import (
    "fmt"
    "strings"
    "wio-utils/cmd/io"
    "wio/internal/config/root"
    "wio/pkg/log"
    "wio/pkg/util"
    "wio/pkg/util/sys"
)

type ModuleData struct {
    Name        string `json:"name"`
    Description string `json:"description"`

    Author  interface{} `json:"author"`
    License interface{} `json:"license"`

    Dependencies  map[string]string `json:"dependencies"`
    ToolchainFile string            `json:"toolchain-file"`
}

type Downloader interface {
    DownloadModule(path, url, reference string) (string, error)
}

var SupportedToolchains = map[string]string{
    "cosa":    "wio-framework-avr-cosa",
    "arduino": "wio-framework-avr-arduino",
}

func DownloadToolchain(toolchainLink string) (string, error) {
    // link must be plain without these accessors
    if strings.Contains(toolchainLink, "https://") || strings.Contains(toolchainLink, "http://") {
        return "", util.Error("toolchain link must be without https or http")
    }

    toolchainName, toolchainRef := func() (string, string) {
        split := strings.Split(toolchainLink, ":")

        if len(split) <= 1 {
            return split[0], ""
        } else {
            return split[0], split[1]
        }
    }()

    var d Downloader

    // this is not a valid name
    if strings.Contains(toolchainName, ".") && strings.Contains(toolchainName, "/") {
        d = GitDownloader{}
    } else {
        d = NpmDownloader{}
    }

    log.Writeln(log.Yellow, "-------------------------------")
    log.Write(log.Cyan, "Verifying/Downloading toolchain module from ")
    log.Write(log.Yellow, toolchainLink)
    log.Writeln(log.Cyan, "...")

    defer func() {
        log.Writeln(log.Yellow, "-------------------------------")
    }()

    if path, err := d.DownloadModule(root.GetToolchainPath(), toolchainName, toolchainRef); err != nil {
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

        configData += fmt.Sprintf("set(WIO_DEP_%s_PATH \"%s\")\n", strings.ToUpper(name), depPath)
        configData += fmt.Sprintf("set(WIO_DEP_%s_VERSION \"%s\")\n", strings.ToUpper(name), depPath)
        returnedData, err := createDepAttributes(depPath, "")
        configData += "\n" + returnedData
        if err != nil {
            return configData, err
        }
    }

    writeConfigData := configData
    writeConfigData += fmt.Sprintf("set(WIO_TOOLCHAIN_FOLDER_PATH \"%s\")\n", root.GetToolchainPath())
    writeConfigData += fmt.Sprintf("set(WIO_CURRENT_TOOLCHAIN_PATH \"%s\")\n", file)

    if err := io.NormalIO.WriteFile(wioConfigPath, []byte(writeConfigData)); err != nil {
        return configData, err
    }

    return configData, nil
}
