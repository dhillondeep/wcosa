#!/bin/bash

TEST_DIR=.wio-tests
BASE_DIR=$(pwd)

echo "-------- WIO TEST SUITE --------"
source $(pwd)/wenv
wmake clean
wmake build
wio -v

_pre() {
    printf "\n"
    echo "-------- BEGIN TEST " $1 " --------"
    cd $BASE_DIR
    rm -rf $TEST_DIR
    mkdir $TEST_DIR
}

_post() {
    cd $BASE_DIR
    rm -rf $TEST_DIR
    echo "-------- END TEST " $1 " --------"
}

_test1() {
    _pre 1
    cd $TEST_DIR
    wio create pkg test-pkg --platform native --framework all
    _post 1
}

_test2() {
    _pre 2
    cd $TEST_DIR
    wio create pkg test-pkg
    _post 2
}

_test3() {
    _pre 3
    cd $TEST_DIR
    wio create app test-app --platform avr --framework cosa --board mega2560
    _post 3
}

_test4() {
    _pre 4
    cp -r tests/project-pkg/pkg-square $TEST_DIR/
    cd $TEST_DIR/pkg-square
    wio create pkg --only-config --platform native
    wio build --all
    _post 4
}

_test1
_test2
_test3
_test4

