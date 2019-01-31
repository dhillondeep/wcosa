// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package main contains the main code for Wio.
// Wio is a tool to make development of embedded system applications easier and simpler.
// It allows for building, testing, and uploading C/C++ applications for Commandline.
package main

import (
    "os"
    "time"
    "wio/internal/cmd/env"
    "wio/internal/executor"

    "wio/internal/cmd"
    "wio/internal/cmd/create"
    "wio/internal/cmd/devices"
    "wio/internal/cmd/pac/install"
    "wio/internal/cmd/pac/publish"
    "wio/internal/cmd/pac/user"
    "wio/internal/cmd/pac/vendor"
    "wio/internal/cmd/run"
    "wio/internal/cmd/upgrade"
    "wio/internal/config/defaults"
    "wio/internal/config/meta"
    "wio/internal/constants"
    "wio/pkg/log"
    "wio/pkg/util/sys"

    "github.com/urfave/cli"
)

var appWideFlags = []cli.Flag{
    cli.BoolFlag{
        Name:  "verbose",
        Usage: "Turns verbose mode on to show detailed logs.",
    },
    cli.BoolFlag{
        Name:  "disable-warnings",
        Usage: "Disables all the warning shown by wio.",
    },
}

var createFlags = []cli.Flag{
    cli.StringFlag{
        Name:  "platform",
        Usage: "Target platform: 'avr', 'native'.",
        Value: "all",
    },
    cli.StringFlag{
        Name:  "framework",
        Usage: "Target framework: 'arduino', 'cosa'.",
        Value: "all",
    },
    cli.StringFlag{
        Name:  "board",
        Usage: "Target boards: e.g. 'uno', 'mega2560', etc.",
        Value: "all",
    },
    cli.StringFlag{
        Name:  "ide",
        Usage: "[clion]",
        Value: "none",
    },
    cli.BoolFlag{
        Name:  "only-config",
        Usage: "Creates only the configuration file (wio.yml).",
    },
    cli.BoolFlag{
        Name:  "header-only",
        Usage: "Specify a header-only package.",
    },
}

var updateFlags = []cli.Flag{
    cli.StringFlag{
        Name:  "ide",
        Usage: "[clion]",
        Value: "none",
    },
    cli.BoolFlag{
        Name:  "full",
        Usage: "Full update and overrides files.",
    },
}

var buildFlags = []cli.Flag{
    cli.BoolFlag{
        Name:  "force",
        Usage: "Forces a full build for targets.",
    },
    cli.BoolFlag{
        Name:  "retool",
        Usage: "Removes existing toolchain and hard resets it.",
    },
    cli.BoolFlag{
        Name:  "all",
        Usage: "Build all available targets.",
    },
}

var cleanFlags = []cli.Flag{
    cli.BoolFlag{
        Name:  "all",
        Usage: "Clean all available targets.",
    },
    cli.BoolFlag{
        Name:  "hard",
        Usage: "Removes build directories.",
    },
}

var runFlags = []cli.Flag{
    cli.StringFlag{
        Name:  "port",
        Usage: "Specify upload port.",
        Value: "none",
    },
    cli.StringFlag{
        Name:  "args",
        Usage: "Arguments passed to executable.",
    },
}

var monitorFlags = []cli.Flag{
    cli.IntFlag{Name: "baud",
        Usage: "Baud rate for the Serial port.",
        Value: defaults.Baud},
    cli.StringFlag{Name: "port",
        Usage: "Serial Port to open.",
        Value: defaults.Port},
    cli.BoolFlag{Name: "gui",
        Usage: "Runs the GUI version of the serial monitor tool."},
}

var devicesListFlags = []cli.Flag{
    cli.BoolFlag{Name: "basic",
        Usage: "Shows only the name of the ports."},
    cli.BoolFlag{Name: "show-all",
        Usage: "Shows all the ports, closed or open (Default: only open devices)."},
}

var envFlags = []cli.Flag{
    cli.BoolFlag{Name: "local",
        Usage: "Creates and updates local environment."},
}

