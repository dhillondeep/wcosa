package downloader

import (
    "io/ioutil"
    "os"
    "wio/pkg/npm/resolve"
    "wio/pkg/util"
    "wio/pkg/util/sys"
)

type NpmDownloader struct{}

func (npmDownloaded NpmDownloader) DownloadModule(path, name, version string) (string, error) {
    info := resolve.NewInfo(path)

    var err error = nil
    if util.IsEmptyString(version) {
        if version, err = info.GetLatest(name); err != nil {
            return "", err
        }
    }

    node := &resolve.Node{Name: name, ConfigVersion: version}

    if err = info.ResolveTree(node); err != nil {
        return "", err
    }

    if err = info.InstallResolved(); err != nil {
        return "", err
    }

    // symlink all the packages under path
    wioPackagesPath := sys.Path(path, sys.Folder, sys.Modules)

    fileInfo, err := ioutil.ReadDir(wioPackagesPath)
    if err != nil {
        return "", err
    }

    for _, file := range fileInfo {
        currFilePath := sys.Path(wioPackagesPath, file.Name())
        newFilePath := sys.Path(path, file.Name())

        if sys.Exists(newFilePath) {
            if err := os.RemoveAll(newFilePath); err != nil {
                return "", err
            }
        }

        if err := os.Symlink(currFilePath, newFilePath); err != nil {
            return "", err
        }
    }

    return sys.Path(path, node.Name) + "__" + version, nil
}
