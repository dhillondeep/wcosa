package resolve

import (
    "strings"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/log"
    "wio/cmd/wio/toolchain/npm/semver"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
)

const (
    Latest = "latest"
)

// Flattens vendor dependency graph and creates a map of remote dependency name and version.
// An entry would like: "packageName__packageVersion: packageVersion"
func FlattenVendorRemoteDependencies(queue *log.Queue, directory string, vendorRemotes map[string]string) error {
    // parent config
    parentConfig, err := utils.ReadWioConfig(directory)
    if err != nil {
        return err
    }

    for parentDependencyName, parentDependency := range parentConfig.GetDependencies() {
        if parentDependency.Vendor {
            currDirectory := io.Path(directory, io.Vendor, parentDependencyName)

            if !utils.PathExists(currDirectory) {
                return errors.Stringf("vendor dependency [%s] from [%s] does not exist",
                    parentDependencyName, parentConfig.Name())
            }

            dependencyConfig, err := utils.ReadWioConfig(currDirectory)
            if err != nil {
                return err
            }

            if dependencyConfig.GetType() == constants.APP {
                log.Warnln(queue, "vendor dependency [%s] is supposed to be a package. Skipping....",
                    parentDependencyName)
                continue
            }

            // go through it's dependencies to get all the remote ones
            for depName, dep := range dependencyConfig.GetDependencies() {
                if dep.Vendor {
                    if err := FlattenVendorRemoteDependencies(queue, currDirectory, vendorRemotes); err != nil {
                        return err
                    }
                } else {
                    vendorRemotes[depName+"__"+dep.Version] = dep.Version
                }
            }
        }
    }

    return nil
}

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

    root := &Node{name: config.Name(), ver: config.Version()}
    if root.resolve = semver.Parse(root.ver); root.resolve == nil {
        return errors.Stringf("project has an invalid version %s", root.ver)
    }

    // remotes from main config file
    deps := config.Dependencies()
    for name, ver := range deps {
        name = strings.Split(name, "__")[0]

        node := &Node{name: name, ver: ver}
        root.deps = append(root.deps, node)
    }

    // remotes from all the vendors
    queue := &log.Queue{}
    vendorRemotes := map[string]string{}
    if err := FlattenVendorRemoteDependencies(queue, i.dir, vendorRemotes); err != nil {
        return err
    }
    for name, ver := range vendorRemotes {
        name = strings.Split(name, "__")[0]

        node := &Node{name: name, ver: ver}
        root.deps = append(root.deps, node)
    }
    log.Writeln(queue)

    // resolve remotes
    for _, dep := range root.deps {
        if err := i.ResolveTree(dep); err != nil {
            return err
        }
    }

    logResolveDone(root)
    return nil
}

func (i *Info) ResolveTree(root *Node) error {
    logResolve(root)

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
