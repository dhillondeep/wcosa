package dependencies

import (
    goerr "errors"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
    "wio/cmd/wio/commands/run/cmake"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/log"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/constants"
)

const (
    REMOTE_NAME     = "node_modules"
    VENDOR_NAME     = "vendor"
    PKG_REMOTE_NAME = "pkg_module"
)

var packageVersions = map[string]string{}     // Keeps track of versions for the packages
var cmakeTargets = map[string]*cmake.Target{} // CMake Target that will be built
var cmakeTargetsLink []cmake.TargetLink       // CMake Target to Link to and from
var cmakeTargetNames = map[string]bool{}      // CMake Target Names. Used to check for unique names

// Stores information about every package that is scanned
type DependencyScanStructure struct {
    Name         string
    Directory    string
    Version      string
    FromVendor   bool
    MainTag      types.PkgTag
    Dependencies types.DependenciesTag
}

// Creates Scan structures for all the scanned packages
func createDependencyScanStructure(name string, depPath string, fromVendor bool) (*DependencyScanStructure, types.DependenciesTag, error) {
    configPath := depPath + io.Sep + "wio.yml"

    config := types.PkgConfig{}
    if err := io.NormalIO.ParseYml(configPath, &config); err != nil {
        return nil, nil, err
    } else {
        if strings.Trim(config.MainTag.Config.WioVersion, " ") == "" {
            wioVer := config.MainTag.Config.WioVersion
            return nil, nil, errors.UnsupportedWioConfigVersion{PackageName: name, Version: wioVer}
        }
        pkg := &DependencyScanStructure{
            Directory:    depPath,
            Name:         config.GetMainTag().GetName(),
            FromVendor:   fromVendor,
            Version:      config.GetMainTag().GetVersion(),
            MainTag:      config.MainTag,
            Dependencies: config.GetDependencies(),
        }

        packageVersions[pkg.Name] = pkg.Version
        return pkg, config.DependenciesTag, nil
    }
}

// Go through all the dependency packages and get information about them
func recursiveDependencyScan(queue *log.Queue, currDirectory string,
    dependencies map[string]*DependencyScanStructure, packageDependencies types.DependenciesTag) error {
    // if directory does not exist, do not do anything
    if !utils.PathExists(currDirectory) {
        log.Verbln(queue, "% does not exist, skipping", currDirectory)
        return nil
    }

    // list all directories
    if dirs, err := ioutil.ReadDir(currDirectory); err != nil {
        return err
    } else if len(dirs) > 0 {
        // directories exist so let's go through each of them
        for _, dir := range dirs {
            // ignore files
            if !dir.IsDir() {
                continue
            }

            // only worry about dependencies in wio.yml file
            if _, exists := packageDependencies[dir.Name()]; !exists {
                continue
            }

            dirPath := currDirectory + io.Sep + dir.Name()
            if !utils.PathExists(dirPath + io.Sep + "wio.yml") {
                return goerr.New("wio.yml file missing")
            }

            var fromVendor = false

            // check if the current directory is for remote or vendor
            if filepath.Base(currDirectory) == VENDOR_NAME {
                fromVendor = true
            }

            // create DependencyScanStructure
            if dependencyPackage, packageDeps, err := createDependencyScanStructure(dir.Name(), dirPath, fromVendor); err != nil {
                return err
            } else {
                dependencyName := dependencyPackage.Name + "__vendor"
                if !dependencyPackage.FromVendor {
                    dependencyName = dependencyPackage.Name + "__" + dependencyPackage.Version
                }

                dependencies[dependencyName] = dependencyPackage

                log.Verbln(queue, "%s package stored as dependency named: %s", dirPath, dependencyName)

                packageDependencies = packageDeps
            }

            if utils.PathExists(dirPath + io.Sep + REMOTE_NAME) {
                // if remote directory exists
                if err := recursiveDependencyScan(queue, dirPath+io.Sep+REMOTE_NAME, dependencies, packageDependencies); err != nil {
                    return err
                }
            } else if utils.PathExists(dirPath + io.Sep + VENDOR_NAME) {
                // if vendor directory exists
                if err := recursiveDependencyScan(queue, dirPath+io.Sep+VENDOR_NAME, dependencies, packageDependencies); err != nil {
                    return err
                }
            }
        }
    } else {
        log.Verbln(queue, "% does not have any package, skipping", currDirectory)
    }

    return nil
}

// When we are building for pkg type, we will copy the files into the remote directory
// This will be picked up while scanning and hence the rest of build process stays the same
func convertPkgToDependency(pkgPath string, projectName string, projectDir string) error {
    if !utils.PathExists(pkgPath) {
        if err := os.MkdirAll(pkgPath, os.ModePerm); err != nil {
            return err
        }
    }

    pkgDir := pkgPath + io.Sep + projectName
    if utils.PathExists(pkgDir) {
        if err := os.RemoveAll(pkgDir); err != nil {
            return err
        }
    }

    // Copy project files
    projectDir += io.Sep
    pkgDir += io.Sep
    for _, file := range []string{"src", "include", "wio.yml"} {
        if err := utils.Copy(projectDir+file, pkgDir+file); err != nil {
            return err
        }
    }
    return nil
}

