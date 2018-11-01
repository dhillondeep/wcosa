package run

import (
    "strings"
    "wio/internal/cmd/devices"
    "wio/internal/types"
    "wio/pkg/util"
    "wio/pkg/util/sys"
)

func getPort(info *runInfo) (string, error) {
    if info.context.IsSet("port") {
        return info.context.String("port"), nil
    }
    ports, err := devices.GetPorts()
    if err != nil {
        return "", err
    }
    serialPort := devices.GetArduinoPort(ports)
    if serialPort == nil {
        return "", util.Error("failed to find Arduino port")
    }
    return serialPort.Name(), nil
}

func portReconfigure(info *runInfo, target types.Target) error {
    // Run check means that executable exists and target is configured
    port, err := getPort(info)
    if err != nil {
        return err
    }
    targetDir := targetPath(info, target)
    data, err := sys.NormalIO.ReadFile(sys.Path(targetDir, "CMakeLists.txt"))
    if err != nil {
        return err
    }
    if !strings.Contains(string(data), port) {
        _, err := configureTargets(info, []types.Target{target})
        if err != nil {
            return err
        }
        return configTarget(binaryPath(info, target))
    }
    return nil
}
