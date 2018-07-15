package install

import (
    "wio/cmd/wio/errors"
    "wio/cmd/wio/toolchain/npm/resolve"
    "wio/cmd/wio/toolchain/npm/semver"
)

func (cmd Cmd) getArgs(info *resolve.Info) (string, string, error) {
    args := cmd.Context.Args()
    switch len(args) {
    case 0:
        return "", "", errors.String("missing package name")
    case 1:
        name := args[0]
        ver, err := info.GetLatest(name)
        if err != nil {
            return "", "", err
        }
        return name, ver, nil
    default:
        name := args[0]
        ver := args[1]
        if ret := semver.MakeQuery(ver); ret == nil {
            return "", "", errors.Stringf("invalid version expression: %s", ver)
        }
        return name, ver, nil
    }
}
