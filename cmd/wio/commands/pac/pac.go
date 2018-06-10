// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands package, which contains all the commands provided by the tool.
// Package manager for wio. It used npm as a backend and pushes packages to that
package pac

import (
    "bytes"
    "errors"
    "github.com/urfave/cli"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"
    "time"
    "wio/cmd/wio/commands"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/utils/io/log"
)

const (
    ADD     = "add"
    RM      = "rm"
    LIST    = "list"
    INFO    = "info"
    PUBLISH = "publish"
    GET     = "get"
    UPDATE  = "update"
    COLLECT = "collect"
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
    directory, err := filepath.Abs(pac.Context.String("dir"))
    commands.RecordError(err, "")

    switch pac.Type {
    case ADD:
        if len(pac.Context.Args()) == 0 {
            commands.RecordError(errors.New("you need to provide at least one package"), "")
        }

        handleAdd(directory, pac.Context.Args(), pac.Context.Bool("vendor"))
    case RM:
        handleRemove(directory, pac.Context.Args(), pac.Context.Bool("A"))
    case LIST:
        handleList(directory)
    case INFO:
        if len(pac.Context.Args()) == 0 {
            commands.RecordError(errors.New("you need to provide a package name"), "")
        } else if len(pac.Context.Args()) > 1 {
            commands.RecordError(errors.New("only one package name is accepted"), "")
        }

        handleInfo(directory, pac.Context.Args()[0])
    case PUBLISH:
        handlePublish(directory)
    case GET:
        handleGet(directory, pac.Context.Bool("clean"))
        break
    case UPDATE:
        handleUpdate(directory)
        break
    case COLLECT:
        handleCollect(directory)
        break
    }
}

// This package.json file will be used to check if dependencies are up to date
func createPrivatePackageJson(directory string, name string, dependencies types.DependenciesTag) {
    log.Norm.Cyan(false, "setting up packages environment ... ")

    // wio path
    wioPath := directory + io.Sep + ".wio"

    // npm config
    npmConfig := types.NpmConfig{}

    // fill all the fields for package.json
    npmConfig.Name = name
    npmConfig.Version = "0.0.1"
    npmConfig.Main = ".wio.js"

    npmConfig.Dependencies = make(types.NpmDependencyTag)

    log.Verb.Verbose(true, "")

    // add dependencies to package.json
    for dependencyName, dependencyValue := range dependencies {
        if !dependencyValue.Vendor {
            npmConfig.Dependencies[dependencyName] = "^" + dependencyValue.Version
        }
    }

    // write package.json file
    if err := io.NormalIO.WriteJson(wioPath+io.Sep+"package.json", &npmConfig); err != nil {
        commands.RecordError(err, "failure")
    }

    log.Norm.Green(true, "success")
}

