package dependencies

import "strings"

func dependencyPackageToCMakeString(dependencyName string, dependencyPackage *DependencyPackage) (string) {
    headerOnlyString := `
add_library({{DEPENDENCY_NAME}} INTERFACE)
target_compile_definitions({{DEPENDENCY_NAME}} INTERFACE __AVR_${FRAMEWORK}__ {{DEPENDENCY_FLAGS}})
target_include_directories({{DEPENDENCY_NAME}} INTERFACE "{{DEPENDENCY_PATH}}/include")
`

    nonHeaderOnlyString := `
file(GLOB_RECURSE SRC_FILES "{{DEPENDENCY_PATH}}/src/*.cpp" "{{DEPENDENCY_PATH}}/src/*.cc" "{{DEPENDENCY_PATH}}/src/*.c")
generate_arduino_library({{DEPENDENCY_NAME}}
	SRCS ${SRC_FILES}
	BOARD ${BOARD})
target_compile_definitions({{DEPENDENCY_NAME}} PRIVATE __AVR_${FRAMEWORK}__ {{DEPENDENCY_FLAGS}})
target_include_directories({{DEPENDENCY_NAME}} PUBLIC "{{DEPENDENCY_PATH}}/include")
target_include_directories({{DEPENDENCY_NAME}} PRIVATE "{{DEPENDENCY_PATH}}/src")`

    var finalString string
    if dependencyPackage.MainTag.HeaderOnly {
        finalString = string(headerOnlyString)
    } else {
        finalString = string(nonHeaderOnlyString)
    }

    finalString = strings.Replace(finalString, "{{DEPENDENCY_NAME}}", dependencyName, -1)
    finalString = strings.Replace(finalString, "{{DEPENDENCY_PATH}}", dependencyPackage.Directory, -1)
    finalString = strings.Replace(finalString, "{{DEPENDENCY_FLAGS}}",
        strings.Join(dependencyPackage.MainTag.GlobalFlags, " "), -1)

    return strings.Trim(strings.Trim(finalString, "\n"), " ")
}
