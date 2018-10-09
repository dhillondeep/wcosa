package util

import (
    "os/exec"
    "wio/pkg/util/sys"
)

const (
    cmakeConfig = "cmake.config"
)

var makeForCmakeGenerators = map[string]map[string]string{
    "windows": {
        "Ninja":               "ninja",
        "Unix Makefiles":      "make",
        "MinGW Makefiles":     "mingw32-make",
        "NMake Makefiles":     "nmake",
        "Borland Makefiles":   "bcc32",
        "MSYS Makefiles":      "make",
        "NMake Makefiles JOM": "jom",
        "Watcom WMake":        "wmake",
    },
    "darwin": {
        "Ninja":          "ninja",
        "Unix Makefiles": "make",
    },
    "linux": {
        "Ninja":          "ninja",
        "Unix Makefiles": "make",
        "Watcom WMake":   "wmake",
    },
}

func isCommandAvailable(name string) bool {
    cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
    if err := cmd.Run(); err != nil {
        return false
    }
    return true
}

func mingwExists() bool {
    return isCommandAvailable("mingw32-make --version")
}

func nmakeExists() bool {
    return isCommandAvailable("nmake ?")
}

func makeExists() bool {
    return isCommandAvailable("make -h")
}

func bcc32Exists() bool {
    return isCommandAvailable("bcc32 -h")
}

func jomExists() bool {
    return isCommandAvailable("jom -h")
}

func wmakeExists() bool {
    return isCommandAvailable("wmake -h")
}

func ninjaExists() bool {
    return isCommandAvailable("ninja -h")
}

func GetCmakeGenerator(buildDir string) (string, error) {
    var generator string
    var err error

    if ninjaExists() {
        generator = "Ninja"
    } else {
        os := sys.GetOS()

        // for windows
        if os == sys.WINDOWS {
            if mingwExists() {
                generator = "MinGW Makefiles"
            } else if nmakeExists() {
                generator = "NMake Makefiles"
            } else if jomExists() {
                generator = "NMake Makefiles JOM"
            } else if makeExists() {
                generator = "Unix Makefiles"
            } else if wmakeExists() {
                generator = "Watcom WMake"
            } else if bcc32Exists() {
                generator = "Borland Makefiles"
            } else {
                err = Error("No build tool found. Please install one of ninja, make, mingw32-make, etc")
            }
        } else {
            if makeExists() {
                generator = "Unix Makefiles"
            } else {
                if os == sys.DARWIN {
                    err = Error("No build tool found. Please install one of ninja, or GNU make")
                } else if wmakeExists() {
                    generator = "Watcom WMake"
                } else {
                    err = Error("No build tool found. Please install one of ninja, GNU make, or wmake")
                }
            }
        }
    }

    if err == nil {
        // build files exist and cmake config file exists
        if sys.Exists(sys.Path(buildDir, cmakeConfig)) {
            currGenerator, _ := GetCurrentGenerator(buildDir)

            if generator != string(currGenerator) {
                return "", Error("Build tools mismatch Before: \"%s\" & Now: \"%s\". "+"wio clean --hard to fix",
                    currGenerator, generator)
            } else {
                return string(currGenerator), nil
            }
        } else if sys.Exists(sys.Path(buildDir, "bin", "CMakeFiles")) {
            // cmake config file doesn't exist but build files exist
            return GetCurrentGenerator(buildDir)
        }

        // cmake config file and build files do not exist
        if err := sys.NormalIO.WriteFile(sys.Path(buildDir, cmakeConfig), []byte(generator)); err != nil {
            return "", err
        }
    }

    return generator, err
}

func GetCurrentGenerator(targetBuildDir string) (string, error) {
    generator, err := sys.NormalIO.ReadFile(sys.Path(targetBuildDir, cmakeConfig))
    if err != nil {
        return "", Error("Build tool cannot be determined. wio clean --hard to fix")
    }

    return string(generator), nil
}

func GetBuildTool(buildDir string) (string, error) {
    generator, err := GetCurrentGenerator(buildDir)
    if err != nil {
        return "", err
    }

    os := sys.GetOS()
    return makeForCmakeGenerators[os][generator], nil
}
