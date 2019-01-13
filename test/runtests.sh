#!/bin/bash

set -e

test_folder="wio-test"
base_folder=$(pwd)
num_tests=23

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
    wio build --all
    wio run
}

_test2() {
    cp -r ./project-pkg/pkg-square ./${test_folder}
    cd ${test_folder}
    rm wio.yml
    wio create pkg --platform native --only-config
    wio build --all
    wio run
}

_test3() {
    cd ./project-pkg/pkg-malloc
    wio clean --hard
    wio update
    wio build --all
    wio run
    wio clean
}

_test4() {
    cd ./project-pkg/pkg-uart
    wio clean --hard
    wio update
    wio build --all
    wio clean --hard
}

_test5() {
    cd ./project-pkg/pkg-list
    wio clean --all --hard
    wio update --verbose
    wio build --all
    wio clean --all --verbose
    wio build --all --disable-warnings
    wio run native-tests
}

_test6() {
    cd ./project-pkg/pkg-trace
    wio clean --hard
    wio update
    wio build --all
}

_test7() {
    cd ./project-app/app-avr
    wio clean --hard
    wio update
    wio build --all
}

_test8() {
    cd ./project-app/app-pipe
    wio clean --hard
    wio update
    wio build --all
    wio run main --args "cat wio.yml"
}

_test9() {
    cd ./project-app/app-stdout
    wio clean --hard
    wio update
    wio build --all
    wio run main --args "15 7"
}

_test10() {
    cd ./project-app/app-alloc/vendor/alloc-one
    wio clean --hard
    wio build --all
    wio run
}

_test11() {
    cd ./project-app/app-alloc/vendor/alloc-two
    wio clean --hard
    wio build --all
    wio run
}

_test12() {
    cd ./project-app/app-alloc
    wio clean --hard
    wio update
    wio build --all
    wio run
}

_test13() {
    wio create pkg --platform native ${test_folder}
    cd ${test_folder}
    wio install react
}

_test14() {
    cd ./project-pkg/pkg-headerOnly
    wio clean native --hard
    wio clean arduino --hard
    wio clean cosa --hard
    wio update
    wio build --all
    wio run native
}

_test15() {
    cd ./project-pkg/pkg-shared/shared/honly
    make
    cd ../../
    wio clean --hard
    wio update
    wio build --all
    export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:$(pwd)/shared/honly
    export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$(pwd)/shared/honly
    wio run
}

_test16() {
    cd ./project-app/app-shared/vendor/pkg-shared/shared/honly
    make
    cd ../../../../
    wio clean --hard
    wio update
    wio build --all
    export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:$(pwd)/vendor/pkg-shared/shared/honly
    export LD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:$(pwd)/vendor/pkg-shared/shared/honly
}

_test17() {
    wio create app --platform native ${test_folder}
    cd ${test_folder}
    wio install wlib-stl@3.0.2
    machine_type=$(uname -m)
    echo "int main(void) {}" >> src/main.cpp
    echo "    definitions:" >> wio.yml
    if [ "$machine_type" == "x86_64" ]; then
        bit_type="64"
        align="3"
    else
        bit_type = "32"
        align="2"
    fi
    echo "    - WLIB_TLSF_ARCH=$bit_type" >> wio.yml
    echo "    - WLIB_TLSF_LOG2_DIV=4" >> wio.yml
    echo "    - WLIB_TLSF_LOG2_ALIGN=$align" >> wio.yml
    echo "    - WLIB_TLSF_LOG2_MAX=10" >> wio.yml
    wio build --all
}

_test18() {
    cd ./project-app/app-pthread
    wio clean --hard
    wio update
    wio build --all
    wio run
}

_test19() {
    cd ./project-app/app-osspecific
    wio clean --hard
    wio update
    wio build --all
    wio run
}

_test20() {
    cd ./project-pkg/pkg-ingest
    wio clean test1 test2 --hard
    wio update
    wio build --all
    wio run test1
    wio run test2
}

_test21() {
    cd ./project-app/app-linker
    wio clean --hard
    wio update
    wio install
    wio build --all
    wio run
}

_test22() {
    cd ./project-pkg/pkg-placeholder
    wio clean --hard
    wio update
    wio build --all
    wio run test-deep
    wio run test-jeff
}

_test23() {
    cd ./project-pkg/pkg-libraries-complex
    wio clean --all --hard
    wio update
    wio build --all
    wio run
}

# Source and build
cd ./../
source ./wenv
./wmake clean
./wmake build
cd ./test

# Remove all build folders
rm -rf ${test_folder}
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