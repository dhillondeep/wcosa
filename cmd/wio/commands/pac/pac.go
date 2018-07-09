// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands package, which contains all the commands provided by the tool.
// Package manager for wio. It used npm as a backend and pushes packages to that
package pac

import (
    goerr "errors"
    "os"
    "strings"
    "wio/cmd/wio/commands/run"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/log"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"

    "github.com/fatih/color"
    "github.com/urfave/cli"
)

const (
    LIST      = "list"
    PUBLISH   = "PUBLISH"
    UNINSTALL = "uninstall"
    INSTALL   = "install"
    COLLECT   = "collect"
)

type Pac struct {
    Context *cli.Context
    Type    string
}

// Get context for the command
func (pac Pac) GetContext() *cli.Context {
    return pac.Context
}

// Executes the libraries command
func (pac Pac) Execute() error {
    // check if valid wio project
    directory, err := os.Getwd()
    if err != nil {
        return err
    }
    if !utils.PathExists(directory + io.Sep + io.Config) {
        return errors.ConfigMissing{}
    }
    log.Info("updating package.json ... ")
    if err := updateNpmConfig(directory, pac.Type == PUBLISH); err != nil {
        log.WriteFailure()
        return err
    }
    log.WriteSuccess()

    switch pac.Type {
    case INSTALL:
        return pac.handleInstall(directory)
    case UNINSTALL:
        return pac.handleUninstall(directory)
    case LIST:
        return handleList(directory)
    case PUBLISH:
        return handlePublish(directory)
    default:
        return goerr.New("invalid pac command")
    }
}

// This handles the install command and uses npm to install packages
func (pac Pac) handleInstall(directory string) error {
    // check install arguments
    installPackage := installArgumentCheck(pac.Context.Args())

    remoteDirectory := io.Path(directory, io.Folder, io.Modules)
    wioPath := io.Path(directory, io.Config)

    // clean npm_modules in .wio folder
    if pac.Context.Bool("clean") {
        log.Info(log.Cyan, "cleaning npm packages ... ")

        if !utils.PathExists(remoteDirectory) || !utils.PathExists(wioPath) {
            log.WriteSuccess()
        } else {
            if err := os.RemoveAll(remoteDirectory); err != nil {
                log.WriteFailure()
                return err
            }
            packageLock := io.Path(directory, io.Folder, "package-lock.json")
            if err := os.RemoveAll(packageLock); err != nil {
                log.WriteFailure()
                return err
            }
            log.WriteSuccess()
        }
    }

    var npmCmdArgs []string

    if installPackage[0] != "all" {
        log.Infoln(log.Blue, "Installing %s", installPackage)
        npmCmdArgs = append(npmCmdArgs, installPackage...)

        if pac.Context.IsSet("save") {
            config, err := utils.ReadWioConfig(directory)
            if err != nil {
                return err
            }

            dependencies := config.GetDependencies()

            if config.GetDependencies() == nil {
                dependencies = types.DependenciesTag{}
            }

            for _, packageGiven := range installPackage {
                strip := strings.Split(packageGiven, "@")

                packageName := strip[0]
                dependencies[packageName] = &types.DependencyTag{
                    Version: "latest",
                    Vendor:  false,
                }

                if len(strip) > 1 {
                    dependencies[packageName].Version = strip[1]
                }
            }

            config.SetDependencies(dependencies)

            log.Write(log.INFO, color.New(color.FgCyan), "saving changes in wio.yml file ... ")
            if err := types.PrettyPrint(config, wioPath); err != nil {
                log.WriteFailure()
                return err
            } else {
                log.WriteSuccess()
            }
        }
    } else {
        if pac.Context.IsSet("save") {
            return goerr.New("--save flag needs at least one dependency specified")
        }

        log.Infoln(log.Blue, "Installing dependencies")

        projectConfig, err := utils.ReadWioConfig(directory)
        if err != nil {
            return err
        }

        for dependencyName, dependency := range projectConfig.GetDependencies() {
            npmCmdArgs = append(npmCmdArgs, dependencyName+"@"+dependency.Version)
        }
    }

    if len(npmCmdArgs) <= 0 {
        log.Writeln(log.NONE, color.New(color.FgGreen), "nothing to do")
    } else {
        // install packages
        if log.IsVerbose() {
            npmCmdArgs = append(npmCmdArgs, "--verbose")
        }
        npmCmdArgs = append([]string{"install"}, npmCmdArgs...)
        return run.Execute(directory+io.Sep+io.Folder, "npm", npmCmdArgs...)
    }
    return nil
}

// This handles the uninstall and removes packages already downloaded
func (pac Pac) handleUninstall(directory string) error {
    // check install arguments
    uninstallPackage, err := uninstallArgumentCheck(pac.Context.Args())
    if err != nil {
        return err
    }

    remoteDirectory := directory + io.Sep + io.Folder + io.Sep + io.Modules
    wioPath := directory + io.Sep + io.Config

    var config types.IConfig
    if pac.Context.IsSet("save") {
        config, err = utils.ReadWioConfig(directory)
        if err != nil {
            return err
        }
    }

    dependencyDeleted := false

    for _, packageGiven := range uninstallPackage {
        log.Write(log.INFO, color.New(color.FgCyan), "uninstalling %s ... ", packageGiven)
        strip := strings.Split(packageGiven, "@")

        packageName := strip[0]

        if !utils.PathExists(remoteDirectory + io.Sep + packageName) {
            log.Writeln(log.INFO, color.New(color.FgYellow), "does not exist")
            continue
        }

        if err := os.RemoveAll(remoteDirectory + io.Sep + packageName); err != nil {
            log.WriteFailure()
            return err
        } else {
            log.WriteSuccess()
        }

        if pac.Context.IsSet("save") {
            if _, exists := config.GetDependencies()[packageName]; exists {
                dependencyDeleted = true
                delete(config.GetDependencies(), packageName)
            }
        }
    }

    if dependencyDeleted {
        log.Info(log.Cyan, "saving changes in wio.yml file ... ")
        if err := types.PrettyPrint(config, wioPath); err != nil {
            log.WriteFailure()
            return err
        } else {
            log.WriteSuccess()
        }
    }
    return nil
}

// This handles the list command to show dependencies of the project
func handleList(directory string) error {
    return run.Execute(directory+io.Sep+io.Folder, "npm", "list")
}

// This handles the publish command and uses npm to publish packages
func handlePublish(directory string) error {
    if err := publishCheck(directory); err != nil {
        return err
    }
    log.Infoln(log.Cyan, "publishing the package to remote server ... ")
    if log.IsVerbose() {
        return run.Execute(directory, "npm", "publish", "--verbose")
    } else {
        return run.Execute(directory, "npm", "publish")
    }
}
