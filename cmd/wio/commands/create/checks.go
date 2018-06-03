package create

import (
    "github.com/urfave/cli"
    "os"
    "wio/cmd/wio/utils/io/log"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/commands"
    "bufio"
    "strings"
)

// This check is used to see if the cli arguments are required length
func performArgumentCheck(args cli.Args, isUpdating bool) {
    command := "create"
    if isUpdating {
        command = "update"
    }

    // check to make sure we are given two arguments (one for directory and one for board)
    if len(args) <= 0 {
        log.Norm.Yellow(true, "Directory needs to be provided for creation/update")
        log.Norm.Yellow(false, "Check: ")
        log.Norm.Cyan(true, "wio " + command + " -h")
        os.Exit(3)
    } else if len(args) == 1 && !isUpdating {
        log.Norm.Yellow(true, "Boards is also needed for creation")

        log.Norm.Yellow(false, "Check: ")
        log.Norm.Cyan(true, "wio " + command + " -h")
        os.Exit(3)
    }
}

// This check is used to see if wio.yml file exists and the directory is valid
func performWioExistsCheck(directory string) {
    if !utils.PathExists(directory) {
        log.Norm.Yellow(true, directory+" : no such path exists")
        os.Exit(3)
    }

    if !utils.PathExists(directory + io.Sep + "wio.yml") {
        log.Norm.Yellow(true, "Not a valid wio project: wio.yml file missing")
        os.Exit(3)
    }
}

func performPreUpdateCheck(directory string, projType string) (bool, error) {
    wioPath := directory + io.Sep + "wio.yml"

    isApp, err := utils.IsAppType(wioPath)
    if err != nil {
        return false, err
    }

    if isApp && projType != "app"  || !isApp && projType == "app" {
        // project is of wrong type we cannot update
        return false, nil
    }

    return true, nil
}

/// This method is a crucial peace of check to make sure people do not lose their work. It makes
/// sure that if people are creating the project when there are files in the folder, they mean it
/// and not doing it by mistake. It will warn them to update instead if they want
func performPreCreateCheck(directory string) (bool, error) {
    if !utils.PathExists(directory) {
        return true, nil
    }

    if status, err := utils.IsEmpty(directory); err != nil {
        return false, err
    } else if status {
        return true, nil
    } else {
        message := `The directory is not empty!!
This action will erase everything and will create a new project.
An alternative is to do: wio update <app type> DIRECTORY
Please type y/yes to indicate creation and anything else to indicate abortion: `
        log.Norm.Cyan(false, message)
        reader := bufio.NewReader(os.Stdin)
        text, err := reader.ReadString('\n')
        commands.RecordError(err, "")

        text = strings.TrimSuffix(strings.ToLower(text), "\n")

        if text == "y" || text == "yes" {
            log.Norm.Write(true, "")
            return true, nil
        } else {
            return false, nil
        }
    }
}
