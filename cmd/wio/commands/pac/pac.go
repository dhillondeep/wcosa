// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands package, which contains all the commands provided by the tool.
// Package manager for wio. It used npm as a backend and pushes packages to that
package pac

import (
    "bufio"
    goerr "errors"
    "github.com/fatih/color"
    "github.com/urfave/cli"
    "io/ioutil"
    "os"
    "os/exec"
    "regexp"
    "strings"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/log"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/commands/run"
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
    if !utils.PathExists(directory + io.Sep + "wio.yml") {
        return errors.ConfigMissing{}
    }

    switch pac.Type {
    case INSTALL:
        return pac.handleInstall(directory)
    case UNINSTALL:
        return pac.handleUninstall(directory)
    case COLLECT:
        return pac.handleCollect(directory)
    case LIST:
        return pac.handleList(directory)
    case PUBLISH:
        return pac.handlePublish(directory)
    default:
        return goerr.New("invalid pac command")
    }
}

// This handles the install command and uses npm to install packages
func (pac Pac) handleInstall(directory string) error {
    // check install arguments
    installPackage := installArgumentCheck(pac.Context.Args())

    remoteDirectory := directory + io.Sep + ".wio" + io.Sep + "node_modules"

    // clean npm_modules in .wio folder
    if pac.Context.Bool("clean") {
        log.Write(log.INFO, color.New(color.FgCyan), "cleaning npm packages ... ")

        if !utils.PathExists(remoteDirectory) || !utils.PathExists(directory + io.Sep + "wio.yml") {
            log.Writeln(log.NONE, color.New(color.FgGreen), "nothing to do")
        } else {
            if err := os.RemoveAll(remoteDirectory); err != nil {
                log.WriteFailure()
                return err
            } else {
                if err := os.RemoveAll(directory + io.Sep + ".wio" + io.Sep + "package-lock.json "); err != nil {
                    log.WriteFailure()
                    return err
                } else {
                    log.Writeln(log.NONE, color.New(color.FgGreen), "success")
                }
            }
        }
    }

    var npmInstallArguments []string

    if installPackage[0] != "all" {
        log.Write(log.INFO, color.New(color.FgCyan), "installing %s ... ", installPackage)
        npmInstallArguments = append(npmInstallArguments, installPackage...)

        if pac.Context.IsSet("save") {
            projectConfig, err := utils.ReadWioConfig(directory + io.Sep + "wio.yml")
            if err != nil {
                return err
            }

            dependencies := projectConfig.GetDependencies()

            if projectConfig.GetDependencies() == nil {
                dependencies = types.DependenciesTag{}
            }

            for _, packageGiven := range installPackage {
                strip := strings.Split(packageGiven, "@")

                packageName := strip[0]
                dependencies[packageName] = &types.DependencyTag{
                    Version: "0.0.1",
                    Vendor:  false,
                }

                if len(strip) > 1 {
                    dependencies[packageName].Version = strip[1]
                }
            }

            projectConfig.SetDependencies(dependencies)

            log.Write(log.INFO, color.New(color.FgCyan), "saving changes in wio.yml file ... ")
            if err := projectConfig.PrettyPrint(directory + io.Sep + "wio.yml"); err != nil {
                log.WriteFailure()
                return err
            } else {
                log.Writeln(log.NONE, color.New(color.FgGreen), "success")
            }
        }
    } else {
        if pac.Context.IsSet("save") {
            return goerr.New("--save flag needs at least one dependency specified")
        }

        log.Write(log.INFO, color.New(color.FgCyan), "installing dependencies ... ")

        projectConfig, err := utils.ReadWioConfig(directory + io.Sep + "wio.yml")
        if err != nil {
            return err
        }

        for dependencyName, dependency := range projectConfig.GetDependencies() {
            npmInstallArguments = append(npmInstallArguments, dependencyName+"@"+dependency.Version)
        }
    }

    if len(npmInstallArguments) <= 0 {
        log.Writeln(log.NONE, color.New(color.FgGreen), "nothing to do")
    } else {
        // install packages
        cmdNpm := exec.Command("npm", "install")

        if log.IsVerbose() {
            npmInstallArguments = append(npmInstallArguments, "--verbose")
        }

        cmdNpm.Args = append([]string{"npm", "install"}, npmInstallArguments...)
        cmdNpm.Dir = directory + io.Sep + ".wio"

        npmStderrReader, err := cmdNpm.StderrPipe()
        if err != nil {
            return err
        }
        npmStdoutReader, err := cmdNpm.StdoutPipe()
        if err != nil {
            return err
        }

        npmStdoutScanner := bufio.NewScanner(npmStdoutReader)
        go func() {
            for npmStdoutScanner.Scan() {
                log.Writeln(log.INFO, color.New(color.Reset), npmStdoutScanner.Text())
            }
        }()

        npmStderrScanner := bufio.NewScanner(npmStderrReader)
        go func() {
            for npmStderrScanner.Scan() {
                line := npmStderrScanner.Text()
                line = strings.Trim(strings.Replace(line, "npm", "wio", -1), " ")

                if strings.Contains(line, "no such file or directory") ||
                    strings.Contains(line, "No description") || strings.Contains(line, "No repository field") ||
                    strings.Contains(line, "No README data") || strings.Contains(line, "No license field") ||
                    strings.Contains(line, "notice created a lockfile as package-lock.json.") {
                    continue
                } else if line == "" {
                    continue
                } else if strings.Contains(line, "debug.log") || strings.Contains(line, "A complete log of this run can be found in:") {
                    continue
                } else {
                    line = strings.Replace(line, "wio", "", -1)
                    line = strings.Replace(line, "verb", "", -1)
                    line = strings.Replace(line, "info", "", -1)
                    line = strings.Replace(line, "WARN ", "", -1)
                    line = strings.Replace(line, "ERR!", "", -1)
                    line = strings.Replace(line, "notarget", "", -1)

                    line = strings.Trim(line, " ")

                    log.Writeln(log.ERR, color.New(color.Reset), line)
                }
            }
        }()

        err = cmdNpm.Start()
        if err != nil {
            return err
        }
        err = cmdNpm.Wait()
        if err != nil {
            return err
        }
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

    remoteDirectory := directory + io.Sep + ".wio" + io.Sep + "node_modules"

    var projectConfig *types.Config
    if pac.Context.IsSet("save") {
        projectConfig, err = utils.ReadWioConfig(directory + io.Sep + "wio.yml")
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
            if _, exists := projectConfig.GetDependencies()[packageName]; exists {
                dependencyDeleted = true
                delete(projectConfig.GetDependencies(), packageName)
            }
        }
    }

    if dependencyDeleted {
        log.Info(log.Cyan, "saving changes in wio.yml file ... ")
        if err := projectConfig.PrettyPrint(directory + io.Sep + "wio.yml"); err != nil {
            log.WriteFailure()
            return err
        } else {
            log.WriteSuccess()
        }
    }
    return nil
}

