// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of run package, which contains all the commands to run the project
// Builds, Uploads, and Executes the project
package run

import (
    "bufio"
    goerr "errors"
    "github.com/fatih/color"
    "github.com/urfave/cli"
    "os"
    "os/exec"
    "strings"
    "wio/cmd/wio/commands/create"
    "wio/cmd/wio/commands/run/cmake"
    "wio/cmd/wio/commands/run/dependencies"
    "wio/cmd/wio/config"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/log"
    "wio/cmd/wio/toolchain"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
)

type Run struct {
    Context *cli.Context
    error
}

// get context for the command
func (run Run) GetContext() *cli.Context {
    return run.Context
}

// Runs the build, upload command (acts as one in all command)
func (run Run) Execute() {
    // perform argument check
    directory := performArgumentCheck(run.Context.Args())

    projectConfig, err := utils.ReadWioConfig(directory + io.Sep + "wio.yml")
    if err != nil {
        log.WriteErrorlnExit(err)
    }

    targetName := run.Context.String("target")
    if targetName == config.ProjectDefaults.DefaultTarget {
        targetName = projectConfig.GetTargets().GetDefaultTarget()
    }

    var target types.Target

    if val, exists := projectConfig.GetTargets().GetTargets()[targetName]; exists {
        target = val
    } else {
        log.WriteErrorlnExit(errors.TargetDoesNotExistError{
            TargetName: targetName,
        })
    }

    targetDirectory := directory + io.Sep + ".wio" + io.Sep + "build" + io.Sep + "targets" + io.Sep + targetName

    // show information about the whole run process
    log.Write(log.INFO, color.New(color.FgYellow), "Platform:             ")
    log.Writeln(log.NONE, nil, projectConfig.GetMainTag().GetCompileOptions().GetPlatform())
    log.Write(log.INFO, color.New(color.FgYellow), "Framework:            ")
    log.Writeln(log.NONE, nil, target.GetFramework())
    log.Write(log.INFO, color.New(color.FgYellow), "Target Name:          ")
    log.Writeln(log.NONE, nil, targetName)
    log.Write(log.INFO, color.New(color.FgYellow), "Target Source:        ")
    log.Writeln(log.NONE, nil, directory+io.Sep+target.GetSrc())

    portToUse := ""
    performUpload := false

    // check if we can perform upload and if we can, choose port
    if projectConfig.GetMainTag().GetCompileOptions().GetPlatform() == create.AVR {
        log.Write(log.INFO, color.New(color.FgYellow), "Board:                ")
        log.Writeln(log.NONE, nil, target.GetBoard())

        // select port if upload is triggered as well
        if run.Context.Bool("upload") {
            performUpload = true
            if !run.Context.IsSet("port") {
                if ports, err := toolchain.GetPorts(); err != nil {
                    log.WriteErrorlnExit(errors.AutomaticPortNotDetectedError{})
                } else {
                    port := toolchain.GetArduinoPort(ports)

                    if port == nil {
                        log.WriteErrorlnExit(errors.AutomaticPortNotDetectedError{})
                    } else {
                        log.Write(log.INFO, color.New(color.FgYellow), "Port (Auto):          ")
                        log.Writeln(log.NONE, nil, port.Port+"\n")
                        portToUse = port.Port
                    }
                }
            } else {
                portToUse = run.Context.String("port")

                log.Write(log.INFO, color.New(color.FgYellow), "Port (Manual):        ")
                log.Writeln(log.NONE, nil, portToUse+"\n")
            }
        } else {
            log.Writeln(log.NONE, nil, "")
        }
    } else {
        log.Writeln(log.NONE, nil, "\n")
        log.WriteErrorln(errors.ActionNotSupportedByPlatform{
            Platform:    projectConfig.GetMainTag().GetCompileOptions().GetPlatform(),
            CommandName: "upload",
            Err:         goerr.New("skipping upload"),
        }, true)
    }

    /////////////////////////////////////////// Main CMakeLists ///////////////////////////////////////

    log.Writeln(log.INFO, color.New(color.FgYellow).Add(color.Underline), "building target")
    log.Write(log.INFO, color.New(color.FgCyan), "creating project CMakeLists file ... ")

    queue := log.GetQueue()

    // create CMakeLists.txt file
    if projectConfig.GetMainTag().GetCompileOptions().GetPlatform() == create.AVR {
        if err := cmake.GenerateAvrMainCMakeLists(projectConfig.GetMainTag().GetName(), directory,
            target.GetBoard(), portToUse, target.GetFramework(),
            targetName, target.GetSrc(), target.GetFlags(), target.GetDefinitions()); err != nil {
            log.Writeln(log.NONE, color.New(color.FgRed), "failure")
            log.PrintQueue(queue, log.TWO_SPACES)
            log.WriteErrorlnExit(err)
        } else {
            log.Writeln(log.NONE, color.New(color.FgGreen), "success")
            log.PrintQueue(queue, log.TWO_SPACES)
        }
    } else {
        err = errors.PlatformNotSupportedError{
            Platform: projectConfig.GetMainTag().GetCompileOptions().GetPlatform(),
        }

        log.WriteErrorlnExit(err)
    }

    ////////////////////////////// Dependency Scan and dependencies.cmake /////////////////////////////////////

    log.Write(log.INFO, color.New(color.FgCyan), "scanning dependencies and creating build files ... ")

    projectConfig.GetMainTag().GetName()

    // scan dependencies and create dependencies.cmake file
    if err := dependencies.CreateCMakeDependencyTargets(queue, projectConfig.GetMainTag().GetName(), directory,
        projectConfig.GetType(), target.GetFlags(), target.GetDefinitions(),
        projectConfig.GetDependencies(), projectConfig.GetMainTag().GetCompileOptions().GetPlatform(),
        projectConfig.GetMainTag().GetVersion()); err != nil {
        log.Writeln(log.NONE, color.New(color.FgRed), "failure")
        log.PrintQueue(queue, log.TWO_SPACES)
        log.WriteErrorlnExit(err)
    } else {
        log.Writeln(log.NONE, color.New(color.FgGreen), "success")
        log.PrintQueue(queue, log.TWO_SPACES)
    }

    //////////////////////////////////////////////// Clean ////////////////////////////////////////////

    // clean build files for the target
    if run.Context.Bool("clean") {
        log.Write(log.INFO, color.New(color.FgCyan), "cleaning build files for target %s ... ", targetName)

        err := os.RemoveAll(targetDirectory)
        if err != nil {
            log.Writeln(log.NONE, color.New(color.FgRed), "failure")
            log.WriteErrorlnExit(errors.DeleteDirectoryError{
                DirName: targetDirectory,
                Err:     err,
            })
        } else {
            log.Writeln(log.NONE, color.New(color.FgGreen), "success")
        }
    }

    /////////////////////////////////////////////// Build ///////////////////////////////////////////

    // create a directory for the target
    if err := os.MkdirAll(targetDirectory, os.ModePerm); err != nil {
        log.WriteErrorlnExit(err)
    }

    log.Write(log.INFO, color.New(color.FgCyan), "generating building files for \"%s\" ... ", targetName)
    log.Writeln(log.NONE, nil, "")

    // build targets cmake
    if err := buildTargetCmake(targetDirectory); err != nil {
        log.Writeln(log.INFO_NONE, color.New(color.FgRed), "failure")
        log.WriteErrorlnExit(err)
    } else {
        log.Writeln(log.INFO_NONE, color.New(color.FgGreen), "success")
    }

    // build targets make
    if err := buildTargetMake(targetDirectory); err != nil {
        log.WriteErrorlnExit(err)
    }

    //////////////////////////////////////////// Upload ////////////////////////////////////////

    // upload if instructed or supported
    if performUpload {
        log.Writeln(log.NONE, nil, "")
        log.Writeln(log.INFO, color.New(color.FgYellow).Add(color.Underline), "uploading target")

        if err := uploadTarget(targetDirectory); err != nil {
            log.WriteErrorlnExit(err)
        }
    }
}

