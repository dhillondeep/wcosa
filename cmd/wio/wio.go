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
    "wio/cmd/wio/errors"
    "wio/cmd/wio/constants"
)

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
                cli.Command{
                    Name:      "pkg",
                    Usage:     "Creates a wio package, intended to be used by other people.",
                    UsageText: "wio create pkg [command options]",
                    Flags: []cli.Flag{
                        cli.BoolFlag{
                            Name:  "header-only",
                            Usage: "Specify a header-only package.",
                        },
                        cli.BoolFlag{
                            Name:  "create-example",
                            Usage: "This will create an example project that user can build and upload.",
                        },
                        cli.BoolFlag{
                            Name:  "only-config",
                            Usage: "Creates only the configuration file (wio.yml).",
                        },
                        cli.BoolFlag{
                            Name:  "no-extras",
                            Usage: "This will restrict wio from creating .gitignore, README.md, etc files.",
                        },
                        cli.BoolFlag{
                            Name:  "disable-warnings",
                            Usage: "Disables all the warning shown by wio.",
                        },
                        cli.BoolFlag{
                            Name:  "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                        },
                    },
                    Subcommands: cli.Commands{
                        cli.Command{
                            Name:      "avr",
                            Usage:     "Creates an AVR package.",
                            UsageText: "wio create pkg avr [directory] [board] [command options]",
                            Flags: []cli.Flag{
                                cli.StringFlag{
                                    Name:  "framework",
                                    Usage: "AVR Framework to use: `Cosa` or `Arduino`",
                                    Value: config.AvrProjectDefaults.Framework,
                                },
                            },
                            Action: func(c *cli.Context) {
                                command = create.Create{Context: c, Type: constants.PKG, Platform: constants.AVR, Update: false}
                            },
                        },
                        cli.Command{
                            Name:      "native",
                            Usage:     "Create native (Desktop) package.",
                            UsageText: "wio create pkg native [directory] [command options]",
                            Flags: []cli.Flag{
                                cli.StringFlag{
                                    Name:  "framework",
                                    Usage: "Currently unused.",
                                    Value: config.NativeProjectDefaults.Framework,
                                },
                            },
                        },
                    },
                },
                cli.Command{
                    Name:      "app",
                    Usage:     "Creates a wio application, intended to be compiled and uploaded to a device.",
                    UsageText: "wio create app [command options]",
                    Subcommands: cli.Commands{
                        cli.Command{
                            Name:      "avr",
                            Usage:     "Creates an AVR application.",
                            UsageText: "wio create app avr [directory] [board] [command options]",
                            Flags: []cli.Flag{
                                cli.StringFlag{Name: "framework",
                                    Usage: "Framework being used for this project. Framework contains the core libraries.",
                                    Value: config.AvrProjectDefaults.Framework},
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
                                command = create.Create{Context: c, Type: constants.APP, Platform: constants.AVR, Update: false}
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
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed."},
                cli.BoolFlag{Name: "disable-warnings",
                    Usage: "Disables all the warning shown by wio"},
                cli.BoolFlag{Name: "no-extras",
                    Usage: "This will restrict wio from creating .gitignore, README.md, etc files."},
                cli.BoolFlag{Name: "config-help",
                    Usage: "Prints help text in the config file."},
            },
            Action: func(c *cli.Context) {
                command = create.Create{Context: c, Update: true}
            },
        },

        {
            Name:      "run",
            Usage:     "Builds, Runs and/or Uploads the project to a device.",
            UsageText: "wio run [directory] [command options]",
            Flags: []cli.Flag{
                cli.StringFlag{Name: "target",
                    Usage: "Builds, Runs and/or uploads a specified target instead of the main/default target.",
                    Value: config.AvrProjectDefaults.DefaultTarget,
                },
                cli.BoolFlag{Name: "clean",
                    Usage: "Clean the project build files before new build is triggered.",
                },
                cli.BoolFlag{Name: "upload",
                    Usage: "Uploads the built target to a device (automatically selected).",
                },
                cli.StringFlag{Name: "port",
                    Usage: "Port to upload the project to, (default: automatically select).",
                    Value: config.AvrProjectDefaults.Port,
                },
                cli.BoolFlag{Name: "build-all",
                    Usage: "Build all the targets specified in wio.yml file.",
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
                            Value: config.AvrProjectDefaults.Baud},
                        cli.StringFlag{Name: "port",
                            Usage: "Serial Port to open.",
                            Value: config.AvrProjectDefaults.Port},
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
        /*
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
        /*},
          },
        */
    }

    app.Action = func(c *cli.Context) error {
        app.Command("help").Run(c)
        return nil
    }

    if err = app.Run(os.Args); err != nil {
        log.WriteErrorlnExit(err)
    }

    defer func() {
        if r := recover(); r != nil {
            fatalError := errors.FatalError{
                Log: r,
            }

            log.WriteErrorlnExit(fatalError)
        }
    }()

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
