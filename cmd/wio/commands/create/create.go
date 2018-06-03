// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create and update command and sub commands provided by the tool.
// Creates, updates and initializes a wio project.
package create

import (
    "github.com/urfave/cli"
    "path/filepath"
    "wio/cmd/wio/utils/io/log"
    "os"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "strings"
    "wio/cmd/wio/commands"
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
func (create Create) GetContext() (*cli.Context) {
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
    createPacket.CreateTemplates = create.Context.Bool("create-templates")
    createPacket.CreateExtras = !create.Context.Bool("no-extras")

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

/// This method is a helper function that prints instructions to the console before "create"
/// or "update" command is executed. This has information about the project and it's ty[e
func prePrint(createPacket *PacketCreate) {
    log.Norm.Yellow(false, "Project Name: ")
    log.Norm.Cyan(true, createPacket.Name)
    log.Norm.Yellow(false, "Project Type: ")
    log.Norm.Cyan(true, createPacket.ProjType)
    log.Norm.Yellow(false, "Project Path: ")
    log.Norm.Cyan(true, createPacket.Directory)
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
                log.Verb.Verbose(true, "tests constrain not approved for: " + path.Entry + " dir")
                moveOnDir = true
                break
            }
        }

        if moveOnDir { continue}


        directoryPath := filepath.Clean(createPacket.Directory+io.Sep+path.Entry)

        if !utils.PathExists(directoryPath) {
            commands.RecordError(os.MkdirAll(directoryPath, os.ModePerm), "failure")
        }

        for _, file := range path.Files {
            toPath := filepath.Clean(directoryPath+io.Sep+file.To)
            moveOnFile := false

            // handle file constrains
            for _, constrain := range file.Constrains {
                if constrain == "ide=clion" && createPacket.Ide != "clion" {
                    log.Verb.Verbose(true, "ide=clion constrain not approved for: " + toPath)
                    moveOnFile = true
                    break
                } else if constrain == "extra" && !createPacket.CreateExtras {
                    log.Verb.Verbose(true, "extra constrain not approved for: " + toPath)
                    moveOnFile = true
                    break
                } else if constrain == "template" && !createPacket.CreateTemplates {
                    log.Verb.Verbose(true, "template constrain not approved for: " + toPath)
                    moveOnFile = true
                    break
                }
            }

            if moveOnFile { continue }

            // handle updates
            if !file.Update && createPacket.Update {
                log.Verb.Verbose(true, "skipping for update: " + toPath)
                continue
            }

            commands.RecordError(io.AssetIO.CopyFile(file.From, toPath, file.Override), "failure")
            log.Verb.Verbose(true, "copied " + file.From + " -> " + toPath)
        }
    }

    log.Norm.Green(true, "success ")
}

