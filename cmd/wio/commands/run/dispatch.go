package run

import (
    "fmt"
    "strings"
    "wio/cmd/wio/commands/run/cmake"
    "wio/cmd/wio/commands/run/dependencies"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/log"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"
)

type dispatchCmakeFunc func(info *runInfo, target *types.Target) error

var dispatchCmakeFuncPlatform = map[string]dispatchCmakeFunc{
    constants.AVR:    dispatchCmakeAvr,
    constants.NATIVE: dispatchCmakeNative,
}
var dispatchCmakeFuncAvrFramework = map[string]dispatchCmakeFunc{
    constants.COSA: dispatchCmakeAvrGeneric,
    //constants.ARDUINO: dispatchCmakeAvrGeneric,
}

func dispatchCmake(info *runInfo, target *types.Target) error {
    platform := strings.ToLower((*target).GetPlatform())
    if _, exists := dispatchCmakeFuncPlatform[platform]; !exists {
        message := fmt.Sprintf("Platform [%s] is not supported", platform)
        return errors.String(message)
    }
    return dispatchCmakeFuncPlatform[platform](info, target)
}

func dispatchCmakeAvr(info *runInfo, target *types.Target) error {
    framework := strings.ToLower((*target).GetFramework())
    if _, exists := dispatchCmakeFuncAvrFramework[framework]; !exists {
        message := fmt.Sprintf("Framework [%s] not supported", framework)
        return errors.String(message)
    }
    return dispatchCmakeFuncAvrFramework[framework](info, target)
}

func dispatchCmakeNative(info *runInfo, target *types.Target) error {
    return dispatchCmakeNativeGeneric(info, target)
}

func dispatchCmakeAvrGeneric(info *runInfo, target *types.Target) error {
    projectName := info.config.GetMainTag().GetName()
    projectPath := info.directory
    port, err := getPort(info)
    if err != nil && info.runType == TypeRun {
        return err
    }
    return cmake.GenerateAvrCmakeLists(target, projectName, projectPath, port)
}

func dispatchCmakeNativeGeneric(info *runInfo, target *types.Target) error {
    projectName := info.config.GetMainTag().GetName()
    projectPath := info.directory
    return cmake.GenerateNativeCmakeLists(target, projectName, projectPath)
}

func dispatchCmakeDependencies(info *runInfo, target *types.Target) error {
    path := info.directory
    queue := log.NewQueue(16)
    return dependencies.CreateCMakeDependencyTargets(info.config, target, path, queue)
}

func dispatchRunTarget(info *runInfo, target *types.Target) error {
    binDir := binaryPath(info, target)
    platform := (*target).GetPlatform()
    switch platform {
    case constants.AVR:
        var err error = nil
        err = portReconfigure(info, target)
        if err == nil {
            err = uploadTarget(binDir)
        }
        return err
    case constants.NATIVE:
        return runTarget(binDir, "."+io.Sep+(*target).GetName())
    default:
        return errors.Stringf("platform [%s] is not supported", platform)
    }
}

func dispatchCanRunTarget(info *runInfo, target *types.Target) bool {
    binDir := binaryPath(info, target)
    platform := (*target).GetPlatform()
    file := binDir + io.Sep + (*target).GetName() + platformExtension(platform)
    exists, _ := io.Exists(file)
    return exists
}
