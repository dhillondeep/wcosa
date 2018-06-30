package run

import (
    "os"
    "os/exec"
    "fmt"
    "runtime"
    sysio "io"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/log"
)

func configTarget(dir string) error {
    return execute(dir, "cmake", "../", "-G", "Unix Makefiles")
}

func buildTarget(dir string) error {
    jobs := runtime.NumCPU() + 2
    jobsFlag := fmt.Sprintf("-j%d", jobs)
    return execute(dir, "make", jobsFlag)
}

func uploadTarget(dir string) error {
    return execute(dir, "make", "upload")
}

func runTarget(dir string, file string) error {
    return execute(dir, file)
}

func cleanTarget(dir string) error {
    return execute(dir, "make", "clean")
}

type targetFunc func(string, chan error)

func configAndBuild(dir string, errChan chan error) {
    log.Verbln(log.Magenta, "Building directory: %s", dir)
    binDir := dir + io.Sep + "bin"
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
    binDir := dir + io.Sep + "bin"
    exists, err := io.Exists(binDir)
    if err != nil {
        errChan <- err
    } else if exists {
        errChan <- cleanTarget(binDir)
    } else {
        errChan <- nil
    }
}

func hardClean(dir string, errChan chan error) {
    log.Verbln(log.Magenta, "Removing directory: %s", dir)
    errChan <- os.RemoveAll(dir)
}

func execute(dir string, name string, args ...string) error {
    cmd := exec.Command(name, args...)
    cmd.Dir = dir
    stdoutIn, err := cmd.StdoutPipe()
    if err != nil {
        return err
    }
    stderrIn, err := cmd.StderrPipe()
    if err != nil {
        return err
    }
    err = cmd.Start()
    if err != nil {
        return err
    }
    go func() { sysio.Copy(os.Stdout, stdoutIn) }()
    go func() { sysio.Copy(os.Stderr, stderrIn) }()
    err = cmd.Wait()
    if err != nil {
        return err
    }
    return nil
}
