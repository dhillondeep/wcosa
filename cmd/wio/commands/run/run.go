// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of run package, which contains all the commands to run the project
// Builds, Uploads, and Executes the project
package run

import (
    "github.com/urfave/cli"
    "os"
    "wio/cmd/wio/commands/run/cmake"
    "wio/cmd/wio/log"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/errors"
)

type Run struct {
    Context *cli.Context
    error
}

type runType int

const (
    type_build  runType = 0
    type_clean  runType = 1
    type_run    runType = 2
    type_upload runType = 3
)

type runInfo struct {
    context *cli.Context
    config  *types.Config

    directory string
    targets   []string

    rtype runType
    jobs  int
}

// get context for the command
func (run Run) GetContext() *cli.Context {
    return run.Context
}

// Runs the build, upload command (acts as one in all command)
func (run Run) Execute() {
    directory, err := os.Getwd()
    if err != nil {
        log.WriteErrorlnExit(err)
    }
    config, err := utils.ReadWioConfig(directory + io.Sep + "wio.yml")
    if err != nil {
        log.WriteErrorlnExit(err)
    }
    targets := run.Context.Args()
    info := runInfo{
        context: run.Context,
        config: config,
        directory: directory,
        targets: targets,
    }
    if err := info.build(); err != nil {
        log.WriteErrorlnExit(err)
    }
}

func (info runInfo) build() error {
    info.rtype = type_build
    targets := make([]types.Target, 0, len(info.targets))
    projectTargets := info.config.GetTargets().GetTargets()

    for _, targetName := range info.targets {
        if _, exists := projectTargets[targetName]; exists {
            projectTargets[targetName].SetName(targetName)
            targets = append(targets, projectTargets[targetName])
        } else {
            log.Warnln("Unrecognized target name: [%s]", targetName)
        }
    }
    if len(info.targets) <= 0 {
        defaultName := info.config.GetTargets().GetDefaultTarget()
        if _, exists := projectTargets[defaultName]; !exists {
            return errors.Stringf("Default target [%s] does not exist", defaultName)
        }
        projectTargets[defaultName].SetName(defaultName)
        targets = append(targets, projectTargets[defaultName])
    }

    log.Info(log.Cyan, "Generating files ... ")
    targetDirs := make([]string, 0, len(targets))
    for _, target := range targets {
        if err := dispatchCmake(&info, &target); err != nil {
            return err
        }
        if err := dispatchCmakeDependencies(&info, &target); err != nil {
            return err
        }
        targetDir := cmake.BuildPath(info.directory)
        targetDir += io.Sep + target.GetName()
        if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
            return err
        }
        targetDirs = append(targetDirs, targetDir)
    }
    log.WriteSuccess()
    log.Infoln(log.Cyan, "Building targets ... ")
    errs := make([]chan error, 0, len(targets))
    for _, targetDir := range targetDirs {
        err := make(chan error)
        go configAndBuild(targetDir, err)
        errs = append(errs, err)
    }
    for _, errchan := range errs {
        if err := <- errchan; err != nil {
            return err
        }
    }
    return nil
}

func (run Run) configure(dir string) error {
    return nil
}
