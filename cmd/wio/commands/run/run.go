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
    "github.com/fatih/color"
    "runtime"
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

func (info runInfo) clean() error {
    return nil
}

func (info runInfo) build() error {
    info.rtype = type_build

    log.Info(log.Cyan, "Reading targets ... ")
    targets, err := getTargetArgs(&info)
    if err != nil {
        log.WriteFailure()
        return err
    }
    log.WriteSuccess()

    log.Info(log.Cyan, "Generating files ... ")
    targetDirs, err := configureTargets(&info, targets)
    if err != nil {
        log.WriteFailure()
        return err
    }
    log.WriteSuccess()

    log.Infoln(log.Cyan.Add(color.Underline), "Building targets")
    log.Infoln(log.Magenta, "Running with JOBS=%d", runtime.NumCPU() + 2)
    errs := asyncBuildTargets(targetDirs, targets)
    return awaitErrors(errs)
}

func getTargetArgs(info *runInfo) ([]types.Target, error) {
    targets := make([]types.Target, 0, len(info.targets))
    projectTargets := info.config.GetTargets().GetTargets()

    if info.context.Bool("all") {
        for name, target := range projectTargets {
            target.SetName(name)
            targets = append(targets, target)
        }
    } else {
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
                return nil, errors.Stringf("Default target [%s] does not exist", defaultName)
            }
            projectTargets[defaultName].SetName(defaultName)
            targets = append(targets, projectTargets[defaultName])
        }
    }
    return targets, nil
}

func configureTargets(info *runInfo, targets []types.Target) ([]string, error) {
    targetDirs := make([]string, 0, len(targets))
    for _, target := range targets {
        if err := dispatchCmake(info, &target); err != nil {
            return nil, err
        }
        if err := dispatchCmakeDependencies(info, &target); err != nil {
            return nil, err
        }
        targetDir := cmake.BuildPath(info.directory)
        targetDir += io.Sep + target.GetName()
        if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
            return nil, err
        }
        targetDirs = append(targetDirs, targetDir)
    }
    return targetDirs, nil
}

func asyncBuildTargets(targetDirs []string, targets []types.Target) []chan error {
    errs := make([]chan error, 0, len(targets))
    for _, targetDir := range targetDirs {
        err := make(chan error)
        go configAndBuild(targetDir, err)
        errs = append(errs, err)
    }
    return errs
}

func awaitErrors(errs []chan error) error {
    for _, errchan := range errs {
        if err := <- errchan; err != nil {
            return err
        }
    }
    return nil
}
