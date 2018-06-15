// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create and update command and sub commands provided by the tool.
// Creates, updates and initializes a wio project.
package create

import (
    "github.com/urfave/cli"
    "os"
    "path/filepath"
    "wio/cmd/wio/errors"
    goerr "errors"
    "wio/cmd/wio/log"
    "github.com/fatih/color"
    "wio/cmd/wio/config"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/types"
    "strings"
)

const (
    AVR = "avr"
)

const (
    APP = "app"
    PKG = "pkg"
)

type Create struct {
    Context *cli.Context
    Type    string
    Platform string
    Update  bool
    error   error
}

// get context for the command
func (create Create) GetContext() *cli.Context {
    return create.Context
}

// Executes the create command
func (create Create) Execute() {
    performArgumentCheck(create.Context, create.Update, create.Platform)

    var directory string
    var err error

    // fetch directory based on the argument and print logs
    if !create.Update {
        directory, err = filepath.Abs(create.Context.Args()[0])
    } else {
        if len(create.Context.Args()) > 0 {
            directory, err = filepath.Abs(create.Context.Args()[0])
        } else {
            directory, err = os.Getwd()

            log.WriteErrorlnExit(err)

            err = errors.ProgrammingArgumentAssumption{
                CommandName:  "update",
                ArgumentName: "directory",
                Err:          goerr.New("directory is not provided so current directory is used: " + directory),
            }
            log.WriteErrorln(err, true)
        }
    }

    if create.Update {
        // this checks if wio.yml file exists for it to update
        performWioExistsCheck(directory)

        // this checks if project is valid state to be updated
        performPreUpdateCheck(directory, &create)
    } else {
        // this checks if directory is empty before create can be triggered
        performPreCreateCheck(directory)
    }

    // handle update
    if create.Update {
        create.handleUpdate(directory)
    } else {
        // handle AVR creation
        if create.Platform == AVR {
            create.handleAVRCreation(directory)
        } else {
            err := errors.PlatformNotSupportedError{
                Platform: create.Platform,
            }

            log.WriteErrorlnExit(err)
        }

    }

}

