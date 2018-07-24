// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of devices package, which contains all the commands related to handling devices
// Runs the serial monitor and list devices
package devices

import (
    "fmt"
    "os"
    "os/signal"
    "strings"
    "syscall"
    "wio/internal/toolchain"
    "wio/pkg/log"
    "wio/pkg/util"

    "github.com/urfave/cli"
    "go.bug.st/serial.v1"
)

type Devices struct {
    Context *cli.Context
    Type    byte
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
func (devices Devices) Execute() error {
    switch devices.Type {
    case MONITOR:
        return HandleMonitor(devices.Context.Int("baud"), devices.Context.IsSet("port"), devices.Context.String("port"))
    case LIST:
        return handlePorts(devices.Context.Bool("basic"), devices.Context.Bool("show-all"))
    default:
        return util.Error("invalid device command")
    }
}

// Provides information abouts ports
func handlePorts(basic bool, showAll bool) error {
    ports, err := toolchain.GetPorts()
    if err != nil {
        return err
    }

    log.Info(log.Cyan, "Num of ports: ")
    log.Infoln("%d\n", len(ports.Ports))

    numOpenPorts := 0
    for _, port := range ports.Ports {
        if port.Product != "None" {
            numOpenPorts++
        }

        if port.Product == "None" && !showAll {
            continue
        }

        log.Infoln(log.Yellow, port.Port)

        if !basic {
            log.Info(log.Cyan, "Product:          ")
            log.Infoln(port.Description)
            log.Info(log.Cyan, "Manufacturer:     ")
            log.Infoln(port.Manufacturer)
            log.Info(log.Cyan, "Serial Number:    ")
            log.Infoln(port.SerialNumber)
            log.Info(log.Cyan, "Hwid:             ")
            log.Infoln(port.Hwid)
            log.Info(log.Cyan, "Vid:              ")
            log.Infoln(port.Vid)
        }

        log.Infoln()
    }

    log.Info(log.Cyan, "Num of open ports: ")
    log.Infoln("%d", numOpenPorts)
    return nil
}

// Opens monitor to see serial data
func HandleMonitor(baud int, portDefined bool, portProvided string) error {
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
            return util.Error("failed to automatically detect AVR port")
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
            return util.Error("invalid baud rate")
        }
    }

    defer serialPort.Close()

    log.Info(log.Cyan, "Wio Serial Monitor")
    log.Info(log.Yellow, "  @  ")
    log.Info(log.Cyan, portToUse)
    log.Info(log.Yellow, "  @  ")
    log.Infoln(log.Cyan, "%d", baud)
    log.Infoln(log.Cyan, "--- Quit: Ctrl+C ---")

    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        log.Infoln("\n--- exit ---")
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
    return nil
}
