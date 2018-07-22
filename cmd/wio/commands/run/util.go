package run

import (
    "wio/cmd/wio/commands/run/cmake"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"
)

func buildPath(info *runInfo) string {
    return cmake.BuildPath(info.directory)
}

func targetPath(info *runInfo, target types.Target) string {
    return io.Path(buildPath(info), target.GetName())
}

func binaryPath(info *runInfo, target types.Target) string {
    return io.Path(targetPath(info, target), constants.BinDir)
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
    case constants.Avr:
        return ".elf"
    case constants.Native:
        return nativeExtension()
    default:
        return ""
    }
}
