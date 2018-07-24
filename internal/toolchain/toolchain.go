package toolchain

import (
    "errors"
    "os/exec"
    "path/filepath"
    "wio/pkg/util/sys"
)

const (
    serialLinux   = "serial/serial-ports-linux"
    serialDarwin  = "serial/serial-ports-mac"
    serialWindows = "serial/serial-ports.exe"
)

var operatingSystem = sys.GetOS()

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

// This is the command to execute PySerial to get ports information
func GetPySerialCommand(args ...string) (*exec.Cmd, error) {
    pySerialPath, err := GetToolchainPath()
    if err != nil {
        return nil, err
    }

    if operatingSystem == sys.LINUX {
        pySerialPath += sys.Sep + serialLinux
    } else if operatingSystem == sys.DARWIN {
        pySerialPath += sys.Sep + serialDarwin
    } else if operatingSystem == sys.WINDOWS {
        pySerialPath += sys.Sep + serialWindows
    } else {
        return nil, errors.New("pyserial not available for this operating system")
    }

    return exec.Command(pySerialPath, args...), nil
}
