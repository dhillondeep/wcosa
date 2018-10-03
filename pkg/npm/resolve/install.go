package resolve

import (
    "io"
    "io/ioutil"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "wio/pkg/npm"
    "wio/pkg/npm/publish"
    "wio/pkg/util"
    "wio/pkg/util/sys"

    "github.com/mholt/archiver"
)

func (i *Info) InstallResolved() error {
    logInstallStart()

    for name, cache := range i.ver {
        for ver, data := range cache {
            if err := i.install(name, ver, data); err != nil {
                return err
            }
        }
    }

    logInstallDone()
    return nil
}

func (i *Info) install(name, ver string, data *npm.Version) error {
    local, err := i.GetPkg(name, ver)
    if err != nil {
        return err
    }
    if local != nil {
        return nil
    }

    file := name + "__" + ver
    tar := sys.Path(i.dir, sys.Folder, sys.Download, file+".tgz")
    if !sys.Exists(tar) {
        url := data.Dist.Tarball
        total, err := contentSize(url)
        if err != nil {
            return err
        }
        cb := &counter{total: total, cb: installCallback(name, ver)}
        if err := download(url, tar, cb); err != nil {
            return err
        }
    }

    // check shashum
    tarData, err := ioutil.ReadFile(tar)
    if err != nil {
        return err
    }
    if sha := publish.Shasum(tarData); sha != data.Dist.Shasum {
        if err := os.RemoveAll(tar); err != nil {
            return err
        }
        return util.Error("expected tar checksum %s", data.Dist.Shasum)
    }

    modules := sys.Path(i.dir, sys.Folder, sys.Modules)
    pkg := sys.Path(modules, "package")
    if err := os.RemoveAll(pkg); err != nil {
        return err
    }
    if err := untar(tar, modules); err != nil {
        return err
    }
    return os.Rename(pkg, sys.Path(modules, file))
}

func download(url string, dst string, cb io.Writer) error {
    if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
        return err
    }
    out, err := os.Create(dst + sys.Temp)
    if err != nil {
        return err
    }
    defer out.Close()
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if _, err := io.Copy(out, io.TeeReader(resp.Body, cb)); err != nil {
        return err
    }
    out.Close()
    resp.Body.Close()
    return os.Rename(dst+sys.Temp, dst)
}

func untar(src string, dst string) error {
    return archiver.TarGz.Open(src, dst)
}

func installCallback(name string, ver string) callback {
    return func(curr uint64, total uint64) {
        logInstall(name, ver, curr, total)
    }
}

func contentSize(url string) (uint64, error) {
    resp, err := http.Head(url)
    if err != nil {
        return 0, err
    }
    if resp.StatusCode != http.StatusOK {
        return 0, util.Error("GET %s returned %d", url, resp.StatusCode)
    }
    str := resp.Header.Get("Content-Length")
    return strconv.ParseUint(str, 10, 64)
}
