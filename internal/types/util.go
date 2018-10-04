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

    text, err := sys.NormalIO.ReadFile(path)
    if err != nil {
        return nil, err
    }

    ret := &ConfigImpl{}
    err = yaml.Unmarshal(text, ret)

    // Check to give error on older versions of wio
    pkgRegex := regexp.MustCompile(`(\s+|)pkg:`)
    appRegex := regexp.MustCompile(`(\s+|)app:`)
    versionRegex := regexp.MustCompile(`wio_version:(|\s+)0.0.[0-3]`)

    // check if wio config is for older versions
    if pkgRegex.MatchString(string(text)) || appRegex.MatchString(string(text)) ||
        versionRegex.MatchString(string(text)) {
        return nil, util.Error("%s: wio version < 0.4.0 is not supported by "+
            "current wio version", dir)
    }

    // check for wio.yml validity
    if ret.GetType() != constants.App && ret.GetType() != constants.Pkg {
        if strings.Trim(ret.GetType(), " ") == "" {
            return nil, util.Error("%s: wio.yml is invalid or cannot be parsed", dir)
        }

        return nil, util.Error("%s: wio.yml has invalid project type: %s", dir, ret.GetType())
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
