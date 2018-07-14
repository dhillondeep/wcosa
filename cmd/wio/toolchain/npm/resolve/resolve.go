package resolve

import (
    "wio/cmd/wio/errors"
    "wio/cmd/wio/toolchain/npm/semver"
)

func (i *Info) ResolveTree(root *Node) error {
	ver, err := i.resolveVer(root.name, root.ver)
	if err != nil {
		return err
	}
	root.resolve = ver
	data := i.GetVersion(root.name, ver.Str())
	for name, ver := data.Dependencies {
		node := &{name: name, ver: ver}
		root.deps = append(root.deps, node)
	}
	for _, node := range root.deps {
		if err := i.ResolveTree(node); err != nil {
			return nil
		}
	}
	return nil
}

func (i *Info) resolveVer(name string, ver string) (*Version, error) {
    if ret := semver.Parse(ver); ret != nil {
		i.StoreVer(name, ret)
        return ret
    }
    query := semver.MakeQuery(ver)
    if query == nil {
        return nil, errors.Stringf("invalid version expression %s", str)
    }
	if ret := i.resolve[name].Find(query); ret != nil {
		return ret, nil
	}
    if ret := query.FindBest(i.GetList(name)); ret != nil {
		i.StoreVer(name, ret)
		return ret, nil
	}
	return nil, errors.Stringf("unable to find suitable version for %s", str)
}
