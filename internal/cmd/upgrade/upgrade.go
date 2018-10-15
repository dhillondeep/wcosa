package upgrade

import (
    "wio/pkg/npm/semver"
    "wio/pkg/toolchain"
    "wio/pkg/util"

    "github.com/urfave/cli"
)

type Upgrade struct {
    Context *cli.Context
}

// get context for the command
func (upgrade Upgrade) GetContext() *cli.Context {
    return upgrade.Context
}

// Runs the upgrade command
func (upgrade Upgrade) Execute() error {
    if upgrade.Context.NArg() <= 0 {
        return toolchain.UpdateWioExecutable(nil)
    } else {
        versionToUpgrade := upgrade.Context.Args().Get(0)
        versionToUpgradeSem := semver.Parse(versionToUpgrade)

        if versionToUpgradeSem == nil {
            return util.Error("%s version is invalid", versionToUpgrade)
        }

        if versionToUpgradeSem.Lt(semver.Parse("0.6.0")) {
            return util.Error("wio can only be upgraded/downgraded to versions >= 0.6.0")
        }

        return toolchain.UpdateWioExecutable(versionToUpgradeSem)
    }
}
