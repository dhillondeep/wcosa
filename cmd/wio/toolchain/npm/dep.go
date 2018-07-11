package npm

type depTreeNode struct {
    name     string
    version  string
    children []*depTreeNode
}

type depTreeInfo struct {
    cache map[string]map[string]*packageVersion
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
    // get the dependencies of the hard version
    var pkgVersion *packageVersion = nil
    if cacheName, exists := info.cache[root.name]; exists {
        if cacheVersion, exists := cacheName[root.version]; exists {
            pkgVersion = cacheVersion
        }
    }
    if pkgVersion == nil {
        pkgVersion, err = findPackage(root.name, root.version)
        if err != nil {
            return err
        }
        if cacheName, exists := info.cache[root.name]; exists {
            cacheName[root.version] = pkgVersion
        } else {
            info.cache[root.name] = map[string]*packageVersion{root.version: pkgVersion}
        }
    }

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
