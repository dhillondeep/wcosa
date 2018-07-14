package resolve

import (
    "wio/cmd/wio/toolchain/npm"
    "wio/cmd/wio/toolchain/npm/client"
)

type DataCache map[string]*npm.Data
type VerCache map[string]map[string]*npm.Version

type Info struct {
    dir  string
    data DataCache
    ver  VerCache

	resolve map[string]semver.List
	lists map[string]semver.List
}

type Node struct {
    name string
    ver  string
    deps []*Node

	resolve *semver.Version
}

func NewInfo(dir string) *Info {
    return &Info{
        dir:  dir,
        data: DataCache{},
        ver:  VerCache{},
    }
}

func (i *Info) getData(name string) *npm.Data {
    if ret, exists := i.data[name]; exists {
        return ret
    }
    return nil
}

func (i *Info) setData(name string, data *npm.Data) {
    i.data[name] = data
}

func (i *Info) getVer(name string, ver string) *npm.Version {
    if data, exists := i.ver[name]; exists {
        if ret, exists := data[ver]; exists {
            return ret
        }
    }
    return nil
}

func (i *Info) setVer(name string, ver string, data *npm.Version) {
    if cache, exists := i.ver[name]; exists {
        cache[ver] = data
    } else {
        i.ver[name] = map[string]*npm.Version{ver: data}
    }
}

func (i *Info) GetData(name string) (*npm.Data, error) {
    if ret := i.getData(name); ret != nil {
        return ret, nil
    }
    ret, err := client.FetchPackageData(name)
    if err != nil {
        return nil, err
    }
    i.setData(name, ret)
    return ret, nil
}

func (i *Info) GetVersion(name string, ver string) (*npm.Version, error) {
    if ret := i.getVer(name, ver); ret != nil {
        return ret, nil
    }
    if data := i.getData(name); data != nil {
        if ret, exists := data.Versions[ver]; exists {
            i.setVer(name, ver, ret)
            return ret
        }
    }
    ret, err := findVersion(name, ver, i.dir)
    if err != nil {
        return nil, err
    }
    if ret != nil {
        i.setVer(name, ver, ret)
        return ret, nil
    }
    ret, err = client.FetchPackageVersion(name, ver)
    if err != nil {
        return nil, err
    }
    i.setVer(name, ver, ret)
    return ret, nil
}

func (i *Info) GetList(name string) (semver.List, error) {
	if ret, exists := i.lists[name]; exists {
		return ret, nil
	}
	data, err := GetData(name)
	if err != nil {
		return err
	}
	vers := data.Versions
	list := make(semver.List, 0, len(vers))
	for ver := range vers {
		list = append(list, semver.Parse(ver))
	}
	list.Sort()
	i.lists[name] = list
	return list
}

func (i *Info) StoreVer(name string, ver *Version) {
	list := i.resolve[name]
	for _, el := range list {
		if el.eq(ver) {
			return
		}
	}
	list = append(list, ver)
	list.Sort()
}
