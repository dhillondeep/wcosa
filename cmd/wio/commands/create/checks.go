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
    "strings"
    "wio/cmd/wio/config"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/log"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/constants"
)

// This check is used to see if the cli arguments are provided and then based on that decide defaults
func performArgumentCheck(context *cli.Context, isUpdating bool, platform string) (string, string) {
    var directory string
    var board string
    var err error

    // check directory
    if len(context.Args()) <= 0 {
        directory, err = os.Getwd()

        log.WriteErrorlnExit(err)

        err = errors.ProgrammingArgumentAssumption{
            CommandName:  "create",
            ArgumentName: "directory",
            Err:          goerr.New("directory is not provided so current directory is used: " + directory),
        }

        log.WriteErrorln(err, true)
    } else {
        directory = context.Args()[0]
    }

    // check board for create
    if !isUpdating && platform == constants.AVR && len(context.Args()) <= 1 {
        err = errors.ProgrammingArgumentAssumption{
            CommandName:  "create",
            ArgumentName: "board",
            Err:          goerr.New("board is not provided so a default board is used: " + config.AvrProjectDefaults.AVRBoard),
        }

        log.WriteErrorln(err, true)
        board = config.AvrProjectDefaults.AVRBoard
    } else if !isUpdating && platform == constants.AVR && len(context.Args()) >= 1 {
        board = context.Args()[1]
    } else {
        board = ""
    }

    return directory, board
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
func performPreUpdateCheck(directory string, create *Create) {
    wioPath := directory + io.Sep + "wio.yml"

    isApp, err := utils.IsAppType(wioPath)
    if err != nil {
        configError := errors.ConfigParsingError{
            Err: goerr.New("project type could not be parsed"),
        }

        log.WriteErrorlnExit(configError)
    }

    if isApp {
        create.Type = constants.APP
    } else {
        create.Type = constants.PKG
    }

    // check the platform
    projectConfig, err := utils.ReadWioConfig(wioPath)
    if err != nil {
        log.WriteErrorlnExit(err)
    } else {
        create.Platform = projectConfig.GetMainTag().GetCompileOptions().GetPlatform()

        if strings.ToLower(create.Platform) != constants.AVR {
            err := errors.PlatformNotSupportedError{
                Platform: create.Platform,
                Err:      goerr.New("update the platform tag before updating the project"),
            }

            log.WriteErrorlnExit(err)
        }
    }
}

/// This method is a crucial peace of check to make sure people do not lose their work. It makes
/// sure that if people are creating the project when there are files in the folder, they mean it
/// and not doing it by mistake. It will warn them to update instead if they want
func performPreCreateCheck(directory string, onlyConfig bool) {
    if onlyConfig {
        if utils.PathExists(directory + io.Sep + "wio.yml") {
            err := errors.OverridePossibilityError{
                Path: directory,
                Err: goerr.New("wio.yml file will be replaced with new config.\n" + errors.Spaces +
                    "Type (y) to indicate creation or anything else otherwise: "),
            }

            log.WriteErrorAndPrompt(err, log.INFO, "y", true)
        }

        return
    }

    if utils.PathExists(directory) {
        if status, err := utils.IsEmpty(directory); err != nil {
            log.WriteErrorlnExit(err)
        } else if !status {
            err := errors.OverridePossibilityError{
                Path: directory,
                Err: goerr.New("files will be replaced with new project files.\n" + errors.Spaces +
                    "Type (y) to indicate creation or anything else otherwise: "),
            }

            log.WriteErrorAndPrompt(err, log.INFO, "y", true)

            // delete all the files
            if err := os.RemoveAll(directory); err != nil {
                deleteError := errors.DeleteDirectoryError{
                    DirName: directory,
                    Err:     err,
                }

                log.WriteErrorlnExit(deleteError)
            }
        }
    }
}
