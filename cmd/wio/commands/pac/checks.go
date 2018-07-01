package pac

import (
    "bytes"
    goerr "errors"
    "net/http"
    "os/exec"
    "wio/cmd/wio/errors"
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
func uninstallArgumentCheck(args []string) ([]string, error) {
    if len(args) <= 0 {
        return nil, goerr.New("provide at least one package to uninstall")
    } else {
        return args, nil
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
func publishCheck(directory string) error {
    status, err := utils.IsAppType(directory + io.Sep + io.Config)
    if err != nil {
        return err
    }
    if status {
        return goerr.New("publish command is only supported for project of pkg type")
    }
    return nil
}

// Checks if dependencies are valid wio packages and if they are already pushed
func dependencyCheck(directory string, name string, version string) error {
    resp, err := http.Get("https://www.npmjs.com/package/" + name + "/v/" + version)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // dependency does not exist
    if resp.StatusCode == 404 {
        return goerr.New("dependency: \"" + name + "\" package does not exist on remote server")
    }

    // verify the version by executing npm info command
    npmInfoCommand := exec.Command("npm", "info", name+"@"+version)
    npmInfoCommand.Dir = directory

    cmdOutOutput := &bytes.Buffer{}
    npmInfoCommand.Stdout = cmdOutOutput

    err = npmInfoCommand.Run()
    if err != nil {
        return errors.CommandStartError{
            CommandName: "npm info",
            Err:         err,
        }
    }

    // version does not exists
    if cmdOutOutput.String() == "" {
        return errors.Stringf("dependency [%s@%s] does not exist", name, version)
    }

    return nil
}
