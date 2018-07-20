package toolchain

import (
    "errors"
    "os/exec"
    "path/filepath"
    "wio/cmd/wio/utils/io"
)

const (
    serialLinux   = "serial/serial-ports-linux"
    serialDarwin  = "serial/serial-ports-mac"
    serialWindows = "serial/serial-ports.exe"
)

var operatingSystem = io.GetOS()

// This returns the path to toolchain directory
func GetToolchainPath() (string, error) {
    executablePath, err := io.NormalIO.GetRoot()
    if err != nil {
        return "", err
    }

    toolchainPath := executablePath + io.Sep + "toolchain"

    if !io.Exists(toolchainPath) {
        toolchainPath, err = filepath.Abs(executablePath + io.Sep + ".." + io.Sep + "toolchain")
        if err != nil {
            return "", err
        }
    }

    return toolchainPath, nil
}

// This is the command to execute PySerial to get ports information
func GetPySerialCommand(args ...string) (*exec.Cmd, error) {
    pySerialPath, err := GetToolchainPath()
    if err != nil {
        return nil, err
    }

    if operatingSystem == io.LINUX {
        pySerialPath += io.Sep + serialLinux
    } else if operatingSystem == io.DARWIN {
        pySerialPath += io.Sep + serialDarwin
    } else if operatingSystem == io.WINDOWS {
        pySerialPath += io.Sep + serialWindows
    } else {
        return nil, errors.New("pyserial not available for this operating system")
    }

    return exec.Command(pySerialPath, args...), nil
}
