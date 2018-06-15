// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package utils contains utilities/files useful throughout the app
// This file contains all the function to manipulate project configuration file

package utils

import (
    "gopkg.in/yaml.v2"
    "strings"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/types"
    "wio/cmd/wio/errors"
    "bufio"
    "regexp"
)

// Write configuration for the project with information on top and nice spacing
func PrettyPrintConfigHelp(projectConfig interface{}, filePath string) error {
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

    // get data
    if ymlData, err = yaml.Marshal(projectConfig); err != nil {
        return err
    }
    if appInfoData, err = io.AssetIO.ReadFile(appInfoPath); err != nil {
        return err
    }
    if pkgInfoData, err = io.AssetIO.ReadFile(pkgInfoPath); err != nil {
        return err
    }
    if targetsInfoData, err = io.AssetIO.ReadFile(targetsInfoPath); err != nil {
        return err
    }
    if dependenciesInfoData, err = io.AssetIO.ReadFile(dependenciesInfoPath); err != nil {
        return err
    }

    finalString := ""
    currentString := strings.Split(string(ymlData), "\n")

    beautify := false
    first := false
    create := true

    for line := range currentString {
        currLine := currentString[line]

        if len(currLine) <= 1 {
            continue
        }

        if strings.Contains(currLine, "app:") && create {
            finalString += string(appInfoData) + "\n"
            create = false
        } else if strings.Contains(currLine, "pkg:") && create {
            finalString += string(pkgInfoData) + "\n"
            create = false
        } else if strings.Contains(currLine, "targets:") {
            finalString += "\n" + string(targetsInfoData) + "\n"
        } else if strings.Contains(currLine, "create:") {
            beautify = true
        } else if strings.Contains(currLine, "dependencies:") {
            beautify = true
            first = false
            finalString += "\n" + string(dependenciesInfoData) + "\n"
        } else if beautify && !first {
            first = true
        } else if !strings.Contains(currLine, "compile_flags:") && beautify {
            simpleString := strings.Trim(currLine, " ")

            if simpleString[len(simpleString)-1] == ':' {
                finalString += "\n"
            }
        }

        finalString += currLine + "\n"
    }

    err = io.NormalIO.WriteFile(filePath, []byte(finalString))

    return err
}

// Write configuration with nice spacing and information
func PrettyPrintConfig(projectConfig types.Config, filePath string) error {
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
            Err: err,
        }
    }
    if pkgInfoData, err = io.AssetIO.ReadFile(pkgInfoPath); err != nil {
        return errors.ReadFileError{
            FileName: pkgInfoPath,
            Err: err,
        }
    }
    if targetsInfoData, err = io.AssetIO.ReadFile(targetsInfoPath); err != nil {
        return errors.ReadFileError{
            FileName: targetsInfoPath,
            Err: err,
        }
    }
    if dependenciesInfoData, err = io.AssetIO.ReadFile(dependenciesInfoPath); err != nil {
        return errors.ReadFileError{
            FileName: dependenciesInfoPath,
            Err: err,
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
    flagsTagPat := regexp.MustCompile(`(^flags:)|((\s| |^\w)flags:(\s+|))`)
    definitionsTagPat := regexp.MustCompile(`(^definitions:)|((\s| |^\w)definitions:(\s+|))`)

    scanner := bufio.NewScanner(strings.NewReader(string(ymlData)))
    for scanner.Scan() {
        line := scanner.Text()

        if appTagPat.MatchString(line) {
            finalStr += string(appInfoData) + "\n"
            finalStr += line
        } else if pkgTagPat.MatchString(line) {
            finalStr += string(pkgInfoData) + "\n"
            finalStr += line
        } else if targetsTagPat.MatchString(line) {
            finalStr += "\n" + string(targetsInfoData) + "\n"
            finalStr += line
        } else if dependenciesTagPat.MatchString(line) {
            finalStr += "\n" + string(dependenciesInfoData) + "\n"
            finalStr += line
        } else if configTagPat.MatchString(line) || compileOptionsTagPat.MatchString(line) ||
            flagsTagPat.MatchString(line) || definitionsTagPat.MatchString(line) || metaTagPat.MatchString(line) {
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
            Err: err,
        }
    }

    return nil
}
