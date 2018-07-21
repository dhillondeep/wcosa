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
)

var cppStandards = map[string]string{
    "c++98": "98",
    "c++03": "98",
    "c++0x": "11",
    "c++11": "11",
    "c++14": "14",
    "c++17": "17",
    "c++2a": "20",
    "c++20": "20",
}

var cStandards = map[string]string{
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
func GetStandard(standard string) (string, string, error) {
    stdString := strings.Trim(standard, " ")
    stdString = strings.ToLower(stdString)
    cppStandard := ""
    cStandard := ""

    for _, std := range strings.Split(stdString, ",") {
        std = strings.Trim(std, " ")
        if val, exists := cppStandards[std]; exists {
            cppStandard = val
        } else if val, exists := cStandards[std]; exists {
            cStandard = val
        } else if val != "" {
            return "", "", errors.Stringf("invalid ISO C/C++ standard [%s]", std)
        }
    }
    if cppStandard == "" {
        cppStandard = "11"
    }
    if cStandard == "" {
        cStandard = "99"
    }
    return cppStandard, cStandard, nil
}

func BuildPath(projectPath string) string {
    return projectPath + io.Sep + io.Folder + io.Sep + constants.TargetDir
}

func generateCmakeLists(templateFile string, buildPath string, values map[string]string) error {
    templatePath := io.Path("templates", "cmake", templateFile+".txt.tpl")
    cmakeListsPath := io.Path(buildPath, "CMakeLists.txt")
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
    targetName string,
    target types.Target,
    projectName string,
    projectPath string,
    cppStandard string,
    cStandard string,
    port string) error {

    flags := target.GetFlags().GetTarget()
    definitions := target.GetDefinitions().GetTarget()
    framework := target.GetFramework()
    buildPath := io.Path(BuildPath(projectPath), targetName)
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
        "PORT":                       port,
        "PLATFORM":                   strings.ToUpper(constants.AVR),
        "FRAMEWORK":                  strings.ToUpper(framework),
        "BOARD":                      target.GetBoard(),
        "TARGET_NAME":                targetName,
        "ENTRY":                      target.GetSource(),
        "TARGET_COMPILE_FLAGS":       strings.Join(flags, " "),
        "TARGET_COMPILE_DEFINITIONS": strings.Join(definitions, " "),
    })
}

func GenerateNativeCmakeLists(
    targetName string,
    target types.Target,
    projectName string,
    projectPath string,
    cppStandard string,
    cStandard string) error {

    flags := target.GetFlags().GetTarget()
    definitions := target.GetDefinitions().GetTarget()
    buildPath := io.Path(BuildPath(projectPath), targetName)
    templateFile := "CMakeListsNative"

    return generateCmakeLists(templateFile, buildPath, map[string]string{
        "PROJECT_PATH":               filepath.ToSlash(projectPath),
        "PROJECT_NAME":               projectName,
        "CPP_STANDARD":               cppStandard,
        "C_STANDARD":                 cStandard,
        "TARGET_NAME":                targetName,
        "PLATFORM":                   strings.ToUpper(constants.NATIVE),
        "FRAMEWORK":                  strings.ToUpper(target.GetFramework()),
        "OS":                         strings.ToUpper(target.GetBoard()),
        "ENTRY":                      target.GetSource(),
        "TARGET_COMPILE_FLAGS":       strings.Join(flags, " "),
        "TARGET_COMPILE_DEFINITIONS": strings.Join(definitions, " "),
    })
}
