#!/bin/bash

set -e

test_folder="wio-test"
base_folder=$(pwd)
num_tests=13

# Check that working directory contains script
if [ ! -f $(pwd)/`basename "${0}"` ]; then
    printf "ERROR: `basename "${0}"` must be run from its directory\n"
    exit 1
fi

# Pre and Post test functions
_pre() {
    printf "======== RUNNING TEST ${1} ========\n"
    cd ${base_folder}
    rm -rf ${test_folder}
}

_post() {
    printf "========== DONE TEST ${1} =========\n"
    printf "\n"
    cd ${base_folder}
    rm -rf ${test_folder}
}

# Functional tests
_test1() {
    cp -r ./project-pkg/pkg-fives ./${test_folder}
    cd ${test_folder}
    rm wio.yml
    wio create pkg --platform native --only-config
    wio build
    wio run
}

_test2() {
    cp -r ./project-pkg/pkg-square ./${test_folder}
    cd ${test_folder}
    rm wio.yml
    wio create pkg --platform native --only-config
    wio build
    wio run
}

_test3() {
    cd ./project-pkg/pkg-malloc
    wio clean --hard
    wio update
    wio build
    wio run
    wio clean
}

_test4() {
    cd ./project-pkg/pkg-uart
    wio clean --hard
    wio update
    wio build
    wio clean --hard
}

_test5() {
    cd ./project-pkg/pkg-list
    wio clean --all --hard
    wio update --verbose
    wio build --all
    wio clean native-tests --verbose
    wio build native-tests --disable-warnings
    wio build avr-tests
    wio run native-tests
}

_test6() {
    cd ./project-pkg/pkg-trace
    wio clean --hard
    wio update
    wio build cosa-tests
}

_test7() {
    cd ./project-app/app-avr
    wio clean --hard
    wio update
    wio build
}

_test8() {
    cd ./project-app/app-pipe
    wio clean --hard
    wio update
    wio build
    wio run main --args "cat wio.yml"
}

_test9() {
    cd ./project-app/app-stdout
    wio clean --hard
    wio update
    wio build
    wio run main --args "15 7"
}

_test10() {
    cd ./project-app/app-alloc/vendor/alloc-one
    wio clean --hard
    wio build
    wio run
}

_test11() {
    cd ./project-app/app-alloc/vendor/alloc-two
    wio clean --hard
    wio build
    wio run
}

_test12() {
    cd ./project-app/app-alloc
    wio clean --hard
    wio update
    wio build
    wio run
}

_test13() {
    wio create pkg --platform native ${test_folder}
    cd ${test_folder}
    wio install react
}

# Source and build
cd ./../
source ./wenv
./wmake clean
./wmake build
cd ./test

# Remove all build folders
find ./ -maxdepth 3 -name ".wio" -type d -exec rm -rf {} \;

# Run all test
_all() {
    for i in `seq 1 ${num_tests}`; do
        test_func="_test${i}"
        _pre "${i}"
        ${test_func}
        _post "${i}"
    done
}

# Run specific test
_test() {
    index="${1}"
    test_func="_test${index}"
    _pre "${index}"
    ${test_func}
    _post "${index}"
}

# Argument parsing
if [ "${#}" -gt 1 ] ; then
    echo "runtests.sh [[index]]"
    echo "Too many arguments"
    exit 1
elif [ "${#}" -eq 1 ] ; then
    _test "${1}"
else
    _all
fi
