// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package io contains helper functions related to io
// This file contains an interface to print output to io in various colors and modes
package log

import (
    "fmt"
    "github.com/fatih/color"
    "strings"
)

var Norm = writer{verbose: true, status: true}
var Verb = writer{verbose: false, status: true}

// user should not touch this
type writer struct {
    verbose bool
    status  bool
}

// Turns verbose mode on. This is the mode when Verbose functions work
func SetVerbose() {
    Verb.verbose = true
}

// This is used to turn normal mode print on and off. This way a silent mode can be implemented
func SetStatus(status bool) {
    Norm.status = status
}

func write(colorFn func(string, ...interface{}) string, newLine bool, text string, a ...interface{}) {
    str := ""

    if strings.Count(text, "%") == len(a) {
        str = fmt.Sprintf(text, a...)
    } else {
        str = text
    }

    if colorFn != nil {
        fmt.Fprintf(color.Output, colorFn(str))
    } else {
        fmt.Printf(text)
    }

    if newLine {
        fmt.Println()
    }
}

// Red is a convenient helper function to print with red foreground.
func (writer writer) Red(newLine bool, text string, a ...interface{}) {
    if !writer.status || !writer.verbose {
        return
    }

    write(color.RedString, newLine, text, a...)
}

// Green is a convenient helper function to print with green foreground.
func (writer writer) Green(newLine bool, text string, a ...interface{}) {
    if !writer.status || !writer.verbose {
        return
    }

    write(color.GreenString, newLine, text, a...)
}

// Yellow is a convenient helper function to print with yellow foreground.
func (writer writer) Yellow(newLine bool, text string, a ...interface{}) {
    if !writer.status || !writer.verbose {
        return
    }

    write(color.YellowString, newLine, text, a...)
}

// Blue is a convenient helper function to print with blue foreground.
func (writer writer) Blue(newLine bool, text string, a ...interface{}) {
    if !writer.status || !writer.verbose {
        return
    }

    write(color.BlueString, newLine, text, a...)
}

// Magenta is a convenient helper function to print with magenta foreground.
func (writer writer) Magenta(newLine bool, text string, a ...interface{}) {
    if !writer.status || !writer.verbose {
        return
    }

    write(color.MagentaString, newLine, text, a...)
}

// Cyan is a convenient helper function to print with cyan foreground.
func (writer writer) Cyan(newLine bool, text string, a ...interface{}) {
    if !writer.status || !writer.verbose {
        return
    }

    write(color.CyanString, newLine, text, a...)
}

// White is a convenient helper function to print with white foreground.
func (writer writer) White(newLine bool, text string, a ...interface{}) {
    if !writer.status || !writer.verbose {
        return
    }

    write(color.WhiteString, newLine, text, a...)
}

// Normal is a convenient helper function to print with default/normal foreground.
func (writer writer) Write(newLine bool, text string, a ...interface{}) {
    if !writer.status || !writer.verbose {
        return
    }

    write(nil, newLine, text, a...)
}

// Special function to be used when using Verbose mode.
// In this mode, color can be set and other verbose default things can be defined
func (writer writer) Verbose(newLine bool, text string, a ...interface{}) {
    if !writer.status || !writer.verbose {
        return
    }

    writer.Write(newLine, text, a...)
}

// This returns true if verbose mode is on and false otherwise
func (writer writer) IsVerbose() bool {
    return writer.verbose
}
