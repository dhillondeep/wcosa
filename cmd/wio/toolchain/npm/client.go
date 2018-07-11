package npm

import (
    "bytes"
    "encoding/json"
    "net/http"
    "time"

    "io"
    "os"
    "strings"
    "wio/cmd/wio/errors"

    "github.com/mholt/archiver"
)

const timeoutSeconds = 10

const (
    registryBaseUrl = "https://registry.npmjs.org"
)

var clientInstance = &http.Client{Timeout: timeoutSeconds * time.Second}

func getJson(client *http.Client, url string, target interface{}) error {
    resp, err := client.Get(url)
    defer resp.Body.Close()
    if err != nil {
        return err
    }
    if resp.StatusCode != http.StatusOK {
        return errors.Stringf("HTTP request to %s returned %d", url, resp.StatusCode)
    }
    return json.NewDecoder(resp.Body).Decode(target)
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

func makePackageUrl(name string) string {
    return urlResolve(registryBaseUrl, name)
}

func getPackageData(name string) (*packageData, error) {
    var data packageData
    url := makePackageUrl(name)
    err := getJson(clientInstance, url, &data)
    if err == nil && data.Error != "" {
        err = errors.String(data.Error)
    }
    return &data, err
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
