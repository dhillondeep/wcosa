// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create and update command and sub commands provided by the tool.
// Creates, updates and initializes a wio project.
package create

import (
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
    info.toLowerCase()

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
    info.printCreateSummary()
}

func (create Create) createProjectStructure(queue *log.Queue, info *createInfo) error {
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

////////////////////////////////////////////// Update /////////////////////////////////////////////////////////

func (create Create) updateApp(directory string, config *types.AppConfig) error {
    return nil
}

func (create Create) updatePackage(directory string, config *types.PkgConfig) error {
    info := &createInfo{
        Directory: directory,
        Type:      constants.PKG,
        Name:      config.MainTag.GetName(),
    }

    log.Info(log.Cyan, "updating package files ... ")
    queue := log.GetQueue()
    if err := create.updateProjectFiles(queue, info); err != nil {
        log.WriteFailure()
        return err
    }
    log.WriteSuccess()
    log.PrintQueue(queue, log.TWO_SPACES)

    log.Info(log.Cyan, "updating wio.yml ... ")
    queue = log.GetQueue()
    if err := updatePackageConfig(queue, config, info); err != nil {
        log.WriteFailure()
        return err
    }
    log.WriteSuccess()
    log.PrintQueue(queue, log.TWO_SPACES)

    log.Writeln()
    log.Info(log.Yellow.Add(color.Underline), "Project update summary")
    log.Info(log.Cyan, "path             ")
    log.Writeln(directory)
    log.Info(log.Cyan, "project type     ")
    log.Writeln("pkg")

    return nil
}

// Update Wio project
func (create Create) handleUpdate(directory string) error {
    cfg, err := utils.ReadWioConfig(directory + io.Sep + "wio.yml")
    if err != nil {
        return err
    }
    switch cfg.Type {
    case types.App:
        err = create.updateApp(directory, cfg.Config.(*types.AppConfig))
        break
    default:
        err = create.updatePackage(directory, cfg.Config.(*types.PkgConfig))
        break
    }
    return err
}

func updatePackageConfig(queue *log.Queue, config *types.PkgConfig, info *createInfo) error {
    // Ensure a minimim wio version is specified
    if strings.Trim(config.MainTag.Config.WioVersion, " ") == "" {
        return errors.String("wio.yml missing `minimum_wio_version`")
    }
    if config.MainTag.GetName() != filepath.Base(info.Directory) {
        log.Warnln(queue, "Base directory different from project name")
    }
    if config.MainTag.CompileOptions.HeaderOnly {
        config.MainTag.Flags.Visibility = "INTERFACE"
        config.MainTag.Definitions.Visibility = "INTERFACE"
    } else {
        config.MainTag.Flags.Visibility = "PRIVATE"
        config.MainTag.Definitions.Visibility = "PRIVATE"
    }
    if strings.Trim(config.MainTag.Meta.Version, " ") == "" {
        config.MainTag.Meta.Version = "0.0.1"
    }
    configPath := info.Directory + io.Sep + "wio.yml"
    return config.PrettyPrint(configPath)
}

/*// Update wio.yml file
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

        updateAVRAppTargets(&appConfig.TargetsTag, directory)

        // make current wio version as the default version if no version is provided
        if strings.Trim(appConfig.MainTag.Config.WioVersion, " ") == "" {
            appConfig.MainTag.Config.WioVersion = config.ProjectMeta.Version
        }
    } else {
    pkgConfig := projectConfig.(*types.PkgConfig)

    pkgConfig.MainTag.Meta.Name = filepath.Base(directory)

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

    if err := utils.PrettyPrintConfig(projectConfig, directory+io.Sep+"wio.yml", false); err != nil {
        return errors.WriteFileError{
            FileName: directory + io.Sep + "wio.yml",
            Err:      err,
        }
    }

    return nil
}*/

// Update AVR project files
func (create Create) updateProjectFiles(queue *log.Queue, info *createInfo) error {
    log.Verb(queue, "reading paths.json file ... ")
    structureData := &StructureConfigData{}
    if err := io.AssetIO.ParseJson("configurations/structure-avr.json", structureData); err != nil {
        log.WriteFailure(queue, log.VERB)
        return err
    } else {
        log.WriteSuccess(queue, log.VERB)
    }
    create.copyProjectAssets(queue, info, structureData.Pkg)
    return nil
}

//////////////////////////////////////////// Helpers //////////////////////////////////////////////////
func (create Create) generateConstraints() (map[string]bool, map[string]bool) {
    context := create.Context
    dirConstraints := map[string]bool{
        "tests":          false,
        "no-header-only": !context.Bool("header-only"),
    }
    fileConstraints := map[string]bool{
        "ide=clion":      false,
        "extra":          !context.Bool("no-extras"),
        "example":        context.Bool("create-example"),
        "no-header-only": !context.Bool("no-header-only"),
    }
    return dirConstraints, fileConstraints
}

func (info createInfo) fillReadMe(queue *log.Queue, readmeFile string) error {
    log.Verb(queue, "filling README file ... ")
    if err := template.IOReplace(readmeFile, map[string]string{
        "PLATFORM":        info.Platform,
        "FRAMEWORK":       info.Framework,
        "BOARD":           info.Board,
        "PROJECT_NAME":    info.Name,
        "PROJECT_VERSION": "0.0.1",
    }); err != nil {
        log.WriteFailure(queue, log.VERB)
        return err
    }
    log.WriteSuccess(queue, log.VERB)
    return nil
}

func (info createInfo) toLowerCase() {
    info.Type = strings.ToLower(info.Type)
    info.Platform = strings.ToLower(info.Platform)
    info.Framework = strings.ToLower(info.Framework)
    info.Board = strings.ToLower(info.Board)
}

// This uses a structure.json file and creates a project structure based on that. It takes in consideration
// all the constrains and copies files. This should be used for creating project for any type of app/pkg
func (create Create) copyProjectAssets(queue *log.Queue, info *createInfo, data StructureTypeData) error {
    dirConstraints, fileConstraints := create.generateConstraints()
    for _, path := range data.Paths {
        directoryPath := filepath.Clean(info.Directory + io.Sep + path.Entry)
        skipDir := false
        log.Verbln(queue, "copying assets to directory: %s", directoryPath)
        // handle directory constraints
        for _, constraint := range path.Constraints {
            _, exists := dirConstraints[constraint]
            if !exists || !dirConstraints[constraint] {
                log.Verbln(queue, "constraint not specified and hence skipping this directory")
                skipDir = true
                break
            }
        }
        if skipDir {
            continue
        }

        if !utils.PathExists(directoryPath) {
            if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
                return err
            }
            log.Verbln(queue, "created directory: %s", directoryPath)
        }

        log.Verbln(queue, "copying asset files for directory: %s", directoryPath)
        for _, file := range path.Files {
            toPath := filepath.Clean(directoryPath + io.Sep + file.To)
            skipFile := false
            // handle file constraints
            for _, constraint := range file.Constraints {
                _, exists := fileConstraints[constraint]
                if !exists || !fileConstraints[constraint] {
                    log.Verbln(queue, "constraint not specified and hence skipping this file")
                    skipFile = true
                    break
                }
            }
            if skipFile {
                continue
            }

            // handle updates
            if !file.Update && create.Update {
                log.Verbln(queue, "project is not updating, hence skipping update for path: %s", toPath)
                continue
            }

            // copy assets
            if err := io.AssetIO.CopyFile(file.From, toPath, file.Override); err != nil {
                return err
            } else {
                log.Verbln(queue, `copied asset file "%s" TO: %s: `, filepath.Base(file.From), toPath)
            }
        }
    }
    return nil
}

// Create wio.yml file for AVR project
func (create Create) fillProjectConfig(queue *log.Queue, info *createInfo) error {
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

func (info createInfo) printCreateSummary() {
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
