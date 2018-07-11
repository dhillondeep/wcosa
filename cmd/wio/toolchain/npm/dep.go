package npm

type depTreeNode struct {
    name     string
    version  string
    children []*depTreeNode
}

type depTreeInfo struct {
    dataCache map[string]*packageData
    versionCache map[string]map[string]*packageVersion
}

func buildDependencyTree(root *depTreeNode, info *depTreeInfo) error {
    query, err := getVersionQuery(root.version)
    if err != nil {
        return err
    }
    // base case reached with ~ or ^ version query
    if query != equal {
        return nil
    }
    // get the version data
    var pkgVersion *packageVersion = nil
    if cachedName, exists := info.versionCache[root.name]; exists {
        if cachedVersion, exists := cachedName[root.version]; exists {
            pkgVersion = cachedVersion
        } else {
            cachedData := info.dataCache[root.name]
            pkgVersion, err = findPackage(cachedData, root.version)
            if err != nil {
                return err
            }
            cachedName[root.version] = pkgVersion
        }
    } else {
        pkgData, err := getPackageData(root.name)
        if err != nil {
            return err
        }
        info.dataCache[root.name] = pkgData
        pkgVersion, err = findPackage(pkgData, root.version)
        if err != nil {
            return err
        }
        info.versionCache[root.name] = map[string]*packageVersion{root.version: pkgVersion}
    }
    // get the dependencies of the hard version
    for depName, depVer := range pkgVersion.Dependencies {
        depNode := &depTreeNode{name: depName, version: depVer}
        root.children = append(root.children, depNode)
    }
    for _, depNode := range root.children {
        // TODO potentially parallel with goroutines
        if err := buildDependencyTree(depNode, info); err != nil {
            return err
        }
    }
    return nil
}
