// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create and update command and sub commands provided by the tool.
// Helper for create and update command to perform various checks
package create

import (
    goerr "errors"
    "github.com/urfave/cli"
    "os"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/log"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
)

// Check directory
func performDirectoryCheck(context *cli.Context) string {
    var directory string
    var err error

    // directory is always the first argument
    if len(context.Args()) <= 0 {
        directory, err = os.Getwd()
        log.WriteErrorlnExit(err)
    } else {
        directory = context.Args()[0]
    }
    return directory
}

// This check is used to see if wio.yml file exists and the directory is valid
func performWioExistsCheck(directory string) {
    if !utils.PathExists(directory) {
        err := errors.PathDoesNotExist{
            Path: directory,
        }

        log.WriteErrorlnExit(err)
    } else if !utils.PathExists(directory + io.Sep + "wio.yml") {
        err := errors.ConfigMissing{}

        log.WriteErrorlnExit(err)
    }
}

// This performs various checks before update can be triggered
func performPreUpdateCheck(directory string, create *Create) {}

/// This method is a crucial piece of check to make sure people do not lose their work. It makes
/// sure that if people are creating the project when there are files in the folder, they mean it
/// and not doing it by mistake. It will warn them to update instead if they want
func performPreCreateCheck(directory string, onlyConfig bool) {
    // Configure existing directory
    if onlyConfig {
        // Check if config exists
        if utils.PathExists(directory + io.Sep + "wio.yml") {
            err := errors.OverridePossibilityError{
                Path: directory,
                Err: goerr.New("wio.yml file will be replaced with new config.\n" + errors.Spaces +
                    "projectType (y) to indicate creation or anything else otherwise: "),
            }

            log.WriteErrorAndPrompt(err, log.INFO, "y", true)
        }

    } else if utils.PathExists(directory) {
        if isEmpty, err := utils.IsEmpty(directory); err != nil {
            log.WriteErrorlnExit(err)
        } else if !isEmpty {
            // directory is not empty
            err := errors.OverridePossibilityError{
                Path: directory,
                Err: goerr.New("files will be replaced with new project files.\n" + errors.Spaces +
                    "projectType (y) to indicate creation or anything else otherwise: "),
            }

            log.WriteErrorAndPrompt(err, log.INFO, "y", true)

            // delete all the files
            if err := utils.RemoveContents(directory); err != nil {
                log.WriteErrorlnExit(err)
            } else {
                log.Writeln(log.VERB, nil, "deleted all the files from: %s", directory)
            }
        }
    }
}