// publishes the package to npm backend
func handlePublish(directory string) {
    log.Norm.Yellow(true, "Publishing this package using npm server")
    log.Norm.Cyan(false, "verifying project structure ... ")

    // read wio.yml file
    pkgConfig := &types.PkgConfig{}

    io.NormalIO.ParseYml(directory+io.Sep+"wio.yml", pkgConfig)

    if pkgConfig.MainTag.Name == "" {
        commands.RecordError(errors.New("push command only works on project of type \"pkg\""), "failure")
    } else {
        log.Norm.Green(true, "success")
    }

    log.Norm.Cyan(false, "packing project files ... ")

    pacDir := directory + io.Sep + ".wio" + io.Sep + "pac" + io.Sep + strconv.Itoa(time.Now().Nanosecond())

    // copy src and include folder to .wio/pac folder
    commands.RecordError(utils.CopyDir(directory+io.Sep+"src",
        pacDir+io.Sep+"src"), "failure")
    commands.RecordError(utils.CopyDir(directory+io.Sep+"include",
        pacDir+io.Sep+"include"), "failure")
    commands.RecordError(utils.CopyFile(directory+io.Sep+"wio.yml",
        pacDir+io.Sep+"wio.yml"), "failure")
    commands.RecordError(utils.CopyFile(directory+io.Sep+"README.md",
        pacDir+io.Sep+"README.md"), "failure")
    commands.RecordError(utils.CopyFile(directory+io.Sep+"LICENSE",
        pacDir+io.Sep+"LICENSE"), "failure")

    log.Norm.Green(true, "success")

    log.Verb.Verbose(false, "creating package manager files ... ")

    // npm config
    npmConfig := types.NpmConfig{}

    // fill all the fields for package.json
    npmConfig.Name = pkgConfig.MainTag.Name
    npmConfig.Version = pkgConfig.MainTag.Version
    npmConfig.Description = pkgConfig.MainTag.Description
    npmConfig.Repository = pkgConfig.MainTag.Repository
    npmConfig.Main = ".wio.js"
    npmConfig.Keywords = utils.AppendIfMissing(pkgConfig.MainTag.Keywords, []string{"c++", "c", "wio"})
    npmConfig.Author = pkgConfig.MainTag.Author
    npmConfig.License = pkgConfig.MainTag.License
    npmConfig.Contributors = pkgConfig.MainTag.Contributors

    npmConfig.Dependencies = make(types.NpmDependencyTag)

    log.Verb.Verbose(true, "")

    // add dependencies to package.json
    for dependencyName, dependencyValue := range pkgConfig.DependenciesTag {
        if !dependencyValue.Vendor {
            dependencyCheck(directory, dependencyName, dependencyValue.Version)

            npmConfig.Dependencies[dependencyName] = "^" + dependencyValue.Version
        }
    }

    // write package.json file
    if err := io.NormalIO.WriteJson(pacDir+io.Sep+"package.json", &npmConfig); err != nil {
        commands.RecordError(err, "failure")
    }

    // write .wio.js file so that this is the entry point
    if err := io.NormalIO.WriteFile(pacDir+io.Sep+".wio.js", []byte("console.log('H"+
        "i!!! Welcome to Wio world!')")); err != nil {
        removeNpmFiles(pacDir)
        commands.RecordError(err, "failure")
    }

    log.Verb.Verbose(true, "success")
    log.Norm.Cyan(false, "running publish command ... ")

    // execute cmake command
    npmPublishCommand := exec.Command("npm", "publish")
    npmPublishCommand.Dir = pacDir

    // Stderr buffer
    cmdErrOutput := &bytes.Buffer{}
    npmPublishCommand.Stderr = cmdErrOutput

    if log.Verb.IsVerbose() {
        npmPublishCommand.Stdout = os.Stdout
    }
    err := npmPublishCommand.Run()
    if err != nil {
        removeNpmFiles(pacDir)

        // this means user needs to login
        if strings.Contains(cmdErrOutput.String(), "npm adduser") ||
            strings.Contains(cmdErrOutput.String(), "npm login") {
            log.Norm.Yellow(true, "\nNPM backend is used for package management")
            log.Norm.Green(true, "Run \"npm adduser\" to login/create account to publish the package")

            os.Exit(2)
        }

        commands.RecordError(err, "failure", strings.Trim(cmdErrOutput.String(), "\n"))
    } else {
        log.Norm.Green(true, "success")
    }

    removeNpmFiles(pacDir)

    log.Verb.Verbose(true, "success")
    log.Norm.Yellow(true, pkgConfig.MainTag.Name+"@"+pkgConfig.MainTag.Version+" published!!")
}

