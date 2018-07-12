package npm

type depTreeNode struct {
    name     string
    version  string
    children []*depTreeNode
}

type depTreeInfo struct {
    baseDir string
    cache   map[string]map[string]*packageVersion
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
    if cacheName, exists := info.cache[root.name]; exists {
        if _, exists := cacheName[root.version]; exists {
            return nil // already resolved
        }
    }
    if pkgVersion == nil {
        pkgVersion, err = getOrFetchVersion(root, info)
        if err != nil {
            return err
        }
        if cacheName, exists := info.cache[root.name]; exists {
            cacheName[root.version] = pkgVersion
        } else {
            info.cache[root.name] = map[string]*packageVersion{root.version: pkgVersion}
        }
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

