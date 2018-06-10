// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package main contains the main code for Wio.
// Wio is a tool to make development of embedded system applications easier and simpler.
// It allows for building, testing, and uploading AVR applications for Commandline.
package main

import (
    "github.com/urfave/cli"
    "os"
    "path/filepath"
    "time"
    "wio/cmd/wio/commands"
    "wio/cmd/wio/commands/build"
    "wio/cmd/wio/commands/clean"
    "wio/cmd/wio/commands/create"
    "wio/cmd/wio/commands/pac"
    "wio/cmd/wio/commands/run"
    "wio/cmd/wio/config"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/utils/io/log"
)

func main() {
    // read help templates
    appHelpText, err := io.AssetIO.ReadFile("cli-helper/app-help.txt")
    commands.RecordError(err, "")

    commandHelpText, err := io.AssetIO.ReadFile("cli-helper/command-help.txt")
    commands.RecordError(err, "")

    subCommandHelpText, err := io.AssetIO.ReadFile("cli-helper/subcommand-help.txt")
    commands.RecordError(err, "")

    // override help templates
    cli.AppHelpTemplate = string(appHelpText)
    cli.CommandHelpTemplate = string(commandHelpText)
    cli.SubcommandHelpTemplate = string(subCommandHelpText)

    // command that will be executed
    var command commands.Command

    app := cli.NewApp()
    app.Name = config.ProjectMeta.Name
    app.Version = config.ProjectMeta.Version
    app.EnableBashCompletion = config.ProjectMeta.EnableBashCompletion
    app.Compiled = time.Now()
    app.Copyright = config.ProjectMeta.Copyright
    app.Usage = config.ProjectMeta.UsageText

    app.Commands = []cli.Command{
        {
            Name:  "create",
            Usage: "Creates and initializes a wio project.",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:      "pkg",
                    Usage:     "Creates a wio package, intended to be used by other people",
                    UsageText: "wio create pkg <DIRECTORY> <BOARD> [command options]",
                    Flags: []cli.Flag{
                        cli.BoolFlag{Name: "header-only",
                            Usage: "This flag can be used to specify that the package is header only"},
                        cli.StringFlag{Name: "framework",
                            Usage: "Framework being used for this project. Framework contains the core libraries",
                            Value: config.ProjectDefaults.Framework},
                        cli.StringFlag{Name: "platform",
                            Usage: "Platform being used for this project. Platform is the type of chip supported (AVR/ ARM)",
                            Value: config.ProjectDefaults.Platform},
                        cli.StringFlag{Name: "ide",
                            Usage: "Creates the project for a specified IDE (CLion, Eclipse, VS Code)",
                            Value: config.ProjectDefaults.Ide},
                        cli.BoolFlag{Name: "create-demo",
                            Usage: "This will create a demo project that user can build and upload"},
                        cli.BoolFlag{Name: "no-extras",
                            Usage: "This will restrict wio from creating .gitignore, README.md, etc files"},
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed"},
                    },
                    Action: func(c *cli.Context) {
                        command = create.Create{Context: c, Type: create.PKG, Update: false}
                    },
                },
                cli.Command{
                    Name:      "app",
                    Usage:     "Creates a wio application, intended to be compiled and uploaded to a device",
                    UsageText: "wio create app <DIRECTORY> <BOARD> [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "framework",
                            Usage: "Framework being used for this project. Framework contains the core libraries",
                            Value: config.ProjectDefaults.Framework},
                        cli.StringFlag{Name: "platform",
                            Usage: "Platform being used for this project. Platform is the type of chip supported (AVR/ ARM)",
                            Value: config.ProjectDefaults.Platform},
                        cli.BoolFlag{Name: "tests",
                            Usage: "Creates a test folder to support unit testing"},
                        cli.StringFlag{Name: "ide",
                            Usage: "Creates the project for a specified IDE (CLion, Eclipse, VS Code)",
                            Value: config.ProjectDefaults.Ide},
                        cli.BoolFlag{Name: "create-demo",
                            Usage: "This will create a demo project that user can build and upload"},
                        cli.BoolFlag{Name: "no-extras",
                            Usage: "This will restrict wio from creating .gitignore, README.md, etc files"},
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed"},
                    },
                    Action: func(c *cli.Context) {
                        command = create.Create{Context: c, Type: create.APP, Update: false}
                    },
                },
            },
        },
        {
            Name:  "update",
            Usage: "Updates the current project and fixes any issues.",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:      "pkg",
                    Usage:     "Updates a wio package, intended to be used by other people",
                    UsageText: "wio update pkg <DIRECTORY> [command options]",
                    Flags: []cli.Flag{
                        cli.BoolFlag{Name: "header-only",
                            Usage: "This flag can be used to specify that the package is header only"},
                        cli.StringFlag{Name: "board",
                            Usage: "Board being used for this project. This will use this board for the update",
                            Value: config.ProjectDefaults.Board},
                        cli.StringFlag{Name: "ide",
                            Usage: "Creates the project for a specified IDE (CLion, Eclipse, VS Code)",
                            Value: config.ProjectDefaults.Ide},
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed"},
                    },
                    Action: func(c *cli.Context) {
                        command = create.Create{Context: c, Type: create.PKG, Update: true}
                    },
                },
                cli.Command{
                    Name:      "app",
                    Usage:     "Updates a wio application, intended to be compiled and uploaded to a device",
                    UsageText: "wio update app <DIRECTORY> [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "board",
                            Usage: "Board being used for this project. This will use this board for the update",
                            Value: config.ProjectDefaults.Board},
                        cli.StringFlag{Name: "framework",
                            Usage: "Framework being used for this project. Framework contains the core libraries",
                            Value: config.ProjectDefaults.Framework},
                        cli.StringFlag{Name: "platform",
                            Usage: "Platform being used for this project. Platform is the type of chip supported (AVR/ ARM)",
                            Value: config.ProjectDefaults.Platform},
                        cli.BoolFlag{Name: "tests",
                            Usage: "Creates a test folder to support unit testing"},
                        cli.StringFlag{Name: "ide",
                            Usage: "Creates the project for a specified IDE (CLion, Eclipse, VS Code)",
                            Value: config.ProjectDefaults.Ide},
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed"},
                    },
                    Action: func(c *cli.Context) {
                        command = create.Create{Context: c, Type: create.APP, Update: true}
                    },
                },
            },
        },
        {
            Name:      "build",
            Usage:     "Builds the wio project.",
            UsageText: "wio build [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "clean",
                    Usage: "Clean the project before building it",
                },
                cli.StringFlag{Name: "target",
                    Usage: "Build a specified target instead of building the default",
                    Value: config.ProjectDefaults.DefaultTarget,
                },
                cli.StringFlag{Name: "dir",
                    Usage: "Directory for the project (default: current working directory)",
                    Value: getCurrDir(),
                },
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                },
            },
            Action: func(c *cli.Context) {
                validateWioProject(c.String("dir"))
                command = build.Build{Context: c}
            },
        },
        {
            Name:      "clean",
            Usage:     "Cleans all the build files for the project.",
            UsageText: "wio clean",
            Flags: []cli.Flag{
                cli.StringFlag{Name: "target",
                    Usage: "Cleans build files for a specified target instead of cleaning all the targets",
                    Value: config.ProjectDefaults.DefaultTarget,
                },
                cli.StringFlag{Name: "dir",
                    Usage: "Directory for the project (default: current working directory)",
                    Value: getCurrDir(),
                },
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                },
            },
            Action: func(c *cli.Context) {
                validateWioProject(c.String("dir"))
                command = clean.Clean{Context: c}
            },
        },
        {
            Name:      "run",
            Usage:     "Builds and Uploads the project to a device (provide port flag to trigger upload)",
            UsageText: "wio run [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "clean",
                    Usage: "Clean the project before building it",
                },
                cli.StringFlag{Name: "target",
                    Usage: "Builds, and uploads a specified target instead of the main/default target",
                    Value: config.ProjectDefaults.DefaultTarget,
                },
                cli.StringFlag{Name: "port",
                    Usage: "Port to upload the project to, (default: automatically select)",
                    Value: config.ProjectDefaults.Port,
                },
                cli.StringFlag{Name: "dir",
                    Usage: "Directory for the project (default: current working directory)",
                    Value: getCurrDir(),
                },
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                },
            },
            Action: func(c *cli.Context) {
                validateWioProject(c.String("dir"))
                command = run.Run{Context: c}
            },
        },
        /*
           {
               Name:      "test",
               Usage:     "Runs unit tests available in the project.",
               UsageText: "wio test",
               Flags: []cli.Flag{
                   cli.BoolFlag{Name: "clean",
                       Usage: "Clean the project before building it",
                   },
                   cli.StringFlag{Name: "port",
                       Usage: "Port to upload the project to, (default: automatically select)",
                       Value: defaults.Port,
                   },
                   cli.StringFlag{Name: "target",
                       Usage: "Builds, and uploads a specified target instead of the main/default target",
                       Value: defaults.Utarget,
                   },
                   cli.BoolFlag{Name: "verbose",
                       Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                   },
               },
               Action: func(c *cli.Context) error {
                   return nil
               },
           },

           {
               Name:      "monitor",
               Usage:     "Runs the serial monitor.",
               UsageText: "wio monitor [command options]",
               Flags: []cli.Flag{
                   cli.BoolFlag{Name: "gui",
                       Usage: "Runs the GUI version of the serial monitor tool",
                   },
                   cli.StringFlag{Name: "port",
                       Usage: "Port to upload the project to, (default: automatically select)",
                       Value: defaults.Port,
                   },
                   cli.BoolFlag{Name: "verbose",
                       Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                   },
               },
               Action: func(c *cli.Context) error {
                   return nil
               },
           },
           {
               Name:      "doctor",
               Usage:     "Guide development tooling and system configurations.",
               UsageText: "wio doctor",
               Action: func(c *cli.Context) error {
                   return nil
               },
           },
           {
               Name:      "analyze",
               Usage:     "Analyzes C/C++ code statically.",
               UsageText: "wio analyze",
               Action: func(c *cli.Context) error {
                   return nil
               },
           },
           {
               Name:      "doxygen",
               Usage:     "Runs doxygen tool to create documentation for the code.",
               UsageText: "wio doxygen",
               Action: func(c *cli.Context) error {
                   return nil
               },
           },
        */
        {
            Name:  "pac",
            Usage: "Package manager for Wio projects.",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:      "add",
                    Usage:     "Add/Update dependencies.",
                    UsageText: "wio pac add [command options]",
                    Flags: []cli.Flag{
                        cli.BoolFlag{Name: "vendor",
                            Usage: "Adds the dependency as vendor",
                        },
                        cli.StringFlag{Name: "dir",
                            Usage: "Directory for the project (default: current working directory)",
                            Value: getCurrDir(),
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                        },
                    },
                    Action: func(c *cli.Context) {
                        validateWioProject(c.String("dir"))
                        command = pac.Pac{Context: c, Type: pac.ADD}
                    },
                },
                cli.Command{
                    Name:      "rm",
                    Usage:     "Remove dependencies.",
                    UsageText: "wio pac rm [command options]",
                    Flags: []cli.Flag{
                        cli.BoolFlag{Name: "A",
                            Usage: "Delete all the dependencies",
                        },
                        cli.StringFlag{Name: "dir",
                            Usage: "Directory for the project (default: current working directory)",
                            Value: getCurrDir(),
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                        },
                    },
                    Action: func(c *cli.Context) {
                        validateWioProject(c.String("dir"))
                        command = pac.Pac{Context: c, Type: pac.RM}
                    },
                },
                cli.Command{
                    Name:      "list",
                    Usage:     "List all the dependencies",
                    UsageText: "wio pac list [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "dir",
                            Usage: "Directory for the project (default: current working directory)",
                            Value: getCurrDir(),
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                        },
                    },
                    Action: func(c *cli.Context) {
                        validateWioProject(c.String("dir"))
                        command = pac.Pac{Context: c, Type: pac.LIST}
                    },
                },
                cli.Command{
                    Name:      "info",
                    Usage:     "Get information about a dependency being used",
                    UsageText: "wio pac info [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "dir",
                            Usage: "Directory for the project (default: current working directory)",
                            Value: getCurrDir(),
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                        },
                    },
                    Action: func(c *cli.Context) {
                        validateWioProject(c.String("dir"))
                        command = pac.Pac{Context: c, Type: pac.INFO}
                    },
                },
                cli.Command{
                    Name:      "publish",
                    Usage:     "Publish the wio package to the package manager site (npm site)",
                    UsageText: "wio pac publish [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "dir",
                            Usage: "Directory for the project (default: current working directory)",
                            Value: getCurrDir(),
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                        },
                    },
                    Action: func(c *cli.Context) {
                        validateWioProject(c.String("dir"))
                        command = pac.Pac{Context: c, Type: pac.PUBLISH}
                    },
                },
                cli.Command{
                    Name:      "get",
                    Usage:     "Gets all the packages mentioned in wio.yml file and vendor folder.",
                    UsageText: "wio pac get [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "dir",
                            Usage: "Directory for the project (default: current working directory)",
                            Value: getCurrDir(),
                        },
                        cli.BoolFlag{Name: "clean",
                            Usage: "Cleans all the current packages and re get all of them.",
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                        },
                    },
                    Action: func(c *cli.Context) {
                        validateWioProject(c.String("dir"))
                        command = pac.Pac{Context: c, Type: pac.GET}
                    },
                },
                /*
                   cli.Command{
                       Name:      "update",
                       Usage:     "Updates all the packages mentioned in wio.yml file and vendor folder.",
                       UsageText: "wio pac update [command options]",
                       Flags: []cli.Flag{
                           cli.StringFlag{Name: "dir",
                               Usage: "Directory for the project (default: current working directory)",
                               Value: getCurrDir(),
                           },
                           cli.BoolFlag{Name: "verbose",
                               Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                           },
                       },
                       Action: func(c *cli.Context) {
                           command = pac.Pac{Context: c, Type: pac.UPDATE}
                       },
                   },
                   cli.Command{
                       Name:      "collect",
                       Usage:     "Creates vendor folder and puts all the packages in that folder.",
                       UsageText: "wio pac collect [command options]",
                       Flags: []cli.Flag{
                           cli.StringFlag{Name: "dir",
                               Usage: "Directory for the project (default: current working directory)",
                               Value: getCurrDir(),
                           },
                           cli.StringFlag{Name: "pkg",
                               Usage: "Packages to collect instead of collecting all of the packages.",
                               Value: "none",
                           },
                           cli.BoolFlag{Name: "verbose",
                               Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                           },
                       },
                       Action: func(c *cli.Context) {
                           command = pac.Pac{Context: c, Type: pac.COLLECT}
                       },
                   },*/
            },
        },
    }

    app.Action = func(c *cli.Context) error {
        app.Command("help").Run(c)
        return nil
    }

    if err = app.Run(os.Args); err != nil {
        panic(err)
    }

    // execute the command
    if command != nil {
        // check if verbose flag is true
        if command.GetContext().Bool("verbose") {
            log.SetVerbose()
        }

        command.Execute()
    }
}

func validateWioProject(directory string) {
    directory, err := filepath.Abs(directory)
    commands.RecordError(err, "")

    if !utils.PathExists(directory) {
        log.Norm.Yellow(true, directory+" : no such path exists")
        os.Exit(3)
    }

    if !utils.PathExists(directory + io.Sep + "wio.yml") {
        log.Norm.Yellow(true, "Not a valid wio project: wio.yml file missing")
        os.Exit(3)
    }
}

// returns the current directory from where wio is being called
func getCurrDir() string {
    directory, err := os.Getwd()
    commands.RecordError(err, "")
    return directory
}
