package cmd

import (
	"github.com/urfave/cli"
)

type Command interface {
	GetContext() *cli.Context
	Execute() error
}

type ExitError struct {
	code int
	error
}

func (exitError ExitError) ExitCode() int {
	return exitError.code
}
