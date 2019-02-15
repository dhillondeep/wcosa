package generate

import (
	"wio/internal/cmd/devices"
	"wio/pkg/util"
)

func GetPort(providedPort string) (string, error) {
	if !util.IsEmptyString(providedPort) && providedPort != "none" {
		return providedPort, nil
	}
	ports, err := devices.GetPorts()
	if err != nil {
		return "", err
	}
	serialPort := devices.GetArduinoPort(ports)
	if serialPort == nil {
		return "", util.Error("failed to detect upload port. Specify an upload port!")
	}
	return serialPort.Name(), nil
}
