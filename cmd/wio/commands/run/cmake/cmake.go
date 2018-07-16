package cmake

import (
    "os"
    "path/filepath"
    "strings"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/utils/template"
)

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
    framework string,
    cppStandard string,
    cStandard string,
    projectName string,
    projectPath string) error {

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
        "FRAMEWORK":                  framework,
        "HARDWARE":                   (*target).GetBoard(),
        "ENTRY":                      (*target).GetSrc(),
        "PLATFORM":                   constants.NATIVE,
        "TARGET_COMPILE_FLAGS":       strings.Join(flags, " "),
        "TARGET_COMPILE_DEFINITIONS": strings.Join(definitions, " "),
    })
}
