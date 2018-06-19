// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of devices package, which contains all the commands related to handling devices
// Runs the serial monitor and list devices
package devices

import (
    goerr "errors"
    "fmt"
    "github.com/fatih/color"
    "github.com/urfave/cli"
    "go.bug.st/serial.v1"
    "os"
    "os/signal"
    "strings"
    "syscall"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/log"
    "wio/cmd/wio/toolchain"
)

type Devices struct {
    Context *cli.Context
    Type    byte
    error
}

// get context for the command
func (devices Devices) GetContext() *cli.Context {
    return devices.Context
}

const (
    LIST    = 0
    MONITOR = 1
)

// Runs the build command when cli build option is provided
func (devices Devices) Execute() {
    switch devices.Type {
    case MONITOR:
        HandleMonitor(devices.Context.Int("baud"), devices.Context.IsSet("port"), devices.Context.String("port"))
        break
    case LIST:
        handlePorts(devices.Context.Bool("basic"), devices.Context.Bool("show-all"))
        break
    }
}

// Provides information abouts ports
func handlePorts(basic bool, showAll bool) {
    ports, err := toolchain.GetPorts()
    if err != nil {
        log.WriteErrorlnExit(goerr.New("port data could not be gathered from the operating system"))
    }

    log.Write(log.INFO, color.New(color.FgCyan), "Num of total ports: ")
    log.Writeln(log.NONE, nil, "%d\n", len(ports.Ports))

    numOpenPorts := 0
    for _, port := range ports.Ports {
        if port.Product == "None" && !showAll {
            continue
        } else {
            numOpenPorts++
        }

        log.Writeln(log.INFO, color.New(color.FgYellow), port.Port)

        if !basic {
            log.Write(log.INFO, color.New(color.FgCyan), "Product:          ")
            log.Writeln(log.NONE, nil, port.Description)
            log.Write(log.INFO, color.New(color.FgCyan), "Manufacturer:     ")
            log.Writeln(log.NONE, nil, port.Manufacturer)
            log.Write(log.INFO, color.New(color.FgCyan), "Serial Number:    ")
            log.Writeln(log.NONE, nil, port.SerialNumber)
            log.Write(log.INFO, color.New(color.FgCyan), "Hwid:             ")
            log.Writeln(log.NONE, nil, port.Hwid)
            log.Write(log.INFO, color.New(color.FgCyan), "Vid:              ")
            log.Writeln(log.NONE, nil, port.Vid)
        }

        log.Writeln(log.NONE, nil, "")
    }

    log.Write(log.INFO, color.New(color.FgCyan), "Num of open ports: ")
    log.Writeln(log.NONE, nil, "%d", numOpenPorts)
}

// Opens monitor to see serial data
func HandleMonitor(baud int, portDefined bool, portProvided string) {
    var port *toolchain.SerialPort

    ports, err := toolchain.GetPorts()
    if err != nil {
        port = nil
    } else {
        port = toolchain.GetArduinoPort(ports)
    }

    portToUse := portProvided

    if !portDefined {
        if port == nil {
            log.WriteErrorlnExit(errors.AutomaticPortNotDetectedError{})
        } else {
            portToUse = port.Port
        }
    }

    // Open the first serial port detected at 9600bps N81
    mode := &serial.Mode{
        BaudRate: baud,
        Parity:   serial.NoParity,
        DataBits: 8,
        StopBits: serial.OneStopBit,
    }
    serialPort, err := serial.Open(portToUse, mode)
    if err != nil {
        if strings.Contains(err.Error(), "Invalid serial port") {
            log.WriteErrorlnExit(goerr.New("invalid baud rate"))
        }
    }

    defer serialPort.Close()

    log.Write(log.INFO, color.New(color.FgCyan), "Wio Serial Monitor")
    log.Write(log.INFO, color.New(color.FgYellow), "  @  ")
    log.Write(log.INFO, color.New(color.FgCyan), portToUse)
    log.Write(log.INFO, color.New(color.FgYellow), "  @  ")
    log.Writeln(log.INFO, color.New(color.FgCyan), "%d", baud)
    log.Writeln(log.INFO, color.New(color.FgCyan), "--- Quit: Ctrl+C ---")

    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        log.Writeln(log.INFO, nil, "\n--- exit ---")
        os.Exit(1)
    }()

    // Read and print the response
    buff := make([]byte, 100)
    for {
        // Reads up to 100 bytes
        n, err := serialPort.Read(buff)
        if err != nil {
            panic(err)
            break
        }
        if n == 0 {
            fmt.Println("\nEOF")
            break
        }
        fmt.Printf("%v", string(buff[:n]))
    }
}
