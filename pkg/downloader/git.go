package downloader

import (
    "fmt"
    "os"
    "wio/pkg/log"
    "wio/pkg/util"
    "wio/pkg/util/sys"

    git "gopkg.in/src-d/go-git.v4"
    "gopkg.in/src-d/go-git.v4/plumbing"
)

type GitDownloader struct{}

const (
    Protocol   = "https"
    DefaultRef = "master"
)

func (gitDownloader GitDownloader) DownloadModule(path, url, reference string, retool bool) (string, error) {
    log.Write(log.Cyan, "Fetching toolchain using git from ")
    log.Write(log.Green, url)
    log.Write(log.Cyan, "... ")

    cloneOptions := &git.CloneOptions{
        URL:               fmt.Sprintf("%s://%s", Protocol, url),
        Progress:          os.Stdout,
        Depth:             1,
        RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
    }

    if !util.IsEmptyString(reference) {
        cloneOptions.ReferenceName = plumbing.ReferenceName(reference)
    } else {
        reference = DefaultRef
    }

    clonePath := fmt.Sprintf("%s__%s", sys.Path(path, url), reference)

    if sys.Exists(clonePath) && !retool {
        log.Writeln(log.Green, "already exists")
        return clonePath, nil
    } else if sys.Exists(clonePath) && retool {
        os.RemoveAll(clonePath)
        log.Writeln(log.Green, "retooling")
    } else {
        log.Writeln(log.Green, "downloading")
    }

    _, err := git.PlainClone(clonePath, false, cloneOptions)

    if err != nil {
        os.RemoveAll(clonePath)
        return "", util.Error("toolchain could not be downloaded, check url and reference")
    }

    moduleData := &ModuleData{}

    if err := sys.NormalIO.ParseJson(sys.Path(clonePath, "package.json"), moduleData); err != nil {
        return "", util.Error("toolchain config error")
    }

    for name, version := range moduleData.Dependencies {
        if _, err := DownloadToolchain(fmt.Sprintf("%s:%s", name, version), retool); err != nil {
            return "", err
        }
    }

    return clonePath, nil
}
