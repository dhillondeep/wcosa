// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands package, which contains all the commands provided by the tool.
// Cleans all the build files for the project
package clean

import (
	"github.com/urfave/cli"
	"os"
	"path/filepath"
	"wio/cmd/wio/commands"
	"wio/cmd/wio/log"
	"wio/cmd/wio/utils"
	"wio/cmd/wio/utils/io"
)

type Clean struct {
	Context *cli.Context
	error
}

// get context for the command
func (clean Clean) GetContext() *cli.Context {
	return clean.Context
}

// Runs the build command when cli build option is provided
func (clean Clean) Execute() {
	RunClean(clean.Context.String("dir"), clean.Context.String("target"),
		clean.Context.IsSet("target"))
}

// This function allows other packages to call the clean command. It contains
// all the main functionality
func RunClean(directoryCli string, targetCli string, targetProvided bool) {
	directory, err := filepath.Abs(directoryCli)
	commands.RecordError(err, "")

	targetToCleanName := targetCli

	log.Norm.Yellow(true, "Cleaning build files")

	if !targetProvided {
		log.Norm.Cyan(false, "cleaning all the targets ... ")

		// clean everything
		if err := os.RemoveAll(directory + io.Sep + ".wio" + io.Sep + "build" + io.Sep + "targets"); err != nil {
			commands.RecordError(err, "failure")
		} else {
			log.Norm.Green(true, "success")
		}
	} else {
		// clean an individual target
		targetPath := directory + io.Sep + ".wio" + io.Sep + "build" + io.Sep + "targets" + io.Sep + targetToCleanName

		if utils.PathExists(targetPath) {
			log.Norm.Cyan(false, "cleaning target: \""+targetToCleanName+"\" ... ")

			// delete the target
			if err := os.RemoveAll(targetPath); err != nil {
				commands.RecordError(err, "failure")
			} else {
				log.Norm.Green(true, "success")
			}

		} else {
			log.Norm.Cyan(false, "no build files exists for target: \""+targetToCleanName+"\" ... ")
			log.Norm.Green(true, "skipped")
		}
	}

	log.Norm.Yellow(true, "Clean command executed successfully!!")
}
