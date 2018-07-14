package resolve

import "wio/cmd/wio/toolchain/npm"

type DataCache map[string]*npm.Data
type VerCache map[string]map[string]*npm.Version

func (c DataCache) TryGet(name string) *npm.Data {
    if data, exists := c[name]; exists {
        return data
    }
    return nil
}

func (c DataCache) Has(name string) bool {
    _, exists := c[name]
    return exists
}

func (c DataCache) Store(name string, data *npm.Data) bool {
    if _, exists := c[name]; !exists {
        c[name] = data
        return true
    }
    return false
}

func (c VerCache) TryGet(name string, ver string) *npm.Version {
    if data, exists := c[name]; exists {
        if version, exists := data[ver]; exists {
            return version
        }
    }
    return nil
}

func (c VerCache) Has(name string, ver string) bool {
    if data, exists := c[name]; exists {
        _, exists = data[ver]
        return exists
    }
    return false
}

func (c VerCache) Store(name string, ver string, version *npm.Version) bool {
    if data, exists := c[name]; exists {
        if _, exists := data[ver]; !exists {
            data[ver] = version
            return true
        }
        return false
    }
    c[name] = map[string]*npm.Version{ver: version}
    return true
}
