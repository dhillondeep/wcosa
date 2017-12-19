set(OS %os)
set(WCOSA_CMAKE_VERSION %cmake-version)
set(WCOSA_PROJECT_NAME %project-name)
set(WCOSA_PATH "%wcosa-path")
set(AVR_PATH "%avr-path")
set(AVR_MCU %avr-mcu)
set(AVR_H_FUSE %avr-h-fuse)
set(AVR_L_FUSE %avr-l-fuse)

cmake_minimum_required(VERSION ${WCOSA_CMAKE_VERSION})

include("wcosa/generic-gcc-avr.cmake")

project(${WCOSA_PROJECT_NAME})

if ((${OS} STREQUAL windows) OR (${OS} STREQUAL cygwin))
      # window avr libraries
      set(CMAKE_FIND_ROOT_PATH "${AVR_PATH}/avr")
else()
      # linux/unix avr libraries
      if(DEFINED ENV{AVR_FIND_ROOT_PATH})
        set(CMAKE_FIND_ROOT_PATH $ENV{AVR_FIND_ROOT_PATH})
    else(DEFINED ENV{AVR_FIND_ROOT_PATH})
      if(EXISTS "/opt/local/avr")
        set(CMAKE_FIND_ROOT_PATH "/opt/local/avr")
      elseif(EXISTS "/usr/avr")
        set(CMAKE_FIND_ROOT_PATH "/usr/avr")
      elseif(EXISTS "/usr/lib/avr")
        set(CMAKE_FIND_ROOT_PATH "/usr/lib/avr")
      elseif(EXISTS "/usr/local/CrossPack-AVR")
        set(CMAKE_FIND_ROOT_PATH "/usr/local/CrossPack-AVR")
      else(EXISTS "/opt/local/avr")
        message(FATAL_ERROR "Please set AVR_FIND_ROOT_PATH in your environment.")
    endif(EXISTS "/opt/local/avr")
endif(DEFINED ENV{AVR_FIND_ROOT_PATH})
endif()

set(CMAKE_FIND_ROOT_PATH_MODE_PROGRAM NEVER)
set(CMAKE_FIND_ROOT_PATH_MODE_LIBRARY ONLY)
set(CMAKE_FIND_ROOT_PATH_MODE_INCLUDE ONLY)
# not added automatically, since CMAKE_SYSTEM_NAME is "generic"
set(CMAKE_SYSTEM_INCLUDE_PATH "${CMAKE_FIND_ROOT_PATH}/include")
set(CMAKE_SYSTEM_LIBRARY_PATH "${CMAKE_FIND_ROOT_PATH}/lib")

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
  set(CMAKE_C_FLAGS_RELEASE "-Os")
  set(CMAKE_CXX_FLAGS_RELEASE "-Os")
endif(CMAKE_BUILD_TYPE MATCHES Release)

if(CMAKE_BUILD_TYPE MATCHES Debug)
   set(CMAKE_C_FLAGS_DEBUG "-O0 -save-temps -g -gdwarf-3 -gstrict-dwarf")
   set(CMAKE_CXX_FLAGS_DEBUG "-O0 -save-temps -g -gdwarf-3 -gstrict-dwarf")
endif(CMAKE_BUILD_TYPE MATCHES Debug)

##########################################################################
# compiler options for all build types
##########################################################################
add_definitions("-DF_CPU=${MCU_SPEED}")
add_definitions("-fpack-struct")
add_definitions("-fshort-enums")
add_definitions("-Wall")
add_definitions("-Werror")
add_definitions("-pedantic")
add_definitions("-pedantic-errors")
add_definitions("-funsigned-char")
add_definitions("-funsigned-bitfields")
add_definitions("-ffunction-sections")
add_definitions("-c")
%custom-compiler-flags

##########################################################################
# include search paths
##########################################################################
%include-directories

##########################################################################
# building library and adding src (do not delete this line)
##########################################################################
file(GLOB_RECURSE LIB_FILES "lib/*.cpp" "lib/*.cc" "lib/*.h" "lib/*.hh")
file(GLOB_RECURSE SRC_FILES "src/*.cpp" "src/*.cc")

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
   if (LIB_FILES)
        ${PROJECT_NAME}-librariers-${AVR_MCU}
    endif ()
)

add_custom_target(UPDATE ALL
    COMMAND python ${WCOSA_PATH}/scripts/handle.py -path ${CMAKE_SOURCE_DIR} -update this
    ${CMAKE} ${CMAKE_CURRENT_SRC_DIR}
)