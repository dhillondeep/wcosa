package npm

import (
    "github.com/deckarep/golang-set"
    "sort"
)

type depTreeNode struct {
    name     string
    version  string
    children []*depTreeNode
}

type depTreeInfo struct {
    baseDir    string
    cache      map[string]map[string]*packageVersion
    data       map[string]*packageData
    unresolved map[string]mapset.Set
    versions   map[string]versionList

    resLeast map[string]string
    resNear  map[string]map[int]string
}

func newTreeInfo(dir string) *depTreeInfo {
    return &depTreeInfo{
        baseDir: dir,
        cache: map[string]map[string]*packageVersion{},
        data: map[string]*packageData{},
        unresolved: map[string]mapset.Set{},
        versions: map[string]versionList,
        resLeast: map[string]string{},
        resNear: map[string]map[int]string{},
    }
}

func buildDependencyTree(root *depTreeNode, info *depTreeInfo) error {
    query, err := getVersionQuery(root.version)
    if err != nil {
        return err
    }

    // version query resolve
    if query != equal {
    }

    // get the version data
    var pkgVersion *packageVersion = nil
    if cacheName, exists := info.cache[root.name]; exists {
        if _, exists := cacheName[root.version]; exists {
            return nil // already resolved
        }
    }
    if pkgVersion == nil {
        pkgVersion, err = getOrFetchVersion(root.name, root.version, info.baseDir)
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

func printTree(node *depTreeNode, level log.Type, pre string) {
    log.Writeln(level, "%s@%s", node.name, node.version)
    for i := 0; i < len(node.children) - 1; i++ {
        log.Write(level, "%s|_ ", pre)
        printTree(node.children[i], level, pre + "|  ")
    }
    if len(node.children) > 0 {
        log.Write(level, "%s\\_ ", pre)
        printTree(node.children[len(node.children) - 1], level, pre + "   ")
    }
}

/*
pkg
|_ wlib-json@1.0.4
|  \_ wlib-wio@1.0.0
\_ wlib-memory@1.0.2
|  |_ wlib-tmp@1.0.0
|  |  \_ wlib-util@1.0.0
|  \_ wlib-malloc@1.0.2
|     \_ wlib-tlsf@1.0.1
\_ wlib-list@1.0.0
*/
