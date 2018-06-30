package cmake

//////////////////////////////////////////////// Dependencies ////////////////////////////////////////

const targetIncludeDirectories = `target_include_directories({{TARGET}} {{VISIBILITY}} {{DIRECTORIES}})`
const targetCompileDefinitions = `target_compile_definitions({{TARGET}} {{VISIBILITY}} {{DEFINITIONS}})`
const targetCompileOptions = `target_compile_options({{TARGET}} {{VISIBILITY}} {{FLAGS}})`
const targetLinkLibraries = `target_link_libraries({{TARGET}} {{VISIBILITY}} {{LIBRARIES}})`

const file = `file({{TYPE}} {{VARIABLE}} {{ARGS}})`
const addLibrary = `add_library({{LIBRARY}} {{TYPE}} {{FILES}})`
const generateArduinoLibrary = `generate_arduino_library({{LIBRARY}} SRCS {{FILES}} BOARD {{BOARD}})`

const addExecutable = `add_executable({{TARGET}} {{FILES}})`
const generateArduinoFirmware = `genreate_arduino_firmware({{TARGET}} SRCS {{FILES}} BOARD {{BOARD}} PORT {{PORT}})`

// This for header only AVR dependency
const avrHeader = `
add_library({{DEPENDENCY_NAME}} INTERFACE)

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    {{DEFINITIONS_VISIBILITY}}
    {{DEPENDENCY_DEFINITIONS}})

target_compile_options(
    {{DEPENDENCY_NAME}}
    {{FLAGS_VISIBILITY}}
    {{DEPENDENCY_FLAGS}})

target_include_directories(
    {{DEPENDENCY_NAME}}
    INTERFACE
    "{{DEPENDENCY_PATH}}/include")
`

// This is for header only AVR dependency
const avrLibrary = `
file(GLOB_RECURSE
    {{DEPENDENCY_NAME}}_files 
    "{{DEPENDENCY_PATH}}/src/*.cpp" 
    "{{DEPENDENCY_PATH}}/src/*.cc"
    "{{DEPENDENCY_PATH}}/src/*.c")

generate_arduino_library(
    {{DEPENDENCY_NAME}}
	SRCS ${{{DEPENDENCY_NAME}}_files}
	BOARD ${BOARD})

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    {{DEFINITIONS_VISIBILITY}} 
    {{DEPENDENCY_DEFINITIONS}})

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    PRIVATE
    WIO_PLATFORM_${PLATFORM}
    WIO_FRAMEWORK_${FRAMEWORK}
    WIO_BOARD_${BOARD})

target_compile_options(
    {{DEPENDENCY_NAME}}
    {{FLAGS_VISIBILITY}} 
    {{DEPENDENCY_FLAGS}})

target_include_directories(
    {{DEPENDENCY_NAME}}
    PUBLIC
    "{{DEPENDENCY_PATH}}/include")

target_include_directories(
    {{DEPENDENCY_NAME}}
    PRIVATE
    "{{DEPENDENCY_PATH}}/src")

target_include_directories(
    {{DEPENDENCY_NAME}}
    PUBLIC
    "{{DEPENDENCY_PATH}}/include")
`

// This for header only desktop dependency
const desktopHeader = `
add_library({{DEPENDENCY_NAME}} INTERFACE)

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    {{DEFINITIONS_VISIBILITY}} 
    {{DEPENDENCY_DEFINITIONS}})

target_compile_options(
    {{DEPENDENCY_NAME}}
    {{FLAGS_VISIBILITY}} 
    {{DEPENDENCY_FLAGS}})

target_include_directories(
    {{DEPENDENCY_NAME}}
    INTERFACE
    "{{DEPENDENCY_PATH}}/include")
`

// This is for header only desktop dependency
const desktopLibrary = `
file(GLOB_RECURSE 
    {{DEPENDENCY_NAME}}_files
    "{{DEPENDENCY_PATH}}/src/*.cpp"
    "{{DEPENDENCY_PATH}}/src/*.cc"
    "{{DEPENDENCY_PATH}}/src/*.c")

add_library(
    {{DEPENDENCY_NAME}}
    STATIC
    ${{{DEPENDENCY_NAME}}_files})

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    {{DEFINITIONS_VISIBILITY}}
    {{DEPENDENCY_DEFINITIONS}})

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    PRIVATE
    WIO_PLATFORM_${PLATFORM}
    WIO_FRAMEWORK_${FRAMEWORK}
    WIO_BOARD_${BOARD})

target_compile_options(
    {{DEPENDENCY_NAME}}
    {{FLAGS_VISIBILITY}}
    {{DEPENDENCY_FLAGS}})

target_include_directories(
    {{DEPENDENCY_NAME}}
    PRIVATE
    "{{DEPENDENCY_PATH}}/src")

target_include_directories(
    {{DEPENDENCY_NAME}}
    PUBLIC
    "{{DEPENDENCY_PATH}}/include")
`

/////////////////////////////////////////////// Linking ////////////////////////////////////////////

// This is for linking dependencies
const linkString = `
target_link_libraries({{LINKER_NAME}} {{LINK_VISIBILITY}} {{DEPENDENCY_NAME}})

`
