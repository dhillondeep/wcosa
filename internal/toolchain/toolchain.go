package toolchain

import (
    "path/filepath"
    "wio/pkg/util/sys"
)

// This returns the path to toolchain directory
func GetToolchainPath() (string, error) {
    executablePath, err := sys.NormalIO.GetRoot()
    if err != nil {
        return "", err
    }

    toolchainPath := executablePath + sys.Sep + "toolchain"

    if !sys.Exists(toolchainPath) {
        toolchainPath, err = filepath.Abs(executablePath + sys.Sep + ".." + sys.Sep + "toolchain")
        if err != nil {
            return "", err
        }
    }

    return toolchainPath, nil
}
