// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package utils contains utilities/files useful throughout the app
// This file contains all the function to manipulate project configuration file

package utils

import (
    "strings"
    "gopkg.in/yaml.v2"
    "wio/cmd/wio/utils/io"
)

// Write configuration for the project with information on top and nice spacing
func PrettyPrintConfigHelp(projectConfig interface{}, filePath string) (error) {
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

// Write configuration with nice spacing
func PrettyPrintConfigSpacing(projectConfig interface{}, filePath string) (error) {
    var ymlData []byte
    var err error

    // get data
    if ymlData, err = yaml.Marshal(projectConfig); err != nil {
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
            create = false
        } else if strings.Contains(currLine, "pkg:") && create {
            create = false
        } else if strings.Contains(currLine, "targets:") {
            finalString += "\n"
        } else if strings.Contains(currLine, "create:") {
            beautify = true
        } else if strings.Contains(currLine, "dependencies:") {
            finalString += "\n"
            beautify = true
            first = false
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
