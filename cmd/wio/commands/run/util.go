package run

import (
    "os"
    "wio/cmd/wio/commands/run/cmake"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"
)

func buildPath(info *runInfo) string {
    return cmake.BuildPath(info.directory)
}

func targetPath(info *runInfo, target types.Target) string {
    return buildPath(info) + io.Sep + target.GetName()
}

func binaryPath(info *runInfo, target types.Target) string {
    return targetPath(info, target) + io.Sep + constants.BinDir
}

func readDirectory(args []string) (string, error) {
    if len(args) <= 0 {
        return os.Getwd()
    }
    return args[0], nil
}

func nativeExtension() string {
    switch io.GetOS() {
    case io.WINDOWS:
        return ".exe"
    default:
        return ""
    }
}

func platformExtension(platform string) string {
    switch platform {
    case constants.AVR:
        return ".elf"
    case constants.NATIVE:
        return nativeExtension()
    default:
        return ""
    }
}
