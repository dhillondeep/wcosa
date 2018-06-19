package toolchain

import (
	"os/exec"
	"wio/cmd/wio/utils/io"
)

const (
	serialLinux  = "serial/serial-ports-linux"
	serialDarwin = "serial/serial-ports-mac"
)

var operatingSystem = io.GetOS()

// This returns the path to toolchain directory
func GetToolchainPath() (string, error) {
	executablePath, err := io.NormalIO.GetRoot()
	if err != nil {
		return "", err
	}

	toolchainPath := executablePath + io.Sep + "toolchain"
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
	} else {
		// TODO change to windows
		pySerialPath += io.Sep + serialDarwin
	}

	return exec.Command(pySerialPath, args...), nil
}
