package commands

import (
    "fmt"
    "github.com/urfave/cli"
    "os"
    "wio/cmd/wio/utils/io/log"
)

type Command interface {
    GetContext() *cli.Context
    Execute()
}

type ExitError struct {
    code int
    error
}

func (exitError ExitError) ExitCode() int {
    return exitError.code
}

// RecordError function allows for error handling with error code and
// nice console error logs
func RecordError(err error, message string, more ...interface{}) {
    if err == nil {
        return
    }

    if message != "" {
        log.Norm.Red(true, message)
    }

    log.Norm.Red(true, "Error Report: ")

    if len(more) > 0 {
        fmt.Fprintln(os.Stderr, more...)
    }

    exitCoder := ExitError{code: 2, error: err}

    cli.HandleExitCoder(exitCoder)
}
