package packages

import (
    "wio/cmd/wio/utils/io"
    "os"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/log"
    "os/user"
    "path/filepath"
)

type GlobalConfigurations struct {
    Platforms struct {
        AvrPath string
        ArmPath string
    }
    Frameworks struct {
        ArduinoAvrPath string
        CosaAVRPath string
    }

}

var operatingSystem = io.GetOS()
var programFiles map[string]string
var packages map[string]string
var platformPaths map[string]string
var frameworkPaths map[string]string

// Set's up the initial filesystem to store all the wio related files
// If it is already setup, it does not do anything
func ResolveProgramFilesPathAndCreate(queue *log.Queue) error {
    executablePath, err := io.NormalIO.GetRoot()
    if err != nil {
        return err
    }

    usr, err := user.Current()
    if err != nil {
        return err
    }

    log.QueueWriteln(queue, log.VERB, nil, "setting up program files path")
    programFiles = map[string]string {
        io.WINDOWS: filepath.Clean(executablePath + io.Sep + ".wio"),
        io.LINUX: filepath.Clean(usr.HomeDir + io.Sep + ".wio"),
        io.DARWIN: filepath.Clean(usr.HomeDir + io.Sep + ".wio"),
    }

    log.QueueWriteln(queue, log.VERB, nil, "setting up program packages path")
    packages = map[string]string {
        io.WINDOWS: filepath.Clean(programFiles[io.WINDOWS] + io.Sep + "packages"),
        io.LINUX: filepath.Clean(programFiles[io.WINDOWS] + io.Sep + "packages"),
        io.DARWIN: filepath.Clean(programFiles[io.WINDOWS] + io.Sep + "packages"),
    }

    log.QueueWriteln(queue, log.VERB, nil, "setting up platform path")
    platformPaths = map[string]string{
        io.WINDOWS: filepath.Clean(packages[io.WINDOWS] + io.Sep + "platforms"),
        io.LINUX: filepath.Clean(packages[io.LINUX] + io.Sep + "platforms"),
        io.DARWIN: filepath.Clean(packages[io.DARWIN] + io.Sep + "platforms"),
    }

    log.QueueWriteln(queue, log.VERB, nil, "setting up framework path")
    frameworkPaths = map[string]string{
        io.WINDOWS: filepath.Clean(packages[io.WINDOWS] + io.Sep + "frameworks"),
        io.LINUX: filepath.Clean(packages[io.LINUX] + io.Sep + "frameworks"),
        io.DARWIN: filepath.Clean(packages[io.DARWIN] + io.Sep + "frameworks"),
    }

    // wio only supports windows, linux and darwin for now
    if operatingSystem != io.WINDOWS && operatingSystem != io.LINUX && operatingSystem != io.DARWIN {
        return errors.OperatingSystemNotSupportedError{
            OperatingSystem: operatingSystem,
        }
    }

    // create platform path
    if !utils.PathExists(platformPaths[operatingSystem]) {
        log.QueueWriteln(queue, log.VERB, nil, "creating path: %s", platformPaths[operatingSystem])
        if err := os.MkdirAll(platformPaths[operatingSystem], os.ModePerm); err != nil {
            return err
        }
    }

    // create framework path
    if !utils.PathExists(frameworkPaths[operatingSystem]) {
        log.QueueWriteln(queue, log.VERB, nil, "creating path: %s", frameworkPaths[operatingSystem])
        if err := os.MkdirAll(frameworkPaths[operatingSystem], os.ModePerm); err != nil {
            return err
        }
    }

    return nil
}
