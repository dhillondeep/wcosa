package login

import (
    "bytes"
    "encoding/json"
    "net/http"
    "wio/pkg/log"
    "wio/pkg/npm/client"
    "wio/pkg/util"
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
        return nil, util.Error("401 invalid npm login")
    }
    if status != http.StatusCreated {
        return nil, util.Error("%d login error", status)
    }
    str, _ := json.MarshalIndent(res, "", Indent)
    log.Verbln("Response:\n%s", str)
    if res.Status != true {
        return nil, util.Error("unknown error while logging in")
    }
    return res, nil
}

func Request(header *Header, body *Body) (*http.Request, error) {
    url := client.UrlResolve(client.BaseUrl, "-", "user", body.Id)
    log.Verbln("\nPUT %s", url)
    str, _ := json.MarshalIndent(body, "", Indent)
    log.Verbln("Body:\n%s", str)
    req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(str)))
    if err != nil {
        return nil, err
    }
    str, _ = json.MarshalIndent(header, "", Indent)
    log.Verbln("Header:\n%s", str)
    req.Header.Set("content-type", header.ContentType)
    return req, nil
}

func GetToken(name, pass, email string) (*Token, error) {
    res, err := Do(name, pass, email)
    if err != nil {
        return nil, err
    }
    return &Token{Value: res.Token}, nil
}
