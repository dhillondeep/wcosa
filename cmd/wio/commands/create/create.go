// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create and update command and sub commands provided by the tool.
// Creates, updates and initializes a wio project.
package create

import (
    goerr "errors"
    "github.com/fatih/color"
    "github.com/urfave/cli"
    "os"
    "path/filepath"
    "strings"
    "wio/cmd/wio/config"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/log"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/utils/template"
)

type Create struct {
    Context *cli.Context
    Type    string
    Update  bool
    error   error
}

type createInfo struct {
    Directory string
    Type      string
    Name      string

    Platform  string
    Framework string
    Board     string

    ConfigOnly bool
    HeaderOnly bool
}

// get context for the command
func (create Create) GetContext() *cli.Context {
    return create.Context
}

// Executes the create command
func (create Create) Execute() {
    directory := performDirectoryCheck(create.Context)

    if create.Update {
        // this checks if wio.yml file exists for it to update
        performWioExistsCheck(directory)
        // this checks if project is valid state to be updated
        performPreUpdateCheck(directory, &create)
        create.handleUpdate(directory)

    } else {
        // this checks if directory is empty before create can be triggered
        performPreCreateCheck(directory, create.Context.Bool("only-config"))
        create.handleCreation(directory)
    }
}

///////////////////////////////////////////// Creation ////////////////////////////////////////////////////////

