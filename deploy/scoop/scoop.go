package main

import (
    "bufio"
    "gopkg.in/src-d/go-git.v4"
    "log"
    "os"
    "path/filepath"
    "strings"
    "wio/internal/config/meta"
    "wio/pkg/util/sys"
    "wio/pkg/util/template"
)

const (
    exec32Name  = "wio_windows_i386"
    exec64Name  = "wio_windows_x86_64"
    extension = "zip"
)

func main() {
    ex, err := os.Executable()
    if err != nil {
        log.Fatal(err)
    }

    currPath := filepath.Dir(ex)

    cloneOptions := &git.CloneOptions{
        URL:               "https://github.com/wio/wio-bucket",
        Progress:          os.Stdout,
        Depth:             1,
        RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
    }

    repoPath := sys.Path(currPath, "wio-bucket")

    if err := os.RemoveAll(repoPath); err != nil {
        log.Fatal(err)
    }

    if _, err := git.PlainClone(repoPath, false, cloneOptions); err != nil {
        log.Fatal(err)
    }

    checkSumsPath := sys.Path(currPath, "../", "bin", "checksums.txt")

    checkSumsFile, err := os.Open(checkSumsPath)
    if err != nil {
       log.Fatal(err)
    }
    defer checkSumsFile.Close()

    scanner := bufio.NewScanner(checkSumsFile)
    var checkSum32 string
    var checkSum64 string
    for scanner.Scan() {
       tokens := strings.Split(scanner.Text(), " ")

       if strings.Trim(tokens[2], " ") == exec32Name+"."+extension {
           checkSum32 = tokens[0]
       } else if strings.Trim(tokens[2], " ") == exec64Name+"."+extension {
           checkSum64 = tokens[0]
       }
    }

    jsonFile := sys.Path(repoPath, "wio.json")

    if err := sys.NormalIO.CopyFile(sys.Path(currPath, "scoop-template.tpl"), jsonFile, true); err != nil {
        log.Fatal(err)
    }

    if err := template.IOReplace(jsonFile, map[string]string{
        "version": meta.Version,
        "exec32bit": exec32Name,
        "exec64bit": exec64Name,
        "extension": extension,
        "checksum32bit": checkSum32,
        "checksum64bit": checkSum64,
    }); err != nil {
        log.Fatal(err)
    }
}
