#!/bin/bash

set -e

# Check that working directory contains script
if [ ! -f $(pwd)/`basename "$0"` ]; then
    printf "ERROR: `basename "$0"` must be run from its directory\n"
    exit 1
fi

mkdir -p $(pwd)/project-pkg/pkg-list/vendor
rm -f $(pwd)/project-pkg/pkg-list/vendor/pkg-malloc
ln -s -f $(pwd)/project-pkg/pkg-malloc $(pwd)/project-pkg/pkg-list/vendor/pkg-malloc

mkdir -p $(pwd)/project-app/app-avr/vendor
rm -f $(pwd)/project-app/app-avr/vendor/pkg-uart
rm -f $(pwd)/project-app/app-avr/vendor/pkg-list
rm -f $(pwd)/project-app/app-avr/vendor/pkg-malloc
ln -s -f $(pwd)/project-pkg/pkg-uart $(pwd)/project-app/app-avr/vendor/pkg-uart
ln -s -f $(pwd)/project-pkg/pkg-list $(pwd)/project-app/app-avr/vendor/pkg-list
ln -s -f $(pwd)/project-pkg/pkg-malloc $(pwd)/project-app/app-avr/vendor/pkg-malloc

mkdir -p $(pwd)/project-pkg/pkg-trace/vendor
rm -f $(pwd)/project-pkg/pkg-trace/vendor/pkg-uart
ln -s -f $(pwd)/project-pkg/pkg-uart $(pwd)/project-pkg/pkg-trace/vendor/pkg-uart
