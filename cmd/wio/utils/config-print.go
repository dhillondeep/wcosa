// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package utils contains utilities/files useful throughout the app
// This file contains all the function to manipulate project configuration file

package utils

import (
    "bufio"
    "gopkg.in/yaml.v2"
    "regexp"
    "strings"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"
)

// Write configuration with nice spacing and information
func PrettyPrintConfig(projectConfig types.Config, filePath string, showHelp bool) error {
    appInfoPath := "templates" + io.Sep + "config" + io.Sep + "app-helper.txt"
    pkgInfoPath := "templates" + io.Sep + "config" + io.Sep + "pkg-helper.txt"
    targetsInfoPath := "templates" + io.Sep + "config" + io.Sep + "targets-helper.txt"
    dependenciesInfoPath := "templates" + io.Sep + "config" + io.Sep + "dependencies-helper.txt"

    var ymlData []byte
    var appInfoData []byte
    var pkgInfoData []byte
    var targetsInfoData []byte
    var dependenciesInfoData []byte
    var err error

    if appInfoData, err = io.AssetIO.ReadFile(appInfoPath); err != nil {
        return errors.ReadFileError{
            FileName: appInfoPath,
            Err:      err,
        }
    }
    if pkgInfoData, err = io.AssetIO.ReadFile(pkgInfoPath); err != nil {
        return errors.ReadFileError{
            FileName: pkgInfoPath,
            Err:      err,
        }
    }
    if targetsInfoData, err = io.AssetIO.ReadFile(targetsInfoPath); err != nil {
        return errors.ReadFileError{
            FileName: targetsInfoPath,
            Err:      err,
        }
    }
    if dependenciesInfoData, err = io.AssetIO.ReadFile(dependenciesInfoPath); err != nil {
        return errors.ReadFileError{
            FileName: dependenciesInfoPath,
            Err:      err,
        }
    }

    // marshall yml data
    if ymlData, err = yaml.Marshal(projectConfig); err != nil {
        marshallError := errors.YamlMarshallError{
            Err: err,
        }
        return marshallError
    }

    finalStr := ""

    // configuration tags
    appTagPat := regexp.MustCompile(`(^app:)|((\s| |^\w)app:(\s+|))`)
    pkgTagPat := regexp.MustCompile(`(^pkg:)|((\s| |^\w)pkg:(\s+|))`)
    targetsTagPat := regexp.MustCompile(`(^targets:)|((\s| |^\w)targets:(\s+|))`)
    dependenciesTagPat := regexp.MustCompile(`(^dependencies:)|((\s| |^\w)dependencies:(\s+|))`)
    configTagPat := regexp.MustCompile(`(^config:)|((\s| |^\w)config:(\s+|))`)
    compileOptionsTagPat := regexp.MustCompile(`(^compile_options:)|((\s| |^\w)compile_options:(\s+|))`)
    metaTagPat := regexp.MustCompile(`(^meta:)|((\s| |^\w)meta:(\s+|))`)

    // empty array
    emptyArrayPat := regexp.MustCompile(`:\s+\[\]`)
    // empty object
    emptyMapPat := regexp.MustCompile(`:\s+\{\}`)
    // empty tag
    emptyTagPat := regexp.MustCompile(`:\s+\n+|:\s+"\s+"|:\s+""|:"\s+"|:""`)
    // board
    boardPat := regexp.MustCompile(`board`)

    scanner := bufio.NewScanner(strings.NewReader(string(ymlData)))
    for scanner.Scan() {
        line := scanner.Text()

        if projectConfig.GetMainTag().GetCompileOptions().GetPlatform() == types.DESKTOP {
            // skip board tags for desktop platform
            if boardPat.MatchString(line) {
                continue
            }
        }

        // ignore empty arrays, objects and tags
        if emptyArrayPat.MatchString(line) || emptyMapPat.MatchString(line) || emptyTagPat.MatchString(line) {
            if !(strings.Contains(line, "global_flags: []") ||
                strings.Contains(line, "target_flags: []") ||
                strings.Contains(line, "pkg_flags: []") ||
                strings.Contains(line, "global_definitions: []") ||
                strings.Contains(line, "target_definitions: []") ||
                strings.Contains(line, "pkg_definitions: []")){
                continue
            }
        }

        if appTagPat.MatchString(line) {
            if showHelp {
                finalStr += string(appInfoData) + "\n"
            }

            finalStr += line
        } else if pkgTagPat.MatchString(line) {
            if showHelp {
                finalStr += string(pkgInfoData) + "\n"
            }

            finalStr += line
        } else if targetsTagPat.MatchString(line) {
            finalStr += "\n"
            if showHelp {
                finalStr += string(targetsInfoData) + "\n"
            }
            finalStr += line
        } else if dependenciesTagPat.MatchString(line) {
            finalStr += "\n"
            if showHelp {
                finalStr += string(dependenciesInfoData) + "\n"
            }
            finalStr += line
        } else if configTagPat.MatchString(line) || compileOptionsTagPat.MatchString(line) ||
            metaTagPat.MatchString(line)  {
            finalStr += "\n"
            finalStr += line
        } else {
            finalStr += line
        }

        finalStr += "\n"
    }

    if err = io.NormalIO.WriteFile(filePath, []byte(finalStr)); err != nil {
        return errors.WriteFileError{
            FileName: filePath,
            Err:      err,
        }
    }

    return nil
}
