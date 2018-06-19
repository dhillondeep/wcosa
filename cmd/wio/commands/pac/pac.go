// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands package, which contains all the commands provided by the tool.
// Package manager for wio. It used npm as a backend and pushes packages to that
package pac

import (
    "github.com/urfave/cli"
    "os"
    "os/exec"
    "wio/cmd/wio/log"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
    "github.com/fatih/color"
    "wio/cmd/wio/errors"
    "bufio"
    "strings"
    goerr "errors"
    "wio/cmd/wio/types"
    "io/ioutil"
    "regexp"
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
    error
}

// Get context for the command
func (pac Pac) GetContext() *cli.Context {
    return pac.Context
}

// Executes the libraries command
func (pac Pac) Execute() {
    // check if valid wio project
    directory, err := os.Getwd()
    if err != nil {
        log.WriteErrorlnExit(err)
    }

    if !utils.PathExists(directory + io.Sep + "wio.yml") {
        log.WriteErrorlnExit(errors.ConfigMissing{})
    }

    switch pac.Type {
    case INSTALL:
        pac.handleInstall(directory)
        break
    case UNINSTALL:
        pac.handleUninstall(directory)
        break
    case COLLECT:
        pac.handleCollect(directory)
        break
    case LIST:
        pac.handleList(directory)
        break
    case PUBLISH:
        pac.handlePublish(directory)
        break
    }
}

// This handles the install command and uses npm to install packages
func (pac Pac) handleInstall(directory string) {
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
                log.Writeln(log.NONE, color.New(color.FgRed), "failure")
                log.WriteErrorlnExit(err)
            } else {
                if err := os.RemoveAll(directory + io.Sep + ".wio" + io.Sep + "package-lock.json "); err != nil {
                    log.Writeln(log.NONE, color.New(color.FgRed), "failure")
                    log.WriteErrorlnExit(err)
                } else {
                    log.Writeln(log.NONE, color.New(color.FgGreen), "success")
                }
            }
        }
    }

    var npmInstallArguments []string

    if installPackage[0] != "all" {
        log.Writeln(log.INFO, color.New(color.FgCyan), "installing %s ...", installPackage)
        npmInstallArguments = append(npmInstallArguments, installPackage...)

        if pac.Context.IsSet("save") {
            projectConfig, err := utils.ReadWioConfig(directory + io.Sep + "wio.yml")
            if err != nil {
                log.WriteErrorlnExit(errors.ConfigParsingError{
                    Err: err,
                })
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

            log.Write(log.INFO, color.New(color.FgCyan), "saving changes in wio.yml file ...")
            if err := utils.PrettyPrintConfig(projectConfig, directory+io.Sep+"wio.yml",
                pac.Context.Bool("config-help")); err != nil {
                log.Writeln(log.INFO, color.New(color.FgRed), "failure")
                log.WriteErrorlnExit(errors.WriteFileError{
                    FileName: directory + io.Sep + "wio.yml",
                    Err:      err,
                })
            } else {
                log.Writeln(log.INFO, color.New(color.FgGreen), "success")
            }
        }
    } else {
        if pac.Context.IsSet("save") {
            log.WriteErrorlnExit(errors.ProgramArgumentsError{
                CommandName:  "install",
                ArgumentName: "--save",
                Err:          goerr.New("--save flag needs atleast one dependency specified"),
            })
        }

        log.Writeln(log.INFO, color.New(color.FgCyan), "installing dependencies ...")

        projectConfig, err := utils.ReadWioConfig(directory + io.Sep + "wio.yml")
        if err != nil {
            log.WriteErrorlnExit(errors.ConfigParsingError{
                Err: err,
            })
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
            log.WriteErrorlnExit(err)
        }
        npmStdoutReader, err := cmdNpm.StdoutPipe()
        if err != nil {
            log.WriteErrorlnExit(err)
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
            log.WriteErrorlnExit(errors.CommandStartError{
                CommandName: "wio install",
            })
        }

        err = cmdNpm.Wait()
        if err != nil {
            log.WriteErrorlnExit(errors.CommandWaitError{
                CommandName: "wio install",
            })
        }
    }
}

