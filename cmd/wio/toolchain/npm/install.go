package npm

import (
    "os"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/types"
    "wio/cmd/wio/log"
)

type versionQuery int

const (
    equal   versionQuery = 0
    atLeast versionQuery = 1
    near    versionQuery = 2
)

func getVersionQuery(versionStr string) (versionQuery, error) {
    if len(versionStr) < 1 {
        return 0, errors.Stringf("invalid version string: %s", versionStr)
    }
    leading := versionStr[0]
    switch leading {
    case '~':
        return near, nil
    case '^':
        return atLeast, nil
    default:
        return equal, nil
    }
}

// Removes `.wio.js` and `package.json` from extracted tarball
func removePackageExtras(pkgDir string) error {
    if err := os.Remove(io.Path(pkgDir, ".wio.js")); err != nil {
        return err
    }
    return os.Remove(io.Path(pkgDir, "package.json"))
}

// The generated folder structure will be
//
// `pkgDir`
//      [packageName]__[packageVersion]
//      [packageName]__[packageVersion]
//      ...
//      [packageName]__[packageVersion]
//          include
//          src
//          wio.yml
//
func installPackages(dir string, config types.IConfig) error {
    deps := config.GetDependencies()
    depNodes := make([]*depTreeNode, 0, len(deps))
    for name, depTag := range deps {
        depNode := &depTreeNode{name: name, version: depTag.Version}
        depNodes = append(depNodes, depNode)
    }
    root := &depTreeNode{
        name: config.Name(),
        version: config.Version(),
        children: depNodes,
    }
    info := newTreeInfo(dir)
    for _, depNode := range root.children {
        if err := buildDependencyTree(depNode, info); err != nil {
            return err
        }
    }
    for name, versions := range info.cache {
        for version := range versions {
            log.Infoln("%s@%s", name, version)
        }
    }

    return nil
}
