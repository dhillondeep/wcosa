// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands package, which contains all the commands provided by the tool.
// Package manager for wio. It used npm as a backend and pushes packages to that
package pac

import (
    "github.com/urfave/cli"
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
}

// Get context for the command
func (pac Pac) GetContext() *cli.Context {
    return pac.Context
}

// Executes the libraries command
func (pac Pac) Execute() error {
    return nil
}
