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
    "time"
    "wio/cmd/wio/commands"
    "wio/cmd/wio/commands/create"
    "wio/cmd/wio/commands/devices"
    "wio/cmd/wio/commands/pac"
    "wio/cmd/wio/commands/run"
    "wio/cmd/wio/config"
    "wio/cmd/wio/log"
    "wio/cmd/wio/utils/io"
)

var packageFlags = []cli.Flag{
    cli.StringFlag{
        Name:  "platform",
        Usage: "Target platform: 'AVR', 'Native', or 'all'",
        Value: "all",
    },
    cli.StringFlag{
        Name:  "framework",
        Usage: "Target framework: 'Arduino', 'Cosa', or 'all'",
        Value: "all",
    },
    cli.StringFlag{
        Name:  "board",
        Usage: "Target boards: e.g. 'uno', 'mega2560', or 'all'",
        Value: "all",
    },
    cli.BoolFlag{
        Name:  "only-config",
        Usage: "Creates only the configuration file (wio.yml).",
    },
    cli.BoolFlag{
        Name:  "header-only",
        Usage: "Specify a header-only package.",
    },
    cli.BoolFlag{
        Name:  "disable-warnings",
        Usage: "Disables all the warning shown by wio.",
    },
    cli.BoolFlag{
        Name:  "verbose",
        Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
    },
    /*cli.BoolFlag{
        Name:  "no-extras",
        Usage: "This will restrict wio from creating .gitignore, README.md, etc files.",
    },
    cli.BoolFlag{
        Name:  "create-example",
        Usage: "This will create an example project that user can build and upload.",
    },*/
}

var updateFlags = []cli.Flag{
    cli.BoolFlag{
        Name:  "verbose",
        Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
    },
    cli.BoolFlag{
        Name:  "disable-warnings",
        Usage: "Disables all the warning shown by wio",
    },
    /*cli.BoolFlag{
        Name:  "no-extras",
        Usage: "This will restrict wio from creating .gitignore, README.md, etc files.",
    },
    /*cli.BoolFlag{
        Name:  "config-help",
        Usage: "Prints help text in the config file.",
    },*/
}

var buildFlags = []cli.Flag{
    cli.BoolFlag{
        Name:  "all",
        Usage: "Build all available targets",
    },
    cli.BoolFlag{
        Name:  "verbose",
        Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
    },
    cli.BoolFlag{
        Name:  "disable-warnings",
        Usage: "Disables all the warning shown by wio",
    },
}

var cleanFlags = []cli.Flag{
    cli.BoolFlag{
        Name: "hard",
        Usage: "Removes build directories",
    },
}

