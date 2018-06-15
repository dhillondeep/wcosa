// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands package, which contains all the commands provided by the tool.
// Runs the serial monitor
package monitor

import (
    "github.com/urfave/cli"
    "wio/cmd/wio/toolchain"
    "bytes"
    "os"
    "wio/cmd/wio/log"
    "encoding/json"
    "strconv"
    "wio/cmd/wio/commands"
    "fmt"
    "go.bug.st/serial.v1"
    "strings"
    "github.com/go-errors/errors"
    "os/signal"
    "syscall"
)

type Monitor struct {
    Context *cli.Context
    Type byte
    error
}

// get context for the command
func (monitor Monitor) GetContext() *cli.Context {
    return monitor.Context
}

type SerialPort struct {
    Port string
    Description string
    Hwid string
    Manufacturer string
    SerialNumber string     `json:"serial-number"`
    Vid string
    Product string
}

type SerialPorts struct {
    Ports []SerialPort
}

const (
    OPEN = 0
    PORTS = 1
)

// Runs the build command when cli build option is provided
func (monitor Monitor) Execute() {
    switch monitor.Type {
    case OPEN:
        HandleMonitor(monitor.Context.Int("baud"), monitor.Context.IsSet("port"), monitor.Context.String("port"))
        break
    case PORTS:
        handlePorts(monitor.Context.Bool("basic"), monitor.Context.Bool("show-all"))
        break
    }
}

// Provides information abouts ports
func handlePorts(basic bool, showAll bool) {
    cmd, err := toolchain.GetPySerialCommand("-get-serial-devices")
    if err != nil {
        commands.RecordError(err, "")
    }

    cmdOutput := &bytes.Buffer{}
    cmd.Stdout = cmdOutput
    cmd.Stderr = os.Stderr
    cmd.Run()

    ports := &SerialPorts{}
    if err := json.Unmarshal([]byte(cmdOutput.String()), ports); err != nil {
        commands.RecordError(err, "")
    }

    log.Norm.Write(true, "Num of total ports: " + strconv.Itoa(len(ports.Ports)))
    log.Norm.Write(true, "")

    numOpenPorts := 0
    for _, port := range ports.Ports {
        if port.Product == "None" && !showAll {
            continue
        } else {
            numOpenPorts++
        }

        log.Norm.Cyan(true, port.Port)

        if !basic {
            log.Norm.Write(true, "Product:          "+port.Product)
            log.Norm.Write(true, "Description:      "+port.Description)
            log.Norm.Write(true, "Manufacturer:     "+port.Manufacturer)
            log.Norm.Write(true, "Serial Number:    "+port.SerialNumber)
            log.Norm.Write(true, "Hwid:             "+port.Hwid)
            log.Norm.Write(true, "Vid:              "+port.Vid)


            log.Norm.Write(true, "")
        }
    }

    if basic {
        log.Norm.Write(true, "")
    }

    log.Norm.Write(true, "Number of open ports: " + strconv.Itoa(numOpenPorts))
}

// Opens monitor to see serial data
func HandleMonitor(baud int, portDefined bool, port string) {
    cmd, err := toolchain.GetPySerialCommand("-get-serial-devices")
    if err != nil {
        commands.RecordError(err, "")
    }

    cmdOutput := &bytes.Buffer{}
    cmd.Stdout = cmdOutput
    cmd.Stderr = os.Stderr
    cmd.Run()

    ports := &SerialPorts{}
    if err := json.Unmarshal([]byte(cmdOutput.String()), ports); err != nil {
        commands.RecordError(err, "")
    }

    var portToUse string

    if portDefined {
        log.Norm.Cyan(false, "Port provided: ")
        log.Norm.Write(true, port)
        portToUse = port
    } else {
        portFound := false
        for _, port := range ports.Ports {
            arduinoStr := "arduino"

            if strings.Contains(strings.ToLower(port.Description), arduinoStr) ||
                strings.Contains(strings.ToLower(port.Product), arduinoStr) ||
                strings.Contains(strings.ToLower(port.Manufacturer), arduinoStr) {
                log.Norm.Cyan(false, "Auto detected port: ")
                log.Norm.Write(true,  port.Port)
                portFound = true
                portToUse = port.Port
            }
        }

        if !portFound {
            commands.RecordError(errors.New("Port cannot be automatically detected. Please specify a port"), "")
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
        commands.RecordError(err, "")
    }

    defer serialPort.Close()

    log.Norm.Cyan(false, "Wio Serial Monitor")
    log.Norm.Yellow(false, "  @  " )
    log.Norm.Write(false, portToUse )
    log.Norm.Yellow(false, "  @  ")
    log.Norm.Write(true,  strconv.Itoa(baud))
    log.Norm.Write(true, "--- Quit: Ctrl+C ---")


    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        log.Norm.Write(true, "\n--- exit ---")
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