// This handles the collect command to collect remote dependencies into vendor folder
func (pac Pac) handleCollect(directory string) error {
    // check install arguments
    collectPackages := collectArgumentCheck(pac.Context.Args())

    remoteDirectory := directory + io.Sep + ".wio" + io.Sep + "node_modules"

    var projectConfig *types.Config
    var err error
    if pac.Context.IsSet("save") {
        projectConfig, err = utils.ReadWioConfig(directory + io.Sep + "wio.yml")
        if err != nil {
            return err
        }
    }

    modified := false

    if collectPackages[0] == "_______all__________" {
        files, err := ioutil.ReadDir(remoteDirectory)
        if err != nil {
            return err
        }

        for _, f := range files {
            collectPackages = append(collectPackages, f.Name())
        }
    }

    for _, packageGiven := range collectPackages {
        // this is put in by verification process
        if packageGiven == "_______all__________" {
            continue
        }

        log.Write(log.INFO, color.New(color.FgCyan), "copying %s to vendor ... ", packageGiven)

        strip := strings.Split(packageGiven, "@")

        packageName := strip[0]

        if !utils.PathExists(remoteDirectory + io.Sep + packageName) {
            log.Infoln(log.Yellow, "remote package does not exist")
            continue
        }

        if utils.PathExists(directory + io.Sep + "vendor" + io.Sep + packageName) {
            log.Infoln(log.Yellow, "package already exists in vendor")
            continue
        }

        if err := utils.CopyDir(remoteDirectory+io.Sep+packageName, directory+io.Sep+"vendor"+io.Sep+packageName); err != nil {
            log.WriteFailure()
            return err
        } else {
            log.WriteSuccess()
        }

        if pac.Context.IsSet("save") {
            if val, exists := projectConfig.GetDependencies()[packageName]; exists {
                val.Vendor = true
                modified = true
            }
        }
    }

    if modified {
        log.Write(log.INFO, color.New(color.FgCyan), "saving changes in wio.yml file ... ")
        if err := projectConfig.PrettyPrint(directory + io.Sep + "wio.yml"); err != nil {
            log.WriteFailure()
            return err
        } else {
            log.WriteSuccess()
        }
    }
    return nil
}

