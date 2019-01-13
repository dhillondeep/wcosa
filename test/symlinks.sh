#!/bin/bash

set -e
dir=$(pwd)

# Check that working directory contains script
if [ ! -f $(pwd)/`basename "$0"` ]; then
    printf "ERROR: `basename "$0"` must be run from its directory\n"
    exit 1
fi

_link() {
    mkdir -p $(dirname "${2}")
    rm -f "${2}"
    ln -s -f "${1}" "${2}"
}

_link ${dir}/project-pkg/pkg-malloc ${dir}/project-pkg/pkg-list/vendor/pkg-malloc

_link ${dir}/project-pkg/pkg-uart ${dir}/project-app/app-avr/vendor/pkg-uart
_link ${dir}/project-pkg/pkg-list ${dir}/project-app/app-avr/vendor/pkg-list
_link ${dir}/project-pkg/pkg-malloc ${dir}/project-app/app-avr/vendor/pkg-malloc

_link ${dir}/project-pkg/pkg-uart ${dir}/project-pkg/pkg-trace/vendor/pkg-uart

_link ${dir}/project-pkg/pkg-malloc ${dir}/project-app/app-alloc/vendor/pkg-malloc
_link ${dir}/project-pkg/pkg-malloc ${dir}/project-app/app-alloc/vendor/alloc-one/vendor/pkg-malloc
_link ${dir}/project-pkg/pkg-malloc ${dir}/project-app/app-alloc/vendor/alloc-two/vendor/pkg-malloc
_link ${dir}/project-pkg/pkg-shared ${dir}/project-app/app-shared/vendor/pkg-shared
_link ${dir}/project-pkg/pkg-ingest ${dir}/project-pkg/pkg-placeholder/vendor/pkg-ingest

