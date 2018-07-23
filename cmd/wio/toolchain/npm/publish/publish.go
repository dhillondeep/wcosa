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
    log.Info(log.Cyan, "Retrieving token ... ")
    token, err := login.LoadToken(dir)
    if err != nil {
        log.WriteFailure()
        return err
    }
    log.WriteSuccess()
    header := NewHeader(token.Value)

    log.Info(log.Cyan, "Zipping package ... ")
    data, err := VersionData(dir, cfg)
    if err != nil {
        log.WriteFailure()
        return err
    }
    if err := GeneratePackage(dir, data); err != nil {
        log.WriteFailure()
        return err
    }
    tarFile := fmt.Sprintf("%s-%s.tgz", data.Name, data.Version)
    tarPath := io.Path(dir, io.Folder, tarFile)
    if err := MakeTar(dir, tarPath); err != nil {
        log.WriteFailure()
        return err
    }
    log.WriteSuccess()

    log.Info(log.Cyan, "Computing shasum ... ")
    tarData, err := ioutil.ReadFile(tarPath)
    if err != nil {
        log.WriteFailure()
        return err
    }
    shasum := Shasum(tarData)
    tarDist := TarEncode(tarData)
    log.WriteSuccess()
    log.Verbln("Data length:    %d", len(tarData))
    log.Verbln("Encoded length: %d", len(tarDist))

    tarUrl := client.UrlResolve(client.BaseUrl, data.Name, "-", tarFile)
    data.Dist = npm.Dist{Shasum: shasum, Tarball: tarUrl}

    payload := &Attachment{
        Type:   "application/octet-stream",
        Data:   tarDist,
        Length: uint64(len(tarData)),
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

    log.Info(log.Cyan, "Sending request .... ")
    res := &Response{}
    status, err := client.GetJson(client.Npm, req, res)
    if err != nil {
        log.WriteFailure()
        return err
    }
    if res.Success != true {
        log.WriteFailure()
        return PublishError{res.Error}
    }
    if status != http.StatusOK {
        log.WriteFailure()
        return HttpFailed{status}
    }
    log.WriteSuccess()
    str, _ = json.MarshalIndent(res, "", login.Indent)
    log.Verbln("Response:\n%s", string(str))
    log.Info(log.Cyan, "Published ")
    log.Info(log.Green, "%s@%s", data.Name, data.Version)
    log.Infoln(log.Cyan, "!")
    return nil
}
