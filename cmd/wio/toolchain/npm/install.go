package npm

import (
    "os"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/utils/io"
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

func findPackage(pkgData *packageData, versionStr string) (*packageVersion, error) {
    name := pkgData.Name
    if len(pkgData.Versions) <= 0 {
        return nil, errors.Stringf("package %s found but no versions exist", name)
    }
    // check for dist tag version
    if distTag, exists := pkgData.DistTags[versionStr]; exists {
        versionStr = distTag
    }
    queryType, err := getVersionQuery(versionStr)
    if err != nil {
        return nil, err
    }
    if queryType == equal {
        if pkgVersion, exists := pkgData.Versions[versionStr]; exists {
            return &pkgVersion, nil
        }
        return nil, errors.Stringf("package %s@%s does not exist", name, versionStr)
    }
    versionStrList := make([]string, 0, len(pkgData.Versions))
    for versionKey := range pkgData.Versions {
        versionStrList = append(versionStrList, versionKey)
    }
    sortedVersions, err := sortedVersionList(versionStrList)
    queryVersion, err := strtover(versionStr[1:])
    if err != nil {
        return nil, err
    }
    var version version
    switch queryType {
    case atLeast:
        version, err = sortedVersions.findAtLeast(queryVersion)
        if err != nil {
            return nil, err
        }
    case near:
        version = sortedVersions.findNearest(queryVersion)
    }
    pkgVersion := pkgData.Versions[vertostr(version)]
    return &pkgVersion, nil
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
        depNode = &depTreeNode{name: name, version: depTag.Version}
        depNodes = append(depNodes, depNode)
    }
    root := &depTreeNode{
        name: config.Name(),
        version: config.Version(),
        children: depNodes,
    }
    info := newTreeInfo(dir)
    for depNode := range root.children {
        if err := buildDependencyTree(depNode, info, false); err != nil {
            return err
        }
    }
    for name, versions := range info.cache {
        for version := range versions {
            log.Infoln("%s@%s", name, version)
        }
    }
}
