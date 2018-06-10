set(CMAKE_VER 3.0.0)
set(PROJECT_NAME {{PROJECT_NAME}})
set(PROJECT_PATH "{{PROJECT_PATH}}")
set(CMAKE_TOOLCHAIN_PATH "{{TOOLCHAIN_PATH}}")
set (CMAKE_MODULE_PATH "${CMAKE_CURRENT_SOURCE_DIR}")
set(DEPENDENCY_FILE dependencies)

# all the paths toolchain can be at (this is because of different package managers)
if (EXISTS "${CMAKE_TOOLCHAIN_PATH}/{{TOOLCHAIN_FILE_REL}}")
    set(CMAKE_TOOLCHAIN_FILE "${CMAKE_TOOLCHAIN_PATH}/{{TOOLCHAIN_FILE_REL}}")
elseif (EXISTS "${CMAKE_TOOLCHAIN_PATH}/../{{TOOLCHAIN_FILE_REL}}")
    set(CMAKE_TOOLCHAIN_FILE "${CMAKE_TOOLCHAIN_PATH}/../{{TOOLCHAIN_FILE_REL}}")
elseif (EXISTS "/usr/share/wio/{{TOOLCHAIN_FILE_REL}}")
    set(CMAKE_TOOLCHAIN_FILE "/usr/share/wio/{{TOOLCHAIN_FILE_REL}}")
else()
    message(FATAL_ERROR "Toolchain cannot be found. Build Halted!")
endif()

# properties
set(TARGET_NAME {{TARGET_NAME}})
set(BOARD {{BOARD}})
set(FRAMEWORK {{FRAMEWORK}})
set(ENTRY {{ENTRY}})

cmake_minimum_required(VERSION ${CMAKE_VERSION})
project(${PROJECT_NAME} C CXX ASM)
cmake_policy(SET CMP0023 OLD)

file(GLOB_RECURSE SRC_FILES "${PROJECT_PATH}/${ENTRY}/*.cpp" "${PROJECT_PATH}/${ENTRY}/*.cc" "${PROJECT_PATH}/${ENTRY}/*.c")
generate_arduino_firmware(${TARGET_NAME}
    SRCS ${SRC_FILES}
    BOARD ${BOARD}
    PORT {{PORT}})
target_compile_definitions(${TARGET_NAME} PRIVATE __AVR_${FRAMEWORK}__ {{TARGET_COMPILE_FLAGS}})
