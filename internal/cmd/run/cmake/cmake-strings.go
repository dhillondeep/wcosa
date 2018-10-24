package cmake

//////////////////////////////////////////////// Dependencies ////////////////////////////////////////

// This for header only AVR dependency
const AvrHeader = `add_library({{DEPENDENCY_NAME}} INTERFACE)

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    INTERFACE
    {{PRIVATE_DEFINITIONS}}
)

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    INTERFACE
    {{PUBLIC_DEFINITIONS}}
)

target_compile_options(
    {{DEPENDENCY_NAME}}
    INTERFACE
    {{DEPENDENCY_FLAGS}}
)

target_include_directories(
    {{DEPENDENCY_NAME}}
    INTERFACE
    "{{DEPENDENCY_PATH}}/include"
)
`

const AvrLibrary = `file(GLOB_RECURSE
    {{DEPENDENCY_NAME}}_files
    "{{DEPENDENCY_PATH}}/src/*.cpp"
    "{{DEPENDENCY_PATH}}/src/*.cc"
    "{{DEPENDENCY_PATH}}/src/*.c")

generate_arduino_library(
    {{DEPENDENCY_NAME}}
    SRCS ${{{DEPENDENCY_NAME}}_files}
    BOARD ${BOARD})

set_target_properties({{DEPENDENCY_NAME}} PROPERTIES CXX_STANDARD {{CXX_STANDARD}})
set_target_properties({{DEPENDENCY_NAME}} PROPERTIES C_STANDARD {{C_STANDARD}})

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    PRIVATE
    {{PRIVATE_DEFINITIONS}}
)

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    PUBLIC
    {{PUBLIC_DEFINITIONS}}
)

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    PRIVATE
    WIO_PLATFORM_${PLATFORM}
    WIO_FRAMEWORK_${FRAMEWORK}
    WIO_BOARD_${BOARD}
)

target_compile_options(
    {{DEPENDENCY_NAME}}
    PUBLIC
    {{DEPENDENCY_FLAGS}}
)

target_include_directories(
    {{DEPENDENCY_NAME}}
    PUBLIC
    "{{DEPENDENCY_PATH}}/include"
)

target_include_directories(
    {{DEPENDENCY_NAME}}
    PRIVATE
    "{{DEPENDENCY_PATH}}/src"
)
`

// This for header only desktop dependency
const DesktopHeader = `add_library({{DEPENDENCY_NAME}} INTERFACE)

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    INTERFACE
    {{PRIVATE_DEFINITIONS}}
)

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    INTERFACE
    {{PUBLIC_DEFINITIONS}}
)

target_compile_options(
    {{DEPENDENCY_NAME}}
    INTERFACE
    {{DEPENDENCY_FLAGS}}
)

target_include_directories(
    {{DEPENDENCY_NAME}}
    INTERFACE
    "{{DEPENDENCY_PATH}}/include"
)
`

const DesktopLibrary = `file(GLOB_RECURSE
    {{DEPENDENCY_NAME}}_files
    "{{DEPENDENCY_PATH}}/src/*.cpp"
    "{{DEPENDENCY_PATH}}/src/*.cc"
    "{{DEPENDENCY_PATH}}/src/*.c")

add_library(
    {{DEPENDENCY_NAME}}
    STATIC
    ${{{DEPENDENCY_NAME}}_files}
)

set_target_properties({{DEPENDENCY_NAME}} PROPERTIES CXX_STANDARD {{CXX_STANDARD}})
set_target_properties({{DEPENDENCY_NAME}} PROPERTIES C_STANDARD {{C_STANDARD}})

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    PRIVATE
    {{PRIVATE_DEFINITIONS}}
)

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    PUBLIC
    {{PUBLIC_DEFINITIONS}}
)

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    PRIVATE
    WIO_PLATFORM_${PLATFORM}
    WIO_FRAMEWORK_${FRAMEWORK}
    WIO_OS_${OS}
)

target_compile_options(
    {{DEPENDENCY_NAME}}
    PUBLIC
    {{DEPENDENCY_FLAGS}}
)

target_include_directories(
    {{DEPENDENCY_NAME}}
    PRIVATE
    "{{DEPENDENCY_PATH}}/src"
)

target_include_directories(
    {{DEPENDENCY_NAME}}
    PUBLIC
    "{{DEPENDENCY_PATH}}/include"
)
`

const LibraryFind = `find_library(
    {{LIB_NAME_VAR}}
    {{LIB_NAME}}
    {{LIB_HINTS}}
    {{LIB_PATHS}}
    {{LIB_REQUIRED}}
)

if ({{LIB_NAME_VAR}}-NOTFOUND)
    message(FATAL_ERROR "{{LIB_NAME}} library not found")
endif()
`

const LibraryInclude = `target_include_directories(
    {{LIB_NAME_VAR}}
    PRIVATE
    "{{LIB_INCLUDE_PATHS}}"
)
`

const LibraryPackageFind = `find_package(
    {{LIB_NAME}}
    {{LIB_VERSION}}
    {{LIB_HINTS}}
    {{LIB_PATHS}}
    COMPONENTS {{LIB_REQUIRED_COMPONENTS}}
    OPTIONAL_COMPONENTS {{LIB_OPTIONAL_COMPONENTS}}
    {{LIB_REQUIRED}}
)

if (NOT {{LIB_NAME}}_FOUND)
    message(FATAL_ERROR "{{LIB_NAME}} library not found")
endif()
`

/////////////////////////////////////////////// Linking ////////////////////////////////////////////

// This is for linking dependencies
const LinkString = `target_link_libraries(
    {{LINK_FROM}}
    {{LINK_VISIBILITY}}
    {{LINK_TO}}
    {{LINKER_FLAGS}}
)
`
