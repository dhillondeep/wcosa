// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create and update command and sub commands provided by the tool.
// Creates, updates and initializes a wio project.
package create

import (
    "path/filepath"
    "wio/cmd/wio/config"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/log"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"

    "github.com/fatih/color"
)

// Creation of AVR projects
func (create Create) createProject(dir string) error {
    info := createInfo{
        context:     create.Context,
        directory:   dir,
        projectType: create.Type,
        name:        filepath.Base(dir),
        platform:    create.Context.String("platform"),
        framework:   create.Context.String("framework"),
        board:       create.Context.String("board"),
        configOnly:  create.Context.Bool("only-config"),
        headerOnly:  create.Context.Bool("header-only"),
        updateOnly:  create.Update,
    }
    info.toLowerCase()

    // Generate project structure
    queue := log.GetQueue()
    if !info.configOnly {
        log.Info(log.Cyan, "creating project structure ... ")
        if err := createStructure(queue, &info); err != nil {
            log.WriteFailure()
            return err
        } else {
            log.WriteSuccess()
        }
        log.PrintQueue(queue, log.TWO_SPACES)
    }

    // Fill configuration file
    queue = log.GetQueue()
    log.Info(log.Cyan, "configuring project files ... ")
    if err := fillConfig(&info); err != nil {
        log.WriteFailure()
        return err
    } else {
        log.WriteSuccess()
    }
    log.PrintQueue(queue, log.TWO_SPACES)

    // print structure summary
    info.printPackageCreateSummary()
    return nil
}

// Copy and generate files for a package project
func createStructure(queue *log.Queue, info *createInfo) error {
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

    dataType := &structureData.Pkg
    if info.projectType == constants.APP {
        dataType = &structureData.App
    }

    if err := copyProjectAssets(subQueue, info, dataType); err != nil {
        log.WriteFailure(queue, log.VERB)
        log.CopyQueue(subQueue, queue, log.FOUR_SPACES)
        return err
    } else {
        log.WriteSuccess(queue, log.VERB)
        log.CopyQueue(subQueue, queue, log.FOUR_SPACES)
    }

    readmeFile := info.directory + io.Sep + "README.md"
    err := info.fillReadMe(queue, readmeFile)

    return err
}

// Generate wio.yml for app project
func fillConfig(info *createInfo) error {
    switch info.projectType {
    case constants.APP:
        return fillAppConfig(info)
    case constants.PKG:
        return fillPackageConfig(info)
    }
    return nil
}

func fillAppConfig(info *createInfo) error {
    config := &types.ConfigImpl{
        Type: constants.APP,
        Info: &types.InfoImpl{
            Name:     info.name,
            Version:  "0.0.1",
            Keywords: []string{"wio"},
            Options: &types.OptionsImpl{
                Default: "main",
                Version: config.ProjectMeta.Version,
            },
        },
        Targets: map[string]*types.TargetImpl{
            "main": &types.TargetImpl{
                Source:    "src",
                Platform:  getPlatform(info.platform),
                Framework: getFramework(info.framework),
                Board:     getBoard(info.board),
            },
        },
    }
    return utils.WriteWioConfig(info.directory, config)
}

// Generate wio.yml for package project
func fillPackageConfig(info *createInfo) error {
    config := &types.ConfigImpl{
        Type: constants.PKG,
        Info: &types.InfoImpl{
            Name:     info.name,
            Version:  "0.0.1",
            Keywords: []string{"wio"},
            Options: &types.OptionsImpl{
                Default: "tests",
                Version: config.ProjectMeta.Version,
            },
        },
        Targets: map[string]*types.TargetImpl{
            "tests": &types.TargetImpl{
                Source:    "tests",
                Platform:  getPlatform(info.platform),
                Framework: getFramework(info.framework),
                Board:     getBoard(info.board),
            },
        },
    }
    return utils.WriteWioConfig(info.directory, config)
}

// Print package creation summary
func (info createInfo) printPackageCreateSummary() {
    log.Writeln()
    log.Infoln(log.Yellow.Add(color.Underline), "Project structure summary")
    if !info.headerOnly {
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
    log.Writeln(info.directory)
    log.Info(log.Cyan, "project type     ")
    log.Writeln(info.projectType)
    log.Info(log.Cyan, "platform         ")
    log.Writeln(info.platform)
    log.Info(log.Cyan, "framework        ")
    log.Writeln(info.framework)
    log.Info(log.Cyan, "board            ")
    log.Writeln(info.board)
}