func main() {
    log.Init()

    // read help templates
    appHelpText, err := io.AssetIO.ReadFile("cli-helper/app-help.txt")
    log.WriteErrorlnExit(err)

    commandHelpText, err := io.AssetIO.ReadFile("cli-helper/command-help.txt")
    log.WriteErrorlnExit(err)

    subCommandHelpText, err := io.AssetIO.ReadFile("cli-helper/subcommand-help.txt")
    log.WriteErrorlnExit(err)

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

                {
                    Name:      "pkg",
                    Usage:     "Creates a wio package.",
                    UsageText: "wio create pkg [command options]",
                    Flags:     packageFlags,
                    Action: func(c *cli.Context) {
                        command = create.Create{Context: c, Update: false}
                    },
                },

                {
                    Name:      "app",
                    Usage:     "Creates a wio application.",
                    UsageText: "wio create app [command options]",
                    Subcommands: cli.Commands{
                        cli.Command{
                            Name:      "avr",
                            Usage:     "Creates an AVR application.",
                            UsageText: "wio create app avr [directory] [board] [command options]",
                            Flags: []cli.Flag{
                                cli.StringFlag{Name: "framework",
                                    Usage: "Framework being used for this project. Framework contains the core libraries.",
                                    Value: config.ProjectDefaults.Framework},
                                cli.BoolFlag{Name: "create-example",
                                    Usage: "This will create an example project that user can build and upload."},
                                cli.BoolFlag{Name: "only-config",
                                    Usage: "Creates only the configuration file (wio.yml)."},
                                cli.BoolFlag{Name: "no-extras",
                                    Usage: "This will restrict wio from creating .gitignore, README.md, etc files."},
                                cli.BoolFlag{Name: "verbose",
                                    Usage: "Turns verbose mode on to show detailed errors and commands being executed."},
                                cli.BoolFlag{Name: "disable-warnings",
                                    Usage: "Disables all the warning shown by wio."},
                            },
                            Action: func(c *cli.Context) {
                                command = create.Create{Context: c, Update: false}
                            },
                        },
                    },
                },
            },
        },

        {
            Name:      "update",
            Usage:     "Updates the current project and fixes any issues.",
            UsageText: "wio update [directory] [command options]",
            Flags:     updateFlags,
            Action: func(c *cli.Context) {
                command = create.Create{Context: c, Update: true}
            },
        },

        {
            Name:      "build",
            Usage:     "Configure and build the project.",
            UsageText: "wio build [targets] [command options]",
            Flags:     buildFlags,
            Action: func(c *cli.Context) {
                command = run.Run{Context: c, RunType: run.TypeBuild}
            },
        },

        {
            Name:  "clean",
            Usage: "Clean project targets",
            Flags: append(buildFlags, cleanFlags...),
            Action: func(c *cli.Context) {
                command = run.Run{Context: c, RunType: run.TypeClean}
            },
        },

        {
            Name:      "run",
            Usage:     "Builds, Runs and/or Uploads the project to a device.",
            UsageText: "wio run [directory] [command options]",
            Flags: []cli.Flag{
                cli.StringFlag{Name: "target",
                    Usage: "Builds, Runs and/or uploads a specified target instead of the main/default target.",
                    Value: config.ProjectDefaults.DefaultTarget,
                },
                cli.BoolFlag{Name: "upload",
                    Usage: "Uploads the built target to a device (automatically selected).",
                },
                cli.StringFlag{Name: "port",
                    Usage: "Port to upload the project to, (default: automatically select).",
                    Value: config.ProjectDefaults.Port,
                },
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                },
                cli.BoolFlag{Name: "disable-warnings",
                    Usage: "Disables all the warning shown by wio.",
                },
            },
            Action: func(c *cli.Context) {
                command = run.Run{Context: c}
            },
        },
        {
            Name:      "devices",
            Usage:     "Handles serial devices connected.",
            UsageText: "wio devices [command options]",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:      "monitor",
                    Usage:     "Opens a Serial monitor.",
                    UsageText: "wio monitor open [command options]",
                    Flags: []cli.Flag{
                        cli.IntFlag{Name: "baud",
                            Usage: "Baud rate for the Serial port.",
                            Value: config.ProjectDefaults.Baud},
                        cli.StringFlag{Name: "port",
                            Usage: "Serial Port to open.",
                            Value: config.ProjectDefaults.Port},
                        cli.BoolFlag{Name: "gui",
                            Usage: "Runs the GUI version of the serial monitor tool"},
                        cli.BoolFlag{Name: "disable-warnings",
                            Usage: "Disables all the warning shown by wio.",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = devices.Devices{Context: c, Type: devices.MONITOR}
                    },
                },
                cli.Command{
                    Name:      "list",
                    Usage:     "Lists all the connected devices/ports and provides information about them.",
                    UsageText: "wio devices list [command options]",
                    Flags: []cli.Flag{
                        cli.BoolFlag{Name: "basic",
                            Usage: "Shows only the name of the ports."},
                        cli.BoolFlag{Name: "show-all",
                            Usage: "Shows all the ports, closed or open (Default: only open devices)."},
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed."},
                        cli.BoolFlag{Name: "disable-warnings",
                            Usage: "Disables all the warning shown by wio.",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = devices.Devices{Context: c, Type: devices.LIST}
                    },
                },
            },
        },
        {
            Name:      "install",
            Usage:     "Install's wio packages from remote server.",
            UsageText: "wio install [package name] [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "save",
                    Usage: "Adds package to wio.yml file and installs it."},
                cli.BoolFlag{Name: "clean",
                    Usage: "Deletes previous packages and installs new ones."},
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed."},
                cli.BoolFlag{Name: "disable-warnings",
                    Usage: "Disables all the warning shown by wio."},
                cli.BoolFlag{Name: "config-help",
                    Usage: "Prints help text in the config file."},
            },
            Action: func(c *cli.Context) {
                command = pac.Pac{Context: c, Type: pac.INSTALL}
            },
        },
        {
            Name:      "uninstall",
            Usage:     "Uninstall's wio packages downloaded.",
            UsageText: "wio uninstall <package name> [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "save",
                    Usage: "Removes package from wio.yml file."},
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed."},
                cli.BoolFlag{Name: "disable-warnings",
                    Usage: "Disables all the warning shown by wio."},
                cli.BoolFlag{Name: "config-help",
                    Usage: "Prints help text in the config file."},
            },
            Action: func(c *cli.Context) {
                command = pac.Pac{Context: c, Type: pac.UNINSTALL}
            },
        },
        {
            Name:      "publish",
            Usage:     "Publishes wio package.",
            UsageText: "wio publish [directory] [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed."},
                cli.BoolFlag{Name: "disable-warnings",
                    Usage: "Disables all the warning shown by wio.",
                },
            },
            Action: func(c *cli.Context) {
                command = pac.Pac{Context: c, Type: pac.PUBLISH}
            },
        },
        {
            Name:      "collect",
            Usage:     "Grabs all the remote packages and stores them in vendor directory.",
            UsageText: "wio collect [package] [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "save",
                    Usage: "Updates packages moved to vendor status to true."},
                cli.BoolFlag{Name: "disable-warnings",
                    Usage: "Disables all the warning shown by wio."},
                cli.BoolFlag{Name: "config-help",
                    Usage: "Prints help text in the config file."},
            },
            Action: func(c *cli.Context) {
                command = pac.Pac{Context: c, Type: pac.COLLECT}
            },
        },
        {
            Name:      "list",
            Usage:     "List all the packages installed.",
            UsageText: "wio list [directory] [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed."},
                cli.BoolFlag{Name: "disable-warnings",
                    Usage: "Disables all the warning shown by wio.",
                },
            },
            Action: func(c *cli.Context) {
                command = pac.Pac{Context: c, Type: pac.LIST}
            },
        },
    }

    app.Action = func(c *cli.Context) error {
        app.Command("help").Run(c)
        return nil
    }
    if err = app.Run(os.Args); err != nil {
        log.WriteErrorlnExit(err)
    }
    //defer func() {
    //    if r := recover(); r != nil {
    //        fatalError := errors.FatalError{
    //            Log: r,
    //        }

    //        log.WriteErrorlnExit(fatalError)
    //    }
    //}()
    // execute the command
    if command != nil {
        // check if verbose flag is true
        if command.GetContext().Bool("verbose") {
            log.SetVerbose()
        }
        if command.GetContext().Bool("disable-warnings") {
            log.DisableWarnings()
        }
        command.Execute()
    }
}