// checks if dependencies are valid wio packages and if they are already pushed
func dependencyCheck(directory string, dependencyName string, dependencyVersion string) {
    log.Verb.Verbose(false, "dependency: checking if "+dependencyName+" package exists ... ")

    resp, err := http.Get("https://www.npmjs.com/package/" + dependencyName + "/v/" + dependencyVersion)
    if err != nil {
        commands.RecordError(err, "failure")
    }

    // dependency does not exist
    if resp.StatusCode == 404 {
        log.Verb.Verbose(true, "failure")
        commands.RecordError(errors.New("dependency: \"" + dependencyName + "\" package does not exist on remote "+
            "server"), "")
    }
    resp.Body.Close()

    log.Verb.Verbose(true, "success")
    log.Verb.Verbose(false, "dependency: checking if " + dependencyName + "@" + dependencyVersion+
        " version exists ... ")

    // verify the version by executing npm info command
    npmInfoCommand := exec.Command("npm", "info", dependencyName+"@"+dependencyVersion)
    npmInfoCommand.Dir = directory

    // Stderr buffer
    cmdErrOutput := &bytes.Buffer{}
    cmdOutOutput := &bytes.Buffer{}
    npmInfoCommand.Stderr = cmdErrOutput
    npmInfoCommand.Stdout = cmdOutOutput

    err = npmInfoCommand.Run()
    if err != nil {
        commands.RecordError(err, "failure", strings.Trim(cmdErrOutput.String(), "\n"))
    }

    // version does not exists
    if cmdOutOutput.String() == "" {
        commands.RecordError(errors.New("dependency: \"" + dependencyName + "@" + dependencyVersion+
            "\" version does not exist"), "failure")
    } else {
        log.Verb.Verbose(true, "success")

        log.Verb.Verbose(false, "dependency: checking if " + dependencyName + "@" + dependencyVersion+
            " is a valid wio package ... ")

        // check if the package is a wio package by checking C, C++ and wio flags
        pat := regexp.MustCompile(`keywords: .*[\r\n]`)
        s := pat.FindString(cmdOutOutput.String())

        // if wio, c and c++ found, this package is a valid wio package
        if strings.Contains(s, "wio") && strings.Contains(s, "c") && strings.Contains(s, "c++") {
            log.Verb.Verbose(true, "success")
        } else {
            commands.RecordError(errors.New("dependency: \"" + dependencyName + "@" + dependencyVersion+
                "\" is not a wio package"), "failure")
        }
    }
}

// removes files used for npm publish
func removeNpmFiles(directory string) {
    log.Verb.Verbose(false, "removing npm files ... ")

    commands.RecordError(os.RemoveAll(directory), "failure")
    log.Verb.Verbose(true, "success")
}

// gets and updates the packages to the versions specified in wio.yml file
func handleGet(directory string, clean bool) {
    log.Norm.Yellow(true, "Getting packages from npm server")

    // clean npm_modules in .wio folder
    if clean {
        log.Norm.Cyan(false, "removing all the pulled packages ... ")

        if err := os.RemoveAll(directory + io.Sep + ".wio" + io.Sep + "node_modules"); err != nil {
            commands.RecordError(err, "failure")
        } else {
            log.Norm.Green(true, "success")
        }
    }

    // read wio.yml file. We gonna use dependencies tag so it does not matter if this is app or pkg
    pkgConfig := &types.PkgConfig{}

    // create private json file
    createPrivatePackageJson(directory, filepath.Base(directory), pkgConfig.DependenciesTag)

    commands.RecordError(io.NormalIO.ParseYml(directory+io.Sep+"wio.yml", pkgConfig), "")

    // add dependencies to package.json
    for dependencyName, dependencyValue := range pkgConfig.DependenciesTag {
        if dependencyValue.Vendor {
            continue
        }

        dependencyCheck(directory, dependencyName, dependencyValue.Version)

        log.Norm.Cyan(false, "pulling " + dependencyName + "@" + dependencyValue.Version+
            " package ... ")

        // execute cmake command
        npmInstallCommand := exec.Command("npm", "install", dependencyName+"@"+dependencyValue.Version)
        npmInstallCommand.Dir = directory + io.Sep + ".wio"

        // Stderr buffer
        cmdErrOutput := &bytes.Buffer{}
        npmInstallCommand.Stderr = cmdErrOutput

        if log.Verb.IsVerbose() {
            log.Norm.Verbose(true, "")
            npmInstallCommand.Stdout = os.Stdout
        }
        err := npmInstallCommand.Run()
        if err != nil {
            commands.RecordError(err, "failure", strings.Trim(cmdErrOutput.String(), "\n"))
        }

        log.Norm.Green(true, "success")
    }

    log.Norm.Yellow(true, "All packages pulled successfully")
}

