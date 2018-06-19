package cmake

import (
    "path/filepath"
    "strings"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/types"
)

// CMake Target information
type CMakeTarget struct {
    TargetName            string
    Path                  string
    Flags                 []string
    Definitions           []string
    FlagsVisibility       string
    DefinitionsVisibility string
    HeaderOnly            bool
}

// CMake Target Link information
type CMakeTargetLink struct {
    From           string
    To             string
    LinkVisibility string
}

// This creates CMake library string that will be used to link libraries
func GenerateAvrDependencyCMakeString(targets map[string]*CMakeTarget, links []CMakeTargetLink) []string {
    cmakeStrings := make([]string, 0)

    for _, target := range targets {
        finalString := avrNonHeaderOnlyString

        if target.HeaderOnly {
            finalString = avrHeaderOnlyString
        }

        finalString = strings.Replace(finalString, "{{DEPENDENCY_NAME}}", target.TargetName, -1)
        finalString = strings.Replace(finalString, "{{DEPENDENCY_PATH}}", target.Path, -1)
        finalString = strings.Replace(finalString, "{{DEPENDENCY_FLAGS}}",
            strings.Join(target.Flags, " "), -1)
        finalString = strings.Replace(finalString, "{{DEPENDENCY_DEFINITIONS}}",
            strings.Join(target.Definitions, " "), -1)
        finalString = strings.Replace(finalString, "{{FLAGS_VISIBILITY}}", target.FlagsVisibility, -1)
        finalString = strings.Replace(finalString, "{{DEFINITIONS_VISIBILITY}}", target.DefinitionsVisibility, -1)

        cmakeStrings = append(cmakeStrings, finalString+"\n")
    }

    for _, link := range links {
        finalString := linkString
        finalString = strings.Replace(finalString, "{{LINKER_NAME}}", link.From, -1)
        finalString = strings.Replace(finalString, "{{DEPENDENCY_NAME}}", link.To, -1)

        finalString = strings.Replace(finalString, "{{LINK_VISIBILITY}}", link.LinkVisibility, -1)

        cmakeStrings = append(cmakeStrings, finalString)
    }

    cmakeStrings = append(cmakeStrings, "")

    return cmakeStrings
}

// THis Creates the main CMakeLists.txt file for AVR app type project
func GenerateAvrMainCMakeLists(appName string, appPath string, board string, port string, framework string,
    targetName string, targetPath string, flags types.TargetFlags, definitions types.TargetDefinitions) error {

    executablePath, err := io.NormalIO.GetRoot()
    if err != nil {
        return err
    }

    toolChainPath := "toolchain/cmake/CosaToolchain.cmake"

    // read the CMakeLists.txt file template
    templateData, err := io.AssetIO.ReadFile("templates/cmake/CMakeListsAVR.txt.tpl")
    if err != nil {
        return err
    }

    templateDataStr := strings.Replace(string(templateData), "{{TOOLCHAIN_PATH}}",
        filepath.ToSlash(executablePath), -1)
    templateDataStr = strings.Replace(templateDataStr, "{{TOOLCHAIN_FILE_REL}}",
        filepath.ToSlash(toolChainPath), -1)
    templateDataStr = strings.Replace(templateDataStr, "{{PROJECT_PATH}}", filepath.ToSlash(appPath), -1)
    templateDataStr = strings.Replace(templateDataStr, "{{PROJECT_NAME}}", appName, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{TARGET_NAME}}", targetName, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{BOARD}}", board, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{PORT}}", port, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{FRAMEWORK}}", strings.Title(framework), -1)
    templateDataStr = strings.Replace(templateDataStr, "{{ENTRY}}", targetPath, 1)
    templateDataStr = strings.Replace(templateDataStr, "{{TARGET_COMPILE_FLAGS}}",
        strings.Join(flags.GetTargetFlags(), " "), -1)
    templateDataStr = strings.Replace(templateDataStr, "{{TARGET_COMPILE_DEFINITIONS}}",
        strings.Join(definitions.GetTargetDefinitions(), " "), -1)

    templateDataStr += "\n\ninclude(${DEPENDENCY_FILE})\n"

    return io.NormalIO.WriteFile(appPath+io.Sep+".wio"+io.Sep+"build"+io.Sep+"CMakeLists.txt",
        []byte(templateDataStr))
}