// This handles the uninstall and removes packages already downloaded
func (pac Pac) handleUninstall(directory string) {
    // check install arguments
    uninstallPackage := uninstallArgumentCheck(pac.Context.Args())

    remoteDirectory := directory + io.Sep + ".wio" + io.Sep + "node_modules"

    var projectConfig types.Config
    var err error
    if pac.Context.IsSet("save") {
        projectConfig, err = utils.ReadWioConfig(directory + io.Sep + "wio.yml")
        if err != nil {
            log.WriteErrorlnExit(errors.ConfigParsingError{
                Err: err,
            })
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
            log.Writeln(log.INFO, color.New(color.FgRed), "failure")
            log.WriteErrorlnExit(errors.DeleteDirectoryError{
                DirName: remoteDirectory + io.Sep + packageName,
                Err:     err,
            })
        } else {
            log.Writeln(log.INFO, color.New(color.FgGreen), "success")
        }

        if pac.Context.IsSet("save") {
            if _, exists := projectConfig.GetDependencies()[packageName]; exists {
                dependencyDeleted = true
                delete(projectConfig.GetDependencies(), packageName)
            }
        }
    }

    if dependencyDeleted {
        log.Write(log.INFO, color.New(color.FgCyan), "saving changes in wio.yml file ...")
        if err := utils.PrettyPrintConfig(projectConfig, directory+io.Sep+"wio.yml",
            pac.Context.Bool("config-help")); err != nil {
            log.Writeln(log.INFO, color.New(color.FgRed), "failure")
            log.WriteErrorlnExit(errors.WriteFileError{
                FileName: directory + io.Sep + "wio.yml",
                Err:      err,
            })
        } else {
            log.Writeln(log.INFO, color.New(color.FgGreen), "success")
        }
    }
}

// This handles the collect command to collect remote dependencies into vendor folder
func (pac Pac) handleCollect(directory string) {
    // check install arguments
    collectPackages := collectArgumentCheck(pac.Context.Args())

    remoteDirectory := directory + io.Sep + ".wio" + io.Sep + "node_modules"

    var projectConfig types.Config
    var err error
    if pac.Context.IsSet("save") {
        projectConfig, err = utils.ReadWioConfig(directory + io.Sep + "wio.yml")
        if err != nil {
            log.WriteErrorlnExit(errors.ConfigParsingError{
                Err: err,
            })
        }
    }

    modified := false

    if collectPackages[0] == "_______all__________" {
        files, err := ioutil.ReadDir(remoteDirectory)
        if err != nil {
            log.WriteErrorlnExit(err)
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
            log.Writeln(log.INFO, color.New(color.FgYellow), "remote package does not exist")
            continue
        }

        if utils.PathExists(directory + io.Sep + "vendor" + io.Sep + packageName) {
            log.Writeln(log.INFO, color.New(color.FgYellow), "package already exists in vendor")
            continue
        }

        if err := utils.CopyDir(remoteDirectory+io.Sep+packageName, directory+io.Sep+"vendor"+io.Sep+packageName); err != nil {
            log.Writeln(log.INFO, color.New(color.FgRed), "failure")
            log.WriteErrorlnExit(err)
        } else {
            log.Writeln(log.INFO, color.New(color.FgGreen), "success")
        }

        if pac.Context.IsSet("save") {
            if val, exists := projectConfig.GetDependencies()[packageName]; exists {
                val.Vendor = true
                modified = true
            }
        }
    }

    if modified {
        log.Write(log.INFO, color.New(color.FgCyan), "saving changes in wio.yml file ...")
        if err := utils.PrettyPrintConfig(projectConfig, directory+io.Sep+"wio.yml",
            pac.Context.Bool("config-help")); err != nil {
            log.Writeln(log.INFO, color.New(color.FgRed), "failure")
            log.WriteErrorlnExit(errors.WriteFileError{
                FileName: directory + io.Sep + "wio.yml",
                Err:      err,
            })
        } else {
            log.Writeln(log.INFO, color.New(color.FgGreen), "success")
        }
    }
}

// This handles the list command to show dependencies of the project
func (pac Pac) handleList(directory string) {
    cmdNpm := exec.Command("npm", "list")
    cmdNpm.Dir = directory + io.Sep + ".wio"
    cmdNpm.Stderr = os.Stderr

    npmStdoutReader, err := cmdNpm.StdoutPipe()
    if err != nil {
        log.WriteErrorlnExit(err)
    }

    npmStdoutScanner := bufio.NewScanner(npmStdoutReader)
    go func() {
        for npmStdoutScanner.Scan() {
            line := npmStdoutScanner.Text()
            if strings.Contains(line, directory) {
                line = directory
            }

            log.Writeln(log.INFO, color.New(color.Reset), line)
        }
    }()

    err = cmdNpm.Start()
    if err != nil {
        log.WriteErrorlnExit(errors.CommandStartError{
            CommandName: "wio install",
        })
    }

    err = cmdNpm.Wait()
    if err != nil {
        log.WriteErrorlnExit(errors.CommandWaitError{
            CommandName: "wio install",
        })
    }
}

