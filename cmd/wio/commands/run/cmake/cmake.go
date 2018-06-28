package cmake

import (
    "path/filepath"
    "strings"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/utils"
    "os"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/utils/template"
)

// This creates the main CMakeLists.txt file for AVR app type project
func GenerateAvrCmakeLists(
    target *types.Target,
    projectName string,
    projectPath string,
    port string) error {

    flags := (*target).GetFlags().GetTargetFlags()
    definitions := (*target).GetDefinitions().GetTargetDefinitions()
    framework := (*target).GetFramework()

    executablePath, err := io.NormalIO.GetRoot()
    if err != nil {
        return err
    }

    var toolChainPath string
    if framework == constants.COSA {
        toolChainPath = "toolchain/cmake/CosaToolchain.cmake"
    } else {
        return errors.FrameworkNotSupportedError{
            Platform:  constants.AVR,
            Framework: framework,
        }
    }

    buildPath := projectPath + io.Sep + ".wio" + io.Sep + "build"
    templatePath := "templates/cmake/CMakeListsAVR.txt.tpl"
    cmakeListsPath := buildPath + io.Sep + "CMakeLists.txt"
    if err := os.MkdirAll(buildPath, os.ModePerm); err != nil {
        return err
    }
    if err := utils.CopyFile(templatePath, cmakeListsPath); err != nil {
        return err
    }
    return template.IOReplace(cmakeListsPath, map[string]string{
        "TOOLCHAIN_PATH":             filepath.ToSlash(executablePath),
        "TOOLCHAIN_FILE_REL":         filepath.ToSlash(toolChainPath),
        "PROJECT_PATH":               filepath.ToSlash(projectPath),
        "PROJECT_NAME":               projectName,
        "FRAMEWORK":                  framework,
        "PORT":                       port,
        "PLATFORM":                   "avr",
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

    flags := (*target).GetFlags().GetTargetFlags()
    definitions := (*target).GetDefinitions().GetTargetDefinitions()

    buildPath := projectPath + io.Sep + ".wio" + io.Sep + "build"
    templatePath := "templates/cmake/CMakeListsNative.txt.tpl"
    cmakeListsPath := buildPath + io.Sep + "CMakeLists.txt"
    if err := os.MkdirAll(buildPath, os.ModePerm); err != nil {
        return err
    }
    if err := utils.CopyFile(templatePath, cmakeListsPath); err != nil {
        return err
    }
    return template.IOReplace(cmakeListsPath, map[string]string{
        "PROJECT_PATH":               filepath.ToSlash(projectPath),
        "PROJECT_NAME":               projectName,
        "TARGET_NAME":                (*target).GetName(),
        "FRAMEWORK":                  (*target).GetFramework(),
        "BOARD":                      (*target).GetBoard(),
        "ENTRY":                      (*target).GetSrc(),
        "PLATFORM":                   "native",
        "TARGET_COMPILE_FLAGS":       strings.Join(flags, " "),
        "TARGET_COMPILE_DEFINITIONS": strings.Join(definitions, " "),
    })
}