func (create Create) handleUpdate(directory string) {
    board := "uno"

    // evaluate project structure and wio.yml file
    log.Write(log.INFO, color.New(color.FgCyan), "evaluating project structure and wio.yml file ... ")
    queue := log.GetQueue()

    if err := create.evaluateAVRProjectFiles(queue); err != nil {
        log.Writeln(log.NONE, color.New(color.FgGreen), "failure")
        log.PrintQueue(queue, log.TWO_SPACES)
        log.WriteErrorlnExit(err)
    } else {
        log.Writeln(log.NONE, color.New(color.FgGreen), "success")
        log.PrintQueue(queue, log.TWO_SPACES)
    }

    // show message on what happened
    log.Write(log.INFO, color.New(color.FgCyan), "updating files ... ")
    queue = log.GetQueue()

    if err := create.updateAVRProjectFiles(queue); err != nil {
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

func (create Create) handleAVRCreation(directory string) {
    board := config.ProjectDefaults.AVRBoard

    if len(create.Context.Args()) > 1 {
        board = create.Context.Args()[1]
    }

    // create project structure
    log.Write(log.INFO, color.New(color.FgCyan), "creating project structure ... ")
    queue := log.GetQueue()

    if err := create.createAVRProjectStructure(queue, directory); err != nil {
        log.Writeln(log.NONE, color.New(color.FgGreen), "failure")
        log.PrintQueue(queue, log.TWO_SPACES)
        log.WriteErrorlnExit(err)
    } else {
        log.Writeln(log.NONE, color.New(color.FgGreen), "success")
        log.PrintQueue(queue, log.TWO_SPACES)
    }

    // fill config file
    log.Write(log.INFO, color.New(color.FgCyan), "configuring project files ... ")
    queue = log.GetQueue()

    if err := create.fillAVRProjectConfig(queue, directory, board); err != nil {
        log.Writeln(log.NONE, color.New(color.FgGreen), "failure")
        log.PrintQueue(queue, log.TWO_SPACES)
        log.WriteErrorlnExit(err)
    } else {
        log.Writeln(log.NONE, color.New(color.FgGreen), "success")
        log.PrintQueue(queue, log.TWO_SPACES)
    }

    // print structure summary
    log.Writeln(log.NONE, nil, "")
    log.Writeln(log.INFO, color.New(color.FgYellow).Add(color.Underline), "Project structure summary")
    if (create.Type == PKG && !create.Context.Bool("header-only")) || create.Type == APP {
        log.Write(log.INFO, color.New(color.FgCyan), "src              ")
        log.Writeln(log.NONE, color.New(color.Reset), "source/non client files go here")
    }

    if create.Type == PKG {
        log.Write(log.INFO, color.New(color.FgCyan), "tests            ")
        log.Writeln(log.NONE, color.New(color.Reset), "source files to test the package go here")
    }

    if create.Type == PKG {
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


func (create Create) createAVRProjectStructure(queue *log.Queue, directory string) (error) {
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

    if create.Type == APP {
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

    log.QueueWrite(queue, log.VERB, color.New(color.Reset), "copying asset files ...")
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
            Err: err,
        }
    } else {
        newReadmeString := strings.Replace(string(data), "{{PLATFORM}}", create.Platform, 1)
        newReadmeString = strings.Replace(newReadmeString, "{{PROJECT_NAME}}", filepath.Base(directory), 1)
        newReadmeString = strings.Replace(newReadmeString, "{{PROJECT_VERSION}}", "0.0.1", 1)

        if err := io.NormalIO.WriteFile(directory + io.Sep + "README.md", []byte(newReadmeString)); err != nil {
            log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
            return  errors.WriteFileError{
                FileName: directory + io.Sep + "README.md",
                Err: err,
            }
        }
    }

    log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")

    return nil
}

func fillMainTagConfiguration(configurations *types.Configurations, board string, platform string, framework string) {
    // supported board, framework and platform and wio version
    configurations.SupportedBoards = append(configurations.SupportedBoards, board)
    configurations.SupportedPlatforms = append(configurations.SupportedPlatforms, AVR)
    configurations.SupportedFrameworks = append(configurations.SupportedFrameworks, framework)
    configurations.WioVersion = config.ProjectMeta.Version
}

func (create Create) fillAVRProjectConfig(queue *log.Queue, directory string, board string) (error) {
    framework := strings.ToLower(create.Context.String("framework"))

    var projectConfig types.Config

    // handle app
    if create.Type == APP {
        log.QueueWrite(queue, log.VERB, nil, "creating config file for application ... ")

        appConfig := &types.AppConfig{}

        appConfig.MainTag.Name = filepath.Base(directory)
        appConfig.MainTag.Ide = config.ProjectDefaults.Ide

        // supported board, framework and platform and wio version
        fillMainTagConfiguration(&appConfig.MainTag.Config, board, AVR, framework)

        appConfig.MainTag.CompileOptions.Platform = AVR

        // create app target
        appConfig.TargetsTag.Targets = map[string]types.AppAVRTarget{
           config.ProjectDefaults.AppTargetName: {
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
        log.QueueWrite(queue, log.VERB, nil, "creating config file for package ... ")

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
        fillMainTagConfiguration(&pkgConfig.MainTag.Config, board, AVR, framework)

        // flags
        pkgConfig.MainTag.Flags.GlobalFlags = []string{}
        pkgConfig.MainTag.Flags.RequiredFlags = []string{}
        pkgConfig.MainTag.Flags.AllowOnlyGlobalFlags = false
        pkgConfig.MainTag.Flags.AllowOnlyRequiredFlags = false

        // definitions
        pkgConfig.MainTag.Definitions.GlobalDefinitions = []string{}
        pkgConfig.MainTag.Definitions.RequiredDefinitions = []string{}
        pkgConfig.MainTag.Definitions.AllowOnlyGlobalDefinitions = false
        pkgConfig.MainTag.Definitions.AllowOnlyRequiredDefinitions = false

        // create pkg target
        pkgConfig.TargetsTag.Targets = map[string]types.PkgAVRTarget{
            config.ProjectDefaults.AppTargetName: {
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


    log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
    log.QueueWrite(queue, log.VERB, nil, "pretty printing wio.yml file ... ")


   if  err := utils.PrettyPrintConfig(projectConfig, directory+io.Sep+"wio.yml"); err != nil {
       log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
       return err
   } else {
       log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
   }

    return  nil
}

func (create Create) evaluateAVRProjectFiles(queue *log.Queue) (error) {
    return nil
}


func (create Create) updateAVRProjectFiles(queue *log.Queue) (error) {
    return nil
}


// This uses a structure.json file and creates a project structure based on that. It takes in consideration
// all the constrains and copies files. This should be used for creating project for any type of app/pkg
func copyProjectAssets(queue *log.Queue, directory string, update bool, structureTypeData StructureTypeData,
    dirConstrainsMap map[string]bool, fileConstrainsMap map[string]bool) (error) {
    for _, path := range structureTypeData.Paths {
        moveOnDir := false

        log.QueueWriteln(queue, log.VERB, nil, "copying assets to directory: " + directory + path.Entry)

        // handle directory constrains
        for _, constrain := range path.Constrains {
            if !dirConstrainsMap[constrain] {
                err := errors.ProjectStructureConstrainError{
                    Constrain: constrain,
                    Path: directory + io.Sep + path.Entry,
                    Err: goerr.New("constrained not specified and hence skipping this directory"),
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
                if !fileConstrainsMap[constrain] {
                    err := errors.ProjectStructureConstrainError{
                        Constrain: constrain,
                        Path: file.From,
                        Err: goerr.New("constrained not specified and hence skipping this file"),
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
                log.QueueWriteln(queue, log.VERB, nil, "project is not updating, hence skipping update for path: "+ toPath)
                continue
            }

            // copy assets
            if err := io.AssetIO.CopyFile(file.From, toPath, file.Override); err != nil{
                return err
            } else {
                log.QueueWriteln(queue, log.VERB, nil,
                    `copied asset file "%s" TO: %s: `, filepath.Base(file.From), directory + io.Sep + file.To)
            }
        }
    }

    return nil
}


/*





    commands.RecordError(err, "")

    createPacket := &PacketCreate{}

    // based on the command line context, get all the important fields
    createPacket.ProjType = create.Type
    createPacket.Update = create.Update
    createPacket.Directory = directory
    createPacket.Name = filepath.Base(directory)
    createPacket.Framework = create.Context.String("framework")
    createPacket.Platform = create.Context.String("platform")
    createPacket.Ide = create.Context.String("ide")
    createPacket.Tests = create.Context.Bool("tests")
    createPacket.CreateDemo = create.Context.Bool("create-demo")
    createPacket.CreateExtras = !create.Context.Bool("no-extras")

    // include header-only flag for pkg
    if createPacket.ProjType == PKG {
        createPacket.HeaderOnly = create.Context.Bool("header-only")
        createPacket.HeaderOnlyFlagSet = create.Context.IsSet("header-only")
    }

    if create.Update {
        performWioExistsCheck(directory)

        // monitor the project structure and files and decide if the update can be performed
        status, err := performPreUpdateCheck(directory, create.Type)
        commands.RecordError(err, "")

        if status {
            // all checks passed we can update
            // board is optional for update
            createPacket.Board = create.Context.String("board")

            // populate project structure and files
            createAndPopulateStructure(createPacket)
            updateProjectSetup(createPacket)
            postPrint(createPacket, true)
        }
    } else {
        // this will check if the directory is empty so that mistakenly work cannot be lost
        status, err := performPreCreateCheck(directory)
        commands.RecordError(err, "")

        if status {
            // we need board for creation
            createPacket.Board = create.Context.Args()[1]

            // populate project structure and files
            createAndPopulateStructure(createPacket)
            initialProjectSetup(createPacket)
            postPrint(createPacket, false)
        }
    }
}

// Creates project structure and based on constrains and other configurations, move template
// and required files
func createAndPopulateStructure(createPacket *PacketCreate) {
    log.Norm.Cyan(false, "creating project structure ... ")
    log.Verb.Verbose(true, "")

    structureData := &StructureConfigData{}

    // read configurationsFile
    commands.RecordError(io.AssetIO.ParseJson("configurations/structure.json", structureData),
        "failure")

    var structureTypeData StructureTypeData

    if createPacket.ProjType == APP {
        structureTypeData = structureData.App
    } else {
        structureTypeData = structureData.Pkg
    }

    for _, path := range structureTypeData.Paths {
        moveOnDir := false

        // handle directory constrains
        for _, constrain := range path.Constrains {
            if constrain == "tests" && !createPacket.Tests {
                log.Verb.Verbose(true, "tests constrain not approved for: "+path.Entry+" dir")
                moveOnDir = true
                break
            } else if constrain == "no-header-only" && createPacket.HeaderOnly {
                log.Verb.Verbose(true, "no-header-only constrain not approved for: "+
                    path.Entry+ " dir")
                moveOnDir = true
                break
            }
        }

        if moveOnDir {
            continue
        }

        directoryPath := filepath.Clean(createPacket.Directory + io.Sep + path.Entry)

        if !utils.PathExists(directoryPath) {
            commands.RecordError(os.MkdirAll(directoryPath, os.ModePerm), "failure")
        }

        for _, file := range path.Files {
            toPath := filepath.Clean(directoryPath + io.Sep + file.To)
            moveOnFile := false

            // handle file constrains
            for _, constrain := range file.Constrains {
                if constrain == "ide=clion" && createPacket.Ide != "clion" {
                    log.Verb.Verbose(true, "ide=clion constrain not approved for: "+toPath)
                    moveOnFile = true
                    break
                } else if constrain == "extra" && !createPacket.CreateExtras {
                    log.Verb.Verbose(true, "extra constrain not approved for: "+toPath)
                    moveOnFile = true
                    break
                } else if constrain == "demo" && !createPacket.CreateDemo {
                    log.Verb.Verbose(true, "demo constrain not approved for: "+toPath)
                    moveOnFile = true
                    break
                } else if constrain == "no-header-only" && createPacket.HeaderOnly {
                    log.Verb.Verbose(true, "no-header-only constrain not approved for: "+toPath)
                    moveOnFile = true
                    break
                } else if constrain == "header-only" && !createPacket.HeaderOnly {
                    log.Verb.Verbose(true, "header-only constrain not approved for: "+toPath)
                    moveOnFile = true
                    break
                }
            }

            if moveOnFile {
                continue
            }

            // handle updates
            if !file.Update && createPacket.Update {
                log.Verb.Verbose(true, "skipping for update: "+toPath)
                continue
            }

            commands.RecordError(io.AssetIO.CopyFile(file.From, toPath, file.Override), "failure")
            log.Verb.Verbose(true, "copied "+file.From+" -> "+toPath)
        }
    }

    log.Norm.Green(true, "success ")
}

/// This is one of the most important step as this sets up the project when update command is used.
/// This also updates the wio.yml file so that it can be fixed and current configurations can be applied.
func updateProjectSetup(createPacket *PacketCreate) {
    log.Norm.Cyan(false, "updating project files ... ")

    // This is used to show src directory warning for header only packages
    showSrcWarning := false

    // we have to copy README file
    var readmeStr string

    var projectConfig interface{}

    if createPacket.ProjType == APP {
        appConfig := &types.AppConfig{}

        commands.RecordError(io.NormalIO.ParseYml(createPacket.Directory+io.Sep+"wio.yml", appConfig),
            "failure")

        // update the name of the project
        appConfig.MainTag.Name = createPacket.Name
        // update the targets to make sure they are valid and there is a default target
        handleAppTargets(&appConfig.TargetsTag, config.ProjectDefaults.Board)
        // set the default board to be from the default target
        createPacket.Board = appConfig.TargetsTag.Targets[appConfig.TargetsTag.DefaultTarget].Board

        // check framework and platform
        checkFrameworkAndPlatform(&appConfig.MainTag.Framework, &appConfig.MainTag.Platform)

        projectConfig = appConfig

        // update readme app content
        readmeStr = getReadmeApp(createPacket.Name, appConfig.MainTag.Platform, appConfig.MainTag.Framework)
    } else {
        pkgConfig := &types.PkgConfig{}

        commands.RecordError(io.NormalIO.ParseYml(createPacket.Directory+io.Sep+"wio.yml", pkgConfig),
            "failure")

        // update the name of the project
        pkgConfig.MainTag.Name = createPacket.Name
        // update the targets to make sure they are valid and there is a default target
        handlePkgTargets(&pkgConfig.TargetsTag, config.ProjectDefaults.Board)
        // set the default board to be from the default target
        createPacket.Board = pkgConfig.TargetsTag.Targets[pkgConfig.TargetsTag.DefaultTarget].Board

        // only change value if the flag is set from cli
        if createPacket.HeaderOnlyFlagSet {
            pkgConfig.MainTag.HeaderOnly = createPacket.HeaderOnly

            if utils.PathExists(createPacket.Directory + io.Sep + "src") {
                showSrcWarning = true
            }
        }

        // check project version
        if pkgConfig.MainTag.Version == "" {
            pkgConfig.MainTag.Version = "0.0.1"
        }

        // make sure boards are updated in yml file
        if !utils.StringInSlice("ALL", pkgConfig.MainTag.Board) {
            pkgConfig.MainTag.Board = []string{"ALL"}
        } else if !utils.StringInSlice(createPacket.Board, pkgConfig.MainTag.Board) {
            pkgConfig.MainTag.Board = append(pkgConfig.MainTag.Board, createPacket.Board)
        }

        // check frameworks and platform
        checkFrameworkArrayAndPlatform(&pkgConfig.MainTag.Framework, &pkgConfig.MainTag.Platform)

        projectConfig = pkgConfig

        // update readme pkg content
        readmeStr = getReadmePkg(createPacket.Name, pkgConfig.MainTag.Platform,
            pkgConfig.MainTag.Framework, pkgConfig.MainTag.Board)
    }

    commands.RecordError(utils.PrettyPrintConfigSpacing(projectConfig, createPacket.Directory+io.Sep+"wio.yml"),
        "failure")

    if utils.PathExists(createPacket.Directory + io.Sep + "README.md") {
        content, err := io.NormalIO.ReadFile(createPacket.Directory + io.Sep + "README.md")
        commands.RecordError(err, "failure")

        // only write if README is empty
        if string(content) == "" || string(content) == "\n" {
            // write readme file
            io.NormalIO.WriteFile(createPacket.Directory+io.Sep+"README.md", []byte(readmeStr))
        }
    } else {
        // readme does not exist so write it
        io.NormalIO.WriteFile(createPacket.Directory+io.Sep+"README.md", []byte(readmeStr))
    }

    log.Norm.Green(true, "success")

    if showSrcWarning {
        // src folder exists and warning is showed
        log.Norm.Magenta(true, "src directory is not needed for header-only packages")
        log.Norm.Magenta(true, "it will ignored by the build system")
    }
}

// This function checks if framework and platform are not empty. It future we can in force in valid
// frameworks and platforms using this
func checkFrameworkAndPlatform(framework *string, platform *string) {
    if *framework == "" {
        *framework = config.ProjectDefaults.Framework
    }

    if *platform == "" {
        *platform = config.ProjectDefaults.Platform
    }
}

// This function is similar to the above but in this case it checks if multiple frameworks are invalid
// and same goes for platform
func checkFrameworkArrayAndPlatform(framework *[]string, platform *string) {
    if len(*framework) == 0 {
        *framework = append(*framework, config.ProjectDefaults.Framework)
    }

    if *platform == "" {
        *platform = config.ProjectDefaults.Platform
    }
}

/// This is one of the most important step as this sets up the project when create command is used.
/// This also fills up the wio.yml file so that default configuration along with user choices
/// are applied.
func initialProjectSetup(createPacket *PacketCreate) {
    log.Norm.Cyan(false, "creating project configuration ... ")

    // we have to copy README file
    var readmeStr string

    var projectConfig interface{}

    if createPacket.ProjType == APP {
        appConfig := &types.AppConfig{}
        defaultTarget := "main"

        // make modifications to the data
        appConfig.MainTag.Ide = createPacket.Ide
        appConfig.MainTag.Platform = createPacket.Platform
        appConfig.MainTag.Framework = createPacket.Framework
        appConfig.MainTag.Name = createPacket.Name
        appConfig.TargetsTag.DefaultTarget = defaultTarget
        targets := make(map[string]types.AppTargetTag, 1)
        appConfig.TargetsTag.Targets = targets

        targetsTag := appConfig.TargetsTag
        handleAppTargets(&targetsTag, createPacket.Board)

        projectConfig = appConfig

        // update readme app content
        readmeStr = getReadmeApp(createPacket.Name, appConfig.MainTag.Platform, appConfig.MainTag.Framework)
    } else {
        pkgConfig := &types.PkgConfig{}
        defaultTarget := "tests"

        // make modifications to the data
        pkgConfig.MainTag.HeaderOnly = createPacket.HeaderOnly
        pkgConfig.MainTag.Ide = createPacket.Ide
        pkgConfig.MainTag.Platform = createPacket.Platform
        pkgConfig.MainTag.Framework = []string{createPacket.Framework}
        pkgConfig.MainTag.Name = createPacket.Name
        pkgConfig.MainTag.Version = "0.0.1"
        pkgConfig.TargetsTag.DefaultTarget = defaultTarget
        targets := make(map[string]types.PkgTargetTag, 0)
        pkgConfig.TargetsTag.Targets = targets
        pkgConfig.MainTag.Board = []string{createPacket.Board}

        targetsTag := pkgConfig.TargetsTag
        handlePkgTargets(&targetsTag, createPacket.Board)

        projectConfig = pkgConfig

        // update readme pkg content
        readmeStr = getReadmePkg(createPacket.Name, pkgConfig.MainTag.Platform,
            pkgConfig.MainTag.Framework, pkgConfig.MainTag.Board)
    }

    commands.RecordError(utils.PrettyPrintConfigHelp(projectConfig, createPacket.Directory+io.Sep+"wio.yml"),
        "failure")

    // write readme file
    commands.RecordError(io.NormalIO.WriteFile(createPacket.Directory+io.Sep+"README.md", []byte(readmeStr)),
        "failure")

    log.Norm.Green(true, "success")
}

// Fill README template and return a string. This is for APP project type
func getReadmeApp(name string, platform string, framework string) string {
    readmeContent, err := io.AssetIO.ReadFile("templates/readme/APP_README.md")
    commands.RecordError(err, "failure")

    readmeStr := strings.Replace(string(readmeContent), "{{PROJECT_NAME}}", name, -1)
    readmeStr = strings.Replace(readmeStr, "{{PLATFORM}}", strings.ToUpper(platform), -1)
    readmeStr = strings.Replace(readmeStr, "{{FRAMEWORK}}", framework, -1)

    return readmeStr
}

// Fill README template and return a string. This is for PKG project type
func getReadmePkg(name string, platform string, framework []string, board []string) string {
    readmeContent, err := io.AssetIO.ReadFile("templates/readme/PKG_README.md")
    commands.RecordError(err, "failure")

    readmeStr := strings.Replace(string(readmeContent), "{{PROJECT_NAME}}", name, -1)
    readmeStr = strings.Replace(readmeStr, "{{PLATFORM}}", strings.ToUpper(platform), -1)
    readmeStr = strings.Replace(readmeStr, "{{FRAMEWORKS}}",
        strings.Join(framework, ","), -1)
    readmeStr = strings.Replace(readmeStr, "{{BOARDS}}",
        strings.Join(board, ","), -1)

    return readmeStr
}

/// This method handles the targets that a user can create and what these targets are
/// in wio.yml file. It targets are not there, it will create a default target. Unless
/// it will keep the targets that are already there. This for targets by App type
func handleAppTargets(targetsTag *types.AppTargetsTag, board string) {
    defaultTarget := types.AppTargetTag{}

    if targetsTag.Targets == nil {
        targetsTag.Targets = make(map[string]types.AppTargetTag)
    }

    if target, ok := targetsTag.Targets[targetsTag.DefaultTarget]; ok {
        defaultTarget.Board = target.Board
        defaultTarget.TargetFlags = target.TargetFlags
        targetsTag.Targets[targetsTag.DefaultTarget] = defaultTarget
    } else {
        defaultTarget.Board = board
        targetsTag.Targets[targetsTag.DefaultTarget] = defaultTarget
    }
}

/// This method handles the targets that a user can create and what these targets are
/// in wio.yml file. It targets are not there, it will create a default target. Unless
/// it will keep the targets that are already there. This for targets by Pkg type
func handlePkgTargets(targetsTag *types.PkgTargetsTag, board string) {
    defaultTarget := types.PkgTargetTag{}

    if targetsTag.Targets == nil {
        targetsTag.Targets = make(map[string]types.PkgTargetTag)
    }

    if target, ok := targetsTag.Targets[targetsTag.DefaultTarget]; ok {
        defaultTarget.Board = target.Board
        defaultTarget.TargetFlags = target.TargetFlags
        defaultTarget.PkgFlags = target.PkgFlags
        targetsTag.Targets[targetsTag.DefaultTarget] = defaultTarget
    } else {
        defaultTarget.Board = board
        targetsTag.Targets[targetsTag.DefaultTarget] = defaultTarget
    }
}
*/
