package resolve

import (
    "wio/cmd/wio/errors"
    "wio/cmd/wio/toolchain/npm/semver"
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

func (i *Info) ResolveTree(root *Node) error {
    if ret := i.GetRes(root.name, root.ver); ret != nil {
        root.resolve = ret
        return nil
    }
    ver, err := i.resolveVer(root.name, root.ver)
    if err != nil {
        return err
    }
    root.resolve = ver
    i.SetRes(root.name, root.ver, ver)
    data, err := i.GetVersion(root.name, ver.Str())
    if err != nil {
        return err
    }
    for name, ver := range data.Dependencies {
        node := &Node{name: name, ver: ver}
        root.deps = append(root.deps, node)
    }
    for _, node := range root.deps {
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