// add and update dependencies from cli
func handleAdd(directory string, args []string, vendor bool) {
    log.Norm.Yellow(true, "Adding/Updating dependencies")

    // read wio.yml file. We gonna use dependencies tag so it does not matter if this is app or pkg
    pkgConfig := &types.PkgConfig{}

    commands.RecordError(io.NormalIO.ParseYml(directory+io.Sep+"wio.yml", pkgConfig), "")

    changed := false

    for _, addArg := range args {
        newDependency := strings.Split(addArg, "@")

        depName := newDependency[0]
        depVersion := "latest"

        if len(newDependency) > 1 {
            depVersion = newDependency[1]
        }

        // check if this dependency exists and check it's version
        if val, ok := pkgConfig.DependenciesTag[depName]; ok {
            // override with vendor
            if !val.Vendor && vendor {
                pkgConfig.DependenciesTag[depName].Vendor = true
                pkgConfig.DependenciesTag[depName].Version = ""

                log.Norm.Cyan(true, "overridden remote dependency by vendor dependency: "+depName)
                changed = true
            } else if vendor {
                log.Norm.Cyan(true, "unchanged vendor dependency: "+depName)
            } else if val.Version != depVersion {
                if val.Vendor {
                    log.Norm.Cyan(true, "delete vendor dependency before updating to remote")
                } else {
                    // change the version since it already exists
                    pkgConfig.DependenciesTag[depName].Version = depVersion

                    log.Norm.Cyan(true, "updated remote dependency: " + depName + "@" + val.Version+
                        "   ->   "+ depName+ "@"+ depVersion)

                    changed = true
                }
            } else {
                log.Norm.Cyan(true, "unchanged remote dependency: "+depName+"@"+val.Version)
            }
        } else {
            if vendor {
                pkgConfig.DependenciesTag[depName] = &types.DependencyTag{Vendor: true}

                log.Norm.Cyan(true, "added vendor dependency: "+depName)
            } else {
                pkgConfig.DependenciesTag[depName] = &types.DependencyTag{Version: depVersion}

                log.Norm.Cyan(true, "added remote dependency: "+depName+"@"+depVersion)
            }

            changed = true
        }
    }

    if !changed {
        log.Norm.Yellow(true, "Dependencies remained unchanged!")
    } else {
        // check if the configuration is app
        if data, err := io.NormalIO.ReadFile(directory + io.Sep + "wio.yml"); err != nil {
            commands.RecordError(err, "")
        } else {
            // if there is "app:" tag use app configuration
            if strings.Contains(string(data), "app:") {
                appConfig := &types.AppConfig{}

                appConfig.DependenciesTag = pkgConfig.DependenciesTag

                if err = io.NormalIO.ParseYml(directory+io.Sep+"wio.yml", appConfig); err != nil {
                    commands.RecordError(err, "")
                }

                commands.RecordError(utils.PrettyPrintConfigSpacing(appConfig, directory+io.Sep+"wio.yml"), "")
            } else {
                commands.RecordError(utils.PrettyPrintConfigSpacing(pkgConfig, directory+io.Sep+"wio.yml"), "")
            }
        }

        log.Norm.Yellow(true, "Provides dependencies added/updated successfully")
    }
}

