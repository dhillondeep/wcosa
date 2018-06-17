// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands package, which contains all the commands provided by the tool.
// Builds the project
package build

import (
	"bytes"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"wio/cmd/wio/commands"
	"wio/cmd/wio/commands/clean"
	"wio/cmd/wio/dependencies"
	"wio/cmd/wio/log"
	"wio/cmd/wio/types"
	"wio/cmd/wio/utils"
	"wio/cmd/wio/utils/io"
)

type Build struct {
	Context *cli.Context
	error
}

// get context for the command
func (build Build) GetContext() *cli.Context {
	return build.Context
}

// Runs the build command when cli build option is provided
func (build Build) Execute() {
	RunBuild(build.Context.String("dir"), build.Context.String("target"),
		build.Context.Bool("clean"), "")
}

// This function allows other packages to call build as well. This is also used when cli build is executed
func RunBuild(directoryCli string, targetCli string, cleanCli bool, port string) string {
	directory, err := filepath.Abs(directoryCli)
	commands.RecordError(err, "")

	cleanStr := cleanCli
	targetToBuildName := targetCli

	configPath := directory + io.Sep + "wio.yml"
	appType, err := utils.IsAppType(configPath)
	commands.RecordError(err, "")

	var projectConfig types.Config

	if appType {
		appConfig := types.AppConfig{}
		commands.RecordError(io.NormalIO.ParseYml(configPath, &appConfig), "")
		projectConfig = appConfig
	} else {
		pkgConfig := types.PkgConfig{}
		commands.RecordError(io.NormalIO.ParseYml(configPath, &pkgConfig), "")
		projectConfig = pkgConfig
	}

	var targetToBuild types.Target

	// if target is default then we choose the default target
	if targetToBuildName == "default" {
		targetToBuildName = projectConfig.GetTargets().GetDefaultTarget()
		targetToBuild = projectConfig.GetTargets().GetTargets()[targetToBuildName]
	} else {
		// choose the target that has been provided (from cli)
		targetToBuild = projectConfig.GetTargets().GetTargets()[targetToBuildName]
	}

	// clean the build files if clean flag is true
	if cleanStr {
		cleanBuildFiles(directory, targetToBuildName)
	}

	log.Norm.Yellow(true, "Building the project")

	flags := targetToBuild.GetFlags()

	projectDependencies := projectConfig.GetDependencies()

	// if we are dealing with package type, make that the dependency
	if !appType {
		pkgConfig := projectConfig.(types.PkgConfig)

		projectDependencies = types.DependenciesTag{pkgConfig.MainTag.Name: &types.DependencyTag{
			Version:         pkgConfig.MainTag.Version,
			Vendor:          false,
			LinkVisibility:  "private",
			DependencyFlags: targetToBuild.GetFlags()["pkg_flags"],
		}}
	}

	// create the target (for now take the first framework and platform)
	createTarget(projectConfig.GetMainTag().GetName(), directory, targetToBuild.GetBoard(), port,
		projectConfig.GetMainTag().GetFrameworks()[0], targetToBuildName,
		flags, projectDependencies, appType, projectConfig.GetMainTag().IsHeaderOnly())

	// build the target
	buildTarget(directory, targetToBuildName)

	// print the ending message
	log.Norm.Yellow(true, "Build successful for \""+projectConfig.GetMainTag().GetName()+"\" project!")

	return targetToBuildName
}

// Scans dependency tree and based on that creates CMake build files
func createTarget(projectName string, directory string, board string, port string, framework string, targetName string,
	projectFlags map[string][]string, projectDependencies types.DependenciesTag, isApp bool, headerOnly bool) {

	log.Verb.Verbose(false, "Generating Dependency Graph ... ")

	// parses all the dependency packages and create a cmake dependency file
	if err := dependencies.CreateCMakeDependencies(projectName, directory, projectFlags, projectDependencies, isApp); err != nil {
		log.Verb.Verbose(true, "failure")
		commands.RecordError(err, "")
	} else {
		log.Verb.Verbose(true, "success")
		commands.RecordError(err, "")
	}

	log.Verb.Verbose(false, "Creating Build Tool files ... ")

	// create main cmake file
	if err := dependencies.CreateMainCMake(projectName, directory, board, port, framework, targetName, projectFlags, isApp); err != nil {
		log.Verb.Verbose(true, "failure")
	} else {
		log.Verb.Verbose(true, "success")
	}
}

// Execute CMake and Make commands to build the project
func buildTarget(directory string, targetName string) {
	targetsDirectory := directory + io.Sep + ".wio" + io.Sep + "build" + io.Sep + "targets"
	targetPath := targetsDirectory + io.Sep + targetName

	log.Norm.Cyan(false, "creating build environment for target: \""+targetName+"\" ... ")

	// create a folder for the target
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		commands.RecordError(err, "failure")
	} else {
		log.Norm.Green(true, "success")
	}

	// execute cmake command
	cmakeCommand := exec.Command("cmake", "../../", "-G", "Unix Makefiles")
	cmakeCommand.Dir = targetPath

	// Stderr buffer
	cmdErrOutput := &bytes.Buffer{}
	cmakeCommand.Stderr = cmdErrOutput

	if log.Verb.IsVerbose() {
		cmakeCommand.Stdout = os.Stdout
	}

	log.Norm.Cyan(false, "generating build files ... ")
	log.Verb.Write(true, "")

	err := cmakeCommand.Run()
	if err != nil {
		commands.RecordError(err, "failure", strings.Trim(cmdErrOutput.String(), "\n"))
	} else {
		log.Norm.Green(true, "success")
	}

	// execute make command
	makeCommand := exec.Command("make")
	makeCommand.Dir = targetPath

	cmdErrOutput = &bytes.Buffer{}
	makeCommand.Stderr = cmdErrOutput

	if log.Verb.IsVerbose() {
		makeCommand.Stdout = os.Stdout
	}

	log.Verb.Write(true, "")
	log.Norm.Cyan(false, "running the build ... ")
	log.Verb.Write(true, "")

	err = makeCommand.Run()
	if err != nil {
		commands.RecordError(err, "failure", strings.Trim(cmdErrOutput.String(), "\n"))
	} else {
		log.Norm.Green(true, "success")
		log.Verb.Write(true, "")
	}
}

// cleans current build files
func cleanBuildFiles(directory string, target string) {
	clean.RunClean(directory, target, true)
	log.Norm.Write(true, "")
}
