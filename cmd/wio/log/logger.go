// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package io contains helper functions related to io
// This file contains an interface to print output to io in various colors and modes
package log

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	NO_SPACES   = ""
	TWO_SPACES   = "  "
	FOUR_SPACES  = "  "
	SIX_SPACES   = "  "
	EIGHT_SPACES = "  "
)

const (
	VERB_NONE = "VERB_NONE"             // Does not Show INFO tag in verbose mode (only activates in verbose mode)
	INFO_NONE = "INFO_NONE"             // Only Shows text in Verbose mode
	NONE      = "NONE"                  // Does not show INFO tag in regular mode
	INFO      = "INFO"                  // Shows like a normal text in regular mode and INFO tag in verbose mode
	VERB      = "VERB"
	ERR       = "ERR"
	WARN      = "WARN"
)

var logTypeColors map[string]*color.Color
var logTypeStream map[string]io.Writer
var createdWriter = writer{verbose: false, warnings: true}

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

// This must be called at the beggining
func Init() {
	logTypeColors = make(map[string]*color.Color)
	logTypeStream = make(map[string]io.Writer)

	logTypeColors[VERB_NONE] = color.New(color.FgWhite).Add(color.BgCyan)
	logTypeStream[VERB_NONE] = colorable.NewColorableStdout()
	logTypeColors[NONE] = color.New(color.FgWhite).Add(color.BgCyan)
	logTypeStream[NONE] = colorable.NewColorableStdout()
	logTypeColors[INFO] = color.New(color.FgWhite).Add(color.BgCyan)
	logTypeStream[INFO] = colorable.NewColorableStdout()
	logTypeColors[ERR] = color.New(color.FgWhite).Add(color.BgRed)
	logTypeStream[ERR] = colorable.NewColorableStderr()
	logTypeColors[WARN] = color.New(color.FgWhite).Add(color.BgYellow)
	logTypeStream[WARN] = colorable.NewColorableStderr()
}

// This provides a queue that can be used to log at different levels
func GetQueue() *Queue {
	return NewQueue(1)
}

// Write Queue
func QueueWrite(queue *Queue, logType string, providedColor *color.Color, message string, a ...interface{}) {
	pushLog(queue, logType, providedColor, message, a...)
}

// Writeln Queue
func QueueWriteln(queue *Queue, logType string, providedColor *color.Color, message string, a ...interface{}) {
	QueueWrite(queue, logType, providedColor, message+"\n", a...)
}

// Copy one queue to another
func CopyQueue(fromQueue *Queue, toQueue *Queue, spaces string) {
	for {
		if fromQueue.count <= 0 {
			break
		} else {
			value := popLog(fromQueue)

			value.text = spaces + value.text

			pat := regexp.MustCompile(`\n[\s]+[a-zA-Z]`)
			findStr := strings.Trim(pat.FindString(value.text), "\n")

			value.text = pat.ReplaceAllString(value.text, "\n"+spaces+findStr)
			pushLog(toQueue, value.logType, value.providedColor, value.text)
		}
	}
}

// Print Queue on the console
func PrintQueue(queue *Queue, spaces string) {
	index := 0

	for {
		if index >= queue.count {
			break
		} else {
			value := popLog(queue)

			value.text = spaces + value.text

			pat := regexp.MustCompile(`\n[\s]+[a-zA-Z]`)
			findStr := strings.Trim(pat.FindString(value.text), "\n")

			value.text = pat.ReplaceAllString(value.text, "\n"+spaces+findStr)
			Write(value.logType, value.providedColor, value.text)
		}
	}
}

// Generic Writeln function
func Writeln(logType string, providedColor *color.Color, message string, a ...interface{}) bool {
	if !showWarnings() && logType == WARN {
		return false
	}

	if Write(logType, providedColor, message, a...) {
        fmt.Println("")
    } else {
        return false
    }

    return true
}

// Generic Write function
func Write(logType string, providedColor *color.Color, message string, a ...interface{}) bool {
    if logType == INFO_NONE && IsVerbose() {
        return false
    }

	if providedColor == nil {
		providedColor = color.New(color.Reset)
	}

	// they only apply in verbose mode
	if (logType == VERB_NONE || logType == VERB) && !IsVerbose() {
		return false
	}

	// only applies when show warnings is enabled
	if !showWarnings() && logType == WARN {
		return false
	}

	// verbose is INFO behind the screen
	if logType == VERB {
		logType = INFO
	}

	if _, ok := logTypeColors[logType]; !ok {
		logType = INFO
	}

	messageColor := color.New(color.Reset)

	if providedColor != nil {
		messageColor = providedColor
	}

	var str string
	if len(a) <= 0 {
		str = fmt.Sprint(message)
	} else {
		str = fmt.Sprintf(message, a...)
	}

	if logType == NONE || logType == VERB_NONE {
		messageColor.Fprintf(logTypeStream[logType], "%s", str)
		return true
	}

	if logType != INFO || IsVerbose() {
		color.New(color.FgHiWhite).Fprintf(logTypeStream[logType], "%s ", "wio")
		logTypeColors[logType].Fprintf(logTypeStream[logType], "%s", strings.ToUpper(logType))
		messageColor.Fprintf(logTypeStream[logType], " %s", str)
	} else if logType == INFO && !IsVerbose() {
		messageColor.Fprintf(logTypeStream[logType], "%s", str)
	}

	return true
}

// Record error to stderr and prints a new line. It also exists the program with an error code
func WriteErrorlnExit(err error) {
	if err == nil {
		return
	}

	Writeln(ERR, color.New(color.Reset), err.Error())
	os.Exit(1)
}

// Record error/warning to stderr and prints a new line
func WriteErrorln(err error, isWarning bool) {
	if err == nil {
		return
	}

	logType := ERR
	if isWarning {
		logType = WARN
	}

	Writeln(logType, color.New(color.Reset), err.Error())
}

// Record error/warning to stderr and prompts user for a choice and based on that decides to exists or not
func WriteErrorAndPrompt(err error, logType string, promptRightAnswer string, caseSensitive bool) {
	if err == nil {
		return
	}

	Write(logType, color.New(color.FgYellow), err.Error())

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	WriteErrorlnExit(err)

	text = strings.TrimSuffix(text, "\n")

	if caseSensitive {
		promptRightAnswer = strings.ToLower(promptRightAnswer)
		text = strings.ToLower(text)
	}

	if text != promptRightAnswer {
		os.Exit(0)
	} else {
		fmt.Fprint(colorable.NewColorableStderr(), "\n")
	}
}

// This returns true if verbose mode is on and false otherwise
func IsVerbose() bool {
	return createdWriter.verbose
}

// This returns true if warnings are enabled
func showWarnings() bool {
	return createdWriter.warnings
}
