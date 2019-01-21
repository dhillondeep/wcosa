package devices

import (
    "strings"

    serial "github.com/dhillondeep/go-serial"
)

func GetPorts() ([]*serial.Info, error) {
    ports, err := serial.ListPorts()
    if err != nil {
        return nil, err
    }

    return ports, nil
}

func GetArduinoPort(ports []*serial.Info) *serial.Info {
    for _, port := range ports {
        arduinoStr := "arduino"

        if strings.Contains(strings.ToLower(port.Description()), arduinoStr) ||
            strings.Contains(strings.ToLower(port.USBProduct()), arduinoStr) ||
            strings.Contains(strings.ToLower(port.USBManufacturer()), arduinoStr) {
            return port
        }
    }

    return nil
}
