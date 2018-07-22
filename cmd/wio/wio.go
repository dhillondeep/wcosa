// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package main contains the main code for Wio.
// Wio is a tool to make development of embedded system applications easier and simpler.
// It allows for building, testing, and uploading AVR applications for Commandline.
package main

import (
    "os"
    "time"
    "wio/cmd/wio/commands"
    "wio/cmd/wio/commands/create"
    "wio/cmd/wio/commands/devices"
    "wio/cmd/wio/commands/pac/install"
    "wio/cmd/wio/commands/pac/vendor"
    "wio/cmd/wio/commands/run"
    "wio/cmd/wio/config"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/log"
    "wio/cmd/wio/utils/io"

    "github.com/urfave/cli"
)

var createFlags = []cli.Flag{
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
        Name:  "verbose",
        Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
    },
    cli.BoolFlag{
        Name:  "disable-warnings",
        Usage: "Disables all the warning shown by wio",
    },
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
}

var buildFlags = []cli.Flag{
    cli.BoolFlag{
        Name:  "all",
        Usage: "Build all available targets",
    },
    cli.StringFlag{
        Name:  "port",
        Usage: "Specify upload port",
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
        Name:  "hard",
        Usage: "Removes build directories",
    },
}

var runFlags = []cli.Flag{
    cli.StringFlag{
        Name:  "port",
        Usage: "Specify upload port",
    },
    cli.StringFlag{
        Name:  "args",
        Usage: "Arguments passed to executable",
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

var command commands.Command
var cmd = []cli.Command{
    {
        Name:  "create",
        Usage: "Creates and initializes a wio project.",
        Subcommands: cli.Commands{
            {
                Name:      "pkg",
                Usage:     "Creates a wio package.",
                UsageText: "wio create pkg [command options]",
                Flags:     createFlags,
                Action: func(c *cli.Context) {
                    command = create.Create{Context: c, Update: false, Type: constants.PKG}
                },
            },
            {
                Name:      "app",
                Usage:     "Creates a wio app.",
                UsageText: "wio create app [command options]",
                Flags:     createFlags,
                Action: func(c *cli.Context) {
                    command = create.Create{Context: c, Update: false, Type: constants.APP}
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
        Flags:     runFlags,
        Action: func(c *cli.Context) {
            command = run.Run{Context: c, RunType: run.TypeRun}
        },
    },

    {

        Name:  "vendor",
        Usage: "Manage locally vendored dependencies.",
        Subcommands: cli.Commands{
            {
                Name:      "add",
                Usage:     "Add a vendored package as a dependency.",
                UsageText: "wio vendor add [package]",
                Action: func(c *cli.Context) {
                    command = vendor.Cmd{Context: c, Op: vendor.Add}
                },
            },
            {
                Name:      "rm",
                Usage:     "Remove a vendor dependency.",
                UsageText: "wio vendor rm [package]",
                Action: func(c *cli.Context) {
                    command = vendor.Cmd{Context: c, Op: vendor.Remove}
                },
            },
        },
    },
    {
        Name:      "install",
        Usage:     "Install packages from remote server.",
        UsageText: "wio install [name] [version]",
        Action: func(c *cli.Context) {
            command = install.Cmd{Context: c}
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
}

func wio() error {
    // read help templates
    appHelpText, err := io.AssetIO.ReadFile("cli-helper/app-help.txt")
    if err != nil {
        return err
    }

    commandHelpText, err := io.AssetIO.ReadFile("cli-helper/command-help.txt")
    if err != nil {
        return err
    }

    subCommandHelpText, err := io.AssetIO.ReadFile("cli-helper/subcommand-help.txt")
    if err != nil {
        return err
    }

    // override help templates
    cli.AppHelpTemplate = string(appHelpText)
    cli.CommandHelpTemplate = string(commandHelpText)
    cli.SubcommandHelpTemplate = string(subCommandHelpText)

    app := cli.NewApp()
    app.Name = config.ProjectMeta.Name
    app.Version = config.ProjectMeta.Version
    app.EnableBashCompletion = config.ProjectMeta.EnableBashCompletion
    app.Compiled = time.Now()
    app.Copyright = config.ProjectMeta.Copyright
    app.Usage = config.ProjectMeta.UsageText
    app.Commands = cmd

    app.Action = func(c *cli.Context) error {
        app.Command("help").Run(c)
        return nil
    }
    if err = app.Run(os.Args); err != nil {
        return err
    }
    // execute the command
    if command != nil {
        // check if verbose flag is true
        if command.GetContext().Bool("verbose") {
            log.SetVerbose()
        }
        if command.GetContext().Bool("disable-warnings") {
            log.DisableWarnings()
        }
        return command.Execute()
    }
    return nil
}

func main() {
    err := wio()
    if err != nil {
        log.Errln(err.Error())
        os.Exit(1)
    }
}
