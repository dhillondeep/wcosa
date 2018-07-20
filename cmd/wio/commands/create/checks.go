// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create and update command and sub commands provided by the tool.
// Helper for create and update command to perform various checks
package create

import (
    goerr "errors"
    "os"
    "path/filepath"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/log"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"

    "github.com/urfave/cli"
)

// Check directory
func performDirectoryCheck(context *cli.Context) (string, error) {
    // directory is always the first argument
    if len(context.Args()) <= 0 {
        return os.Getwd()
    }
    return filepath.Abs(context.Args()[0])
}

// This check is used to see if wio.yml file exists and the directory is valid
func performWioExistsCheck(directory string) error {
    if !io.Exists(directory) {
        return errors.PathDoesNotExist{Path: directory}
    } else if !io.Exists(directory + io.Sep + io.Config) {
        return errors.ConfigMissing{}
    }
    return nil
}

// This performs various checks before update can be triggered
func performPreUpdateCheck(directory string, create *Create) error {
    return nil
}

/// This method is a crucial piece of check to make sure people do not lose their work. It makes
/// sure that if people are creating the project when there are files in the folder, they mean it
/// and not doing it by mistake. It will warn them to update instead if they want
func performPreCreateCheck(directory string, onlyConfig bool) error {
    // Configure existing directory
    if onlyConfig {
        // Check if config exists
        if io.Exists(directory + io.Sep + io.Config) {
            promptMsg := "Override existing wio.yml?"
            yes, err := log.PromptYes(promptMsg)
            if err != nil {
                return err
            }
            if !yes {
                return goerr.New("project config already exists")
            }
        }

    } else if io.Exists(directory) {
        if isEmpty, err := utils.IsEmpty(directory); err != nil {
            return err
        } else if !isEmpty {
            // directory is not empty
            promptMsg := "Directory is not empty; erase contents?"
            yes, err := log.PromptYes(promptMsg)
            if err != nil {
                return err
            }
            if !yes {
                return goerr.New("working directory not empty")
            }

            // delete all the files
            if err := utils.RemoveContents(directory); err != nil {
                return err
            }
        }
    }
    return nil
}