// This handles the list command to show dependencies of the project
func (pac Pac) handleList(directory string) error {
    return run.Execute(directory+io.Sep+".wio", "npm", "list")
}

// This handles the publish command and uses npm to publish packages
func (pac Pac) handlePublish(directory string) error {
    if err := publishCheck(directory); err != nil {
        return err
    }

    // read wio.yml file
    pkgConfig := &types.PkgConfig{}

    if err := io.NormalIO.ParseYml(directory+io.Sep+"wio.yml", pkgConfig); err != nil {
        return err
    }

    log.Info(log.Cyan, "checking files and packing them ... ")

    // npm config
    meta := pkgConfig.MainTag.Meta
    npmConfig := types.NpmConfig{
        Name:         pkgConfig.GetMainTag().GetName(),
        Version:      pkgConfig.GetMainTag().GetVersion(),
        Description:  meta.Description,
        Repository:   meta.Repository,
        Main:         ".wio.js",
        Keywords:     meta.Keywords,
        Author:       meta.Author,
        License:      meta.License,
        Contributors: meta.Contributors,
    }

    // fill all the fields for package.json
    versionPat := regexp.MustCompile(`[0-9]+.[0-9]+.[0-9]+`)
    stringPat := regexp.MustCompile(`[\w"]+`)

    // verify tag values
    if !stringPat.MatchString(npmConfig.Author) {
        log.WriteFailure()
        return goerr.New("author must be specified for a package")
    }
    if !stringPat.MatchString(npmConfig.Description) {
        log.WriteFailure()
        return goerr.New("description must be specified for a package")
    }
    if !versionPat.MatchString(npmConfig.Version) {
        log.WriteFailure()
        return goerr.New("package does not have a valid version")
    }
    if !stringPat.MatchString(npmConfig.License) {
        npmConfig.License = "MIT"
    }
    log.WriteSuccess()

    npmConfig.Dependencies = make(types.NpmDependencyTag)

    // add dependencies to package.json
    for dependencyName, dependencyValue := range pkgConfig.DependenciesTag {
        if !dependencyValue.Vendor {
            if err := dependencyCheck(directory, dependencyName, dependencyValue.Version); err != nil {
                log.WriteFailure()
                return err
            }

            npmConfig.Dependencies[dependencyName] = dependencyValue.Version
        }
    }

    // write package.json file
    if err := io.NormalIO.WriteJson(directory+io.Sep+"package.json", &npmConfig); err != nil {
        log.WriteFailure()
        return err
    }
    log.WriteSuccess()

    log.Writeln(log.INFO, color.New(color.FgCyan), "publishing the package to remote server ... ")

    // execute cmake command
    if log.IsVerbose() {
        return run.Execute(directory, "npm", "publish", "--verbose")
    } else {
        return run.Execute(directory, "npm", "publish")
    }
}
