package generate

import (
    "fmt"
    "os"
    "strings"
    "wio/internal/cmd/run/cmake"
    "wio/internal/cmd/run/dependencies"
    "wio/internal/constants"
    "wio/internal/env"
    "wio/internal/types"
    "wio/internal/utils"
    "wio/pkg/downloader"
    "wio/pkg/log"
    "wio/pkg/util"
    "wio/pkg/util/sys"
    "wio/pkg/util/template"
)

var supportedPlatforms = map[string]bool{
    constants.Avr:    true,
    constants.Native: true,
}

type InfoGenerate struct {
    Config types.Config

    Directory   string
    ProjectType string
    Port        string

    Retool      bool
    NoDepCreate bool
}

func HardwareFile(info *InfoGenerate, target types.Target) error {
    targetPath := sys.Path(utils.BuildPath(info.Directory), target.GetName())

    templateFilePath := sys.Path("templates", "cmake", "hardware.cmake.tpl")
    hardwareFilePath := sys.Path(targetPath, "hardware.cmake")
    if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
        return err
    }

    if err := sys.AssetIO.CopyFile(templateFilePath, hardwareFilePath, true); err != nil {
        return err
    }
    hardware := ""
    switch target.GetPlatform() {
    case constants.Avr:
        hardware = target.GetBoard()
        break
    case constants.Native:
        hardware = env.GetOS()
    }

    return template.IOReplace(hardwareFilePath, map[string]string{
        "UPLOAD_PORT":     info.Port,
        "TARGET_HARDWARE": hardware,
    })
}

func CMakeListsFile(info *InfoGenerate, target types.Target) error {
    platform := strings.ToLower(target.GetPlatform())

    // this means platform was not specified at all
    if util.IsEmptyString(platform) {
        message := fmt.Sprintf("No Platform specified by the [%s] target", target.GetName())
        return util.Error(message)
    }

    if _, exists := supportedPlatforms[platform]; !exists {
        message := fmt.Sprintf("Platform [%s] is not supported", platform)
        return util.Error(message)
    }

    var toolchainPath string
    var err error
    if platform != constants.Native {
        if toolchainPath, err = downloader.DownloadToolchain(target.GetFramework(), info.Retool); err != nil {
            return err
        }
    }

    projectName := info.Config.GetName()
    projectPath := info.Directory

    cppStandard, cStandard, err := cmake.GetStandard(info.Config.GetInfo().GetOptions().GetStandard())
    if err != nil {
        return err
    }

    return cmake.GenerateCmakeLists(toolchainPath, target, projectName, projectPath, cppStandard, cStandard)
}

func DependenciesFile(info *InfoGenerate, target types.Target) error {
    cmakePath := sys.Path(utils.BuildPath(info.Directory), target.GetName())
    if err := os.MkdirAll(cmakePath, os.ModePerm); err != nil {
        return err
    }

    cmakePath = sys.Path(cmakePath, "dependencies.cmake")

    if !sys.Exists(cmakePath) {
        if _, err := os.Create(cmakePath); err != nil {
            return err
        }
    }

    if !info.NoDepCreate {
        buildTargets, libraryTargets, err := dependencies.CreateBuildTargets(info.Directory, target)
        if err != nil {
            return err
        } else {
            err := dependencies.GenerateCMakeDependencies(cmakePath, target.GetPlatform(), buildTargets, libraryTargets)
            if err != nil {
                return err
            }
        }

        log.Verbln()
    }

    return nil
}
