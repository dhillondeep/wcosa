package toolchain

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "wio/internal/config"
    "wio/pkg/log"
    "wio/pkg/util"
    "wio/pkg/util/sys"
)

func GetFrameworkAsset(platform, framework, version string) (*assetRelease, error) {
    repoName := fmt.Sprintf("framework-%s-%s", platform, framework)

    return getRemoteAsset(repoName, framework, version)
}

func DownloadFramework(platform, framework string, frameworkAsset *assetRelease) error {
    log.Write(log.Cyan, "Downloading ")
    log.Write(log.Green, framework)
    log.Write(log.Cyan, " framework for ")
    log.Write(log.Green, platform)
    log.Write(log.Cyan, " platform...")
    frameworkPath := sys.Path(config.GetToolchainPath(), platform, frameworkAsset.LocalName)

    if sys.Exists(frameworkPath) {
        log.Writeln(log.Green, " already exists")
        return nil
    } else {
        log.Writeln()
    }

    if err := os.MkdirAll(frameworkPath, os.ModePerm); err != nil {
        return err
    }

    frameworkTarPath := sys.Path(frameworkPath, frameworkAsset.LocalName+".tar.gz")

    if err := downloadAsset(frameworkAsset, frameworkTarPath); err != nil {
        return err
    }

    err := extractTarball(frameworkTarPath, frameworkPath,
        fmt.Sprintf("Extracting %s ", frameworkAsset.LocalName+".tar.gz"))
    if err != nil {
        log.WriteFailure()
        return err
    }

    if err := os.Remove(frameworkTarPath); err != nil {
        log.WriteFailure()
        return err
    }

    log.WriteSuccess()

    // download requirements
    frameworkConfig := &FrameworkConfig{}

    // open framework.json file
    if err := sys.NormalIO.ParseJson(sys.Path(frameworkPath, FrameworkConfigName), frameworkConfig); err != nil {
        return err
    }

    for _, requirement := range frameworkConfig.Requirements {
        log.Write(log.Cyan, "Downloading requirement %s...", requirement)

        decodeName := strings.Split(requirement, "@")
        reqNameProvided := decodeName[0]
        var reqVersion string

        if len(decodeName) > 1 {
            reqVersion = decodeName[1]
        }

        reqLocalName := fmt.Sprintf("%s@%s", reqNameProvided, reqVersion)
        reqLocalPath := sys.Path(config.GetToolchainPath(), reqLocalName)
        reqTarballPath := sys.Path(reqLocalPath + ".tar.gz")

        if sys.Exists(reqLocalPath) {
            log.Writeln(log.Green, " already exists")
            return nil
        } else {
            log.Writeln()
        }

        reqAsset, err := getRemoteAsset(reqNameProvided, reqLocalName, reqVersion)
        if err != nil {
            return err
        }
        if err := downloadAsset(reqAsset, reqTarballPath); err != nil {
            return err
        }

        err = extractTarball(reqTarballPath, reqLocalPath,
            fmt.Sprintf("Extracting %s ", reqLocalName+".tar.gz"))
        if err != nil {
            log.WriteFailure()
            return err
        }

        if err := os.Remove(reqTarballPath); err != nil {
            log.WriteFailure()
            return err
        }

        log.WriteSuccess()
    }

    return nil
}

func GetToolchainPath(platform, framework string) (string, error) {
    var frameworkPath string

    frameworkDecode := strings.Split(framework, "@")
    frameworkName := frameworkDecode[0]
    var frameworkVersion string
    if len(frameworkDecode) > 1 {
        frameworkVersion = frameworkDecode[1]
    }

    if util.IsEmptyString(frameworkVersion) {
        paths, err := filepath.Glob(sys.Path(config.GetToolchainPath(), platform, frameworkName+"@*"))
        if err != nil {
            return "", err
        }

        if len(paths) <= 0 {
            return "", util.Error("toolchain not found")
        }

        frameworkPath = paths[0]
    } else {
        frameworkPath = sys.Path(config.GetToolchainPath(), platform, frameworkName+"@"+frameworkVersion)
    }

    frameworkConfig := &FrameworkConfig{}

    // open framework.json file
    sys.NormalIO.ParseJson(sys.Path(frameworkPath, FrameworkConfigName), frameworkConfig)

    return sys.Path(frameworkPath, frameworkConfig.ToolchainFile), nil
}
