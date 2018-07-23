package publish

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "wio/cmd/wio/log"
    "wio/cmd/wio/toolchain/npm"
    "wio/cmd/wio/toolchain/npm/client"
    "wio/cmd/wio/toolchain/npm/login"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"
)

func Do(dir string, cfg types.Config) error {
    token, err := login.LoadToken(dir)
    if err != nil {
        return err
    }
    header := NewHeader(token.Value)

    data, err := VersionData(dir, cfg)
    if err != nil {
        return err
    }
    if err := GeneratePackage(dir, data); err != nil {
        return err
    }
    tarFile := fmt.Sprintf("%s-%s.tgz", data.Name, data.Version)
    tarPath := io.Path(dir, io.Folder, tarFile)
    if err := MakeTar(dir, tarPath); err != nil {
        return err
    }

    tarData, err := ioutil.ReadFile(tarPath)
    if err != nil {
        return err
    }
    shasum := Shasum(tarData)
    tarDist := TarEncode(tarData)

    tarUrl := client.UrlResolve(client.BaseUrl, data.Name, "-", tarFile)
    data.Dist = npm.Dist{Shasum: shasum, Tarball: tarUrl}

    payload := &Attachment{
        Type:   "application/octet-stream",
        Data:   tarDist,
        Length: 1024,
    }
    body := &Data{
        Id:          data.Name,
        Name:        data.Name,
        Description: data.Description,
        Readme:      data.Readme,

        DistTags:    map[string]string{"latest": data.Version},
        Versions:    map[string]*npm.Version{data.Version: data},
        Attachments: map[string]*Attachment{tarFile: payload},
    }

    url := client.UrlResolve(client.BaseUrl, data.Name)
    log.Verbln("PUT %s", url)
    str, _ := json.MarshalIndent(header, "", login.Indent)
    log.Verbln("Header:\n%s", string(str))
    str, _ = json.MarshalIndent(body, "", login.Indent)
    log.Verbln("Body:\n%s", string(str))
    req, err := http.NewRequest("PUT", url, bytes.NewBuffer(str))
    if err != nil {
        return err
    }
    req.Header.Set("authorization", header.Authorization)
    req.Header.Set("content-type", header.ContentType)
    req.Header.Set("npm-session", header.NpmSession)

    res := &Response{}
    status, err := client.GetJson(client.Npm, req, res)
    if err != nil {
        return err
    }
    if status != http.StatusOK {
        return HttpFailed{status}
    }
    if res.Success != true {
        return UnknownError{}
    }
    str, _ = json.MarshalIndent(res, "", login.Indent)
    log.Verbln("Response:\n%s", string(str))
    return nil
}
