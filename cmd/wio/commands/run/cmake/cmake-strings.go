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
const AvrHeader = `
add_library({{DEPENDENCY_NAME}} INTERFACE)

set_property(TARGET {{DEPENDENCY_NAME}} PROPERTY CXX_STANDARD {{CXX_STANDARD}})
set_property(TARGET {{DEPENDENCY_NAME}} PROPERTY C_STANDARD {{C_STANDARD}})

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    INTERFACE
    {{PRIVATE_DEFINITIONS}})

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    INTERFACE
    {{PUBLIC_DEFINITIONS}})

target_compile_options(
    {{DEPENDENCY_NAME}}
    INTERFACE
    {{DEPENDENCY_FLAGS}})

target_include_directories(
    {{DEPENDENCY_NAME}}
    INTERFACE
    "{{DEPENDENCY_PATH}}/include")
`

const AvrLibrary = `
file(GLOB_RECURSE
    {{DEPENDENCY_NAME}}_files
    "{{DEPENDENCY_PATH}}/src/*.cpp"
    "{{DEPENDENCY_PATH}}/src/*.cc"
    "{{DEPENDENCY_PATH}}/src/*.c")

generate_arduino_library(
    {{DEPENDENCY_NAME}}
    SRCS ${{{DEPENDENCY_NAME}}_files}
    BOARD ${BOARD})

set_property(TARGET {{DEPENDENCY_NAME}} PROPERTY CXX_STANDARD {{CXX_STANDARD}})
set_property(TARGET {{DEPENDENCY_NAME}} PROPERTY C_STANDARD {{C_STANDARD}})

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    PRIVATE
    {{PRIVATE_DEFINITIONS}})

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    PUBLIC
    {{PUBLIC_DEFINITIONS}})

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    PRIVATE
    WIO_PLATFORM_${PLATFORM}
    WIO_FRAMEWORK_${FRAMEWORK}
    WIO_BOARD_${BOARD})

target_compile_options(
    {{DEPENDENCY_NAME}}
    PUBLIC
    {{DEPENDENCY_FLAGS}})

target_include_directories(
    {{DEPENDENCY_NAME}}
    PUBLIC
    "{{DEPENDENCY_PATH}}/include")

target_include_directories(
    {{DEPENDENCY_NAME}}
    PRIVATE
    "{{DEPENDENCY_PATH}}/src")
`

// This for header only desktop dependency
const DesktopHeader = `
add_library({{DEPENDENCY_NAME}} INTERFACE)

set_property(TARGET {{DEPENDENCY_NAME}} PROPERTY CXX_STANDARD {{CXX_STANDARD}})
set_property(TARGET {{DEPENDENCY_NAME}} PROPERTY C_STANDARD {{C_STANDARD}})

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    INTERFACE
    {{PRIVATE_DEFINITIONS_DEFINITIONS}})

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    INTERFACE
    {{PUBLIC_DEFINITIONS_DEFINITIONS}})

target_compile_options(
    {{DEPENDENCY_NAME}}
    INTERFACE
    {{DEPENDENCY_FLAGS}})

target_include_directories(
    {{DEPENDENCY_NAME}}
    INTERFACE
    "{{DEPENDENCY_PATH}}/include")
`

const DesktopLibrary = `
file(GLOB_RECURSE
    {{DEPENDENCY_NAME}}_files
    "{{DEPENDENCY_PATH}}/src/*.cpp"
    "{{DEPENDENCY_PATH}}/src/*.cc"
    "{{DEPENDENCY_PATH}}/src/*.c")

add_library(
    {{DEPENDENCY_NAME}}
    STATIC
    ${{{DEPENDENCY_NAME}}_files})

set_property(TARGET {{DEPENDENCY_NAME}} PROPERTY CXX_STANDARD {{CXX_STANDARD}})
set_property(TARGET {{DEPENDENCY_NAME}} PROPERTY C_STANDARD {{C_STANDARD}})

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    PRIVATE
    {{PRIVATE_DEFINITIONS}})

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    PUBLIC
    {{PUBLIC_DEFINITIONS}})

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    PRIVATE
    WIO_PLATFORM_${PLATFORM}
    WIO_FRAMEWORK_${FRAMEWORK}
    WIO_OS_${OS})

target_compile_options(
    {{DEPENDENCY_NAME}}
    PUBLIC
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
const LinkString = `
target_link_libraries({{LINKER_NAME}} {{LINK_VISIBILITY}} {{DEPENDENCY_NAME}})`
