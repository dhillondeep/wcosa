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
    "strings"
    "wio/cmd/wio/commands"
    "wio/cmd/wio/config"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/utils/io/log"
)

const (
    APP = "app"
    PKG = "pkg"
)

type Create struct {
    Context *cli.Context
    Type    string
    Update  bool
    error   error
}

// get context for the command
func (create Create) GetContext() *cli.Context {
    return create.Context
}

// Executes the create command
func (create Create) Execute() {
    performArgumentCheck(create.Context.Args(), create.Update)

    // fetch directory based on the argument
    directory, err := filepath.Abs(create.Context.Args()[0])
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

/// This method prints next steps for any type of create/update command. This will help user
/// decide what they can do next
func postPrint(createPacket *PacketCreate, update bool) {
    if !update {
        log.Norm.Yellow(true, "Project Creation finished!! Summary: ")
    } else {
        log.Norm.Yellow(true, "Project Update finished!! Summary: ")
    }

    log.Norm.Cyan(false, "project name: ")
    log.Norm.Cyan(true, createPacket.Name)
    log.Norm.Cyan(false, "project type: ")
    log.Norm.Cyan(true, createPacket.ProjType)
    log.Norm.Cyan(false, "project path: ")
    log.Norm.Cyan(true, createPacket.Directory)

    log.Norm.Yellow(true, "Check following commands:")

    log.Norm.Cyan(true, "`wio build -h`")
    log.Norm.Cyan(true, "`wio run -h`")
    log.Norm.Cyan(true, "`wio pac -h`")

    if createPacket.Tests {
        log.Norm.Cyan(true, "`wio test -h`")
    }
}
