package run

import (
    "wio/cmd/wio/toolchain"
    "wio/cmd/wio/errors"
)

func getPort(info *runInfo) (string, error) {
    if info.context.IsSet("port") {
        return info.context.String("port"), nil
    }
    ports, err := toolchain.GetPorts()
    if err != nil {
        return "", err
    }
    serialPort := toolchain.GetArduinoPort(ports)
    if serialPort == nil {
        return "", errors.String("Failed to find Arduino port")
    }
    return serialPort.Port, nil
}
