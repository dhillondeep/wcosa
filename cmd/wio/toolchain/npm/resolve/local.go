package resolve

import (
    "io/ioutil"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
)

func findLocalConfigs(root string) ([]string, error) {
    paths := []string{
        io.Path(root, io.Vendor),
        io.Path(root, io.Folder, io.Modules),
    }
    var ret []string
    for _, path := range paths {
        if !io.Exists(path) {
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
            dir := io.Path(path, info.Name())
            if io.Exists(io.Path(dir, io.Config)) {
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
        return nil, errors.Stringf("config %s has wrong name", path)
    }
    if config.GetVersion() != ver {
        if strict {
            return nil, errors.Stringf("config %s has wrong version", path)
        } else {
            return nil, nil
        }
    }
    return config, nil
}

func tryGetConfig(path string) (types.Config, error) {
    wioPath := io.Path(path, io.Config)
    if !io.Exists(wioPath) {
        return nil, nil
    }
    config, err := utils.ReadWioConfig(path)
    if err != nil {
        return nil, err
    }
    if config.GetType() == constants.APP {
        return nil, errors.Stringf("config %s is supposed to be package")
    }
    return config, nil
}