var upgradeFlags = []cli.Flag{
    cli.BoolFlag{Name: "force",
        Usage: "Overrides all the restrictions and forces an update."},
}

var command cmd.Command
var commands = []cli.Command{
    {
        Name:      "create",
        Usage:     "Creates and initializes a wio project.",
        UsageText: "wio create <subcommand> [arguments]",
        Subcommands: cli.Commands{
            {
                Name:      "pkg",
                Usage:     "Creates a wio package.",
                UsageText: "wio create pkg [directory] [command options]",
                Flags:     append(createFlags, appWideFlags...),
                Action: func(c *cli.Context) {
                    command = create.Create{Context: c, Update: false, Type: constants.Pkg}
                },
            },
            {
                Name:      "app",
                Usage:     "Creates a wio app.",
                UsageText: "wio create app [directory] [command options]",
                Flags:     append(createFlags, appWideFlags...),
                Action: func(c *cli.Context) {
                    command = create.Create{Context: c, Update: false, Type: constants.App}
                },
            },
        },
    },
    {
        Name:      "update",
        Usage:     "Updates the current project and fixes any issues.",
        UsageText: "wio update [directory] [command options]",
        Flags:     append(updateFlags, appWideFlags...),
        Action: func(c *cli.Context) {
            command = create.Create{Context: c, Update: true}
        },
    },
    {
        Name:      "build",
        Usage:     "Configure and build the project.",
        UsageText: "wio build [targets...] [command options]",
        Flags:     append(buildFlags, appWideFlags...),
        Action: func(c *cli.Context) {
            command = run.Run{Context: c, RunType: run.TypeBuild}
        },
    },
    {
        Name:      "clean",
        Usage:     "Clean project targets.",
        UsageText: "wio clean [targets...] [command options]",
        Flags:     append(cleanFlags, appWideFlags...),
        Action: func(c *cli.Context) {
            command = run.Run{Context: c, RunType: run.TypeClean}
        },
    },
    {
        Name:      "run",
        Usage:     "Builds, Runs and/or Uploads the project to a device.",
        UsageText: "wio run [targets...] [command options]",
        Flags:     append(runFlags, appWideFlags...),
        Action: func(c *cli.Context) {
            command = run.Run{Context: c, RunType: run.TypeRun}
        },
    },
    {
        Name:      "vendor",
        Usage:     "Manage locally vendored dependencies.",
        UsageText: "wio vendor <subcommand> [command options]",
        Subcommands: cli.Commands{
            {
                Name:      "add",
                Usage:     "Add a vendored package as a dependency.",
                UsageText: "wio vendor add [package]",
                Flags: appWideFlags,
                Action: func(c *cli.Context) {
                    command = vendor.Cmd{Context: c, Op: vendor.Add}
                },
            },
            {
                Name:      "rm",
                Usage:     "Remove a vendor dependency.",
                UsageText: "wio vendor rm [package]",
                Flags: appWideFlags,
                Action: func(c *cli.Context) {
                    command = vendor.Cmd{Context: c, Op: vendor.Remove}
                },
            },
        },
    },
    {
        Name:      "install",
        Usage:     "Install packages from remote server.",
        UsageText: "wio install [packages...]",
        Flags:     appWideFlags,
        Action: func(c *cli.Context) {
            command = install.Cmd{Context: c}
        },
    },
    {
        Name:      "login",
        Usage:     "Login to the registry.",
        UsageText: "wio login",
        Flags: appWideFlags,
        Action: func(c *cli.Context) {
            command = user.Login{Context: c}
        },
    },
    {
        Name:      "logout",
        Usage:     "Logout from registry account.",
        UsageText: "wio logout",
        Flags: appWideFlags,
        Action: func(c *cli.Context) {
            command = user.Logout{Context: c}
        },
    },
    {
        Name:      "publish",
        Usage:     "Publish package to registry.",
        UsageText: "wio publish",
        Flags:     appWideFlags,
        Action: func(c *cli.Context) {
            command = publish.Cmd{Context: c}
        },
    },
    {
        Name:      "devices",
        Usage:     "Handles serial devices connected.",
        UsageText: "wio devices <subcommand> [command options]",
        Subcommands: cli.Commands{
            cli.Command{
                Name:      "monitor",
                Usage:     "Opens a Serial monitor.",
                UsageText: "wio devices monitor [command options]",
                Flags:     append(monitorFlags, appWideFlags...),
                Action: func(c *cli.Context) {
                    command = devices.Devices{Context: c, Type: devices.MONITOR}
                },
            },
            cli.Command{
                Name:      "list",
                Usage:     "Lists all the connected devices/ports and provides information about them.",
                UsageText: "wio devices list [command options]",
                Flags:     append(devicesListFlags, appWideFlags...),
                Action: func(c *cli.Context) {
                    command = devices.Devices{Context: c, Type: devices.LIST}
                },
            },
        },
    },
    {
        Name:      "monitor",
        Usage:     "Opens a Serial monitor.",
        UsageText: "wio monitor [command options]",
        Flags:     append(monitorFlags, appWideFlags...),
        Action: func(c *cli.Context) {
            command = devices.Devices{Context: c, Type: devices.MONITOR}
        },
    },
    {
        Name:      "env",
        Usage:     "Wio global environment variables.",
        UsageText: "wio env [command options]",
        Action: func(c *cli.Context) {
            command = env.Env{Context: c, Command: env.VIEW}
        },
        Subcommands: cli.Commands{
            cli.Command{
                Name:      "reset",
                Usage:     "Resets environment variables to default",
                UsageText: "wio env reset [command options]",
                Flags:     append(envFlags, appWideFlags...),
                Action: func(c *cli.Context) {
                    command = env.Env{Context: c, Command: env.RESET}
                },
            },
            cli.Command{
                Name:      "set",
                Usage:     "Modifies the environment variable or adds a new one (name=value or name).",
                UsageText: "wio env set [vars...] [command options]",
                Flags:     append(envFlags, appWideFlags...),
                Action: func(c *cli.Context) {
                    command = env.Env{Context: c, Command: env.SET}
                },
            },
            cli.Command{
                Name:      "unset",
                Usage:     "Removes the environment variable.",
                UsageText: "wio env unset [vars...] [command options]",
                Flags:     append(envFlags, appWideFlags...),
                Action: func(c *cli.Context) {
                    command = env.Env{Context: c, Command: env.UNSET}
                },
            },
        },
    },
    {
        Name:      "upgrade",
        Usage:     "Upgrades wio to a specific version or latest version.",
        UsageText: "wio upgrade [version]",
        Flags: append(upgradeFlags, appWideFlags...),
        Action: func(c *cli.Context) {
            command = upgrade.Upgrade{Context: c}
        },
    },
}

func wio() error {
    appHelpText, err := sys.AssetIO.ReadFile("cli-helper/app-help.txt")
    if err != nil {
        return err
    }

    commandHelpText, err := sys.AssetIO.ReadFile("cli-helper/command-help.txt")
    if err != nil {
        return err
    }

    subCommandHelpText, err := sys.AssetIO.ReadFile("cli-helper/subcommand-help.txt")
    if err != nil {
        return err
    }

    // override help templates
    cli.AppHelpTemplate = string(appHelpText)
    cli.CommandHelpTemplate = string(commandHelpText)
    cli.SubcommandHelpTemplate = string(subCommandHelpText)

    app := cli.NewApp()
    app.Name = meta.Name
    app.Version = meta.Version
    app.Compiled = time.Now()
    app.Copyright = meta.Copyright
    app.Usage = meta.UsageText
    app.Commands = commands

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
    errorHandle := func(err error) {
        if err != nil {
            log.Errln(err.Error())
            os.Exit(1)
        }
    }

    // startup
    errorHandle(executor.ExecuteStartup())

    // wio stuff
    errorHandle(wio())
}
