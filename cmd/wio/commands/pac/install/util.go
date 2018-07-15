package install

import (
    "strings"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/toolchain/npm/resolve"
    "wio/cmd/wio/toolchain/npm/semver"
)

func (cmd Cmd) getArgs(info *resolve.Info) (name string, ver string, err error) {
    args := cmd.Context.Args()
    switch len(args) {
    case 0:
        err = errors.String("missing package name")

    case 1:
        if strings.Contains(args[0], "@") {
            args = strings.Split(args[0], "@")
            goto TwoArgs
        }
        name = args[0]
        ver, err = info.GetLatest(name)
        break

    TwoArgs:
        fallthrough
    default:
        name = args[0]
        ver = args[1]
        if semver.IsValid(ver) {
            exists := false
            exists, err = info.Exists(name, ver)
            if err == nil && !exists {
                err = errors.Stringf("version %s does not exist", ver)
            }
        } else if ret := semver.MakeQuery(ver); ret == nil {
            err = errors.Stringf("invalid version expression: %s", ver)
        }
    }
    return
}
