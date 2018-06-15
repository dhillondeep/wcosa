package create

import (
    "github.com/urfave/cli"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/errors"
    goerr "errors"
    "wio/cmd/wio/config"
    "wio/cmd/wio/log"
    "strings"
)

// This check is used to see if the cli arguments are required length
func performArgumentCheck(context *cli.Context, isUpdating bool, platform string) {
    // check to make sure we are given two arguments (one for directory and one for board)
    if len(context.Args()) <= 0 && !isUpdating {
        err := errors.ProgramArgumentsError{
            CommandName: "create",
            ArgumentName: "directory",
            Err: goerr.New("directory needs to be provided for project creation"),
        }

        log.WriteErrorlnExit(err)
    } else if !isUpdating && platform == AVR && len(context.Args()) <= 1 {
        // check if board is provided and if it is not give a warning
        err := errors.ProgrammingArgumentAssumption{
            CommandName: "create",
            ArgumentName: "board",
            Err: goerr.New("board is not provided so a default board is used: " + config.ProjectDefaults.AVRBoard),
        }

        log.WriteErrorln(err, true)
    }
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
        create.Type = APP
    } else {
        create.Type = PKG
    }

    // check the platform
    projectConfig, err := utils.ReadWioConfig(wioPath)
    if err != nil {
        log.WriteErrorlnExit(err)
    } else {
        create.Platform = projectConfig.GetMainTag().GetCompileOptions().GetPlatform()

        if strings.ToLower(create.Platform) != AVR {
            err := errors.PlatformNotSupportedError{
                Platform: create.Platform,
                Err: goerr.New("update the platform tag before updating the project"),
            }

            log.WriteErrorlnExit(err)
        }
    }
}

/// This method is a crucial peace of check to make sure people do not lose their work. It makes
/// sure that if people are creating the project when there are files in the folder, they mean it
/// and not doing it by mistake. It will warn them to update instead if they want
func performPreCreateCheck(directory string) {
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
        }
    }
}
