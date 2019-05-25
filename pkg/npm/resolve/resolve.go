package resolve

import (
	"fmt"
	"github.com/hashicorp/go-getter"
	"os"
	"strings"
	"wio/internal/constants"
	"wio/internal/types"
	"wio/pkg/log"
	"wio/pkg/npm/semver"
	"wio/pkg/util"
	"wio/pkg/util/sys"

	s "github.com/blang/semver"
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
	return list.Last().String(), nil
}

func (i *Info) Exists(name string, ver string) (bool, error) {
	if ret := semver.Parse(ver); ret == nil {
		return false, util.Error("invalid version %s", ver)
	}
	data, err := i.GetData(name)
	if err != nil {
		return false, err
	}
	_, exists := data.Versions[ver]
	return exists, nil
}

func (i *Info) customUrlResolve(name string, dep types.Dependency, customPath string, root *Node, install bool) error {
	givenVer := semver.Parse(dep.GetVersion())
	if givenVer == nil {
		return util.Error("%s dependency version %s specified is not valid", name, dep.GetVersion())
	}

	dst := sys.Path(customPath, name) + "__" + givenVer.String()
	subDir := "/"

	if install && !sys.Exists(dst) {
		// add options to the source url
		url := dep.GetUrl().GetName()

		if !util.IsEmptyString(dep.GetUrl().GetDir()) {
			subDir = dep.GetUrl().GetDir()
			url += "//" + subDir
		}

		options := dep.GetUrl().GetOptions()
		if options != nil {
			url += "?"
			for name, value := range options {
				url += name + "=" + value + "&"
			}
			url = strings.Trim(strings.Trim(url, "&"), "?")
		}

		log.Info(log.Cyan, "Installing ")
		log.Info(log.Green, "%s@%s ", name, givenVer.String())
		log.Info(log.Cyan, "from ")
		log.Info(log.Green, "%s [%s]", dep.GetUrl().GetName(), subDir)
		log.Info(log.Cyan, "... ")

		if err := getter.Get(dst, url, func(client *getter.Client) error {
			client.Pwd = i.dir
			return nil
		}); err != nil {
			log.WriteFailure()
			return err
		} else {
			log.WriteSuccess()
		}
	}

	node := &Node{Name: name, ConfigVersion: dep.GetVersion(), Vendor: false, CustomUrl: true}
	root.Dependencies = append(root.Dependencies, node)

	if sys.Exists(dst) {
		// read wio.yml config
		config, err := types.ReadWioConfig(dst, true)
		if err != nil {
			return err
		}

		downloadVer := semver.Parse(config.GetVersion())
		if downloadVer == nil {
			return util.Error("%s dependency cannot have invalid version: %s", name, config.GetVersion())
		}

		if downloadVer.NE(*givenVer) {
			err = os.RemoveAll(dst)
			return util.Error("%s version mismatch between specified and downloaded: %s != %s",
				name, givenVer.String(), downloadVer.String())
		}

		i.SetPkg(name, config.GetVersion(), &Package{
			Vendor:  false,
			Path:    dst,
			Config:  config,
			Version: nil,
		})

		return i.resolveRemote(config, node, install)
	}

	return nil
}

func (i *Info) createNodesAndFetch(deps map[string]types.Dependency, root *Node, install bool) error {
	customPath := sys.Path(i.dir, sys.WioFolder, sys.Modules, sys.Custom)
	if err := os.MkdirAll(customPath, os.ModePerm); err != nil {
		return err
	}

	for name, dep := range deps {
		if dep == nil {
			return util.Error("%s dependency cannot be empty", name)
		}

		// custom url is not provided
		if dep.GetUrl() == nil {
			node := &Node{Name: name, ConfigVersion: dep.GetVersion(), Vendor: dep.IsVendor()}
			root.Dependencies = append(root.Dependencies, node)
		} else {
			if err := i.customUrlResolve(name, dep, customPath, root, install); err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *Info) resolveRemote(config types.Config, root *Node, install bool) error {
	if err := i.createNodesAndFetch(config.GetDependencies(), root, install); err != nil {
		return err
	}

	for _, dep := range root.Dependencies {
		if !dep.CustomUrl {
			if err := i.ResolveTree(dep, install); err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *Info) ResolveRemote(config types.Config, install bool) error {
	logResolveStart(config)

	if err := i.LoadLocal(); err != nil {
		return err
	}
	i.root = &Node{
		Name:            config.GetName(),
		ConfigVersion:   config.GetVersion(),
		ResolvedVersion: semver.Parse(config.GetVersion()),
	}
	if i.root.ResolvedVersion == nil {
		return util.Error("project has invalid version %s", i.root.ConfigVersion)
	}

	// adds pkg config for the initial package
	if config.GetType() == constants.Pkg {
		i.SetPkg(i.root.Name, i.root.ResolvedVersion.String(), &Package{
			Vendor: false,
			Path:   i.dir,
			Config: config,
		})
	}

	if err := i.resolveRemote(config, i.root, install); err != nil {
		return err
	}

	logResolveDone(i.root)
	return nil
}

func (i *Info) ResolveTree(root *Node, install bool) error {
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

	if root.Vendor {
		pkg, err := i.GetPkg(root.Name, ver.String())
		if err != nil {
			return err
		}
		if err := i.createNodesAndFetch(pkg.Config.GetDependencies(), root, install); err != nil {
			return err
		}
	} else {
		data, err := i.GetVersion(root.Name, ver.String())
		if err != nil {
			return err
		}
		for name, ver := range data.Dependencies {
			node := &Node{Name: name, ConfigVersion: ver, Vendor: false}
			root.Dependencies = append(root.Dependencies, node)
		}
	}

	for _, node := range root.Dependencies {
		if !node.CustomUrl {
			if err := i.ResolveTree(node, install); err != nil {
				return err
			}
		}
	}
	return nil
}

func (i *Info) resolveVer(name string, ver string) (*s.Version, error) {
	if ret := semver.Parse(ver); ret != nil {
		i.StoreVer(name, ret)
		return ret, nil
	}

	query := semver.MakeQuery(ver)
	if query == nil {
		return nil, util.Error("invalid version expression %s", ver)
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
	} else {
		fmt.Println("")
	}
	return nil, util.Error("unable to find suitable version for %s", ver)
}
