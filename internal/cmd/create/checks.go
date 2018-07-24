// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create and update command and sub commands provided by the tool.
// Helper for create and update command to perform various checks
package create

import (
    "os"
    "path/filepath"
    "wio/pkg/log"
    "wio/pkg/util"
    "wio/pkg/util/sys"

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
    if !sys.Exists(directory) {
        return util.Error("path does not exist: %s", directory)
    } else if !sys.Exists(sys.Path(directory, sys.Config)) {
        return util.Error("path is missing wio.yml: %s", directory)
    }
    return nil
}

/// This method is a crucial piece of check to make sure people do not lose their work. It makes
/// sure that if people are creating the project when there are files in the folder, they mean it
/// and not doing it by mistake. It will warn them to update instead if they want
func performPreCreateCheck(directory string, onlyConfig bool) error {
    // Configure existing directory
    if onlyConfig {
        // Check if config exists
        if sys.Exists(sys.Path(directory, sys.Config)) {
            promptMsg := "Override existing wio.yml?"
            yes, err := log.PromptYes(promptMsg)
            if err != nil {
                return err
            }
            if !yes {
                return util.Error("project config already exists")
            }
        }

    } else if sys.Exists(directory) {
        if isEmpty, err := util.IsEmpty(directory); err != nil {
            return err
        } else if !isEmpty {
            // directory is not empty
            promptMsg := "Directory is not empty; erase contents?"
            yes, err := log.PromptYes(promptMsg)
            if err != nil {
                return err
            }
            if !yes {
                return util.Error("working directory not empty")
            }

            // delete all the files
            if err := util.RemoveContents(directory); err != nil {
                return err
            }
        }
    }
    return nil
}
