package resolve

import (
    "wio/cmd/wio/constants"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/toolchain/npm/semver"
    "wio/cmd/wio/types"
)

const (
    Latest = "latest"
)

func (i *Info) GetLatest(name string) (string, error) {
    data, err := i.GetData(name)
    if err != nil {
        return "", err
    }
    if ver, exists := data.DistTags[Latest]; exists {
        return ver, nil
    }
    list, err := i.GetList(name)
    if err != nil {
        return "", err
    }
    return list.Last().Str(), nil
}

func (i *Info) Exists(name string, ver string) (bool, error) {
    if ret := semver.Parse(ver); ret == nil {
        return false, errors.Stringf("invalid version %s", ver)
    }
    data, err := i.GetData(name)
    if err != nil {
        return false, err
    }
    _, exists := data.Versions[ver]
    return exists, nil
}

func (i *Info) ResolveRemote(config types.IConfig) error {
    logResolveStart(config)

    if err := i.LoadLocal(); err != nil {
        return err
    }
    i.root = &Node{
        Name:            config.Name(),
        ConfigVersion:   config.Version(),
        ResolvedVersion: semver.Parse(config.Version()),
    }
    if i.root.ResolvedVersion == nil {
        return errors.Stringf("project has invalid version %s", i.root.ConfigVersion)
    }

    // adds pkg config for the initial package
    if config.GetType() == constants.PKG {
        i.SetPkg(i.root.Name, i.root.ResolvedVersion.Str(), &Package{
            Vendor: false,
            Path:   i.dir,
            Config: config.(*types.PkgConfig),
        })
    }

    deps := config.Dependencies()
    for name, ver := range deps {
        node := &Node{Name: name, ConfigVersion: ver}
        i.root.Dependencies = append(i.root.Dependencies, node)
    }
    for _, dep := range i.root.Dependencies {
        if err := i.ResolveTree(dep); err != nil {
            return err
        }
    }

    logResolveDone(i.root)
    return nil
}

func (i *Info) ResolveTree(root *Node) error {
    logResolve(root)

    if ret := i.GetRes(root.Name, root.ConfigVersion); ret != nil {
        root.ResolvedVersion = ret
        return nil
    }
    ver, err := i.resolveVer(root.Name, root.ConfigVersion)
    if err != nil {
        return err
    }
    root.ResolvedVersion = ver
    i.SetRes(root.Name, root.ConfigVersion, ver)
    data, err := i.GetVersion(root.Name, ver.Str())
    if err != nil {
        return err
    }
    for name, ver := range data.Dependencies {
        node := &Node{Name: name, ConfigVersion: ver}
        root.Dependencies = append(root.Dependencies, node)
    }
    for _, node := range root.Dependencies {
        if err := i.ResolveTree(node); err != nil {
            return err
        }
    }
    return nil
}

func (i *Info) resolveVer(name string, ver string) (*semver.Version, error) {
    if ret := semver.Parse(ver); ret != nil {
        i.StoreVer(name, ret)
        return ret, nil
    }
    query := semver.MakeQuery(ver)
    if query == nil {
        return nil, errors.Stringf("invalid version expression %s", ver)
    }
    if ret := i.resolve[name].Find(query); ret != nil {
        return ret, nil
    }
    list, err := i.GetList(name)
    if err != nil {
        return nil, err
    }
    if ret := query.FindBest(list); ret != nil {
        i.StoreVer(name, ret)
        return ret, nil
    }
    return nil, errors.Stringf("unable to find suitable version for %s", ver)
}
