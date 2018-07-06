package run

import (
    "strings"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/toolchain"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"
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
        return "", errors.String("failed to find Arduino port")
    }
    return serialPort.Port, nil
}

func portReconfigure(info *runInfo, target *types.Target) error {
    // Run check means that executable exists and target is configured
    port, err := getPort(info)
    if err != nil {
        return err
    }
    targetDir := targetPath(info, target)
    data, err := io.NormalIO.ReadFile(targetDir + io.Sep + "CMakeLists.txt")
    if err != nil {
        return err
    }
    if !strings.Contains(string(data), port) {
        _, err := configureTargets(info, []types.Target{*target})
        if err != nil {
            return err
        }
        return configTarget(binaryPath(info, target))
    }
    return nil
}
