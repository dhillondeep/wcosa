package resolve

import (
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/toolchain/npm"
    util "wio/cmd/wio/utils/io"
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
    local, err := tryFindConfig(name, ver, i.dir)
    if err != nil {
        return err
    }
    if local != nil {
        return nil
    }

    file := name + "__" + ver
    dst := util.Path(i.dir, util.Folder, util.Download, file+".tgz")
    if util.Exists(dst) {
        return nil
    }

    url := data.Dist.Tarball
    total, err := contentSize(url)
    if err != nil {
        return err
    }
    cb := &counter{total: total, cb: installCallback(name, ver)}
    return download(url, dst, cb)
}

func download(url string, dst string, cb io.Writer) error {
    if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
        return err
    }
    out, err := os.Create(dst + util.Temp)
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
    return os.Rename(dst+util.Temp, dst)
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
