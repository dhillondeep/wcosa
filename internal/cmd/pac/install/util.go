package install

import (
	"strings"
	"wio/pkg/npm/resolve"
	"wio/pkg/npm/semver"
	"wio/pkg/util"
)

func (cmd Cmd) getArgs(info *resolve.Info, versionCheck bool) (name string, ver string, err error) {
	args := cmd.Context.Args()
	switch len(args) {
	case 0:
		err = util.Error("missing package name")

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
		if versionCheck {
			if semver.Parse(ver) != nil {
				exists := false
				exists, err = info.Exists(name, ver)
				if err == nil && !exists {
					err = util.Error("version %s does not exist", ver)
				}
			} else if ret := semver.MakeQuery(ver); ret == nil {
				err = util.Error("invalid version expression: %s", ver)
			}
		}
	}
	return
}
