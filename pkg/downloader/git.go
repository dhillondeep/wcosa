package downloader

import (
    "fmt"
    "os"
    "wio-utils/cmd/io"
    "wio/pkg/util"
    "wio/pkg/util/sys"

    git "gopkg.in/src-d/go-git.v4"
    "gopkg.in/src-d/go-git.v4/plumbing"
)

type GitDownloader struct{}

func (gitDownloader GitDownloader) DownloadModule(path, url, reference string) (string, error) {
    if !util.IsCommandAvailable("git", "--version") {
        return "", util.Error("In order to use dev version of toolchain, git must be installed")
    }

    cloneOptions := &git.CloneOptions{
        URL:               "https://" + url,
        Progress:          os.Stdout,
        Depth:             1,
        RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
    }

    if !util.IsEmptyString(reference) {
        cloneOptions.ReferenceName = plumbing.ReferenceName(reference)
    } else {
        reference = "default"
    }

    clonePath := sys.Path(path, url) + "__" + reference

    if sys.Exists(clonePath) {
        return clonePath, nil
    }

    _, err := git.PlainClone(clonePath, false, cloneOptions)

    if err != nil {
        os.RemoveAll(clonePath)
        return "", util.Error("toolchain could not be downloaded, check url and reference")
    }

    moduleData := &ModuleData{}

    io.NormalIO.ParseJson(sys.Path(clonePath, "package.json"), moduleData)

    for name, version := range moduleData.Dependencies {
        DownloadToolchain(fmt.Sprintf("%s:%s", name, version))
    }

    return clonePath, nil
}
