package devices

import (
	"go.bug.st/serial.v1"
)

func GetPorts() ([]string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, err
	}

	return ports, nil
}