// This executes cmake command to generate build files
func buildTargetCmake(buildDirectory string) error {
    cmakeCommand := exec.Command("cmake", "../../", "-G", "Unix Makefiles")
    cmakeCommand.Dir = buildDirectory

    cmakeStderrReader, err := cmakeCommand.StderrPipe()
    if err != nil {
        return err
    }
    cmakeStdoutReader, err := cmakeCommand.StdoutPipe()
    if err != nil {
        return err
    }

    cmakeStdoutScanner := bufio.NewScanner(cmakeStdoutReader)
    go func() {
        for cmakeStdoutScanner.Scan() {
            log.Writeln(log.VERB, nil, cmakeStdoutScanner.Text())
        }
    }()

    // we use deprecated feature so ignore the warning
    cmakeWarning := `CMake Deprecation Warning at CMakeLists.txt:27 (cmake_policy):
  The OLD behavior for policy CMP0023 will be removed from a future version
  of CMake.

  The cmake-policies(7) manual explains that the OLD behaviors of all
  policies are deprecated and that a policy should be set to OLD only under
  specific short-term circumstances.  Projects should be ported to the NEW
  behavior and not rely on setting a policy to OLD.

`
    cmakeStderrScanner := bufio.NewScanner(cmakeStderrReader)
    go func() {
        for cmakeStderrScanner.Scan() {
            line := cmakeStderrScanner.Text()

            if strings.Contains(cmakeWarning, line) {
                continue
            } else {
                log.Writeln(log.ERR, color.New(color.FgRed), line)
            }
        }
    }()

    err = cmakeCommand.Start()
    if err != nil {
        return errors.CommandStartError{
            CommandName: "cmake",
        }
    }

    err = cmakeCommand.Wait()
    if err != nil {
        return errors.CommandWaitError{
            CommandName: "cmake",
        }
    }

    return nil
}

