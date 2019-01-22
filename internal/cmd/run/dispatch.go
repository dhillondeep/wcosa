package run

import (
    "fmt"
    "os"
    "strings"
    "wio/internal/cmd/run/cmake"
    "wio/internal/cmd/run/dependencies"
    "wio/internal/constants"
    "wio/internal/types"
    "wio/pkg/downloader"
    "wio/pkg/log"
    "wio/pkg/util"
    "wio/pkg/util/sys"
)

type dispatchCmakeFunc func(info *runInfo, target types.Target, toolchainPath string) error

var dispatchCmakeFuncPlatform = map[string]dispatchCmakeFunc{
    constants.Avr:    dispatchCmakeAvr,
    constants.Native: dispatchCmakeNative,
}

func dispatchCmake(info *runInfo, target types.Target) error {
    platform := strings.ToLower(target.GetPlatform())

    // this means platform was not specified at all
    if util.IsEmptyString(platform) {
        message := fmt.Sprintf("No Platform specified by the [%s] target", target.GetName())
        return util.Error(message)
    }

    if _, exists := dispatchCmakeFuncPlatform[platform]; !exists {
        message := fmt.Sprintf("Platform [%s] is not supported", platform)
        return util.Error(message)
    }

    var toolchainPath string
    var err error
    if platform != constants.Native {
        defer func() {
            log.Writeln(log.Yellow, "-------------------------------")
        }()

        log.Writeln(log.Yellow, "-------------------------------")

        if toolchainPath, err = downloader.DownloadToolchain(target.GetFramework(), info.retool); err != nil {
            return err
        }
    }

    return dispatchCmakeFuncPlatform[platform](info, target, toolchainPath)
}

func dispatchCmakeAvr(info *runInfo, target types.Target, toolchainPath string) error {
    board := strings.Trim(strings.ToUpper(target.GetBoard()), " ")

    // this means board was not specified at all
    if board == "" {
        message := fmt.Sprintf("No Board specified by the [%s] target", target.GetName())
        return util.Error(message)
    }

    projectName := info.config.GetName()
    projectPath := info.directory
    port, err := getPort(info)
    if err != nil && info.runType == TypeRun {
        return err
    }
    cppStandard, cStandard, err := cmake.GetStandard(info.config.GetInfo().GetOptions().GetStandard())
    if err != nil {
        return err
    }

    return cmake.GenerateAvrCmakeLists(toolchainPath, target, projectName, projectPath, cppStandard, cStandard, port)
}

func dispatchCmakeNative(info *runInfo, target types.Target, _ string) error {
    projectName := info.config.GetName()
    projectPath := info.directory

    cppStandard, cStandard, err := cmake.GetStandard(info.config.GetInfo().GetOptions().GetStandard())
    if err != nil {
        return err
    }

    return cmake.GenerateNativeCmakeLists(target, projectName, projectPath, cppStandard, cStandard)
}

func dispatchCmakeDependencies(info *runInfo, target types.Target) error {
    cmakePath := sys.Path(cmake.BuildPath(info.directory), target.GetName())
    if err := os.MkdirAll(cmakePath, os.ModePerm); err != nil {
        return err
    }
    cmakePath = sys.Path(cmakePath, "dependencies.cmake")

    buildTargets, libraryTargets, err := dependencies.CreateBuildTargets(info.directory, target)
    if err != nil {
        return err
    } else {
        err := dependencies.GenerateCMakeDependencies(cmakePath, target.GetPlatform(), buildTargets, libraryTargets)
        if err != nil {
            return err
        }
    }

    log.Verbln()
    return err
}

func dispatchRunTarget(info *runInfo, target types.Target) error {
    binDir := binaryPath(info, target)
    platform := target.GetPlatform()
    switch platform {
    case constants.Avr:
        var err error = nil
        err = portReconfigure(info, target)
        if err == nil {
            err = uploadTarget(binDir)
        }
        return err
    case constants.Native:
        args := info.context.String("args")
        return runTarget(info.directory, sys.Path(binDir, target.GetName()), args)
    default:
        return util.Error("Platform [%s] is not supported", platform)
    }
}

func dispatchCanRunTarget(info *runInfo, target types.Target) bool {
    binDir := binaryPath(info, target)
    platform := target.GetPlatform()
    file := sys.Path(binDir, target.GetName(), platformExtension(platform))
    return sys.Exists(file)
}
