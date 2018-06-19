package toolchain

import (
    "bytes"
    "os"
    "encoding/json"
    "strings"
)

type SerialPort struct {
    Port         string
    Description  string
    Hwid         string
    Manufacturer string
    SerialNumber string `json:"serial-number"`
    Vid          string
    Product      string
}

type SerialPorts struct {
    Ports []SerialPort
}

func GetPorts() (*SerialPorts, error) {
    cmd, err := GetPySerialCommand("-get-serial-devices")
    if err != nil {
        return nil, err
    }

    cmdOutput := &bytes.Buffer{}
    cmd.Stdout = cmdOutput
    cmd.Stderr = os.Stderr
    cmd.Run()

    ports := &SerialPorts{}
    if err := json.Unmarshal([]byte(cmdOutput.String()), ports); err != nil {
        return nil, err
    }

    return ports, nil
}

func GetArduinoPort(ports *SerialPorts) *SerialPort {
    for _, port := range ports.Ports {
        arduinoStr := "arduino"

        if strings.Contains(strings.ToLower(port.Description), arduinoStr) ||
            strings.Contains(strings.ToLower(port.Product), arduinoStr) ||
            strings.Contains(strings.ToLower(port.Manufacturer), arduinoStr) {
                return &port
        }
    }

    return nil
}
