package run

import (
    "fmt"
    "os"
    "os/exec"
    "runtime"
    "strings"
    "wio/pkg/log"
    "wio/pkg/util"
    "wio/pkg/util/sys"
)

func configTarget(dir string) error {
    cmakeGenerator, err := util.GetCmakeGenerator(sys.Path(dir, "../"))
    if err != nil {
        return err
    }

    return Execute(dir, "cmake", "../", "-G", cmakeGenerator)
}

func buildTarget(dir string) error {
    buildTool, err := util.GetBuildTool(sys.Path(dir, "../"))
    if err != nil {
        return err
    }

    if buildTool == "make" {
        jobs := runtime.NumCPU() + 2
        jobsFlag := fmt.Sprintf("-j%d", jobs)
        return Execute(dir, buildTool, jobsFlag)
    }

    return Execute(dir, buildTool)
}

func uploadTarget(dir string) error {
    buildTool, err := util.GetBuildTool(sys.Path(dir, "../"))
    if err != nil {
        return err
    }

    return Execute(dir, buildTool, "upload")
}

func runTarget(dir, file, args string) error {
    var argv []string
    if args != "" {
        argv = strings.Split(args, " ")
    }

    return Execute(dir, file, argv...)
}

func cleanTarget(dir string) error {
    buildTool, err := util.GetBuildTool(sys.Path(dir, "../"))
    if err != nil {
        return err
    }

    return Execute(dir, buildTool, "clean")
}

type targetFunc func(string, chan error)

func configAndBuild(dir string, errChan chan error) {
    log.Verbln(log.Magenta, "Building directory: %s", dir)
    binDir := sys.Path(dir, "bin")
    if err := os.MkdirAll(binDir, os.ModePerm); err != nil {
        errChan <- err
    } else if err := configTarget(binDir); err != nil {
        errChan <- err
    } else {
        errChan <- buildTarget(binDir)
    }
}

func cleanIfExists(dir string, errChan chan error) {
    log.Verbln(log.Magenta, "Cleaning directory: %s", dir)
    binDir := sys.Path(dir, "bin")
    exists := sys.Exists(binDir)
    if exists {
        errChan <- cleanTarget(binDir)
    } else {
        errChan <- nil
    }
}

func hardClean(dir string, errChan chan error) {
    log.Verbln(log.Magenta, "Removing directory: %s", dir)
    errChan <- os.RemoveAll(dir)
}

func Execute(dir string, name string, args ...string) error {
    cmd := exec.Command(name, args...)
    cmd.Dir = dir
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}
