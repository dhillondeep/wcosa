package pac

import (
    "bytes"
    goerr "errors"
    "github.com/fatih/color"
    "net/http"
    "os/exec"
    "regexp"
    "strings"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/log"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
)

// Checks arguments to verify what to install
func installArgumentCheck(args []string) []string {
    if len(args) <= 0 {
        return []string{"all"}
    } else {
        return args
    }
}

// Checks arguments to verify what to uninstall
func uninstallArgumentCheck(args []string) []string {
    if len(args) <= 0 {
        log.WriteErrorlnExit(errors.ProgramArgumentsError{
            CommandName:  "uninstall",
            ArgumentName: "package name",
            Err:          goerr.New("atleast one package must be provided"),
        })
        return nil
    } else {
        return args
    }
}

// checks arguments to verify what to collect
func collectArgumentCheck(args []string) []string {
    if len(args) <= 0 {
        return []string{"_______all__________"}
    } else {
        return args
    }
}

// checks arguments to verify what to publish
func publishCheck(directory string) {
    status, err := utils.IsAppType(directory + io.Sep + "wio.yml")
    if err != nil {
        log.WriteErrorlnExit(err)
    }

    if status {
        log.WriteErrorlnExit(goerr.New("publish command is only supported for project of pkg type"))
    }
}

// Checks if dependencies are valid wio packages and if they are already pushed
func dependencyCheck(queue *log.Queue, directory string, dependencyName string, dependencyVersion string) error {
    log.QueueWrite(queue, log.VERB, nil, "dependency: checking if "+dependencyName+" package exists ... ")

    resp, err := http.Get("https://www.npmjs.com/package/" + dependencyName + "/v/" + dependencyVersion)
    if err != nil {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
        return err
    }
    defer resp.Body.Close()

    // dependency does not exist
    if resp.StatusCode == 404 {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
        return goerr.New("dependency: \"" + dependencyName + "\" package does not exist on remote server")
    } else {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
    }

    log.QueueWrite(queue, log.VERB, nil, "dependency: checking if " + dependencyName + "@" + dependencyVersion + " "+
        "version exists ... ")

    // verify the version by executing npm info command
    npmInfoCommand := exec.Command("npm", "info", dependencyName+"@"+dependencyVersion)
    npmInfoCommand.Dir = directory

    cmdOutOutput := &bytes.Buffer{}
    npmInfoCommand.Stdout = cmdOutOutput

    err = npmInfoCommand.Run()
    if err != nil {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
        return errors.CommandStartError{
            CommandName: "npm info",
            Err:         err,
        }
    }

    // version does not exists
    if cmdOutOutput.String() == "" {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
        return goerr.New("dependency: \"" + dependencyName + "@" + dependencyVersion +
            "\" version does not exist")
    } else {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
        log.QueueWrite(queue, log.VERB, nil, "dependency: checking if " + dependencyName + "@" + dependencyVersion+
            " is a valid wio package ... ")

        // check if the package is a wio package by checking C, C++ and wio flags
        pat := regexp.MustCompile(`keywords: .*[\r\n]`)
        s := pat.FindString(cmdOutOutput.String())

        // if wio, c and c++ found, this package is a valid wio package
        if strings.Contains(s, "wio") && strings.Contains(s, "c") && strings.Contains(s, "c++") &&
            strings.Contains(s, "pkg") && strings.Contains(s, "iot") {
            log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
        } else {
            log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
            return goerr.New("dependency: \"" + dependencyName + "@" + dependencyVersion +
                "\" is not a wio package")
        }
    }

    return nil
}
