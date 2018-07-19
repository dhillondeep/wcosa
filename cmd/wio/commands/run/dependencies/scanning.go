package dependencies

import (
    "io/ioutil"
    "strings"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/toolchain/npm/semver"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
)

type DependencyInfo struct {
    Name         string
    Directory    string
    Version      *semver.Version
    Vendor       bool
    MainTag      types.MainTag
    Dependencies types.DependenciesTag
}

// scans a directory and creates a DependencyInfo for each valid package and add that to given map
func scanDependencies(directory string, dependencies map[string]*DependencyInfo) error {
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
                return errors.Stringf("%s dependency named \"%s\" is an application not a package",
                    from, currDir.Name())
            }

            depVersion := semver.Parse(config.Version())
            if depVersion == nil {
                return errors.Stringf("%s dependency named \"%s\" does not have a valid version: %s",
                    from, currDir.Name(), config.Version())
            }

            dependencies[config.Name()+"__"+SemverVersionToString(depVersion)] = &DependencyInfo{
                Name:         config.Name(),
                Directory:    currDirPath,
                Version:      depVersion,
                Vendor:       !isRemoteDependency,
                MainTag:      config.GetMainTag(),
                Dependencies: config.GetDependencies(),
            }
        }
    }

    return nil
}

// This gathers dependency info for the whole project. Vendor packages override remote packages if they have
// the same name and version
func GatherDependencies(projectPath string, projectConfig types.IConfig) (map[string]*DependencyInfo, error) {
    dependencies := make(map[string]*DependencyInfo)

    // gather from remote folder
    err := scanDependencies(io.Path(projectPath, io.Folder, io.Modules), dependencies)
    if err != nil {
        return nil, err
    }

    // gather from vendor folder (will override if names are the same
    err = scanDependencies(io.Path(projectPath, io.Vendor), dependencies)
    if err != nil {
        return nil, err
    }

    // if package add that as an dependency
    if projectConfig.GetType() == constants.PKG {
        pkgVersion := semver.Parse(projectConfig.Version())
        if pkgVersion == nil {
            return nil, errors.Stringf("Project package does not have a valid version: %s",
                projectConfig.Version())
        }

        dependencies[projectConfig.Name()+"__"+SemverVersionToString(pkgVersion)] = &DependencyInfo{
            Name:         projectConfig.Name(),
            Directory:    projectPath,
            Version:      pkgVersion,
            Vendor:       false,
            MainTag:      projectConfig.GetMainTag(),
            Dependencies: projectConfig.GetDependencies(),
        }
    }

    return dependencies, nil
}
