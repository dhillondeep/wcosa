// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of run package, which contains all the commands to run the project
// Builds, Uploads, and Executes the project
package run

import (
    "os"
    "runtime"
    "wio/internal/cmd/generate"
    "wio/internal/types"
    "wio/pkg/log"
    "wio/pkg/util"
    "wio/pkg/util/sys"

    "github.com/urfave/cli"
)

type Type int

type Run struct {
    Context *cli.Context
    RunType Type
}

const (
    TypeBuild Type = 0
    TypeClean Type = 1
    TypeRun   Type = 2
)

type runInfo struct {
    context *cli.Context
    config  types.Config

    directory   string
    projectType string
    headerOnly  bool
    targets     []string
    port        string

    runType Type
    jobs    int

    force  bool
    retool bool
}

type runExecuteFunc func(*runInfo, []types.Target) error

var runFuncs = []runExecuteFunc{
    (*runInfo).build,
    (*runInfo).clean,
    (*runInfo).run,
}

// get context for the command
func (run Run) GetContext() *cli.Context {
    return run.Context
}

// Runs the build, upload command (acts as one in all command)
func (run Run) Execute() error {
    directory, err := os.Getwd()
    if err != nil {
        return err
    }
    config, err := types.ReadWioConfig(directory, true)
    if err != nil {
        return err
    }
    targets := run.Context.Args()
    info := runInfo{
        context:     run.Context,
        config:      config,
        directory:   directory,
        projectType: config.GetType(),
        headerOnly:  config.GetInfo().GetOptions().GetIsHeaderOnly(),
        targets:     targets,
        port:        run.Context.String("port"),
        force:       run.Context.Bool("force"),
        retool:      run.Context.Bool("retool"),
    }
    if err := info.execute(run.RunType); err != nil {
        return err
    }
    return nil
}

func (info *runInfo) execute(runType Type) error {
    info.runType = runType

    log.Info(log.Cyan, "Reading targets... ")
    targets, err := getTargetArgs(info)
    if err != nil {
        log.WriteFailure()
        return err
    }
    log.WriteSuccess()

    return runFuncs[info.runType](info, targets)
}

func (info *runInfo) clean(targets []types.Target) error {
    targetDirs := make([]string, 0, len(targets))
    for _, target := range targets {
        targetDirs = append(targetDirs, targetPath(info, target))

        wioTimePath := sys.Path(targetPath(info, target), "wio.time")

        if sys.Exists(wioTimePath) {
            if err := os.Remove(wioTimePath); err != nil {
                return err
            }
        }
    }

    doHardClean := info.context.Bool("hard")

    log.Infoln(log.Magenta, "Running " + func() string {
        if doHardClean {
            return "hard "
        } else {
            return ""
        }
    }() +  "clean with JOBS=%d", runtime.NumCPU()+2)

    errs := asyncCleanTargets(targetDirs, doHardClean)
    if err := awaitErrors(errs); err != nil {
        return err
    }

    log.Infoln(log.Green, "|> Done!")
    return nil
}

func (info *runInfo) build(targets []types.Target) error {
    targetDirs, err := configureTargetsBuildFiles(info, targets)
    if err != nil {
        return err
    }

    log.Infoln(log.Magenta, "Running build with JOBS=%d", runtime.NumCPU()+2)
    errs := asyncBuildTargets(targetDirs)
    return awaitErrors(errs)
}

func (info *runInfo) run(targets []types.Target) error {
    target := targets[0]
    log.Info(log.Cyan, "Target: ")
    log.Infoln(log.Magenta, target.GetName())
    if !dispatchCanRunTarget(info, target) {
        if err := info.build(targets[:1]); err != nil {
            return err
        }
    }
    return dispatchRunTarget(info, target)
}

func getTargetArgs(info *runInfo) ([]types.Target, error) {
    targets := make([]types.Target, 0, len(info.targets))
    projectTargets := info.config.GetTargets()

    if info.context.Bool("all") {
        for name, target := range projectTargets {
            target.SetName(name)
            targets = append(targets, target)
        }
    } else {
        for _, name := range info.targets {
            if _, exists := projectTargets[name]; exists {
                projectTargets[name].SetName(name)
                targets = append(targets, projectTargets[name])
            } else {
                return nil, util.Error("unrecognized target %s", name)
            }
        }
        if len(info.targets) <= 0 {
            defaultName := info.config.GetInfo().GetOptions().GetDefault()
            if defaultName == "" {
                return nil, util.Error("no default target specified")
            }
            if _, exists := projectTargets[defaultName]; !exists {
                return nil, util.Error("default target %s does not exist", defaultName)
            }
            projectTargets[defaultName].SetName(defaultName)
            targets = append(targets, projectTargets[defaultName])
        }
    }
    return targets, nil
}

func configureTargetsBuildFiles(info *runInfo, targets []types.Target) ([]string, error) {
    targetDirs := make([]string, 0, len(targets))

    infoGen := &generate.InfoGenerate{
        Config:      info.config,
        Directory:   info.directory,
        ProjectType: info.projectType,
        Port:        info.port,
    }

    for _, target := range targets {
        buildStatus, err := shouldCreateBuildFiles(info.directory, target.GetName())
        if err != nil {
            return nil, err
        }

        if info.retool || info.force || buildStatus {
            log.Infoln(log.Cyan, "Generating CMake build files for target %s", target.GetName())

            if err := generate.CMakeListsFile(infoGen, target); err != nil {
                return nil, err
            }
            if err := generate.DependenciesFile(infoGen, target); err != nil {
                return nil, err
            }
            if err := generate.HardwareFile(infoGen, target); err != nil {
                return nil, err
            }
        } else {
            log.Infoln(log.Cyan, "Build files up to date for target %s", target.GetName())
        }

        targetDir := targetPath(info, target)
        if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
            return nil, err
        }

        targetDirs = append(targetDirs, targetDir)
    }
    return targetDirs, nil
}

func asyncBuildTargets(targetDirs []string) []chan error {
    var function targetFunc = configAndBuild
    return function.asyncApply(targetDirs)
}

func asyncCleanTargets(targetDirs []string, hard bool) []chan error {
    var function targetFunc = cleanIfExists
    if hard {
        function = hardClean
    }
    return function.asyncApply(targetDirs)
}

func (function targetFunc) asyncApply(targetDirs []string) []chan error {
    errs := make([]chan error, 0, len(targetDirs))
    for _, targetDir := range targetDirs {
        err := make(chan error)
        go function(targetDir, err)
        errs = append(errs, err)
    }
    return errs
}

func awaitErrors(errs []chan error) error {
    for _, errChan := range errs {
        if err := <-errChan; err != nil {
            return err
        }
    }
    return nil
}
