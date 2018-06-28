package template

import (
    "strings"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/errors"
)

func IOReplace(path string, values map[string]string) error {
    data, err := io.NormalIO.ReadFile(path)
    if nil != err {
        return errors.ReadFileError{FileName: path, Err: err}
    }
    result := Replace(string(data), values)
    err = io.NormalIO.WriteFile(path, []byte(result))
    if nil != err {
        return errors.WriteFileError{FileName: path, Err: err}
    }
    return nil
}

func Replace(template string, values map[string]string) string {
    for match, replace := range values {
        template = strings.Replace(template, "{{"+match+"}}", replace, -1)
    }
    return template
}
