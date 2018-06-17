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

func GetToolchainPath() (string, error) {
	executablePath, err := io.NormalIO.GetRoot()
	if err != nil {
		return "", err
	}

	toolchainPath := executablePath + io.Sep + "toolchain"
	return toolchainPath, nil
}

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

//# all the paths toolchain can be at (this is because of different package managers)
//if (EXISTS "${CMAKE_TOOLCHAIN_PATH}/{{TOOLCHAIN_FILE_REL}}")
//set(CMAKE_TOOLCHAIN_FILE "${CMAKE_TOOLCHAIN_PATH}/{{TOOLCHAIN_FILE_REL}}")
//elseif (EXISTS "${CMAKE_TOOLCHAIN_PATH}/../{{TOOLCHAIN_FILE_REL}}")
//set(CMAKE_TOOLCHAIN_FILE "${CMAKE_TOOLCHAIN_PATH}/../{{TOOLCHAIN_FILE_REL}}")
//elseif (EXISTS "/usr/share/wio/{{TOOLCHAIN_FILE_REL}}")
//set(CMAKE_TOOLCHAIN_FILE "/usr/share/wio/{{TOOLCHAIN_FILE_REL}}")
//else()
//message(FATAL_ERROR "Toolchain cannot be found. Build Halted!")
//endif()
