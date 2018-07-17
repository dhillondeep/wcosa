package cmake

import (
    "os"
    "path/filepath"
    "strings"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/utils/template"

    "github.com/thoas/go-funk"
)

var validCppStandard = map[string]string{
    "c++98": "98",
    "c++03": "98",
    "c++11": "11",
    "c++14": "14",
    "c++17": "17",
    "c++2a": "20",
    "c++20": "20",
}

var validCStandard = map[string]string{
    "iso9899:1990": "90",
    "c90":          "90",
    "iso9899:1999": "99",
    "c99":          "99",
    "iso9899:2011": "11",
    "c11":          "11",
    "iso9899:2017": "11",
    "c17":          "11",
    "gnu90":        "90",
    "gnu99":        "99",
    "gnu11":        "11",
}

// Reads standard provided by user and returns CMake standard
func GetStandard(target *types.Target) (string, string, error) {
    standard := strings.ToLower(strings.Trim((*target).GetStandard(), " "))
    var cppStandard string
    var cStandard string

    for _, currStandard := range strings.Split(standard, ",") {
        currStandard = strings.Trim(currStandard, " ")

        if funk.Contains(validCppStandard, currStandard) {
            cppStandard = validCppStandard[currStandard]
            cStandard = validCStandard["c11"]
        } else if funk.Contains(validCppStandard, currStandard) {
            cStandard = validCStandard[currStandard]
            cppStandard = validCppStandard["c++11"]
        } else {
            return "", "", errors.Stringf("[%s] standard is not valid for ISO standard", currStandard)
        }
    }

    return cppStandard, cStandard, nil
}

func BuildPath(projectPath string) string {
    return projectPath + io.Sep + io.Folder + io.Sep + constants.TargetDir
}

func generateCmakeLists(
    templateFile string,
    buildPath string,
    values map[string]string) error {

    templatePath := "templates/cmake/" + templateFile + ".txt.tpl"
    cmakeListsPath := buildPath + io.Sep + "CMakeLists.txt"
    if err := os.MkdirAll(buildPath, os.ModePerm); err != nil {
        return err
    }
    if err := io.AssetIO.CopyFile(templatePath, cmakeListsPath, true); err != nil {
        return err
    }
    return template.IOReplace(cmakeListsPath, values)
}

// This creates the main CMakeLists.txt file for AVR app type project
func GenerateAvrCmakeLists(
    toolchainPath string,
    target *types.Target,
    projectName string,
    projectPath string,
    port string) error {

    cppStandard, cStandard, err := GetStandard(target)
    if err != nil {
        return err
    }

    flags := (*target).GetFlags().GetTargetFlags()
    definitions := (*target).GetDefinitions().GetTargetDefinitions()
    framework := (*target).GetFramework()
    buildPath := BuildPath(projectPath) + io.Sep + (*target).GetName()
    templateFile := "CMakeListsAVR"
    executablePath, err := io.NormalIO.GetRoot()
    if err != nil {
        return err
    }

    return generateCmakeLists(templateFile, buildPath, map[string]string{
        "TOOLCHAIN_PATH":             filepath.ToSlash(executablePath),
        "TOOLCHAIN_FILE_REL":         filepath.ToSlash(toolchainPath),
        "PROJECT_PATH":               filepath.ToSlash(projectPath),
        "PROJECT_NAME":               projectName,
        "CPP_STANDARD":               cppStandard,
        "C_STANDARD":                 cStandard,
        "FRAMEWORK":                  framework,
        "PORT":                       port,
        "PLATFORM":                   constants.AVR,
        "TARGET_NAME":                (*target).GetName(),
        "BOARD":                      (*target).GetBoard(),
        "ENTRY":                      (*target).GetSrc(),
        "TARGET_COMPILE_FLAGS":       strings.Join(flags, " "),
        "TARGET_COMPILE_DEFINITIONS": strings.Join(definitions, " "),
    })
}

func GenerateNativeCmakeLists(
    target *types.Target,
    projectName string,
    projectPath string) error {

    cppStandard, cStandard, err := GetStandard(target)
    if err != nil {
        return err
    }

    flags := (*target).GetFlags().GetTargetFlags()
    definitions := (*target).GetDefinitions().GetTargetDefinitions()
    buildPath := BuildPath(projectPath) + io.Sep + (*target).GetName()
    templateFile := "CMakeListsNative"

    return generateCmakeLists(templateFile, buildPath, map[string]string{
        "PROJECT_PATH":               filepath.ToSlash(projectPath),
        "PROJECT_NAME":               projectName,
        "CPP_STANDARD":               cppStandard,
        "C_STANDARD":                 cStandard,
        "TARGET_NAME":                (*target).GetName(),
        "FRAMEWORK":                  (*target).GetFramework(),
        "HARDWARE":                   (*target).GetBoard(),
        "ENTRY":                      (*target).GetSrc(),
        "PLATFORM":                   constants.NATIVE,
        "TARGET_COMPILE_FLAGS":       strings.Join(flags, " "),
        "TARGET_COMPILE_DEFINITIONS": strings.Join(definitions, " "),
    })
}