// This executes make command to build the project
func buildTargetMake(buildDirectory string) error {
    makeCommand := exec.Command("make")
    makeCommand.Dir = buildDirectory

    if !log.IsVerbose() {
        // this version is used if verbose mode is not on

        makeCommand.Stdout = os.Stdout
        makeCommand.Stderr = os.Stderr

        if err := makeCommand.Run(); err != nil {
            return errors.CommandWaitError{
                CommandName: "make",
            }
        }

        return nil
    }

    // this version is used if verbose mode is on

    makeStderrReader, err := makeCommand.StderrPipe()
    if err != nil {
        return err
    }
    makeStdoutReader, err := makeCommand.StdoutPipe()
    if err != nil {
        return err
    }

    makeStdoutScanner := bufio.NewScanner(makeStdoutReader)
    go func() {
        for makeStdoutScanner.Scan() {
            // add some colors for logging
            if strings.Contains(makeStdoutScanner.Text(), "Scanning dependencies") {
                log.Writeln(log.INFO, color.New(color.FgGreen), makeStdoutScanner.Text())
            } else if strings.Contains(makeStdoutScanner.Text(), "Built target") {
                log.Writeln(log.INFO, color.New(color.FgGreen), makeStdoutScanner.Text())
            } else {
                log.Writeln(log.INFO, color.New(color.Reset), makeStdoutScanner.Text())
            }
        }
    }()

    makeStderrScanner := bufio.NewScanner(makeStderrReader)
    go func() {
        for makeStderrScanner.Scan() {
            log.Writeln(log.ERR, color.New(color.FgRed), makeStderrScanner.Text())
        }
    }()

    err = makeCommand.Start()
    if err != nil {
        return errors.CommandStartError{
            CommandName: "make",
        }
    }

    err = makeCommand.Wait()
    if err != nil {
        return errors.CommandWaitError{
            CommandName: "make",
        }
    }

    return nil
}

// This executes make upload command to upload the target
func uploadTarget(targetDirectory string) error {
    // execute make upload command
    makeCommand := exec.Command("make", "upload")
    makeCommand.Dir = targetDirectory

    makeStderrReader, err := makeCommand.StderrPipe()
    if err != nil {
        return err
    }
    makeStdoutReader, err := makeCommand.StdoutPipe()
    if err != nil {
        return err
    }

    makeStdoutScanner := bufio.NewScanner(makeStdoutReader)
    go func() {
        for makeStdoutScanner.Scan() {
            log.Writeln(log.INFO, color.New(color.Reset), makeStdoutScanner.Text())
        }
    }()

    makeStderrScanner := bufio.NewScanner(makeStderrReader)
    go func() {
        for makeStderrScanner.Scan() {
            if !log.IsVerbose() {
                log.Writeln(log.INFO, color.New(color.FgGreen), makeStderrScanner.Text())
            } else {
                log.Writeln(log.ERR, color.New(color.FgRed), makeStderrScanner.Text())
            }
        }
    }()

    err = makeCommand.Start()
    if err != nil {
        return errors.CommandStartError{
            CommandName: "make",
        }
    }

    err = makeCommand.Wait()
    if err != nil {
        return errors.CommandWaitError{
            CommandName: "make",
        }
    }

    return nil
}
