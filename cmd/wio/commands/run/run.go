// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of run package, which contains all the commands to run the project
// Builds, Uploads, and Executes the project
package run

import (
    goerr "errors"
    "github.com/fatih/color"
    "github.com/urfave/cli"
    "os"
    sysio "io"
    "os/exec"
    "wio/cmd/wio/commands/run/cmake"
    "wio/cmd/wio/commands/run/dependencies"
    cfg "wio/cmd/wio/config"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/log"
    "wio/cmd/wio/toolchain"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
    "bytes"
)

type Run struct {
    Context *cli.Context
    error
}

type runInfo struct {
    context *cli.Context
    config  *types.Config

    directory string
    targets   []string
}

// get context for the command
func (run Run) GetContext() *cli.Context {
    return run.Context
}

// Runs the build, upload command (acts as one in all command)
func (run Run) Execute() {
    directory, err := readDirectory(run.Context.Args())
    if err != nil {
        log.WriteErrorlnExit(err)
    }
    config, err := utils.ReadWioConfig(directory + io.Sep + "wio.yml")
    if err != nil {
        log.WriteErrorlnExit(err)
    }


    targetName := run.Context.String("target")
    if targetName == cfg.ProjectDefaults.DefaultTarget {
        targetName = config.GetTargets().GetDefaultTarget()
    }

    var target types.Target

    if val, exists := config.GetTargets().GetTargets()[targetName]; exists {
        target = val
    } else {
        log.WriteErrorlnExit(errors.TargetDoesNotExistError{
            TargetName: targetName,
        })
    }

    targetDirectory := directory + io.Sep + ".wio" + io.Sep + "build" + io.Sep + "targets" + io.Sep + targetName

    // show information about the whole run process
    log.Write(log.INFO, color.New(color.FgYellow), "Platform:             ")
    log.Writeln(log.NONE, nil, config.GetMainTag().GetCompileOptions().GetPlatform())
    log.Write(log.INFO, color.New(color.FgYellow), "Framework:            ")
    log.Writeln(log.NONE, nil, target.GetFramework())
    log.Write(log.INFO, color.New(color.FgYellow), "Target Name:          ")
    log.Writeln(log.NONE, nil, targetName)
    log.Write(log.INFO, color.New(color.FgYellow), "Target Source:        ")
    log.Writeln(log.NONE, nil, directory+io.Sep+target.GetSrc())

    portToUse := ""
    performUpload := false

    // check if we can perform upload and if we can, choose port
    if config.GetMainTag().GetCompileOptions().GetPlatform() == constants.AVR {
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
            Platform:    config.GetMainTag().GetCompileOptions().GetPlatform(),
            CommandName: "upload",
            Err:         goerr.New("skipping upload"),
        }, true)
    }

    /////////////////////////////////////////// Main CMakeLists ///////////////////////////////////////

    log.Writeln(log.INFO, color.New(color.FgYellow).Add(color.Underline), "building target")
    log.Write(log.INFO, color.New(color.FgCyan), "creating project CMakeLists file ... ")

    queue := log.GetQueue()

    // create CMakeLists.txt file
    if config.GetMainTag().GetCompileOptions().GetPlatform() == constants.AVR {
        if err := cmake.GenerateAvrCmakeLists(config.GetMainTag().GetName(), directory,
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
            Platform: config.GetMainTag().GetCompileOptions().GetPlatform(),
        }

        log.WriteErrorlnExit(err)
    }

    ////////////////////////////// Dependency Scan and dependencies.cmake /////////////////////////////////////

    log.Write(log.INFO, color.New(color.FgCyan), "scanning dependencies and creating build files ... ")

    config.GetMainTag().GetName()

    // scan dependencies and create dependencies.cmake file
    if err := dependencies.CreateCMakeDependencyTargets(queue, config.GetMainTag().GetName(), directory,
        config.GetType(), target.GetFlags(), target.GetDefinitions(),
        config.GetDependencies(), config.GetMainTag().GetCompileOptions().GetPlatform(),
        config.GetMainTag().GetVersion()); err != nil {
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
    if err := configTarget(targetDirectory); err != nil {
        log.Writeln(log.INFO_NONE, color.New(color.FgRed), "failure")
        log.WriteErrorlnExit(err)
    } else {
        log.Writeln(log.INFO_NONE, color.New(color.FgGreen), "success")
    }

    // build targets make
    if err := buildTarget(targetDirectory); err != nil {
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

func (info runInfo) build() error {
    targets := make([]types.Target, 0, len(info.targets))
    projectTargets := info.config.GetTargets().GetTargets()
    for _, targetName := range info.targets {
        if _, exist := projectTargets[targetName]; exist {
            targets = append(targets, projectTargets[targetName])
        } else {
            log.Warnln("Unrecognized target name: %s", targetName)
        }
    }
    if len(info.targets) <= 0 {
        defaultName := info.config.GetTargets().GetDefaultTarget()
        targets = append(targets, projectTargets[defaultName])
    }
    
}

func (run Run) configure(dir string) error {

}

func configTarget(dir string) error {
    return execute(dir, "cmake", "../../", "-G", "Unix Makefiles")
}

func buildTarget(dir string) error {
    return execute(dir, "make")
}

func uploadTarget(dir string) error {
    return execute(dir, "make", "upload")
}

func cleanTarget(dir string) error {
    return execute(dir, "make", "clean")
}

func execute(dir string, name string, args ...string) error {
    cmd := exec.Command(name, args...)
    cmd.Dir = dir
    var stdoutBuf bytes.Buffer
    var stderrBuf bytes.Buffer
    stdoutIn, err := cmd.StdoutPipe()
    if err != nil {
        return err
    }
    stderrIn, err := cmd.StderrPipe()
    if err != nil {
        return err
    }
    var errStdout error
    var errStderr error
    stdout := sysio.MultiWriter(os.Stdout, &stdoutBuf)
    stderr := sysio.MultiWriter(os.Stderr, &stderrBuf)
    err = cmd.Start()
    if err != nil {
        return err
    }
    go func() { _, errStdout = sysio.Copy(stdout, stdoutIn) }()
    go func() { _, errStderr = sysio.Copy(stderr, stderrIn) }()
    err = cmd.Wait()
    if err != nil {
        return err
    }
    if errStdout != nil {
        return errStdout
    }
    if errStderr != nil {
        return errStderr
    }
    return nil
}
