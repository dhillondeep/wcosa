package publish

import (
    "fmt"
    "wio/cmd/wio/toolchain/npm"
)

type Attachment struct {
    Type   string `json:"content-type"`
    Data   string `json:"data"`
    Length uint64 `json:"length"`
}

type Data struct {
    Id          string `json:"_id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Readme      string `json:"readme"`

    DistTags    map[string]string       `json:"dist-tags"`
    Versions    map[string]*npm.Version `json:"versions"`
    Attachments map[string]*Attachment  `json:"_attachments"`
}

type Response struct {
    Success bool   `json:"success"`
    Error   string `json:"error"`
}

type Header struct {
    Authorization string `json:"authorization"`
    ContentType   string `json:"content-type"`
    NpmSession    string `json:"npm-session"`
}

func NewHeader(token string) *Header {
    return &Header{
        Authorization: fmt.Sprintf("Bearer %s", token),
        ContentType:   "application/json",
        NpmSession:    GenerateSession(),
    }
}