// This handles the publish command and uses npm to publish packages
func (pac Pac) handlePublish(directory string) {
    publishCheck(directory)

    // read wio.yml file
    pkgConfig := &types.PkgConfig{}

    if err := io.NormalIO.ParseYml(directory+io.Sep+"wio.yml", pkgConfig); err != nil {
        log.WriteErrorlnExit(errors.ConfigParsingError{
            Err: err})
    }

    log.Write(log.INFO, color.New(color.FgCyan), "checking files and packing them ... ")

    queue := log.GetQueue()

    log.QueueWrite(queue, log.VERB, nil, "creating package.json file ... ")

    // npm config
    npmConfig := types.NpmConfig{}

    // fill all the fields for package.json
    npmConfig.Name = pkgConfig.MainTag.Meta.Name
    npmConfig.Version = pkgConfig.MainTag.Meta.Version
    npmConfig.Description = pkgConfig.MainTag.Meta.Description
    npmConfig.Repository = pkgConfig.MainTag.Meta.Repository
    npmConfig.Main = ".wio.js"
    npmConfig.Keywords = utils.AppendIfMissing(pkgConfig.MainTag.Meta.Keywords, []string{"c++", "c", "wio", "pkg", "iot"})
    npmConfig.Author = pkgConfig.MainTag.Meta.Author
    npmConfig.License = pkgConfig.MainTag.Meta.License
    npmConfig.Contributors = pkgConfig.MainTag.Meta.Contributors

    versionPat := regexp.MustCompile(`[0-9]+.[0-9]+.[0-9]+`)
    stringPat := regexp.MustCompile(`[\w"]+`)

    // verify tag values
    if !stringPat.MatchString(npmConfig.Author) {
        log.Writeln(log.NONE, color.New(color.FgRed), "failure")
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
        log.PrintQueue(queue, log.TWO_SPACES)
        log.WriteErrorlnExit(goerr.New("author must be specified for a package"))
    }
    if !stringPat.MatchString(npmConfig.Description){
        log.Writeln(log.NONE, color.New(color.FgRed), "failure")
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
        log.PrintQueue(queue, log.TWO_SPACES)
        log.WriteErrorlnExit(goerr.New("description must be specified for a package"))
    }
    if !versionPat.MatchString(npmConfig.Version) {
        log.Writeln(log.NONE, color.New(color.FgRed), "failure")
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
        log.PrintQueue(queue, log.TWO_SPACES)
        log.WriteErrorlnExit(goerr.New("package does not have a valid version"))
    }
    if !stringPat.MatchString(npmConfig.License) {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgYellow), "warning")
        log.WriteErrorln(errors.ProgrammingArgumentAssumption{
            CommandName: "wio publish",
            ArgumentName: "LICENSE",
            Err: goerr.New("LICENSE was not provided, and is assumed to be MIT"),
        }, true)
        npmConfig.License = "MIT"
    }

    npmConfig.Dependencies = make(types.NpmDependencyTag)

    subQueue := log.GetQueue()

    // add dependencies to package.json
    for dependencyName, dependencyValue := range pkgConfig.DependenciesTag {
        if !dependencyValue.Vendor {
            if err := dependencyCheck(subQueue, directory, dependencyName, dependencyValue.Version); err != nil {
                log.Writeln(log.NONE, color.New(color.FgRed), "failure")
                log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
                log.CopyQueue(subQueue, queue, log.TWO_SPACES)
                log.PrintQueue(queue, log.TWO_SPACES)
                log.WriteErrorlnExit(err)
            }

            npmConfig.Dependencies[dependencyName] = dependencyValue.Version
        }
    }

    // write package.json file
    if err := io.NormalIO.WriteJson(directory+io.Sep+"package.json", &npmConfig); err != nil {
        log.Writeln(log.NONE, color.New(color.FgRed), "failure")
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
        log.PrintQueue(queue, log.TWO_SPACES)
        log.WriteErrorlnExit(errors.WriteFileError{
            FileName: "package.json",
        })
    } else {
        log.Writeln(log.NONE, color.New(color.FgGreen), "success")
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
        log.PrintQueue(queue, log.TWO_SPACES)
    }

    log.Writeln(log.INFO, color.New(color.FgCyan), "publishing the package to remote server ... ")

    var cmdNpm *exec.Cmd

    // execute cmake command
    if !log.IsVerbose() {
        cmdNpm = exec.Command("npm", "publish")
    } else {
        cmdNpm = exec.Command("npm", "publish", "--verbose")
    }

    cmdNpm.Dir = directory

    npmStderrReader, err := cmdNpm.StderrPipe()
    if err != nil {
        log.WriteErrorlnExit(err)
    }
    npmStdoutReader, err := cmdNpm.StdoutPipe()
    if err != nil {
        log.WriteErrorlnExit(err)
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

            if line == "" {
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

                log.Writeln(log.INFO, color.New(color.Reset), line)
            }
        }
    }()

    err = cmdNpm.Start()
    if err != nil {
        log.WriteErrorlnExit(errors.CommandStartError{
            CommandName: "wio publish",
        })
    }

    err = cmdNpm.Wait()
    if err != nil {
        log.WriteErrorlnExit(errors.CommandWaitError{
            CommandName: "wio publish",
        })
    }
}
