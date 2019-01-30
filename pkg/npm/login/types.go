package login

import (
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "os"
    "strings"
    "time"
    "wio/pkg/util/sys"
)

type Body struct {
    Id       string   `json:"_id"`
    Name     string   `json:"name"`
    Password string   `json:"password"`
    Email    string   `json:"email"`
    Type     string   `json:"type"`
    Roles    []string `json:"roles"`
    Date     string   `json:"date"`
}

type Header struct {
    ContentType string `json:"content-type"`
}

type Response struct {
    Status bool   `json:"ok"`
    Id     string `json:"id"`
    Rev    string `json:"rev"`
    Token  string `json:"token"`
}

type Token struct {
    Value string `json:"token"`
}

const TimeFormat = "2006-01-01 15:04:05.000"
const IdPrefix = "org.couchdb.user:"
const Indent = "    "

func DateFormat(t time.Time) string {
    ret := t.Format(TimeFormat)
    parts := strings.Split(ret, " ")
    return fmt.Sprintf("%sT%sZ", parts[0], parts[1])
}

func ReqHeader() *Header {
    return &Header{ContentType: "application/json"}
}

func ReqBody(user, pass, email string) *Body {
    return &Body{
        Id:       IdPrefix + user,
        Name:     user,
        Password: pass,
        Email:    email,
        Type:     "user",
        Roles:    []string{},
        Date:     DateFormat(time.Now()),
    }
}

func (t *Token) Save(dir string) error {
    path := sys.Path(dir, sys.WioFolder)
    if err := os.MkdirAll(path, os.ModePerm); err != nil {
        return err
    }
    path = sys.Path(path, "token.json")
    str, _ := json.MarshalIndent(t, "", Indent)
    return ioutil.WriteFile(path, []byte(str), os.ModePerm)
}

func LoadToken(dir string) (*Token, error) {
    path := sys.Path(dir, sys.WioFolder, "token.json")
    if !sys.Exists(path) {
        return nil, errors.New("not logged in")
    }
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }
    ret := &Token{}
    if err := json.Unmarshal(data, ret); err != nil {
        return nil, err
    }
    return ret, nil
}
