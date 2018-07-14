package client

import (
    "bytes"
    "encoding/json"
    "net/http"
    "time"

    "io"
    "os"
    "strings"
    "wio/cmd/wio/errors"

    "wio/cmd/wio/toolchain/npm"

    "github.com/mholt/archiver"
)

const timeoutSeconds = 10

const (
    registryBaseUrl = "https://registry.npmjs.org"
)

var npmClient = &http.Client{Timeout: timeoutSeconds * time.Second}

func getJson(client *http.Client, url string, target interface{}) (int, error) {
    resp, err := client.Get(url)
    defer resp.Body.Close()
    if err != nil {
        return 0, err
    }
    if resp.StatusCode != http.StatusOK {
        return resp.StatusCode, nil
    }
    return http.StatusOK, json.NewDecoder(resp.Body).Decode(target)
}

func findFirstSlash(value string) int {
    i := 0
    for ; i < len(value) && value[i] == '/'; i++ {
    }
    return i
}

func findLastSlash(value string) int {
    i := len(value) - 1
    for ; i >= 0 && value[i] == '/'; i-- {
    }
    return i
}

func urlResolve(values ...string) string {
    var buffer bytes.Buffer
    for _, value := range values {
        i := findFirstSlash(value)
        j := findLastSlash(value)
        buffer.WriteString(value[i : j+1])
        buffer.WriteString("/")
    }
    result := buffer.String()
    return result[:len(result)-1]
}

func FetchPackageData(name string) (*npm.Data, error) {
    var data npm.Data
    url := urlResolve(registryBaseUrl, name)
    status, err := getJson(npmClient, url, &data)
    if err != nil {
        return nil, err
    }
    if status == http.StatusNotFound {
        return nil, errors.Stringf("package not found: %s", name)
    }
    if status != http.StatusOK {
        return nil, errors.Stringf("registry GET (%s) returned %d", url, status)
    }
    return &data, nil
}

func FetchPackageVersion(name string, versionStr string) (*npm.Version, error) {
    // assumes `versionStr` is a hard version
    var version npm.Version
    url := urlResolve(registryBaseUrl, name, versionStr)
    status, err := getJson(npmClient, url, &version)
    if err != nil {
        return nil, err
    }
    if status == http.StatusNotFound {
        return nil, errors.Stringf("package not found: %s@%s", name, versionStr)
    }
    if status != http.StatusOK {
        return nil, errors.Stringf("registry GET (%s) returned %d", url, status)
    }
    return &version, nil
}

func downloadTarball(url string, dest string) error {
    if !strings.HasSuffix(url, ".tgz") {
        return errors.Stringf("invalid tarball URL: %s", url)
    }
    if !strings.HasSuffix(dest, ".tgz") {
        return errors.Stringf("invalid tarball path: %s", dest)
    }
    out, err := os.Create(dest)
    defer out.Close()
    if err != nil {
        return err
    }
    resp, err := http.Get(url)
    defer resp.Body.Close()
    if err != nil {
        return err
    }
    _, err = io.Copy(out, resp.Body)
    return err
}

func untar(src string, dest string) error {
    return archiver.TarGz.Open(src, dest)
}
