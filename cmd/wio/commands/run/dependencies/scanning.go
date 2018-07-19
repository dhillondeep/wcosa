package dependencies

import (
    "wio/cmd/wio/types"
    "io/ioutil"
    "strings"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/constants"
)

type DependencyInfo struct {
    Name         string
    Directory    string
    Version      string
    Vendor       bool
    MainTag      types.MainTag
    Dependencies types.DependenciesTag
}

func ScanDependencies(directory string, dependencies map[string]*DependencyInfo) error {
    if dirs, err := ioutil.ReadDir(directory); err != nil {
        return err
    } else if len(dirs) > 0 {
        for _, currDir := range dirs {
            // ignore files
            if currDir.Mode().IsRegular() {
                continue
            }

            currDirPath := io.Path(directory, currDir.Name())
            isRemoteDependency := strings.Contains(directory, io.Path(io.Folder, io.Modules))
            from := "vendor"
            if isRemoteDependency {
                from = "remote"
            }

            if !utils.PathExists(io.Path(currDirPath, io.Config)) {
                return errors.Stringf("%s dependency named \"%s\" does not contain a wio.yml file",
                    from, currDir.Name())
            }
            config, err := utils.ReadWioConfig(currDirPath)
            if err != nil {
                return err
            }

            // confirm if config is a package type
            if config.GetType() == constants.APP {
                errors.Stringf("%s dependency named \"%s\" is an application not a package",
                    from, currDir.Name())
            }

            dependencies[config.Name() + "__" + config.Version()] = &DependencyInfo{
                Name: config.Name(),
                Version: config.Version(),
                Vendor: !isRemoteDependency,
                MainTag: config.GetMainTag(),
                Dependencies: config.GetDependencies(),
            }
        }
    }

    return nil
}
