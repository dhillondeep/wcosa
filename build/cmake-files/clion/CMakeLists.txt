set(WCOSA_CMAKE_VERSION %cmake-version)
set(WCOSA_PROJECT_NAME %project-name)

cmake_minimum_required(VERSION ${WCOSA_CMAKE_VERSION})
project(${WCOSA_PROJECT_NAME})

include("CMakeListsPrivate.txt")

add_custom_target(
    WCOSA_BUILD ALL
    COMMAND ${WCOSA_CMD} build
)

add_custom_target(
    WCOSA_UPDATE_ALL ALL
    COMMAND ${WCOSA_CMD} update
)

add_custom_target(
    WCOSA_UPLOAD ALL
    COMMAND ${WCOSA_CMD} upload
)

add_executable(${PROJECT_NAME} ${SRC_FILES})