package types

import (
    "bufio"
    "fmt"
    "regexp"
    "strings"
    "time"
    "wio/internal/config/meta"
    "wio/internal/constants"
    "wio/pkg/npm/semver"
    "wio/pkg/util"
    "wio/pkg/util/sys"

    "gopkg.in/yaml.v2"
)

var configCategoryTag = regexp.MustCompile(`^[\s\w]*[a-z]:\s*$`)

func resolveVariables(line string, projectPath string) string {
    return strings.Replace(line, "$(PROJECT_PATH)", projectPath, -1)
}

// resolves operating system specific flags and definitions
func resolveOSSpecific(line string) (string, error) {
    osName := sys.GetOS()

    const windowsStr = "$windows("
    const darwinStr = "$darwin("
    const linuxStr = "$linux("

    diffOsSystemMatch := regexp.MustCompile(
        strings.Replace(`\$(darwin|windows|linux)\(([^,])*\)`, osName, "", 1))

    i := 0
    for end := time.Now().Add(time.Second); ; {
        if i&0x0f == 0 { // Check in every 16th iteration
            if time.Now().After(end) {
                return "", util.Error("error occured while fulfulling OS variables")
            }
        }

        if strings.Contains(line, windowsStr) || strings.Contains(line, darwinStr) || strings.Contains(line, linuxStr) {
            line = diffOsSystemMatch.ReplaceAllString(line, "")
            line = strings.Replace(line, ",", "", -1)
            line = strings.TrimSuffix(line, " ")

            if strings.Contains(line, fmt.Sprintf("$%s(", osName)) {
                line = strings.Replace(line, "$"+osName+"(", "", 1)
                line = line[:len(line)-1]
            }
        } else {
            break
        }

        i++
    }

    return line, nil
}

func fulfillVariables(data []byte, projectPath string) ([]byte, error) {
    scanner := bufio.NewScanner(strings.NewReader(string(data)))
    newStr := ""

    for scanner.Scan() {
        if configCategoryTag.MatchString(scanner.Text()) {
            newStr += scanner.Text() + "\n"
            continue
        }

        // resolve operating system constrains
        osStr, err := resolveOSSpecific(scanner.Text())
        if err != nil {
            return nil, err
        }

        // resolve PROJECT_PATH variable
        variableStr := resolveVariables(osStr, projectPath)

        newStr += variableStr + "\n"
    }

    return []byte(newStr), nil
}

func ReadWioConfig(dir string, fulfill bool) (Config, error) {
    path := sys.Path(dir, sys.Config)
    if !sys.Exists(path) {
        return nil, util.Error("path does not contain a wio.yml: %s", dir)
    }

    text, err := sys.NormalIO.ReadFile(path)
    if err != nil {
        return nil, err
    }

    if fulfill {
        text, err = fulfillVariables(text, dir)
        if err != nil {
            return nil, err
        }
    }

    ret := &ConfigImpl{}
    if err = yaml.Unmarshal(text, ret); err != nil {
        return nil, err
    }

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
    } else if semver.Parse(ret.GetInfo().GetOptions().GetWioVersion()) == nil {
        return nil, util.Error("%s: wio.yml does not have a valid version", dir)
    }

    currVersion := semver.Parse(meta.Version)
    needVersion := semver.Parse(ret.Info.GetOptions().GetWioVersion())

    if currVersion.Lt(needVersion) {
        return nil, util.Error("%s: current wio version is %s but wio.yml needs it >= %s", dir,
            currVersion.Str(), needVersion.Str())
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