// Creation of AVR projects
func (create Create) handleCreation(dir string) {
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
    infoLowerCase(&info)

    // Generate project structure
    queue := log.GetQueue()
    if !info.ConfigOnly {
        log.Info(log.Cyan, "creating project structure ... ")
        if err := create.createProjectStructure(queue, &info); err != nil {
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
    if err := create.fillProjectConfig(queue, &info); err != nil {
        log.WriteFailure()
        log.WriteErrorlnExit(err)
    } else {
        log.WriteSuccess()
    }
    log.PrintQueue(queue, log.TWO_SPACES)

    // print structure summary
    log.Writeln()
    log.Info(log.Yellow.Add(color.Underline), "Project structure summary")
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
    log.Writeln(create.Type)
    log.Info(log.Cyan, "platform         ")
    log.Writeln(info.Platform)
    log.Info(log.Cyan, "framework        ")
    log.Writeln(info.Framework)
    log.Info(log.Cyan, "board            ")
    log.Writeln(info.Board)
}

// AVR project structure creation
func (create Create) createProjectStructure(queue *log.Queue, info *createInfo) error {
    log.QueueWrite(queue, log.VERB, log.Default, "reading paths.json file ... ")
    structureData := &StructureConfigData{}

    // read configurationsFile
    if err := io.AssetIO.ParseJson("configurations/structure-avr.json", structureData); err != nil {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
        return err
    } else {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
    }

    var structureTypeData StructureTypeData

    if create.Type == constants.APP {
        structureTypeData = structureData.App
    } else {
        structureTypeData = structureData.Pkg
    }

    // constrains for AVR project
    dirConstrainsMap := map[string]bool{}
    dirConstrainsMap["tests"] = false
    dirConstrainsMap["no-header-only"] = !create.Context.Bool("header-only")

    fileConstrainsMap := map[string]bool{}
    fileConstrainsMap["ide=clion"] = false
    fileConstrainsMap["extra"] = !create.Context.Bool("no-extras")
    fileConstrainsMap["example"] = create.Context.Bool("create-example")
    fileConstrainsMap["no-header-only"] = !create.Context.Bool("header-only")

    log.QueueWrite(queue, log.VERB, log.Default, "copying asset files ... ")
    subQueue := log.GetQueue()

    if err := copyProjectAssets(subQueue, directory, create.Update, structureTypeData, dirConstrainsMap, fileConstrainsMap); err != nil {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
        log.CopyQueue(subQueue, queue, log.FOUR_SPACES)
        return err
    } else {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
        log.CopyQueue(subQueue, queue, log.FOUR_SPACES)
    }

    readmeFile := info.Directory + io.Sep + "README.md"
    err := fillReadMe(queue, readmeFile, info)

    return err
}

// Create wio.yml file for AVR project
func (create Create) fillProjectConfig(
    queue *log.Queue, directory string,
    platform string, framework string, board string) error {

    var projectConfig types.Config

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
    log.QueueWrite(queue, log.VERB, nil, "creating config file for package ... ")

    pkgConfig := &types.PkgConfig{}

    pkgConfig.MainTag.Ide = config.ProjectDefaults.Ide

    // package meta information
    pkgConfig.MainTag.Meta.Name = filepath.Base(directory)
    pkgConfig.MainTag.Meta.Version = "0.0.1"
    pkgConfig.MainTag.Meta.License = "MIT"
    pkgConfig.MainTag.Meta.Keywords = []string{platform, "c", "c++", "wio", framework}
    pkgConfig.MainTag.Meta.Description = "A wio " + platform + " " + create.Type + " using " + framework + " framework"

    pkgConfig.MainTag.CompileOptions.HeaderOnly = create.Context.Bool("header-only")
    pkgConfig.MainTag.CompileOptions.Platform = platform

    // supported board, framework and platform and wio version
    fillMainTagConfiguration(&pkgConfig.MainTag.Config, []string{board}, platform, []string{framework})

    // flags
    pkgConfig.MainTag.Flags.GlobalFlags = []string{}
    pkgConfig.MainTag.Flags.RequiredFlags = []string{}
    pkgConfig.MainTag.Flags.AllowOnlyGlobalFlags = false
    pkgConfig.MainTag.Flags.AllowOnlyRequiredFlags = false

    if pkgConfig.MainTag.CompileOptions.HeaderOnly {
        pkgConfig.MainTag.Flags.Visibility = "INTERFACE"
    } else {
        pkgConfig.MainTag.Flags.Visibility = "PRIVATE"
    }

    // definitions
    pkgConfig.MainTag.Definitions.GlobalDefinitions = []string{}
    pkgConfig.MainTag.Definitions.RequiredDefinitions = []string{}
    pkgConfig.MainTag.Definitions.AllowOnlyGlobalDefinitions = false
    pkgConfig.MainTag.Definitions.AllowOnlyRequiredDefinitions = false

    if pkgConfig.MainTag.CompileOptions.HeaderOnly {
        pkgConfig.MainTag.Definitions.Visibility = "INTERFACE"
    } else {
        pkgConfig.MainTag.Definitions.Visibility = "PRIVATE"
    }

    // create pkg target
    pkgConfig.TargetsTag.DefaultTarget = config.ProjectDefaults.PkgTargetName
    pkgConfig.TargetsTag.Targets = map[string]types.PkgAVRTarget{
        config.ProjectDefaults.PkgTargetName: {
            Src:       config.ProjectDefaults.PkgTargetName,
            Framework: framework,
            Board:     board,
            Flags: types.PkgTargetFlags{
                GlobalFlags: []string{},
                TargetFlags: []string{},
                PkgFlags:    []string{},
            },
        },
    }

    projectConfig = pkgConfig

    log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
    log.QueueWrite(queue, log.VERB, nil, "pretty printing wio.yml file ... ")

    if err := utils.PrettyPrintConfig(projectConfig, directory+io.Sep+"wio.yml", true); err != nil {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
        return err
    } else {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
    }

    return nil
}

////////////////////////////////////////////// Update /////////////////////////////////////////////////////////

// Update Wio project
func (create Create) handleUpdate(directory string) {
    board := "uno"

    projectConfig, err := utils.ReadWioConfig(directory + io.Sep + "wio.yml")
    if err != nil {
        log.WriteErrorlnExit(err)
    }

    platform := projectConfig.GetMainTag().GetCompileOptions().GetPlatform()

    if platform == constants.AVR {
        // update AVR project files
        log.Info(log.Cyan, "updating files for AVR platform ... ")
        queue := log.GetQueue()

        if err := create.updateAVRProjectFiles(queue, directory); err != nil {
            log.Writeln(log.NONE, color.New(color.FgGreen), "failure")
            log.PrintQueue(queue, log.TWO_SPACES)
            log.WriteErrorlnExit(err)
        } else {
            log.Writeln(log.NONE, color.New(color.FgGreen), "success")
            log.PrintQueue(queue, log.TWO_SPACES)
        }
    }

    // update wio.yml file
    log.Info(log.Cyan, "updating wio.yml file ... ")
    queue := log.GetQueue()

    if err := create.updateConfig(queue, projectConfig, directory); err != nil {
        log.Writeln(log.NONE, color.New(color.FgGreen), "failure")
        log.PrintQueue(queue, log.TWO_SPACES)
        log.WriteErrorlnExit(err)
    } else {
        log.Writeln(log.NONE, color.New(color.FgGreen), "success")
        log.PrintQueue(queue, log.TWO_SPACES)
    }

    // print update summary
    log.Writeln(log.NONE, nil, "")
    log.Writeln(log.INFO, color.New(color.FgYellow).Add(color.Underline), "Project update summary")
    log.Info(log.Cyan, "path             ")
    log.Writeln(log.NONE, directory)
    log.Info(log.Cyan, "project type     ")
    log.Writeln(log.NONE, create.Type)
    log.Info(log.Cyan, "platform         ")
    log.Writeln(log.NONE, platform)
    log.Info(log.Cyan, "board            ")
    log.Writeln(log.NONE, board)
}

// Update wio.yml file
func (create Create) updateConfig(queue *log.Queue, projectConfig types.Config, directory string) error {
    isApp, err := utils.IsAppType(directory + io.Sep + "wio.yml")
    if err != nil {
        return errors.ConfigParsingError{
            Err: err,
        }
    }

    if isApp {
        appConfig := projectConfig.(*types.AppConfig)

        appConfig.MainTag.Name = filepath.Base(directory)

        //////////////////////////////////////////// Targets //////////////////////////////////////////////////
        updateAVRAppTargets(&appConfig.TargetsTag, directory)

        // make current wio version as the default version if no version is provided
        if strings.Trim(appConfig.MainTag.Config.WioVersion, " ") == "" {
            appConfig.MainTag.Config.WioVersion = config.ProjectMeta.Version
        }
    } else {
        pkgConfig := projectConfig.(*types.PkgConfig)

        pkgConfig.MainTag.Meta.Name = filepath.Base(directory)

        //////////////////////////////////////////// Targets //////////////////////////////////////////////////
        updateAVRPkgTargets(&pkgConfig.TargetsTag, directory)

        // make sure framework is AVR
        pkgConfig.MainTag.CompileOptions.Platform = constants.AVR

        // make current wio version as the default version if no version is provided
        if strings.Trim(pkgConfig.MainTag.Config.WioVersion, " ") == "" {
            pkgConfig.MainTag.Config.WioVersion = config.ProjectMeta.Version
        }

        // keywords
        pkgConfig.MainTag.Meta.Keywords = utils.AppendIfMissing(pkgConfig.MainTag.Meta.Keywords,
            []string{"wio", "c", "c++"})
        pkgConfig.MainTag.Meta.Keywords = utils.AppendIfMissing(pkgConfig.MainTag.Meta.Keywords,
            pkgConfig.MainTag.Config.SupportedFrameworks)
        pkgConfig.MainTag.Meta.Keywords = utils.AppendIfMissing(pkgConfig.MainTag.Meta.Keywords,
            pkgConfig.MainTag.Config.SupportedPlatforms)
        pkgConfig.MainTag.Meta.Keywords = utils.AppendIfMissing(pkgConfig.MainTag.Meta.Keywords,
            pkgConfig.MainTag.Config.SupportedBoards)

        // flags
        if strings.Trim(pkgConfig.MainTag.Flags.Visibility, " ") == "" {
            if pkgConfig.MainTag.CompileOptions.HeaderOnly {
                pkgConfig.MainTag.Flags.Visibility = "INTERFACE"
            } else {
                pkgConfig.MainTag.Flags.Visibility = "PRIVATE"
            }
        }

        // definitions
        if strings.Trim(pkgConfig.MainTag.Definitions.Visibility, " ") == "" {
            if pkgConfig.MainTag.CompileOptions.HeaderOnly {
                pkgConfig.MainTag.Definitions.Visibility = "INTERFACE"
            } else {
                pkgConfig.MainTag.Definitions.Visibility = "PRIVATE"
            }
        }

        // version
        if strings.Trim(pkgConfig.MainTag.Meta.Version, " ") == "" {
            pkgConfig.MainTag.Meta.Version = "0.0.1"
        }
    }

    if err := utils.PrettyPrintConfig(projectConfig, directory+io.Sep+"wio.yml", create.Context.Bool("config-help")); err != nil {
        return errors.WriteFileError{
            FileName: directory + io.Sep + "wio.yml",
            Err:      err,
        }
    }

    return nil
}

// Update AVR project files
func (create Create) updateAVRProjectFiles(queue *log.Queue, directory string) error {
    log.QueueWrite(queue, log.VERB, log.Default, "reading paths.json file ... ")

    isApp, err := utils.IsAppType(directory + io.Sep + "wio.yml")
    if err != nil {
        return errors.ConfigParsingError{
            Err: err,
        }
    }

    structureData := &StructureConfigData{}

    // read configurationsFile
    if err := io.AssetIO.ParseJson("configurations/structure-avr.json", structureData); err != nil {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
        return err
    } else {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
    }

    dirConstrainsMap := map[string]bool{}
    dirConstrainsMap["tests"] = false

    fileConstrainsMap := map[string]bool{}
    fileConstrainsMap["extra"] = !create.Context.Bool("no-extras")

    if isApp {
        // update AVR application
        copyProjectAssets(queue, directory, true, structureData.App, dirConstrainsMap, fileConstrainsMap)
    } else {
        // update AVR package
        copyProjectAssets(queue, directory, true, structureData.Pkg, dirConstrainsMap, fileConstrainsMap)
    }

    return nil
}

func fillReadMe(queue *log.Queue, readmeFile string, info *createInfo) error {
    log.QueueWrite(queue, log.VERB, log.Default, "filling README file ... ")
    if nil != template.IOReplace(readmeFile, map[string]string{
        "PLATFORM":        info.Platform,
        "FRAMEWORK":       info.Framework,
        "BOARD":           info.Board,
        "PROJECT_NAME":    info.Name,
        "PROJECT_VERSION": "0.0.1",
    }) {
        log.WriteFailure(queue, log.VERB_NONE)
    }
    log.WriteSuccess(queue, log.VERB_NONE)
    return nil
}

func infoLowerCase(info *createInfo) {
    info.Type = strings.ToLower(info.Type)
    info.Platform = strings.ToLower(info.Platform)
    info.Framework = strings.ToLower(info.Framework)
    info.Board = strings.ToLower(info.Board)
}

// Updates configurations for the project to specify supported platforms, frameworks, and boards
func fillMainTagConfiguration(configurations *types.Configurations, boards []string, platform string, frameworks []string) {
    // supported boards, frameworks and platform and wio version
    if utils.ContainsNoCase(configurations.SupportedBoards, "all") {
        configurations.SupportedBoards = []string{"all"}
    } else {
        configurations.SupportedBoards = append(configurations.SupportedBoards, boards...)
    }
    if utils.ContainsNoCase(configurations.SupportedFrameworks, "all") {
        configurations.SupportedFrameworks = []string{"all"}
    } else {
        configurations.SupportedFrameworks = append(configurations.SupportedFrameworks, frameworks...)
    }
    if utils.ContainsNoCase(configurations.SupportedPlatforms, "all") {
        configurations.SupportedPlatforms = []string{"all"}
    } else {
        configurations.SupportedPlatforms = append(configurations.SupportedPlatforms, platform)
    }
    configurations.WioVersion = config.ProjectMeta.Version
}

// This uses a structure.json file and creates a project structure based on that. It takes in consideration
// all the constrains and copies files. This should be used for creating project for any type of app/pkg
func copyProjectAssets(
    queue *log.Queue, directory string, update bool, structureTypeData StructureTypeData,
    dirConstraintMap map[string]bool, fileConstraintMap map[string]bool) error {
    for _, path := range structureTypeData.Paths {
        skipDir := false

        log.QueueWriteln(queue, log.VERB, nil, "copying assets to directory: "+directory+path.Entry)

        // handle directory constrains
        for _, constraint := range path.Constrains {
            _, exists := dirConstraintMap[constraint]
            constrainProvided := exists && dirConstraintMap[constraint]

            if !constrainProvided {
                err := errors.ProjectStructureConstraintError{
                    Constraint: constraint,
                    Path:       directory + io.Sep + path.Entry,
                    Err:        goerr.New("constrained not specified and hence skipping this directory"),
                }

                log.QueueWriteln(queue, log.VERB, nil, err.Error())
                skipDir = true
                break
            }
        }

        if skipDir {
            continue
        }

        directoryPath := filepath.Clean(directory + io.Sep + path.Entry)

        if !utils.PathExists(directoryPath) {
            if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
                return err
            } else {
                log.QueueWriteln(queue, log.VERB, nil, "created directory: %s", directoryPath)
            }
        }

        for _, file := range path.Files {
            log.QueueWriteln(queue, log.VERB, nil, "copying asset files for directory: %s", directoryPath)

            toPath := filepath.Clean(directoryPath + io.Sep + file.To)
            skipFile := false

            // handle file constrains
            for _, constraint := range file.Constraints {
                _, exists := fileConstraintMap[constraint]
                constraintProvided := exists && fileConstraintMap[constraint]

                if !constraintProvided {
                    err := errors.ProjectStructureConstraintError{
                        Constraint: constraint,
                        Path:       file.From,
                        Err:        goerr.New("constraint not specified and hence skipping this file"),
                    }

                    log.QueueWriteln(queue, log.VERB, nil, err.Error())
                    skipFile = true
                    break
                }
            }

            if skipFile {
                continue
            }

            // handle updates
            if !file.Update && update {
                log.QueueWriteln(queue, log.VERB, nil, "project is not updating, hence skipping update for path: "+toPath)
                continue
            }

            // copy assets
            if err := io.AssetIO.CopyFile(file.From, toPath, file.Override); err != nil {
                return err
            } else {
                log.QueueWriteln(queue, log.VERB, nil,
                    `copied asset file "%s" TO: %s: `, filepath.Base(file.From), directory+io.Sep+file.To)
            }
        }
    }

    return nil
}
