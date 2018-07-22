package utils

import (
    "bufio"
    "regexp"
    "strings"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"

    "gopkg.in/yaml.v2"
)

func ReadWioConfig(dir string) (types.Config, error) {
    path := io.Path(dir, io.Config)
    if !io.Exists(path) {
        return nil, errors.Stringf("path does not contain a wio.yml: %s", dir)
    }
    ret := &types.ConfigImpl{}
    err := io.NormalIO.ParseYml(path, ret)
    return ret, err
}

func WriteWioConfig(dir string, config types.Config) error {
    return prettyPrintHelp(config, io.Path(dir, io.Config))
}

// Write configuration with nice spacing and information
func prettyPrintHelp(config types.Config, filePath string) error {
    var ymlData []byte
    var err error

    // marshall yml data
    if ymlData, err = yaml.Marshal(config); err != nil {
        marshallError := errors.YamlMarshallError{
            Err: err,
        }
        return marshallError
    }

    finalStr := ""

    // configuration tags
    projectTagPat := regexp.MustCompile(`(^project:)|((\s| |^\w)project:(\s+|))`)
    targetsTagPat := regexp.MustCompile(`(^targets:)|((\s| |^\w)targets:(\s+|))`)
    dependenciesTagPat := regexp.MustCompile(`(^dependencies:)|((\s| |^\w)dependencies:(\s+|))`)

    scanner := bufio.NewScanner(strings.NewReader(string(ymlData)))
    for scanner.Scan() {
        line := scanner.Text()

        if projectTagPat.MatchString(line) || targetsTagPat.MatchString(line) || dependenciesTagPat.MatchString(line) {
            finalStr += "\n" + line
        } else {
            finalStr += line
        }

        finalStr += "\n"
    }

    if err = io.NormalIO.WriteFile(filePath, []byte(finalStr)); err != nil {
        return errors.WriteFileError{
            FileName: filePath,
            Err:      err,
        }
    }

    return nil
}
