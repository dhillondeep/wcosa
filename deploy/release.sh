#!/bin/bash

root_dir=$(cd -P -- "$(dirname -- "$0")" && pwd -P)
root_dir_name=$(basename "$root_dir")
cd "$root_dir"

# build a release
cd gen
go build build.go
./build
cd ../

# release git and npm
npm run release

# brew release
cd brew
go build brew.go
./brew
cd ../

# scoop release
cd scoop
go build scoop.go
./scoop
cd ../
