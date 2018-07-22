package login

import (
    "bytes"
    "encoding/json"
    "net/http"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/log"
    "wio/cmd/wio/toolchain/npm/client"
)

func Do(name, pass, email string) (*Response, error) {
    header := ReqHeader()
    body := ReqBody(name, pass, email)
    req, err := Request(header, body)
    if err != nil {
        return nil, err
    }
    res := &Response{}
    status, err := client.GetJson(client.Npm, req, res)
    if err != nil {
        return nil, err
    }
    if status == http.StatusUnauthorized {
        return nil, errors.String("401 invalid npm login")
    }
    if status != http.StatusCreated {
        return nil, errors.Stringf("%d status code", status)
    }
    str, _ := json.MarshalIndent(res, "", Indent)
    log.Infoln("Response:\n%s", str)
    return res, nil
}

func Request(header *Header, body *Body) (*http.Request, error) {
    url := client.UrlResolve(client.BaseUrl, "-", "user", body.Id)
    log.Infoln("PUT %s", url)
    str, _ := json.MarshalIndent(body, "", Indent)
    log.Infoln("Body:\n%s", str)
    req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(str)))
    if err != nil {
        return nil, err
    }
    str, _ = json.MarshalIndent(header, "", Indent)
    log.Infoln("Header:\n%s", str)
    req.Header.Set("content-type", header.ContentType)
    return req, nil
}

func GetToken(name, pass, email string) (*Token, error) {
    res, err := Do(name, pass, email)
    if err != nil {
        return nil, err
    }
    return &Token{Value: res.Token}
}
