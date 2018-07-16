package run

import (
    "fmt"
    "regexp"
    "strings"
    "wio/cmd/wio/commands/run/cmake"
    "wio/cmd/wio/commands/run/dependencies"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/log"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"

    "github.com/thoas/go-funk"
)

type dispatchCmakeFunc func(info *runInfo, target *types.Target) error

var dispatchCmakeFuncPlatform = map[string]dispatchCmakeFunc{
    constants.AVR:    dispatchCmakeAvr,
    constants.NATIVE: dispatchCmakeNative,
}
var dispatchCmakeFuncAvrFramework = map[string]dispatchCmakeFunc{
    constants.COSA:    dispatchCmakeAvrCosa,
    constants.ARDUINO: dispatchCmakeAvrArduino,
}

func dispatchCmake(info *runInfo, target *types.Target) error {
    platform := strings.ToLower((*target).GetPlatform())

    // this means platform was not specified at all
    if strings.Trim(platform, " ") == "" {
        message := fmt.Sprintf("No Platform specified by the [%s] target", (*target).GetName())
        return errors.String(message)
    }

    if _, exists := dispatchCmakeFuncPlatform[platform]; !exists {
        message := fmt.Sprintf("Platform [%s] is not supported", platform)
        return errors.String(message)
    }
    return dispatchCmakeFuncPlatform[platform](info, target)
}

func dispatchCmakeAvr(info *runInfo, target *types.Target) error {
    framework := strings.Trim(strings.ToLower((*target).GetFramework()), " ")
    board := strings.Trim(strings.ToUpper((*target).GetBoard()), " ")

    // this means framework was not specified at all
    if framework == "" {
        message := fmt.Sprintf("No Framework specified by the [%s] target. Try one of %s",
            (*target).GetName(), funk.Keys(dispatchCmakeFuncAvrFramework))
        return errors.String(message)
    }

    // this means board was not specified at all
    if board == "" {
        message := fmt.Sprintf("No Board specified by the [%s] target", (*target).GetName())
        return errors.String(message)
    }

    if _, exists := dispatchCmakeFuncAvrFramework[framework]; !exists {
        message := fmt.Sprintf("Framework [%s] not supported", framework)
        return errors.String(message)
    }
    return dispatchCmakeFuncAvrFramework[framework](info, target)
}

func dispatchCmakeNative(info *runInfo, target *types.Target) error {
    // perform a check for frameworks
    frameworkProvided := strings.Trim(strings.ToLower((*target).GetFramework()), " ")

    if frameworkProvided == "" {
        return errors.Stringf("Framework not provided. You can try c++11, c, or c++11,c11, etc")
    }

    return dispatchCmakeNativeGeneric(info, target)
}

func dispatchCmakeAvrCosa(info *runInfo, target *types.Target) error {
    projectName := info.config.GetMainTag().GetName()
    projectPath := info.directory
    port, err := getPort(info)
    if err != nil && info.runType == TypeRun {
        return err
    }
    return cmake.GenerateAvrCmakeLists("toolchain/cmake/CosaToolchain.cmake", target,
        projectName, projectPath, port)
}

func dispatchCmakeAvrArduino(info *runInfo, target *types.Target) error {
    projectName := info.config.GetMainTag().GetName()
    projectPath := info.directory
    port, err := getPort(info)
    if err != nil && info.runType == TypeRun {
        return err
    }
    return cmake.GenerateAvrCmakeLists("toolchain/cmake/ArduinoToolchain.cmake", target,
        projectName, projectPath, port)
}

func dispatchCmakeNativeGeneric(info *runInfo, target *types.Target) error {
    projectName := info.config.GetMainTag().GetName()
    projectPath := info.directory

    frameworks := strings.Split((*target).GetFramework(), ",")

    cDetectionPat := regexp.MustCompile(`c[0-9][0-9]`)
    cppDetectionPat := regexp.MustCompile(`c\+\+[0-9][0-9]`)

    var cppStandard string
    var cStandard string

    for _, framework := range frameworks {
        if cppDetectionPat.MatchString(framework) {
            cppStandard = strings.Split(framework, "++")[1]
            framework = strings.Replace(framework, "++", "PP", 1) // use CPP instead of C++
        } else if framework == "c" {
            cStandard = "11"
        } else if cDetectionPat.MatchString(framework) {
            cStandard = strings.Split(framework, "c")[1]
        } else {
            return errors.Stringf("Framework [%s] is not valid. You can try c++11, c, or c++11,c11", framework)
        }
    }

    return cmake.GenerateNativeCmakeLists(target, strings.Join(frameworks, ""), cppStandard, cStandard,
        projectName, projectPath)
}

func dispatchCmakeDependencies(info *runInfo, target *types.Target) error {
    path := info.directory
    queue := log.NewQueue(16)
    err := dependencies.CreateCMakeDependencyTargets(info.config, target, path, queue)
    log.Verbln()
    log.PrintQueue(queue, log.TWO_SPACES)
    return err
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
        return errors.Stringf("Platform [%s] is not supported", platform)
    }
}

func dispatchCanRunTarget(info *runInfo, target *types.Target) bool {
    binDir := binaryPath(info, target)
    platform := (*target).GetPlatform()
    file := binDir + io.Sep + (*target).GetName() + platformExtension(platform)
    return io.Exists(file)
}
