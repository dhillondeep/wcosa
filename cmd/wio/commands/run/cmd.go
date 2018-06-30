package run

import (
    "os/exec"
    "os"
    sysio "io"
    "wio/cmd/wio/utils/io"
)

func configTarget(dir string) error {
    return execute(dir, "cmake", "../", "-G", "Unix Makefiles")
}

func buildTarget(dir string) error {
    return execute(dir, "make")
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
        return
    }
    if err := configTarget(binDir); err != nil {
        errchan <- err
        return
    }
    errchan <- buildTarget(binDir)
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
