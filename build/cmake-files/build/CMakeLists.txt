set(OS %os)
set(WCOSA_CMAKE_VERSION %cmake-version)
set(WCOSA_PROJECT_NAME %project-name)
set(WCOSA_PATH "%wcosa-path")
set(WCOSA_CMD "python ${WCOSA_PATH}/main.py")
set(AVR_TOOL_PATH "%avr-tool-path")
set(AVR_PATH "%avr-path")
set(AVR_MCU %avr-mcu)
set(AVR_F_CPU %avr-f-cpu)
set(AVR_H_CPU %avr-h-fuse)
set(AVR_L_CPU %avr-l-fuse)
set(AVR_CC "%cmake-c-flags-compiler")
set(AVR_CXX "%cmake-cxx-flags-compiler")

cmake_minimum_required(VERSION ${WCOSA_CMAKE_VERSION})

include("generic-gcc-avr.cmake")

project(${WCOSA_PROJECT_NAME})

set(CMAKE_CXX_FLAGS ${AVR_CXX})
set(CMAKE_C_FLAGS ${AVR_CC})

set(CMAKE_FIND_ROOT_PATH_MODE_PROGRAM NEVER)
set(CMAKE_FIND_ROOT_PATH_MODE_LIBRARY ONLY)
set(CMAKE_FIND_ROOT_PATH_MODE_INCLUDE ONLY)
# not added automatically, since CMAKE_SYSTEM_NAME is "generic"
set(CMAKE_SYSTEM_INCLUDE_PATH "${AVR_PATH}/include")
set(CMAKE_SYSTEM_LIBRARY_PATH "${AVR_PATH}/lib")

##########################################################################
# set build type
##########################################################################
if(NOT CMAKE_BUILD_TYPE)
   set(CMAKE_BUILD_TYPE Release)
endif(NOT CMAKE_BUILD_TYPE)

##########################################################################
# status messages for generating
##########################################################################
message(STATUS "Set CMAKE_FIND_ROOT_PATH to ${CMAKE_FIND_ROOT_PATH}")
message(STATUS "Set CMAKE_SYSTEM_INCLUDE_PATH to ${CMAKE_SYSTEM_INCLUDE_PATH}")
message(STATUS "Set CMAKE_SYSTEM_LIBRARY_PATH to ${CMAKE_SYSTEM_LIBRARY_PATH}")

if(CMAKE_BUILD_TYPE MATCHES Release)
  set(CMAKE_C_FLAGS_RELEASE "%cmake-c-flags-release")
  set(CMAKE_CXX_FLAGS_RELEASE "%cmake-cxx-flags-release")
endif(CMAKE_BUILD_TYPE MATCHES Release)

if(CMAKE_BUILD_TYPE MATCHES Debug)
   set(CMAKE_C_FLAGS_DEBUG "%cmake-c-flags-debug")
   set(CMAKE_CXX_FLAGS_DEBUG "%cmake-cxx-flags-debug")
endif(CMAKE_BUILD_TYPE MATCHES Debug)

##########################################################################
# compiler options for all build types
##########################################################################
add_definitions("-DF_CPU=${AVR_F_CPU}")
add_definitions(%custom-definitions)

##########################################################################
# include search paths
##########################################################################
include_directories("${AVR_PATH}/include")
include_directories("${AVR_PATH}/lib")
%custom-directories-include

##########################################################################
# building library and adding src (do not delete this line)
##########################################################################
file(GLOB_RECURSE LIB_FILES "../../lib/*.cpp" "../../lib/*.cc" "../../lib/*.h" "../../lib/*.hh")
file(GLOB_RECURSE SRC_FILES "../../src/*.cpp" "../../src/*.cc")

if (LIB_FILES)
    add_avr_library(
        "${PROJECT_NAME}-libraries"
        ${LIB_FILES}
    )
endif()

add_avr_executable(
   ${PROJECT_NAME}
   ${SRC_FILES}
)

target_link_libraries(
   ${PROJECT_NAME}-${AVR_MCU}.elf
)

if (LIB_FILES)
   target_link_libraries(
        ${PROJECT_NAME}-libraries-${AVR_MCU}
)
endif ()