func CreateCMakeDependencyTargets(
    config types.IConfig,
    target *types.Target,
    projectPath string,
    queue *log.Queue) error {

    projectName := config.GetMainTag().GetName()
    projectType := config.GetType()
    projectDependencies := config.GetDependencies()
    pkgVersion := config.GetMainTag().GetVersion()
    projectFlags := (*target).GetFlags()
    projectDefinitions := (*target).GetDefinitions()

    remotePackagesPath := projectPath + io.Sep + ".wio" + io.Sep + REMOTE_NAME
    vendorPackagesPath := projectPath + io.Sep + VENDOR_NAME

    scannedDependencies := map[string]*DependencyScanStructure{}

    if projectType == constants.PKG {
        if projectDependencies == nil {
            projectDependencies = types.DependenciesTag{}
        }

        // convert pkg project into tests dependency
        projectDependencies[projectName] = &types.DependencyTag{
            Version:        pkgVersion,
            Flags:          projectFlags.GetPkgFlags(),
            Definitions:    projectDefinitions.GetPkgDefinitions(),
            LinkVisibility: "PRIVATE",
        }

        packageDependencyPath := projectPath + io.Sep + ".wio" + io.Sep + PKG_REMOTE_NAME

        log.Verb(queue, "converting project to a dependency to be used in tests ... ")

        if err := convertPkgToDependency(packageDependencyPath, projectName, projectPath); err != nil {
            log.WriteFailure(queue, log.VERB)
            return err
        }
        log.WriteSuccess(queue, log.VERB)

        log.Verb(queue, "recursively scanning package dependency at path "+ "%s ...", packageDependencyPath)
        subQueue := log.GetQueue()
        if err := recursiveDependencyScan(subQueue, packageDependencyPath, scannedDependencies, projectDependencies); err != nil {
            log.WriteFailure(queue, log.VERB)
            log.CopyQueue(subQueue, queue, log.TWO_SPACES)
            return err
        } else {
            log.WriteSuccess(queue, log.VERB)
            log.CopyQueue(subQueue, queue, log.TWO_SPACES)
        }
    }

    log.Verb(queue, log.VERB, nil, "recursively scanning remote dependencies at path: %s ... ", remotePackagesPath)
    subQueue := log.GetQueue()

    if err := recursiveDependencyScan(subQueue, remotePackagesPath, scannedDependencies, projectDependencies); err != nil {
        log.WriteFailure(queue, log.VERB)
        log.CopyQueue(subQueue, queue, log.TWO_SPACES)
        return err
    }
    log.WriteSuccess(queue, log.VERB)
    log.CopyQueue(subQueue, queue, log.TWO_SPACES)

    log.Verb(queue, "recursively scanning vendor dependencies at path: %s ... ", vendorPackagesPath)
    subQueue = log.GetQueue()

    if err := recursiveDependencyScan(subQueue, vendorPackagesPath, scannedDependencies, projectDependencies); err != nil {
        log.WriteFailure(queue, log.VERB)
        log.CopyQueue(subQueue, queue, log.TWO_SPACES)
        return err
    }
    log.WriteSuccess(queue, log.VERB)
    log.CopyQueue(subQueue, queue, log.TWO_SPACES)

    parentTarget := `${TARGET_NAME}`

    // go through all the direct dependencies and create a cmake targets
    for dependencyName, projectDependency := range projectDependencies {
        var dependencyTargetName string
        var dependencyTarget *DependencyScanStructure

        fullName := dependencyName + "@" + packageVersions[dependencyName]

        if projectDependency.Vendor {
            dependencyTargetName = dependencyName + "__vendor"
        } else {
            dependencyTargetName = dependencyName + "__" + packageVersions[dependencyName]
        }

        if dependencyTarget = scannedDependencies[dependencyTargetName]; dependencyTarget == nil {
            return errors.DependencyDoesNotExistError{
                DependencyName: fullName,
                Vendor:         projectDependency.Vendor,
            }
        }

        log.Verb(queue, "creating cmake target for %s ... ", fullName)
        subQueue := log.GetQueue()

        requiredFlags, requiredDefinitions, err := CreateCMakeTargets(subQueue, parentTarget, false,
            fullName, dependencyTargetName, dependencyTarget, projectFlags.GetGlobalFlags(),
            projectDefinitions.GetGlobalDefinitions(), projectDependency, projectDependency)
        if err != nil {
            log.WriteFailure(queue, log.VERB)
            log.CopyQueue(subQueue, queue, log.TWO_SPACES)
            return err
        }
        log.WriteSuccess(queue, log.VERB)
        log.CopyQueue(subQueue, queue, log.TWO_SPACES)

        log.Verb(queue, "recursively creating cmake targets for %s dependencies ... ", fullName)
        subQueue = log.GetQueue()

        if err := traverseDependencies(subQueue, dependencyTargetName,
            dependencyTarget.MainTag.GetCompileOptions().IsHeaderOnly(), scannedDependencies,
            dependencyTarget.Dependencies, projectFlags.GetGlobalFlags(), requiredFlags,
            projectDefinitions.GetGlobalDefinitions(), requiredDefinitions,
            projectDependency); err != nil {
            log.WriteFailure(queue, log.VERB)
            log.CopyQueue(subQueue, queue, log.TWO_SPACES)
            return err
        }
        log.WriteSuccess(queue, log.VERB)
        log.CopyQueue(subQueue, queue, log.TWO_SPACES)
    }

    cmakePath := cmake.BuildPath(projectPath) + io.Sep + (*target).GetName()
    cmakePath += io.Sep + "dependencies.cmake"

    platform := (*target).GetPlatform()
    avrCmake := cmake.GenerateDependencies(platform, cmakeTargets, cmakeTargetsLink)
    return io.NormalIO.WriteFile(cmakePath, []byte(strings.Join(avrCmake, "\n")))
}
