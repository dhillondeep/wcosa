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
		if devices.Context.NArg() == 0 {
			return util.Error("please provide a serial port")
		}
		return HandleMonitor(devices.Context.Int("baud"), devices.Context.Args().Get(0))
	case LIST:
		return handlePorts()
	default:
		return util.Error("invalid device command")
	}
}

// Provides information abouts ports
func handlePorts() error {
	ports, err := GetPorts()
	if err != nil {
		return err
	}

	log.Info(log.Cyan, "Num of ports: ")
	log.Infoln("%d\n", len(ports))

	for _, port := range ports {
		log.Infoln(port)
	}

	return nil
}

// Opens monitor to see serial data
func HandleMonitor(baud int, portProvided string) error {

	// Open the first serial port detected at 9600bps N81
	mode := &serial.Mode{
		BaudRate: baud,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	serialPort, err := serial.Open(portProvided, mode)
	if err != nil {
		if strings.Contains(err.Error(), "Invalid serial port") {
			return util.Error("invalid baud rate")
		}
	}

	defer serialPort.Close()

	log.Info(log.Cyan, "Wio Serial Monitor")
	log.Info(log.Yellow, "  @  ")
	log.Info(log.Cyan, portProvided)
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

	// When invalid port is read, serial library panics and hence that panic needs to be caught
	defer func() {
		if recover() != nil {
			log.Errln("%s port is not valid or cannot be used", portProvided)
			os.Exit(1)
		}
	}()

	// Read and print the response
	buff := make([]byte, 100)
	for {
		// Reads up to 100 bytes
		n, err := serialPort.Read(buff)

		if err != nil {
			return err
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
