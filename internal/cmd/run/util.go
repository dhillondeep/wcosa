package run

import (
    "wio/internal/cmd/run/cmake"
    "wio/internal/constants"
    "wio/internal/types"
    "wio/pkg/util/sys"
)

func buildPath(info *runInfo) string {
    return cmake.BuildPath(info.directory)
}

func targetPath(info *runInfo, target types.Target) string {
    return sys.Path(buildPath(info), target.GetName())
}

func binaryPath(info *runInfo, target types.Target) string {
    return sys.Path(targetPath(info, target), constants.BinDir)
}

func nativeExtension() string {
    switch sys.GetOS() {
    case sys.WINDOWS:
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
