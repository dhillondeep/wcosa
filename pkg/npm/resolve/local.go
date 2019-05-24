package resolve

import (
	"path/filepath"
	"wio/internal/constants"
	"wio/internal/types"
	"wio/pkg/util"
	"wio/pkg/util/sys"
)

func findLocalConfigs(root string) ([]string, error) {
	paths := []string{
		sys.Path(root, sys.Vendor),
		sys.Path(root, sys.WioFolder, sys.Modules),
		sys.Path(root, sys.WioFolder, sys.Modules, sys.Custom),
	}
	var ret []string
	for _, path := range paths {
		matches, err := filepath.Glob(path)
		if err != nil {
			return nil, err
		}

		for _, match := range matches {
			status, err := util.IsDir(match)
			if err != nil {
				return nil, err
			}

			if !status {
				continue
			}
			dir := sys.Path(match)
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
