package run

import (
    "os"
    "os/exec"
    "fmt"
    "runtime"
    sysio "io"
    "wio/cmd/wio/utils/io"
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

func cleanTarget(dir string) error {
    return execute(dir, "make", "clean")
}

func configAndBuild(dir string, errchan chan error) {
    binDir := dir + io.Sep + "bin"
    if err := os.MkdirAll(binDir, os.ModePerm); err != nil {
        errchan <- err
    } else if err := configTarget(binDir); err != nil {
        errchan <- err
    } else {
        errchan <- buildTarget(binDir)
    }
}

func cleanIfExists(dir string, errchan chan error) {
    binDir := dir + io.Sep + "bin"
    exists, err := io.Exists(binDir)
    if err != nil {
        errchan <- err
    } else if exists {
        errchan <- cleanTarget(dir)
    } else {
        errchan <- nil
    }
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
