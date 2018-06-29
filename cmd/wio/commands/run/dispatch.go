package run

import (
    "wio/cmd/wio/types"
    "wio/cmd/wio/errors"
    "fmt"
    "wio/cmd/wio/commands/run/cmake"
    "strings"
    "wio/cmd/wio/log"
    "wio/cmd/wio/commands/run/dependencies"
)

type cmakeFunc func(info *runInfo, target *types.Target) error

var dispatchCmakeFuncPlatform = map[string]cmakeFunc{
    "avr":    dispatchCmakeAvr,
    "native": dispatchCmakeNative,
}
var dispatchCmakeFuncAvrFramework = map[string]cmakeFunc{
    "cosa": dispatchCmakeAvrGeneric,
    //"arduino": dispatchCmakeAvrGeneric,
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
    if err != nil && info.rtype == type_upload {
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