// remove dependencies from cli
func handleRemove(directory string, args []string, all bool) {
    log.Norm.Yellow(true, "Removing dependencies")

    // read wio.yml file. We gonna use dependencies tag so it does not matter if this is app or pkg
    pkgConfig := &types.PkgConfig{}

    commands.RecordError(io.NormalIO.ParseYml(directory+io.Sep+"wio.yml", pkgConfig), "")

    changed := false

    if len(pkgConfig.DependenciesTag) == 0 {
        log.Norm.Cyan(true, "project has No dependencies")
    } else if all {
        pkgConfig.DependenciesTag = make(map[string]*types.DependencyTag)

        log.Norm.Cyan(true, "deleted All the dependencies")
        changed = true
    } else {
        for _, depName := range args {
            // check if this dependency exists and check it's version
            if _, ok := pkgConfig.DependenciesTag[depName]; ok {
                delete(pkgConfig.DependenciesTag, depName)

                log.Norm.Cyan(true, "deleted "+depName)
                changed = true
            } else {
                log.Norm.Cyan(true, "no such dependency of name: \""+depName+"\", skipping!")
            }
        }
    }

    if !changed {
        log.Norm.Yellow(true, "Dependencies remained unchanged!")
    } else {
        // check if the configuration is app
        if data, err := io.NormalIO.ReadFile(directory + io.Sep + "wio.yml"); err != nil {
            commands.RecordError(err, "")
        } else {
            // if there is "app:" tag use app configuration
            if strings.Contains(string(data), "app:") {
                appConfig := &types.AppConfig{}

                if err = io.NormalIO.ParseYml(directory+io.Sep+"wio.yml", appConfig); err != nil {
                    commands.RecordError(err, "")
                }

                appConfig.DependenciesTag = pkgConfig.DependenciesTag

                commands.RecordError(utils.PrettyPrintConfigSpacing(appConfig, directory+io.Sep+"wio.yml"), "")
            } else {
                commands.RecordError(utils.PrettyPrintConfigSpacing(pkgConfig, directory+io.Sep+"wio.yml"), "")
            }
        }

        log.Norm.Yellow(true, "Provided dependencies removed successfully")
    }
}

// list all the project dependencies
func handleList(directory string) {
    // read wio.yml file. We gonna use dependencies tag so it does not matter if this is app or pkg
    pkgConfig := &types.PkgConfig{}

    commands.RecordError(io.NormalIO.ParseYml(directory+io.Sep+"wio.yml", pkgConfig), "")

    if len(pkgConfig.DependenciesTag) == 0 {
        log.Norm.Yellow(true, "This project has no dependencies")
    } else {
        log.Norm.Yellow(true, "Project dependencies: ")
        for key, value := range pkgConfig.DependenciesTag {
            if value.Vendor {
                log.Norm.Cyan(false, "vendor: ")
                log.Norm.Cyan(true, key)
            } else {
                log.Norm.Cyan(false, "remote: ")
                log.Norm.Cyan(true, key+"@"+value.Version)
            }
        }
    }
}

// provide information about one individual package
func handleInfo(directory string, depName string) {
    // read wio.yml file. We gonna use dependencies tag so it does not matter if this is app or pkg
    pkgConfig := &types.PkgConfig{}

    commands.RecordError(io.NormalIO.ParseYml(directory+io.Sep+"wio.yml", pkgConfig), "")

    if len(pkgConfig.DependenciesTag) == 0 {
        log.Norm.Yellow(true, "This project has no dependencies")
    } else {
        if val, ok := pkgConfig.DependenciesTag[depName]; ok {
            log.Norm.Yellow(true, depName+": ")
            log.Norm.Cyan(true, "is vendor: "+strconv.FormatBool(val.Vendor))

            if !val.Vendor {
                log.Norm.Cyan(true, "version: "+val.Version)
            }
            log.Norm.Cyan(true, "compile flags: ["+strings.Join(val.DependencyFlags, ",")+"]")
        }
    }
}

func handleUpdate(directory string) {
    // TODO future versions
    // do npm outdated and see if something needs to be updated
    // then read package.json file and update the version of wio dependency
}

func handleCollect(directory string) {
    // TODO future versions
}
