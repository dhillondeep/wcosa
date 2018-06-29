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
    "os"
    "regexp"
    "strings"
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

// This must be called at the beggining
func Init() {}

// This provides a queue that can be used to log at different levels
func GetQueue() *Queue {
    return NewQueue(5)
}

// Write Queue
func QueueWrite(queue *Queue, logType Type, providedColor *color.Color, message string, a ...interface{}) {
    pushLog(queue, logType, providedColor, message, a...)
}

// Writeln Queue
func QueueWriteln(queue *Queue, logType Type, providedColor *color.Color, message string, a ...interface{}) {
    QueueWrite(queue, logType, providedColor, message+"\n", a...)
}

// Copy one queue to another
func CopyQueue(fromQueue *Queue, toQueue *Queue, spaces Indentation) {
    for {
        if len(*fromQueue) <= 0 {
            break
        } else {
            value := popLog(fromQueue)

            value.text = string(spaces) + value.text

            pat := regexp.MustCompile(`\n[\s]+[a-zA-Z]`)
            findStr := strings.Trim(pat.FindString(value.text), "\n")

            value.text = pat.ReplaceAllString(value.text, "\n"+string(spaces)+findStr)
            pushLog(toQueue, value.logType, value.providedColor, value.text)
        }
    }
}

// Print Queue on the console with a set indentation
func PrintQueue(queue *Queue, spaces Indentation) {
    index := 0

    for {
        if index >= len(*queue) {
            break
        } else {
            value := popLog(queue)

            value.text = string(spaces) + value.text

            pat := regexp.MustCompile(`\n[\s]+[a-zA-Z]`)
            findStr := strings.Trim(pat.FindString(value.text), "\n")

            value.text = pat.ReplaceAllString(value.text, "\n"+string(spaces)+findStr)
            Write(value.logType, value.providedColor, value.text)
        }
    }
}

// Generic Writeln function
func Writeln(args ...interface{}) bool {
    return Write(append(args, true)...)
}

// Generic Write function
func Write(args ...interface{}) bool {
    var queue *Queue = nil
    logType := NONE
    logColor := Default
    message := ""
    newline := false
    printfArgs := make([]interface{}, 0, len(args))
    for _, arg := range args {
        switch val := arg.(type) {
        case Type:
            logType = val
            break
        case *color.Color:
            logColor = val
            break
        case string:
            if "" == message {
                message = val
            } else {
                printfArgs = append(printfArgs, val)
            }
            break
        case *Queue:
            queue = val
            break
        case bool:
            newline = val
            break
        case error:
            message = val.Error()
            break
        default:
            break
        }
    }
    if newline {
        message = message + "\n"
    }
    if nil != queue {
        pushLog(queue, logType, logColor, message, printfArgs...)
        return true
    }
    return write(logType, logColor, message, printfArgs...)
}

func write(logType Type, providedColor *color.Color, message string, a ...interface{}) bool {
    if (logType == VERB && !IsVerbose()) || (logType == WARN && !showWarnings()) {
        return false
    }
    if providedColor == nil {
        providedColor = Default
    }

    // verbose is INFO behind the screen
    if logType == VERB {
        logType = INFO
    }

    // invalid log type defaults to NONE
    if logType >= NUM_TYPES {
        logType = NONE
    }

    str := fmt.Sprintf(message, a...)

    outStream := logOut
    if logType == WARN || logType == ERR {
        logTypeColors[logType].Fprintf(logOut, "%s", logTypeTags[logType])
        str = " " + str
        outStream = logErr
    }
    providedColor.Fprintf(outStream, "%s", str)
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
func WriteErrorAndPrompt(err error, logType Type, promptRightAnswer string, caseSensitive bool) {
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

// Shorthands
func Info(args ...interface{}) {
    Write(append(args, INFO)...)
}

func Infoln(args ...interface{}) {
    Writeln(append(args, INFO)...)
}

func Verb(args ...interface{}) {
    Write(append(args, VERB)...)
}

func Verbln(args ...interface{}) {
    Writeln(append(args, VERB)...)
}

func Warn(args ...interface{}) {
    Write(append(args, WARN, Yellow))
}

func Warnln(args ...interface{}) {
    Writeln(append(args, WARN, Yellow)...)
}

func Err(args ...interface{}) {
    Write(append(args, ERR, Red)...)
}

func Errln(args ...interface{}) {
    Writeln(append(args, ERR, Red)...)
}

func ErrExit(args ...interface{}) {
    Err(args...)
    os.Exit(1)
}

func ErrlnExit(args ...interface{}) {
    Errln(args...)
    os.Exit(1)
}

func WriteSuccess(args ...interface{}) {
    Writeln(append(args, Green, "success")...)
}

func WriteFailure(args ...interface{}) {
    Writeln(append(args, Red, "failure")...)
}

// This returns true if verbose mode is on and false otherwise
func IsVerbose() bool {
    return createdWriter.verbose
}

// This returns true if warnings are enabled
func showWarnings() bool {
    return createdWriter.warnings
}
