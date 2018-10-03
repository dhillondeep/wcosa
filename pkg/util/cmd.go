package util

import (
    "os/exec"
	"wio/pkg/util/sys"
)

func mingwExists() bool {
	cmd := exec.Command("mingw32-make", "--version")
	return nil == cmd.Run()
}

func nmakeExists() bool {
	cmd := exec.Command("nmake", "/?")
	return nil == cmd.Run()
}

func ninjaExists() bool {
	cmd := exec.Command("ninja", "-h")
	return nil == cmd.Run()
}

func GetCmakeGenerator() string {
	if ninjaExists() {
		return "Ninja"
	}
	if sys.GetOS() != sys.WINDOWS {
		return "Unix Makefiles"
	}
	if mingwExists() {
		return "MinGW Makefiles"
	}
	if nmakeExists() {
		return "NMake Makefiles"
	}
	return "Unix Makefiles"
}

func GetMake() string {
	if ninjaExists() {
		return "ninja"
	}
	if sys.GetOS() != sys.WINDOWS {
		return "make"
	}
	if mingwExists() {
		return "mingw32-make"
	}
	if nmakeExists() {
		return "nmake"
	}
	return "make"
}
