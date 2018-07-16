// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package io contains helper functions related to io
// This file contains an interface to print output to io in various colors and modes
package log

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "wio/cmd/wio/utils"

    "github.com/fatih/color"
    "github.com/mattn/go-colorable"
)

type Indentation string

const (
    NO_SPACES   Indentation = ""
    TWO_SPACES  Indentation = "  "
    FOUR_SPACES Indentation = "    "
)

// Log type levels
type Type int

const (
    NONE      Type = 0
    INFO      Type = 1
    INFO_NONE Type = 1
    VERB      Type = 2
    VERB_NONE Type = 2
    WARN      Type = 3
    ERR       Type = 4
    NUM_TYPES      = 5
)

// Colors
var White = color.New(color.FgWhite)
var Cyan = color.New(color.FgCyan)
var Green = color.New(color.FgGreen)
var Yellow = color.New(color.FgYellow)
var Red = color.New(color.FgRed)
var Magenta = color.New(color.FgMagenta)
var Blue = color.New(color.FgHiBlue)
var Default = color.New(color.Reset)

// Log colors and streams
var logTypeColors = [NUM_TYPES]*color.Color{
    White.Add(color.BgCyan),
    White.Add(color.BgCyan),
    White.Add(color.BgCyan),
    White.Add(color.BgYellow),
    White.Add(color.BgRed),
}
var logTypeTags = [NUM_TYPES]string{"NONE", "INFO", "VERB", "WARN", "ERR"}
var createdWriter = writer{verbose: false, warnings: true}
var logOut = colorable.NewColorableStdout()
var logErr = colorable.NewColorableStderr()

// user should not touch this
type writer struct {
    verbose  bool
    warnings bool
}

// Turns verbose mode on. This is the mode when Verbose functions work
func SetVerbose() {
    createdWriter.verbose = true
}

// Disable all the warning shown by wio
func DisableWarnings() {
    createdWriter.warnings = false
}

// Generic Write function
func Write(args ...interface{}) bool {
    a := GetArgs(args...)
    if a.newline {
        a.message = a.message + "\n"
    }
    if nil != a.queue {
        pushLog(a)
        return true
    }
    return write(a)
}

func write(a *Args) bool {
    if a.level == VERB && !IsVerbose() {
        return false
    }
    if a.level == WARN && !showWarnings() {
        return false
    }
    if a.color == nil {
        a.color = Default
    }
    // verbose is INFO behind the screen
    if a.level == VERB {
        a.level = INFO
    }
    // invalid log type defaults to NONE
    if a.level >= NUM_TYPES {
        a.level = NONE
    }

    str := fmt.Sprintf(a.message, a.args...)
    buf := buffer{}
    if a.level == WARN || a.level == ERR {
        logTypeColors[a.level].Fprintf(&buf, "%s", logTypeTags[a.level])
        str = " " + str
    }
    a.color.Fprintf(&buf, "%s", str)

    data := []byte(buf)
    if a.writer != nil {
        a.writer.Write(data)
    } else if a.level == WARN || a.level == ERR {
        logErr.Write(data)
    } else {
        logOut.Write(data)
    }
    return true
}

// Record error/warning to stderr and prompts user for a choice and based on that decides to exists or not
var yesValues = []string{"y", "ye", "yes", "oui"}

func PromptYes(promptMsg string) (bool, error) {
    Info(Yellow, promptMsg+" (y/N): ")

    reader := bufio.NewReader(os.Stdin)
    text, err := reader.ReadString('\n')
    if err != nil {
        return false, err
    }
    text = strings.Trim(text, "\n")
    text = strings.Trim(text, "\r")
    text = strings.Trim(text, " ")
    text = strings.ToLower(text)
    return utils.ContainsNoCase(yesValues, text), nil
}

// This returns true if verbose mode is on and false otherwise
func IsVerbose() bool {
    return createdWriter.verbose
}

// This returns true if warnings are enabled
func showWarnings() bool {
    return createdWriter.warnings
}
