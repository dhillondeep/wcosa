package dependencies

import (
    "path/filepath"
    "strings"
    "wio/cmd/wio/utils/io"
)

const headerOnlyString = `add_library({{DEPENDENCY_NAME}} INTERFACE)
target_compile_definitions({{DEPENDENCY_NAME}} INTERFACE __AVR_${FRAMEWORK}__ {{DEPENDENCY_FLAGS}})
target_include_directories({{DEPENDENCY_NAME}} INTERFACE "{{DEPENDENCY_PATH}}/include")`

const nonHeaderOnlyString = `file(GLOB_RECURSE SRC_FILES "{{DEPENDENCY_PATH}}/src/*.cpp" "{{DEPENDENCY_PATH}}/src/*.cc" "{{DEPENDENCY_PATH}}/src/*.c")
generate_arduino_library({{DEPENDENCY_NAME}}
	SRCS ${SRC_FILES}
	BOARD ${BOARD})
target_compile_definitions({{DEPENDENCY_NAME}} PRIVATE __AVR_${FRAMEWORK}__ {{DEPENDENCY_FLAGS}})
target_include_directories({{DEPENDENCY_NAME}} PUBLIC "{{DEPENDENCY_PATH}}/include")
target_include_directories({{DEPENDENCY_NAME}} PRIVATE "{{DEPENDENCY_PATH}}/src")`

const linkString = `target_link_libraries({{LINKER_NAME}} {{VISIBILITY}} {{DEPENDENCY_NAME}})`

// creates
func generateAvrDependencyCMakeString(targets map[string]*CMakeTarget, links []CMakeTargetLink) []string {
    cmakeStrings := make([]string, 0)

    for _, target := range targets {
        finalString := nonHeaderOnlyString

        if target.HeaderOnly {
            finalString = headerOnlyString
        }

        finalString = strings.Replace(finalString, "{{DEPENDENCY_NAME}}", target.TargetName, -1)
        finalString = strings.Replace(finalString, "{{DEPENDENCY_PATH}}", target.Path, -1)
        finalString = strings.Replace(finalString, "{{DEPENDENCY_FLAGS}}", strings.Join(target.Flags, " "), -1)

        cmakeStrings = append(cmakeStrings, finalString+"\n")
    }

    for _, link := range links {
        finalString := linkString
        finalString = strings.Replace(finalString, "{{LINKER_NAME}}", link.From, -1)
        finalString = strings.Replace(finalString, "{{DEPENDENCY_NAME}}", link.To, -1)

        finalString = strings.Replace(finalString, "{{VISIBILITY}}", link.LinkVisibility, -1)

        cmakeStrings = append(cmakeStrings, finalString)
    }

    cmakeStrings = append(cmakeStrings, "")

    return cmakeStrings
}

// Creates the main CMakeLists.txt file for AVR app type project
func generateAvrMainCMakeLists(appName string, appPath string, board string, port string, framework string, target string,
    flags map[string][]string, isAPP bool) error {

    executablePath, err := io.NormalIO.GetRoot()
    if err != nil {
        return err
    }

    entry := "src"
    if !isAPP {
        entry = "tests"
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
    templateDataStr = strings.Replace(templateDataStr, "{{TARGET_NAME}}", target, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{BOARD}}", board, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{PORT}}", port, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{FRAMEWORK}}", strings.ToUpper(framework), -1)
    templateDataStr = strings.Replace(templateDataStr, "{{ENTRY}}", entry, 1)
    templateDataStr = strings.Replace(templateDataStr, "{{TARGET_COMPILE_FLAGS}}",
        strings.Join(flags["target_compile_flags"], " "), -1)
    templateDataStr += "\n\ninclude(${DEPENDENCY_FILE})\n"

    return io.NormalIO.WriteFile(appPath+io.Sep+".wio"+io.Sep+"build"+io.Sep+"CMakeLists.txt",
        []byte(templateDataStr))
}
