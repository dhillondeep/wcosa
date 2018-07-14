#!/bin/bash

mkdir -p $(pwd)/project-pkg/pkg-list/vendor
rm -f $(pwd)/project-pkg/pkg-list/vendor/pkg-malloc
ln -s -f $(pwd)/project-pkg/pkg-malloc $(pwd)/project-pkg/pkg-list/vendor/pkg-malloc

mkdir -p $(pwd)/project-app/app-avr/vendor
rm -f $(pwd)/project-app/app-avr/vendor/pkg-uart
rm -f $(pwd)/project-app/app-avr/vendor/pkg-list
ln -s -f $(pwd)/project-pkg/pkg-uart $(pwd)/project-app/app-avr/vendor/pkg-uart
ln -s -f $(pwd)/project-pkg/pkg-list $(pwd)/project-app/app-avr/vendor/pkg-list

mkdir -p $(pwd)/project-pkg/pkg-trace/vendor
rm -f $(pwd)/project-pkg/pkg-trace/vendor/pkg-uart
ln -s -f $(pwd)/project-pkg/pkg-uart $(pwd)/project-pkg/pkg-trace/vendor/pkg-uart
