#!/bin/bash

set -e

test_folder="wio-test"
base_folder=$(pwd)
num_tests=6

# Check that working directory contains script
if [ ! -f $(pwd)/`basename "$0"` ]; then
    printf "ERROR: `basename "$0"` must be run from its directory\n"
    exit 1
fi

# Pre and Post test functions
_pre() {
    printf "======== RUNNING TEST $1 ========\n"
    cd $base_folder
    rm -rf $test_folder
}

_post() {
    printf "========== DONE TEST $1 =========\n"
    printf "\n"
    cd $base_folder
    rm -rf $test_folder
}

# Functional tests
_test1() {
    cp -r ./project-pkg/pkg-fives ./$test_folder
    cd $test_folder
    rm wio.yml
    wio create pkg --platform native --only-config
    wio build
    wio run
}

_test2() {
    cp -r ./project-pkg/pkg-square ./$test_folder
    cd $test_folder
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
    wio clean --hard
    wio update
    wio build --all
    wio clean
    wio build native-tests
    wio build avr-tests
    wio run native-tests
}

_test6() {
    cd ./project-pkg/pkg-trace
    wio clean --hard
    wio update
    wio build cosa-tests
}

# Source and build
cd ./../
source ./wenv
./wmake clean
./wmake build
cd ./tests

# Run each test
for i in `seq 1 $num_tests`; do
    test_func="_test$i"
    _pre "$i"
    $test_func
    _post "$i"
done
