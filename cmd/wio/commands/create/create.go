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
)

const (
    AVR = "avr"
)

type Create struct {
    Context  *cli.Context
    Type     string
    Platform string
    Update   bool
    error    error
}

// get context for the command
func (create Create) GetContext() *cli.Context {
    return create.Context
}

// Executes the create command
func (create Create) Execute() {
    directory, board := performArgumentCheck(create.Context, create.Update, create.Platform)

    if create.Update {
        // this checks if wio.yml file exists for it to update
        performWioExistsCheck(directory)

        // this checks if project is valid state to be updated
        performPreUpdateCheck(directory, &create)

        create.handleUpdate(directory)

    } else {
        // this checks if directory is empty before create can be triggered
        performPreCreateCheck(directory, create.Context.Bool("only-config"))

        // handle AVR creation
        if create.Platform == AVR {
            create.handleAVRCreation(directory, board)
        } else {
            err := errors.PlatformNotSupportedError{
                Platform: create.Platform,
            }

            log.WriteErrorlnExit(err)
        }
    }
}

///////////////////////////////////////////// Creation ////////////////////////////////////////////////////////

// Creation of AVR projects
func (create Create) handleAVRCreation(directory string, board string) {
    onlyConfig := create.Context.Bool("only-config")

    queue := log.GetQueue()

    if !onlyConfig {
        // create project structure
        log.Write(log.INFO, color.New(color.FgCyan), "creating project structure ... ")

        if err := create.createAVRProjectStructure(queue, directory); err != nil {
            log.Writeln(log.NONE, color.New(color.FgGreen), "failure")
            log.PrintQueue(queue, log.TWO_SPACES)
            log.WriteErrorlnExit(err)
        } else {
            log.Writeln(log.NONE, color.New(color.FgGreen), "success")
            log.PrintQueue(queue, log.TWO_SPACES)
        }
    }

    if !onlyConfig {
        // fill config file
        log.Write(log.INFO, color.New(color.FgCyan), "configuring project files ... ")
    }
    queue = log.GetQueue()

    if err := create.fillAVRProjectConfig(queue, directory, board, onlyConfig); err != nil {
        if !onlyConfig {
            log.Writeln(log.NONE, color.New(color.FgGreen), "failure")
            log.PrintQueue(queue, log.TWO_SPACES)
            log.WriteErrorlnExit(err)
        } else {
            log.PrintQueue(queue, "")
            log.WriteErrorlnExit(err)
        }
    } else {
        if !onlyConfig {
            log.Writeln(log.NONE, color.New(color.FgGreen), "success")
            log.PrintQueue(queue, log.TWO_SPACES)
        } else {
            log.PrintQueue(queue, "")
        }
    }

    // print structure summary
    log.Writeln(log.NONE, nil, "")
    log.Writeln(log.INFO, color.New(color.FgYellow).Add(color.Underline), "Project structure summary")
    if (create.Type == types.PKG && !create.Context.Bool("header-only")) || create.Type == types.APP {
        log.Write(log.INFO, color.New(color.FgCyan), "src              ")
        log.Writeln(log.NONE, color.New(color.Reset), "source/non client files go here")
    }

    if create.Type == types.PKG {
        log.Write(log.INFO, color.New(color.FgCyan), "tests            ")
        log.Writeln(log.NONE, color.New(color.Reset), "source files to test the package go here")
    }

    if create.Type == types.PKG {
        log.Write(log.INFO, color.New(color.FgCyan), "include          ")
        log.Writeln(log.NONE, color.New(color.Reset), "client headers for the package go here")
    }

    // print project summary
    log.Writeln(log.NONE, nil, "")
    log.Writeln(log.INFO, color.New(color.FgYellow).Add(color.Underline), "Project creation summary")
    log.Write(log.INFO, color.New(color.FgCyan), "path             ")
    log.Writeln(log.NONE, color.New(color.Reset), directory)
    log.Write(log.INFO, color.New(color.FgCyan), "project type     ")
    log.Writeln(log.NONE, color.New(color.Reset), create.Type)
    log.Write(log.INFO, color.New(color.FgCyan), "platform         ")
    log.Writeln(log.NONE, color.New(color.Reset), create.Platform)
    log.Write(log.INFO, color.New(color.FgCyan), "board            ")
    log.Writeln(log.NONE, color.New(color.Reset), board)
}

