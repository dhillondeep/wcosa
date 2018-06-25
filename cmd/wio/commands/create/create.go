// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create and update command and sub commands provided by the tool.
// Creates, updates and initializes a wio project.
package create

import (
    "github.com/fatih/color"
    "path/filepath"
    "wio/cmd/wio/config"
    "wio/cmd/wio/log"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/constants"
)
// Creation of AVR projects
func (create Create) createPackageProject(dir string) {
    info := createInfo{
        Directory:  dir,
        Type:       constants.PKG,
        Name:       filepath.Base(dir),
        Platform:   create.Context.String("platform"),
        Framework:  create.Context.String("framework"),
        Board:      create.Context.String("board"),
        ConfigOnly: create.Context.Bool("only-config"),
        HeaderOnly: create.Context.Bool("header-only"),
    }
    info.toLowerCase()

    // Generate project structure
    queue := log.GetQueue()
    if !info.ConfigOnly {
        log.Info(log.Cyan, "creating project structure ... ")
        if err := create.createPackageStructure(queue, &info); err != nil {
            log.WriteFailure()
            log.WriteErrorlnExit(err)
        } else {
            log.WriteSuccess()
        }
        log.PrintQueue(queue, log.TWO_SPACES)
    }

    // Fill configuration file
    queue = log.GetQueue()
    log.Info(log.Cyan, "configuring project files ... ")
    if err := create.fillPackageConfig(queue, &info); err != nil {
        log.WriteFailure()
        log.WriteErrorlnExit(err)
    } else {
        log.WriteSuccess()
    }
    log.PrintQueue(queue, log.TWO_SPACES)

    // print structure summary
    info.printPackageCreateSummary()
}


// Copy and generate files for a package project
func (create Create) createPackageStructure(queue *log.Queue, info *createInfo) error {
    log.Verb(queue, "reading paths.json file ... ")
    structureData := &StructureConfigData{}

    // read configurationsFile
    if err := io.AssetIO.ParseJson("configurations/structure-avr.json", structureData); err != nil {
        log.WriteFailure(queue, log.VERB)
        return err
    } else {
        log.WriteSuccess(queue, log.VERB)
    }

    log.Verb(queue, "copying asset files ... ")
    subQueue := log.GetQueue()

    if err := create.copyProjectAssets(subQueue, info, structureData.Pkg); err != nil {
        log.WriteFailure(queue, log.VERB)
        log.CopyQueue(subQueue, queue, log.FOUR_SPACES)
        return err
    } else {
        log.WriteSuccess(queue, log.VERB)
        log.CopyQueue(subQueue, queue, log.FOUR_SPACES)
    }

    readmeFile := info.Directory + io.Sep + "README.md"
    err := info.fillReadMe(queue, readmeFile)

    return err
}

// Generate wio.yml for package project
func (create Create) fillPackageConfig(queue *log.Queue, info *createInfo) error {
    /*// handle app
    if create.Type == constants.APP {
        log.QueueWrite(queue, log.INFO, nil, "creating config file for application ... ")

        appConfig := &types.AppConfig{}
        appConfig.MainTag.Name = filepath.Base(directory)
        appConfig.MainTag.Ide = config.ProjectDefaults.Ide

        // supported board, framework and platform and wio version
        fillMainTagConfiguration(&appConfig.MainTag.Config, []string{board}, constants.AVR, []string{framework})

        appConfig.MainTag.CompileOptions.Platform = constants.AVR

        // create app target
        appConfig.TargetsTag.DefaultTarget = config.ProjectDefaults.AppTargetName
        appConfig.TargetsTag.Targets = map[string]types.AppAVRTarget{
            config.ProjectDefaults.AppTargetName: {
                Src:       "src",
                Framework: framework,
                Board:     board,
                Flags: types.AppTargetFlags{
                    GlobalFlags: []string{},
                    TargetFlags: []string{},
                },
            },
        }

        projectConfig = appConfig
    } else {*/
    log.Verb(queue, "creating config file for package ... ")
    visibility := "PRIVATE"
    if info.HeaderOnly {
        visibility = "INTERFACE"
    }
    target := config.ProjectDefaults.PkgTargetName
    projectConfig := &types.PkgConfig{
        MainTag: types.PkgTag{
            Ide: config.ProjectDefaults.Ide,
            Meta: types.PackageMeta{
                Name:     info.Name,
                Version:  "0.0.1",
                License:  "MIT",
                Keywords: []string{info.Platform, info.Framework, "wio"},
            },
            CompileOptions: types.PkgCompileOptions{
                HeaderOnly: info.HeaderOnly,
                Platform:   info.Platform,
            },
            Config: types.Configurations{
                WioVersion: config.ProjectMeta.Version,
                SupportedPlatforms:  []string{info.Platform},
                SupportedFrameworks: []string{info.Framework},
                SupportedBoards:     []string{info.Board},
            },
            Flags:       types.Flags{Visibility: visibility},
            Definitions: types.Definitions{Visibility: visibility},
        },
        TargetsTag: types.PkgAVRTargets{
            DefaultTarget: target,
            Targets: map[string]types.PkgAVRTarget{
                target: {
                    Src:       target,
                    Platform:  info.Platform,
                    Framework: info.Framework,
                    Board:     info.Board,
                },
            },
        },
    }
    log.WriteSuccess(queue, log.VERB)
    log.Verb(queue, "pretty printing wio.yml file ... ")
    wioYmlPath := info.Directory + io.Sep + "wio.yml"
    if err := projectConfig.PrettyPrint(wioYmlPath); err != nil {
        log.WriteFailure(queue, log.VERB)
        return err
    }
    log.WriteSuccess(queue, log.VERB)
    return nil
}

// Print package creation summary
func (info createInfo) printPackageCreateSummary() {
    log.Writeln()
    log.Infoln(log.Yellow.Add(color.Underline), "Project structure summary")
    if !info.HeaderOnly {
        log.Info(log.Cyan, "src              ")
        log.Writeln("source/non client files")
    }
    log.Info(log.Cyan, "tests            ")
    log.Writeln("source files for test target")
    log.Info(log.Cyan, "include          ")
    log.Writeln("public headers for the package")

    // print project summary
    log.Writeln()
    log.Infoln(log.Yellow.Add(color.Underline), "Project creation summary")
    log.Info(log.Cyan, "path             ")
    log.Writeln(info.Directory)
    log.Info(log.Cyan, "project type     ")
    log.Writeln("pkg")
    log.Info(log.Cyan, "platform         ")
    log.Writeln(info.Platform)
    log.Info(log.Cyan, "framework        ")
    log.Writeln(info.Framework)
    log.Info(log.Cyan, "board            ")
    log.Writeln(info.Board)
}
