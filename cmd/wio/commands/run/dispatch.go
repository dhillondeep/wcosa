package run

import (
    "wio/cmd/wio/types"
    "wio/cmd/wio/errors"
    "fmt"
    "wio/cmd/wio/commands/run/cmake"
    "strings"
)

type cmakeFunc func(info *runInfo, target *types.Target) error
var dispatchCmakeFuncPlatform = map[string]cmakeFunc{
    "avr": dispatchCmakeAvr,
    "native": dispatchCmakeNative,
}
var dispatchCmakeFuncFramework = map[string]cmakeFunc{
    "cosa": dispatchCmakeAvrGeneric,
    "arduino": dispatchCmakeAvrGeneric,
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

}

func dispatchCmakeNative(info *runInfo, target *types.Target) error {

}

func dispatchCmakeAvrGeneric(info *runInfo, target *types.Target) error {

}