// AVR project structure creation
func (create Create) createAVRProjectStructure(queue *log.Queue, directory string) error {
    log.QueueWrite(queue, log.VERB, color.New(color.Reset), "reading paths.json file ... ")

    structureData := &StructureConfigData{}

    // read configurationsFile
    if err := io.AssetIO.ParseJson("configurations/structure-avr.json", structureData); err != nil {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
        return err
    } else {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
    }

    var structureTypeData StructureTypeData

    if create.Type == types.APP {
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

    log.QueueWrite(queue, log.VERB, color.New(color.Reset), "copying asset files ... ")
    subQueue := log.GetQueue()

    if err := copyProjectAssets(subQueue, directory, create.Update, structureTypeData, dirConstrainsMap, fileConstrainsMap); err != nil {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
        log.CopyQueue(subQueue, queue, log.FOUR_SPACES)
        return err
    } else {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
        log.CopyQueue(subQueue, queue, log.FOUR_SPACES)
    }

    log.QueueWrite(queue, log.VERB, color.New(color.Reset), "filling README file ... ")
    if data, err := io.NormalIO.ReadFile(directory + io.Sep + "README.md"); err != nil {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
        return errors.ReadFileError{
            FileName: directory + io.Sep + "wio.yml",
            Err:      err,
        }
    } else {
        newReadmeString := strings.Replace(string(data), "{{PLATFORM}}", create.Platform, 1)
        newReadmeString = strings.Replace(newReadmeString, "{{PROJECT_NAME}}", filepath.Base(directory), 1)
        newReadmeString = strings.Replace(newReadmeString, "{{PROJECT_VERSION}}", "0.0.1", 1)

        if err := io.NormalIO.WriteFile(directory+io.Sep+"README.md", []byte(newReadmeString)); err != nil {
            log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
            return errors.WriteFileError{
                FileName: directory + io.Sep + "README.md",
                Err:      err,
            }
        }
    }

    log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")

    return nil
}

// Create wio.yml file for AVR project
func (create Create) fillAVRProjectConfig(queue *log.Queue, directory string, board string, onlyConfig bool) error {
    framework := strings.ToLower(create.Context.String("framework"))

    var projectConfig types.Config

    // handle app
    if create.Type == types.APP {
        if !onlyConfig {
            log.QueueWrite(queue, log.VERB, nil, "creating config file for application ... ")
        } else {
            log.QueueWrite(queue, log.INFO, nil, "creating config file for application ... ")
        }

        appConfig := &types.AppConfig{}

        appConfig.MainTag.Name = filepath.Base(directory)
        appConfig.MainTag.Ide = config.ProjectDefaults.Ide

        // supported board, framework and platform and wio version
        fillMainTagConfiguration(&appConfig.MainTag.Config, []string{board}, AVR, []string{framework})

        appConfig.MainTag.CompileOptions.Platform = AVR

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
    } else {
        if !onlyConfig {
            log.QueueWrite(queue, log.VERB, nil, "creating config file for package ... ")
        } else {
            log.QueueWrite(queue, log.INFO, nil, "creating config file for package ... ")
        }

        pkgConfig := &types.PkgConfig{}

        pkgConfig.MainTag.Ide = config.ProjectDefaults.Ide

        // package meta information
        pkgConfig.MainTag.Meta.Name = filepath.Base(directory)
        pkgConfig.MainTag.Meta.Version = "0.0.1"
        pkgConfig.MainTag.Meta.License = "MIT"
        pkgConfig.MainTag.Meta.Keywords = []string{AVR, "c", "c++", "wio", framework}
        pkgConfig.MainTag.Meta.Description = "A wio " + AVR + " " + create.Type + " using " + framework + " framework"

        pkgConfig.MainTag.CompileOptions.HeaderOnly = create.Context.Bool("header-only")
        pkgConfig.MainTag.CompileOptions.Platform = AVR

        // supported board, framework and platform and wio version
        fillMainTagConfiguration(&pkgConfig.MainTag.Config, []string{board}, AVR, []string{framework})

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
                Src:       "tests",
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
    }

    if !onlyConfig {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
        log.QueueWrite(queue, log.VERB, nil, "pretty printing wio.yml file ... ")
    } else {
        log.QueueWriteln(queue, log.NONE, color.New(color.FgGreen), "success")
        log.QueueWrite(queue, log.INFO, nil, "pretty printing wio.yml file ... ")
    }

    if err := utils.PrettyPrintConfig(projectConfig, directory+io.Sep+"wio.yml", true); err != nil {
        if !onlyConfig {
            log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
        } else {
            log.QueueWriteln(queue, log.NONE, color.New(color.FgRed), "failure")
        }
        return err
    } else {
        if !onlyConfig {
            log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
        } else {
            log.QueueWriteln(queue, log.NONE, color.New(color.FgGreen), "success")
        }
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

    if platform == AVR {
        // update AVR project files
        log.Write(log.INFO, color.New(color.FgCyan), "updating files for AVR platform ... ")
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
    log.Write(log.INFO, color.New(color.FgCyan), "updating wio.yml file ... ")
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
    log.Write(log.INFO, color.New(color.FgCyan), "path             ")
    log.Writeln(log.NONE, color.New(color.Reset), directory)
    log.Write(log.INFO, color.New(color.FgCyan), "project type     ")
    log.Writeln(log.NONE, color.New(color.Reset), create.Type)
    log.Write(log.INFO, color.New(color.FgCyan), "platform         ")
    log.Writeln(log.NONE, color.New(color.Reset), create.Platform)
    log.Write(log.INFO, color.New(color.FgCyan), "board            ")
    log.Writeln(log.NONE, color.New(color.Reset), board)
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
        pkgConfig.MainTag.CompileOptions.Platform = AVR

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
    log.QueueWrite(queue, log.VERB, color.New(color.Reset), "reading paths.json file ... ")

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

// Updates configurations for the project to specify supported platforms, frameworks, and boards
func fillMainTagConfiguration(configurations *types.Configurations, board []string, platform string, framework []string) {
    // supported board, framework and platform and wio version
    configurations.SupportedBoards = append(configurations.SupportedBoards, board...)
    configurations.SupportedPlatforms = append(configurations.SupportedPlatforms, platform)
    configurations.SupportedFrameworks = append(configurations.SupportedFrameworks, framework...)
    configurations.WioVersion = config.ProjectMeta.Version
}

// This uses a structure.json file and creates a project structure based on that. It takes in consideration
// all the constrains and copies files. This should be used for creating project for any type of app/pkg
func copyProjectAssets(queue *log.Queue, directory string, update bool, structureTypeData StructureTypeData,
    dirConstrainsMap map[string]bool, fileConstrainsMap map[string]bool) error {
    for _, path := range structureTypeData.Paths {
        moveOnDir := false

        log.QueueWriteln(queue, log.VERB, nil, "copying assets to directory: "+directory+path.Entry)

        // handle directory constrains
        for _, constrain := range path.Constrains {
            constrainProvided := false

            if _, exists := dirConstrainsMap[constrain]; exists {
                if dirConstrainsMap[constrain] {
                    constrainProvided = true
                } else {
                    constrainProvided = false
                }
            }

            if !constrainProvided {
                err := errors.ProjectStructureConstrainError{
                    Constrain: constrain,
                    Path:      directory + io.Sep + path.Entry,
                    Err:       goerr.New("constrained not specified and hence skipping this directory"),
                }

                log.QueueWriteln(queue, log.VERB, nil, err.Error())
                moveOnDir = true
                break
            }
        }

        if moveOnDir {
            continue
        }

        directoryPath := filepath.Clean(directory + io.Sep + path.Entry)

        if !utils.PathExists(directoryPath) {
            if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
                return err
            } else {
                log.QueueWriteln(queue, log.VERB, nil,
                    "created directory: %s", directoryPath)
            }
        }

        for _, file := range path.Files {
            log.QueueWriteln(queue, log.VERB, nil, "copying asset files for directory: %s", directoryPath)

            toPath := filepath.Clean(directoryPath + io.Sep + file.To)
            moveOnFile := false

            // handle file constrains
            for _, constrain := range file.Constrains {
                constrainProvided := false

                if _, exists := fileConstrainsMap[constrain]; exists {
                    if fileConstrainsMap[constrain] {
                        constrainProvided = true
                    } else {
                        constrainProvided = false
                    }
                }

                if !constrainProvided {
                    err := errors.ProjectStructureConstrainError{
                        Constrain: constrain,
                        Path:      file.From,
                        Err:       goerr.New("constrained not specified and hence skipping this file"),
                    }

                    log.QueueWriteln(queue, log.VERB, nil, err.Error())
                    moveOnFile = true
                    break
                }
            }

            if moveOnFile {
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
