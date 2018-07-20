package resolve

import (
	"wio/cmd/wio/utils/io"
    "wio/cmd/wio/toolchain/npm"
    "wio/cmd/wio/toolchain/npm/client"
    "wio/cmd/wio/toolchain/npm/semver"
    "wio/cmd/wio/types"
)

type DataCache map[string]*npm.Data
type VerCache map[string]map[string]*npm.Version
type ResCache map[string]map[string]*semver.Version
type PkgCache map[string]map[string]*Package
type ListMap map[string]semver.List

type Info struct {
    dir  string
    data DataCache
    ver  VerCache
    res  ResCache
	pkg  PkgCache

    resolve ListMap
    lists   ListMap

	root *Node
}

type Node struct {
    Name            string
    ConfigVersion   string
    ResolvedVersion *semver.Version
    Dependencies    []*Node
}

type Package struct {
	Vendor bool
	Path string
	Config *types.PkgConfig
	Version *npm.Version
}

func NewInfo(dir string) *Info {
    return &Info{
        dir:     dir,
        data:    DataCache{},
        ver:     VerCache{},
        res:     ResCache{},
		pkg:     PkgCache{},
        resolve: ListMap{},
        lists:   ListMap{},
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

func (i *Info) SetRes(name string, query string, ver *semver.Version) {
    if data, exists := i.res[name]; exists {
        data[query] = ver
    } else {
        i.res[name] = map[string]*semver.Version{query: ver}
    }
}

func (i *Info) GetRes(name string, query string) *semver.Version {
    if data, exists := i.res[name]; exists {
        if ret, exists := data[query]; exists {
            return ret
        }
    }
    return nil
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

func (i *Info) GetVersion(name, ver string) (*npm.Version, error) {
    if ret := i.getVer(name, ver); ret != nil {
        return ret, nil
    }
    if data := i.getData(name); data != nil {
        if ret, exists := data.Versions[ver]; exists {
            i.setVer(name, ver, &ret)
            return &ret, nil
        }
    }
    ret, err := i.GetLocalVersion(name, ver)
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
    data, err := i.GetData(name)
    if err != nil {
        return nil, err
    }
    vers := data.Versions
    list := make(semver.List, 0, len(vers))
    for ver := range vers {
        parse := semver.Parse(ver)
        if parse != nil {
            list = append(list, semver.Parse(ver))
        }
    }
    list.Sort()
    i.lists[name] = list
    return list, nil
}

func (i *Info) StoreVer(name string, ver *semver.Version) {
    i.resolve[name] = i.resolve[name].Insert(ver)
}

func (i *Info) GetLocalVersion(name, ver string) (*npm.Version, error) {
	pkg, err := i.GetPkg(name, ver)
	if err != nil {
		return nil, err
	}
	if pkg == nil {
		return nil, nil
	}
	return pkg.Version, nil
}

func (i *Info) GetPkg(name, ver string) (*Package, error) {
	if data, exists := i.pkg[name]; exists {
		if pkg, exists := data[ver]; exists {
			return pkg, nil
		}
	}

	vendor := []bool{true, true, false}
	strict := []bool{false, true, true}
	paths := []string{
		io.Path(i.dir, io.Vendor, name),
		io.Path(i.dir, io.Vendor, name+"__"+ver),
		io.Path(i.dir, io.Folder, io.Modules, name+"__"+ver),
	}
	for n, path := range paths {
		ret, err := tryFindConfig(name, ver, path, strict[n])
		if err != nil {
			return nil, err
		}
		if ret == nil {
			continue
		}
		pkg := &Package{Vendor: vendor[n], Path: path, Config: ret}
		pkg.Version = &npm.Version{
			Name: ret.Name(),
			Version: ret.Version(),
			Dependencies: ret.Dependencies(),
		}
		i.SetPkg(name, ver, pkg)
		return pkg, nil
	}
	return nil, nil
}

func (i *Info) SetPkg(name, ver string, pkg *Package) {
	if data, exists := i.pkg[name]; exists {
		data[ver] = pkg
	} else {
		i.pkg[name] = map[string]*Package{ver: pkg}
	}
}
