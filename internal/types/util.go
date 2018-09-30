package types

import (
    "bufio"
    "regexp"
    "strings"
    "wio/internal/constants"
    "wio/pkg/util"
    "wio/pkg/util/sys"

    "gopkg.in/yaml.v2"
)

func ReadWioConfig(dir string) (Config, error) {
    path := sys.Path(dir, sys.Config)
    if !sys.Exists(path) {
        return nil, util.Error("path does not contain a wio.yml: %s", dir)
    }
    ret := &ConfigImpl{}
    err := sys.NormalIO.ParseYml(path, ret)

    // check if it is a valid config
    if ret.GetType() != constants.App && ret.GetType() != constants.Pkg {
        return nil, util.Error("wio.yml has invalid type: %s", ret.GetType())
    }
    return ret, err
}

func WriteWioConfig(dir string, config Config) error {
    return prettyPrintHelp(config, sys.Path(dir, sys.Config))
}

// Write configuration with nice spacing and information
func prettyPrintHelp(config Config, filePath string) error {
    var ymlData []byte
    var err error

    // marshall yml data
    if ymlData, err = yaml.Marshal(config); err != nil {
        return err
    }

    finalStr := ""

    // configuration tags
    projectTagPat := regexp.MustCompile(`(^project:)|((\s| |^\w)project:(\s+|))`)
    targetsTagPat := regexp.MustCompile(`(^targets:)|((\s| |^\w)targets:(\s+|))`)
    dependenciesTagPat := regexp.MustCompile(`(^dependencies:)|((\s| |^\w)dependencies:(\s+|))`)
    librariesTagPat := regexp.MustCompile(`(^libraries:)|((\s| |^\w)libraries:(\s+|))`)

    scanner := bufio.NewScanner(strings.NewReader(string(ymlData)))
    for scanner.Scan() {
        line := scanner.Text()

        if projectTagPat.MatchString(line) || targetsTagPat.MatchString(line) ||
            dependenciesTagPat.MatchString(line) || librariesTagPat.MatchString(line) {
            finalStr += "\n" + line
        } else {
            finalStr += line
        }

        finalStr += "\n"
    }

    finalStr += "\n"

    if err = sys.NormalIO.WriteFile(filePath, []byte(finalStr)); err != nil {
        return err
    }

    return nil
}
