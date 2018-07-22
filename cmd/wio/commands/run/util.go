package run

import (
    "io/ioutil"
    "os"
    "path/filepath"
    "wio/cmd/wio/commands/run/cmake"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"
)

var sourceExtensions = map[string]bool{
    ".cpp": true,
    ".c":   true,
    ".cc":  true,
}

// perform a check that at least one executable file is in target directory
func sourceFilesExist(directory string) (bool, error) {
    files, err := ioutil.ReadDir(directory)
    if err != nil {
        return false, err
    } else {
        sourceFileExists := false
        for _, f := range files {
            if _, exists := sourceExtensions[filepath.Ext(io.Path(directory, f.Name()))]; exists {
                sourceFileExists = true
                break
            }
        }

        return sourceFileExists, nil
    }
}

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
