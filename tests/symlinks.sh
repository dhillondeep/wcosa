#!/bin/bash

mkdir -p $(pwd)/project-pkg/pkg-list/vendor
ln -s -f $(pwd)/project-pkg/pkg-malloc $(pwd)/project-pkg/pkg-list/vendor/pkg-malloc
