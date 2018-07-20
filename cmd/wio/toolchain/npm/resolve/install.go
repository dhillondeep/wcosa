package resolve

import (
    sysio "io"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/toolchain/npm"
    "wio/cmd/wio/utils/io"

    "wio/cmd/wio/toolchain/npm/semver"

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

func (i *Info) install(name string, ver string, data *npm.Version) error {
    local, err := tryFindConfig(&Node{Name: name, ResolvedVersion: semver.Parse(ver)}, i.dir)
    if err != nil {
        return err
    }
    if local != nil {
        return nil
    }

    file := name + "__" + ver
    tar := io.Path(i.dir, io.Folder, io.Download, file+".tgz")
    if !io.Exists(tar) {
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

    modules := io.Path(i.dir, io.Folder, io.Modules)
    pkg := io.Path(modules, "package")
    if err := os.RemoveAll(pkg); err != nil {
        return err
    }
    if err := untar(tar, modules); err != nil {
        return err
    }
    return os.Rename(pkg, io.Path(modules, file))
}

func download(url string, dst string, cb sysio.Writer) error {
    if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
        return err
    }
    out, err := os.Create(dst + io.Temp)
    if err != nil {
        return err
    }
    defer out.Close()
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if _, err := sysio.Copy(out, sysio.TeeReader(resp.Body, cb)); err != nil {
        return err
    }
    return os.Rename(dst+io.Temp, dst)
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
        return 0, errors.Stringf("GET %s returned %d", url, resp.StatusCode)
    }
    str := resp.Header.Get("Content-Length")
    return strconv.ParseUint(str, 10, 64)
}
