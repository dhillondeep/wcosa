package resolve

import (
	"io/ioutil"
	"wio/internal/constants"
	"wio/internal/types"
	"wio/pkg/util"
	"wio/pkg/util/sys"
)

func findLocalConfigs(root string) ([]string, error) {
	paths := []string{
		sys.Path(root, sys.Vendor),
		sys.Path(root, sys.WioFolder, sys.Modules),
	}
	var ret []string
	for _, path := range paths {
		if !sys.Exists(path) {
			continue
		}
		infos, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, info := range infos {
			if !info.IsDir() {
				continue
			}
			dir := sys.Path(path, info.Name())
			if sys.Exists(sys.Path(dir, sys.Config)) {
				ret = append(ret, dir)
			}
		}
	}
	return ret, nil
}

func tryFindConfig(name, ver, path string, strict bool) (types.Config, error) {
	config, err := tryGetConfig(path)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}
	if config.GetName() != name {
		return nil, util.Error("config %s has wrong name", path)
	}
	if config.GetVersion() != ver {
		if strict {
			return nil, util.Error("config %s has wrong version", path)
		} else {
			return nil, nil
		}
	}
	return config, nil
}

func tryGetConfig(path string) (types.Config, error) {
	wioPath := sys.Path(path, sys.Config)
	if !sys.Exists(wioPath) {
		return nil, nil
	}
	config, err := types.ReadWioConfig(path, true)
	if err != nil {
		return nil, err
	}
	if config.GetType() == constants.App {
		return nil, util.Error("config %s is supposed to be package", config.GetName())
	}
	return config, nil
}
