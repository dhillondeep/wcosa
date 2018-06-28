package cmake

//////////////////////////////////////////////// Dependencies ////////////////////////////////////////

// This for header only AVR dependency
const avrHeaderOnlyString = `
add_library({{DEPENDENCY_NAME}} INTERFACE)

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    {{DEFINITIONS_VISIBILITY}}
    WIO_FRAMEWORK_${FRAMEWORK} 
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
const avrNonHeaderOnlyString = `
file(GLOB_RECURSE
    SRC_FILES_{{DEPENDENCY_NAME}} 
    "{{DEPENDENCY_PATH}}/src/*.cpp" 
    "{{DEPENDENCY_PATH}}/src/*.cc"
    "{{DEPENDENCY_PATH}}/src/*.c")

generate_arduino_library(
    {{DEPENDENCY_NAME}}
	SRCS ${SRC_FILES_{{DEPENDENCY_NAME}}}
	BOARD ${BOARD})

target_compile_definitions(
    {{DEPENDENCY_NAME}}
    {{DEFINITIONS_VISIBILITY}} 
    WIO_FRAMEWORK_${FRAMEWORK} 
    {{DEPENDENCY_DEFINITIONS}})

target_compile_options(
    {{DEPENDENCY_NAME}}
    {{FLAGS_VISIBILITY}} 
    {{DEPENDENCY_FLAGS}})

target_include_directories(
    {{DEPENDENCY_NAME}}
    PRIVATE
    "{{DEPENDENCY_PATH}}/include")
`

// This for header only desktop dependency
const desktopHeaderOnlyString = `
add_library({{DEPENDENCY_NAME}} INTERFACE)

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
    INTERFACE
    "{{DEPENDENCY_PATH}}/include")
`

// This is for header only desktop dependency
const desktopNonHeaderOnlyString = `
file(GLOB_RECURSE 
    SRC_FILES__{{DEPENDENCY_NAME}}
    "{{DEPENDENCY_PATH}}/src/*.cpp"
    "{{DEPENDENCY_PATH}}/src/*.cc"
    "{{DEPENDENCY_PATH}}/src/*.c")

add_library(
    {{DEPENDENCY_NAME}}
    STATIC
    ${SRC_FILES_{{DEPENDENCY_NAME}}})

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
`

/////////////////////////////////////////////// Linking ////////////////////////////////////////////

// This is for linking dependencies
const linkString = `target_link_libraries({{LINKER_NAME}} {{LINK_VISIBILITY}} {{DEPENDENCY_NAME}})`