/// This is one of the most important step as this sets up the project when update command is used.
/// This also updates the wio.yml file so that it can be fixed and current configurations can be applied.
func updateProjectSetup(createPacket *PacketCreate) {
    log.Norm.Cyan(false, "updating project files ... ")

    // get default configuration values
    defaults := types.DConfig{}
    commands.RecordError(io.AssetIO.ParseYml("config/defaults.yml", &defaults), "failure")

    // we have to copy README file
    var readmeStr string

    var config interface{}

    if createPacket.ProjType == APP {
        projectConfig := &types.AppConfig{}

        commands.RecordError(io.NormalIO.ParseYml(createPacket.Directory+io.Sep+"wio.yml", projectConfig),
            "failure")

        // update the name of the project
        projectConfig.MainTag.Name = createPacket.Name
        // update the targets to make sure they are valid and there is a default target
        handleAppTargets(&projectConfig.TargetsTag, defaults.Board)
        // set the default board to be from the default target
        createPacket.Board = projectConfig.TargetsTag.Targets[projectConfig.TargetsTag.DefaultTarget].Board

        // check framework and platform
        checkFrameworkAndPlatform(&projectConfig.MainTag.Framework,
            &projectConfig.MainTag.Platform, &defaults)

        config = projectConfig

        // update readme app content
        readmeStr = getReadmeApp(createPacket.Name, projectConfig.MainTag.Platform, projectConfig.MainTag.Framework)
    } else {
        projectConfig := &types.PkgConfig{}

        commands.RecordError(io.NormalIO.ParseYml(createPacket.Directory+io.Sep+"wio.yml", projectConfig),
            "failure")

        // update the name of the project
        projectConfig.MainTag.Name = createPacket.Name
        // update the targets to make sure they are valid and there is a default target
        handlePkgTargets(&projectConfig.TargetsTag, defaults.Board)
        // set the default board to be from the default target
        createPacket.Board = projectConfig.TargetsTag.Targets[projectConfig.TargetsTag.DefaultTarget].Board

        // check project version
        if projectConfig.MainTag.Version == "" {
            projectConfig.MainTag.Version = "0.0.1"
        }

        // make sure boards are updated in yml file
        if !utils.StringInSlice("ALL", projectConfig.MainTag.Board) {
            projectConfig.MainTag.Board = []string{"ALL"}
        } else if !utils.StringInSlice(createPacket.Board, projectConfig.MainTag.Board) {
            projectConfig.MainTag.Board = append(projectConfig.MainTag.Board, createPacket.Board)
        }

        // check frameworks and platform
        checkFrameworkArrayAndPlatform(&projectConfig.MainTag.Framework,
            &projectConfig.MainTag.Platform, &defaults)

        config = projectConfig

        // update readme pkg content
        readmeStr = getReadmePkg(createPacket.Name, projectConfig.MainTag.Platform,
            projectConfig.MainTag.Framework, projectConfig.MainTag.Board)
    }

    commands.RecordError(utils.PrettyPrintConfigSpacing(config, createPacket.Directory+io.Sep+"wio.yml"),
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
}

// This function checks if framework and platform are not empty. It future we can in force in valid
// frameworks and platforms using this
func checkFrameworkAndPlatform(framework *string, platform *string, defaults *types.DConfig) {
    if *framework == "" {
        *framework = defaults.Framework
    }

    if *platform == "" {
        *platform = defaults.Platform
    }
}

// This function is similar to the above but in this case it checks if multiple frameworks are invalid
// and same goes for platform
func checkFrameworkArrayAndPlatform(framework *[]string, platform *string, defaults *types.DConfig) {
    if len(*framework) == 0 {
        *framework = append(*framework, defaults.Framework)
    }

    if *platform == "" {
        *platform = defaults.Platform
    }
}

/// This is one of the most important step as this sets up the project when create command is used.
/// This also fills up the wio.yml file so that default configuration along with user choices
/// are applied.
func initialProjectSetup(createPacket *PacketCreate) {
    log.Norm.Cyan(false, "creating project configuration ... ")

    // we have to copy README file
    var readmeStr string

    var config interface{}

    if createPacket.ProjType == APP {
        projectConfig := &types.AppConfig{}
        defaultTarget := "main"

        // make modifications to the data
        projectConfig.MainTag.Ide = createPacket.Ide
        projectConfig.MainTag.Platform = createPacket.Platform
        projectConfig.MainTag.Framework = createPacket.Framework
        projectConfig.MainTag.Name = createPacket.Name
        projectConfig.TargetsTag.DefaultTarget = defaultTarget
        targets := make(map[string]types.AppTargetTag, 1)
        projectConfig.TargetsTag.Targets = targets

        targetsTag := projectConfig.TargetsTag
        handleAppTargets(&targetsTag, createPacket.Board)

        config = projectConfig

        // update readme app content
        readmeStr = getReadmeApp(createPacket.Name, projectConfig.MainTag.Platform, projectConfig.MainTag.Framework)
    } else {
        projectConfig := &types.PkgConfig{}
        defaultTarget := "tests"

        // make modifications to the data
        projectConfig.MainTag.Ide = createPacket.Ide
        projectConfig.MainTag.Platform = createPacket.Platform
        projectConfig.MainTag.Framework = []string{createPacket.Framework}
        projectConfig.MainTag.Name = createPacket.Name
        projectConfig.MainTag.Version = "0.0.1"
        projectConfig.TargetsTag.DefaultTarget = defaultTarget
        targets := make(map[string]types.PkgTargetTag, 0)
        projectConfig.TargetsTag.Targets = targets
        projectConfig.MainTag.Board = []string{createPacket.Board}

        targetsTag := projectConfig.TargetsTag
        handlePkgTargets(&targetsTag, createPacket.Board)

        config = projectConfig

        // update readme pkg content
        readmeStr = getReadmePkg(createPacket.Name, projectConfig.MainTag.Platform,
            projectConfig.MainTag.Framework, projectConfig.MainTag.Board)
    }

    commands.RecordError(utils.PrettyPrintConfigHelp(config, createPacket.Directory+io.Sep+"wio.yml"),
        "failure")

    // write readme file
    commands.RecordError(io.NormalIO.WriteFile(createPacket.Directory+io.Sep+"README.md", []byte(readmeStr)),
        "failure")

    log.Norm.Green(true, "success")
}

// Fill README template and return a string. This is for APP project type
func getReadmeApp(name string, platform string, framework string) (string) {
    readmeContent, err := io.AssetIO.ReadFile("templates/readme/APP_README.md")
    commands.RecordError(err, "failure")

    readmeStr := strings.Replace(string(readmeContent), "{{PROJECT_NAME}}", name, -1)
    readmeStr = strings.Replace(readmeStr, "{{PLATFORM}}", strings.ToUpper(platform), -1)
    readmeStr = strings.Replace(readmeStr, "{{FRAMEWORK}}", framework, -1)

    return readmeStr
}

// Fill README template and return a string. This is for PKG project type
func getReadmePkg(name string, platform string, framework []string, board []string) (string) {
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
        defaultTarget.TargetCompileFlags = target.TargetCompileFlags
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
        defaultTarget.TargetCompileFlags = target.TargetCompileFlags
        defaultTarget.PkgCompileFlags = target.PkgCompileFlags
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
