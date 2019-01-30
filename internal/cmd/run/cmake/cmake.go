package cmake

import (
    "os"
    "path/filepath"
    "strings"
    "wio/internal/constants"
    "wio/internal/types"
    "wio/internal/utils"
    "wio/pkg/downloader"
    "wio/pkg/util"
    "wio/pkg/util/sys"
    "wio/pkg/util/template"
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
            return "", "", util.Error("invalid ISO C/C++ standard [%s]", std)
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

func generateCmakeLists(templateFile string, buildPath string, values map[string]string) error {
    templatePath := sys.Path("templates", "cmake", templateFile+".txt.tpl")
    cmakeListsPath := sys.Path(buildPath, "CMakeLists.txt")
    if err := os.MkdirAll(buildPath, os.ModePerm); err != nil {
        return err
    }
    if err := sys.AssetIO.CopyFile(templatePath, cmakeListsPath, true); err != nil {
        return err
    }
    return template.IOReplace(cmakeListsPath, values)
}

var templateFiles = map[string]string{
    constants.Avr:    "CMakeListsAVR",
    constants.Native: "CMakeListsNative",
}

func GenerateCmakeLists(
    toolchainPath string,
    target types.Target,
    projectName string,
    projectPath string,
    cppStandard string,
    cStandard string) error {

    flags := target.GetFlags().GetTarget()
    definitions := target.GetDefinitions().GetTarget()
    buildPath := sys.Path(utils.BuildPath(projectPath), target.GetName())

    moduleData := &downloader.ModuleData{}
    packageInfoPath := sys.Path(toolchainPath, "package.json")
    if sys.Exists(packageInfoPath) {
        if err := sys.NormalIO.ParseJson(packageInfoPath, moduleData); err != nil {
            return err
        }
    }

    return generateCmakeLists(templateFiles[target.GetPlatform()], buildPath, map[string]string{
        "TOOLCHAIN_FILE":             sys.Path(filepath.ToSlash(toolchainPath), moduleData.ToolchainFile),
        "PROJECT_PATH":               filepath.ToSlash(projectPath),
        "PROJECT_NAME":               projectName,
        "CPP_STANDARD":               cppStandard,
        "C_STANDARD":                 cStandard,
        "PLATFORM":                   strings.ToUpper(target.GetPlatform()),
        "FRAMEWORK":                  strings.ToUpper(target.GetFramework()),
        "TARGET_NAME":                target.GetName(),
        "ENTRY":                      target.GetSource(),
        "TARGET_COMPILE_FLAGS":       strings.Join(flags, " "),
        "TARGET_COMPILE_DEFINITIONS": strings.Join(definitions, " "),
        "TARGET_LINK_LIBRARIES": func() string {
            if len(target.GetLinkerFlags()) > 0 {
                return "\n" + template.Replace(LinkString, map[string]string{
                    "LINK_FROM":       "${TARGET_NAME}",
                    "LINK_VISIBILITY": "PRIVATE",
                    "LINK_TO":         "# no dep",
                    "LINKER_FLAGS":    strings.Join(target.GetLinkerFlags(), " "),
                })
            } else {
                return ""
            }
        }(),
    })
}
