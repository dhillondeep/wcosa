package run

import (
    "os"
    "wio/cmd/wio/log"
    "wio/cmd/wio/errors"
    goerr "errors"
)

func performArgumentCheck(args []string) string {
    var directory string
    var err error

    // check directory
    if len(args) <= 0 {
        directory, err = os.Getwd()

        log.WriteErrorlnExit(err)

        err = errors.ProgrammingArgumentAssumption{
            CommandName:  "create",
            ArgumentName: "directory",
            Err:          goerr.New("directory is not provided so current directory is used: " + directory),
        }

        log.WriteErrorln(err, true)
    } else {
        directory = args[0]
    }

    return directory
}
