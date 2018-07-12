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
    unresolved map[string]mapset.Set
    versions   map[string]versionList
}

func buildDependencyTree(root *depTreeNode, info *depTreeInfo) error {
    query, err := getVersionQuery(root.version)
    if err != nil {
        return err
    }
    // base case reached with ~ or ^ version query
    if query != equal {
        if unresolved, exists := info.unresolved[root.name]; exists {
            unresolved.Add(root.version)
        } else {
            info.unresolved[root.name] = mapset.NewSet(root.version)
        }
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

func resolveVersionQueries(info *depTreeInfo) error {
    // collect available versions into sorted `versionList`
    for name, versions := range info.cache {
        list := make(versionList, 0, len(versions))
        for version := range versions {
            versionVal, _ := strtover(version) // err unlikely
            list = append(list, versionVal)
        }
        sort.Sort(list)
        info.versions[name] = list
    }
    // for `atLeast` queries only need to resolve highest
    // for `near` queries resolve once per unique major version
    for name, unresolved := range info.unresolved {
        atLeastVer, nearVers := trimVersionQueries(unresolved)
        resAtLeast, resNear, err := resolveQueries(name, atLeastVer, nearVers, info.versions[name])
        if err != nil {
            return err
        }
        newVers := make(versionList, 0, len(resNear) + 1)
        for _, nearVer := range resNear {
            if nearVer.major != resAtLeast.major {
                newVers = append(newVers, nearVer)
            }
        }
        newVers = append(newVers, resAtLeast)
        for _, newVer := range newVers {
            verStr := vertostr(newVer)
            pkgVer, err := getOrFetchVersion(name, verStr, info.baseDir)
            if err != nil {
                return err
            }
            if cacheName, exists := info.cache[name]; exists {
                cacheName[verStr] = pkgVer
            } else {
                info.cache[name] = map[string]*packageVersion{verStr: pkgVer}
            }
        }
    }
    return nil
}

func trimVersionQueries(unresolved mapset.Set) (version, mapset.Set) {
    atLeastVer := version{0, 0, 0}
    nearVers := mapset.NewSet()
    for ver := range unresolved.Iter() {
        verStr := ver.(string)
        verVal, _ := strtover(verStr[1:])   // err unlikely
        query, _ := getVersionQuery(verStr) // err unlikely
        switch query {
        case atLeast:
            if versionLess(atLeastVer, verVal) {
                atLeastVer = verVal
            }
        case near:
            nearVers.Add(verVal.major)
        }
    }
    return atLeastVer, nearVers
}

func resolveQueries(
    name string,
    atLeastVer version,
    nearVers mapset.Set,
    available versionList) (version, versionList, error) {

    didAtLeast := false
    var atLeast version

    majorRemain := make([]int, 0, nearVers.Cardinality())
    resolvedNear := make(versionList, 0, nearVers.Cardinality())

    if available != nil && len(available) > 0 {
        ver, err := available.findAtLeast(atLeastVer)
        if err == nil {
            didAtLeast = true
            atLeast = ver
        }

        for major := range nearVers.Iter() {
            ver, err := available.highestMajor(major.(int))
            if err == nil {
                resolvedNear = append(resolvedNear, ver)
            } else {
                majorRemain = append(majorRemain, major.(int))
            }
        }
    }

    // atLeast resolved and no leftover near
    if didAtLeast && len(majorRemain) == 0 {
        return atLeast, resolvedNear, nil
    }

    // need to query package data
    data, err := fetchPackageData(name)
    if err != nil {
        return version{}, nil, err
    }
    remoteVers := make(versionList, 0, len(data.Versions))
    for ver := range data.Versions {
        verVal, _ := strtover(ver) // err unlikely
        remoteVers = append(remoteVers, verVal)
    }
    sort.Sort(remoteVers)

    // resolve remaining
    if !didAtLeast {
        ver, err := remoteVers.findAtLeast(atLeastVer)
        if err != nil {
            return version{}, nil, err
        }
        atLeast = ver
    }
    if len(majorRemain) > 0 {
        for major := range majorRemain {
            ver, err := remoteVers.highestMajor(major)
            if err != nil {
                return version{}, nil, err
            }
            resolvedNear = append(resolvedNear, ver)
        }
    }
    return atLeast, resolvedNear, nil
}