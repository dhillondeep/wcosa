package run

import (
    "wio/cmd/wio/log"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"
    "os"
    "wio/cmd/wio/utils"
    "github.com/fatih/color"
    "io/ioutil"
    "path/filepath"
    "wio/cmd/wio/errors"
    goerr "errors"
    "strings"
    "wio/cmd/wio/commands/create"
)

const (
    REMOTE_NAME     = "node_modules"
    VENDOR_NAME     = "vendor"
    PKG_REMOTE_NAME = "pkg_module"
)

var packageVersions = map[string]string{}    /* Keeps track of versions for the packages */
var cmakeTargets = map[string]*CMakeTarget{} /* CMake Target that will be built */
var cmakeTargetsLink []CMakeTargetLink       /* CMake Target to Link to and from */
var cmakeTargetNames = map[string]bool{}     /* CMake Target Names. Used to check for unique names */

// CMake Target information
type CMakeTarget struct {
    TargetName            string
    Path                  string
    Flags                 []string
    Definitions           []string
    FlagsVisibility       string
    DefinitionsVisibility string
    HeaderOnly            bool
}

// CMake Target Link information
type CMakeTargetLink struct {
    From           string
    To             string
    LinkVisibility string
}

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
func createDependencyScanStructure(depPath string, fromVendor bool) (*DependencyScanStructure, error) {
    wioPath := depPath + io.Sep + "wio.yml"
    wioObject := types.PkgConfig{}
    dependencyPackage := &DependencyScanStructure{}

    if err := io.NormalIO.ParseYml(wioPath, &wioObject); err != nil {
        return nil, errors.ConfigParsingError{
            Err: err,
        }
    } else {
        dependencyPackage.Directory = depPath
        dependencyPackage.Name = wioObject.MainTag.Meta.Name
        dependencyPackage.FromVendor = fromVendor
        dependencyPackage.Version = wioObject.MainTag.Meta.Version
        dependencyPackage.MainTag = wioObject.MainTag
        dependencyPackage.Dependencies = wioObject.DependenciesTag
        packageVersions[dependencyPackage.Name] = dependencyPackage.Version

        return dependencyPackage, nil
    }
}

// Go through all the dependency packages and get information about them
func recursiveDependencyScan(queue *log.Queue, currDirectory string, dependencies map[string]*DependencyScanStructure) error {
    // if directory does not exist, do not do anything
    if !utils.PathExists(currDirectory) {
        log.QueueWriteln(queue, log.VERB, nil, "% does not exist, skipping", currDirectory)
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

            dirPath := currDirectory + io.Sep + dir.Name()

            if !utils.PathExists(dirPath + io.Sep + "wio.yml") {
                return errors.NotValidWioProjectError{
                    Directory: dirPath,
                    Err:       goerr.New("wio.yml file missing"),
                }
            }

            var fromVendor = false

            // check if the current directory is for remote or vendor
            if filepath.Base(currDirectory) == VENDOR_NAME {
                fromVendor = true
            }

            // create DependencyScanStructure
            if dependencyPackage, err := createDependencyScanStructure(dirPath, fromVendor); err != nil {
                return err
            } else {
                dependencyName := dependencyPackage.Name + "__vendor"
                if !dependencyPackage.FromVendor {
                    dependencyName = dependencyPackage.Name + "__" + dependencyPackage.Version
                }

                dependencies[dependencyName] = dependencyPackage

                log.QueueWriteln(queue, log.VERB, nil, "%s package stored as dependency named: %s",
                    dirPath, dependencyName)
            }

            if utils.PathExists(dirPath + io.Sep + REMOTE_NAME) {
                // if remote directory exists
                if err := recursiveDependencyScan(queue, dirPath+io.Sep+REMOTE_NAME, dependencies); err != nil {
                    return err
                }
            } else if utils.PathExists(dirPath + io.Sep + VENDOR_NAME) {
                // if vendor directory exists
                if err := recursiveDependencyScan(queue, dirPath+io.Sep+VENDOR_NAME, dependencies); err != nil {
                    return err
                }
            }
        }
    } else {
        log.QueueWriteln(queue, log.VERB, nil, "% does not have any package, skipping", currDirectory)
    }

    return nil
}

// When we are building for pkg type, we will copy the files into the remote directory
// This will be picked up while scanning and hence the rest of build process stays the same
func convertPkgToDependency(packageDependencyPath string, projectName string, projectDirectory string) error {
    if !utils.PathExists(packageDependencyPath) {
        if err := os.MkdirAll(packageDependencyPath, os.ModePerm); err != nil {
            return err
        }
    }

    if utils.PathExists(packageDependencyPath + io.Sep + projectName) {
        if err := os.RemoveAll(packageDependencyPath + io.Sep + projectName); err != nil {
            return err
        }
    }

    // copy src directory
    if err := utils.CopyDir(projectDirectory+io.Sep+"src", packageDependencyPath+io.Sep+projectName+io.Sep+"src"); err != nil {
        return err
    }

    // copy include directory
    if err := utils.CopyDir(projectDirectory+io.Sep+"include", packageDependencyPath+io.Sep+projectName+io.Sep+"include"); err != nil {
        return err
    }

    // copy wio.yml file
    if err := utils.CopyFile(projectDirectory+io.Sep+"wio.yml", packageDependencyPath+io.Sep+projectName+io.Sep+"wio.yml"); err != nil {
        return err
    }

    return nil
}

func createCMakeDependencyTargets(queue *log.Queue, projectName string, projectDirectory string, projectType string,
    projectFlags types.TargetFlags, projectDefinitions types.TargetDefinitions, projectDependencies types.DependenciesTag,
    platform string) error {
    remotePackagesPath := projectDirectory + io.Sep + ".wio" + io.Sep + REMOTE_NAME
    vendorPackagesPath := projectDirectory + io.Sep + VENDOR_NAME

    scannedDependencies := map[string]*DependencyScanStructure{}

    if projectType == types.PKG {
        packageDependencyPath := projectDirectory + io.Sep + ".wio" + io.Sep + PKG_REMOTE_NAME

        log.QueueWrite(queue, log.VERB, nil, "converting project to a dependency to be used in tests ...")

        if err := convertPkgToDependency(packageDependencyPath, projectName, projectDirectory); err != nil {
            log.QueueWrite(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
            return err
        } else {
            log.QueueWrite(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
        }

        log.QueueWrite(queue, log.VERB, nil, "recursively scanning package dependency at path "+
            "%s ...", packageDependencyPath)

        subQueue := log.GetQueue()

        if err := recursiveDependencyScan(subQueue, packageDependencyPath, scannedDependencies); err != nil {
            log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
            log.CopyQueue(subQueue, queue, log.TWO_SPACES)
            return err
        } else {
            log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
            log.CopyQueue(subQueue, queue, log.TWO_SPACES)
        }
    }

    log.QueueWrite(queue, log.VERB, nil, "recursively scanning remote dependencies at path: %s ...", remotePackagesPath)
    subQueue := log.GetQueue()

    if err := recursiveDependencyScan(subQueue, remotePackagesPath, scannedDependencies); err != nil {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
        log.CopyQueue(subQueue, queue, log.TWO_SPACES)
        return err
    } else {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
        log.CopyQueue(subQueue, queue, log.TWO_SPACES)
    }

    log.QueueWrite(queue, log.VERB, nil, "recursively scanning vendor dependencies at path: %s ...",
        vendorPackagesPath)
    subQueue = log.GetQueue()

    if err := recursiveDependencyScan(subQueue, vendorPackagesPath, scannedDependencies); err != nil {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
        log.CopyQueue(subQueue, queue, log.TWO_SPACES)
        return err
    } else {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
        log.CopyQueue(subQueue, queue, log.TWO_SPACES)
    }

    parentTarget := `${TARGET_NAME}`

    // go through all the direct dependencies and create a cmake targets
    for projectDependencyName, projectDependency := range projectDependencies {
        var dependencyTargetName string
        var dependencyTarget *DependencyScanStructure

        dependencyNameToUseForLogs := projectDependencyName + "@" + packageVersions[projectDependencyName]

        if projectDependency.Vendor {
            dependencyTargetName = projectDependencyName + "__vendor"
        } else {
            dependencyTargetName = projectDependencyName + "__" + packageVersions[projectDependencyName]
        }

        if dependencyTarget = scannedDependencies[dependencyTargetName]; dependencyTarget == nil {
            return errors.DependencyDoesNotExistError{
                DependencyName: dependencyNameToUseForLogs,
                Vendor:         projectDependency.Vendor,
            }
        }

        log.QueueWrite(queue, log.VERB, nil, "creating cmake target for %s ...", dependencyNameToUseForLogs)
        subQueue := log.GetQueue()

        requiredFlags, requiredDefinitions, err := createCMakeTargets(subQueue, parentTarget, false,
            dependencyNameToUseForLogs, dependencyTargetName, dependencyTarget, projectFlags.GetGlobalFlags(),
            projectDefinitions.GetGlobalDefinitions(), projectDependency, projectDependency)
        if err != nil {
            log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
            log.CopyQueue(subQueue, queue, log.TWO_SPACES)
            return err
        } else {
            log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
            log.CopyQueue(subQueue, queue, log.TWO_SPACES)
        }

        log.QueueWrite(queue, log.VERB, nil, "recursively creating cmake targets for %s dependencies ...", dependencyNameToUseForLogs)
        subQueue = log.GetQueue()

        if err := recursivelyGoThroughTransDependencies(subQueue, dependencyTargetName,
            dependencyTarget.MainTag.GetCompileOptions().IsHeaderOnly(), scannedDependencies,
            dependencyTarget.Dependencies, projectFlags.GetGlobalFlags(), requiredFlags,
            projectDefinitions.GetGlobalDefinitions(), requiredDefinitions,
            projectDependency); err != nil {
            log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
            log.CopyQueue(subQueue, queue, log.TWO_SPACES)
            return err
        } else {
            log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
            log.CopyQueue(subQueue, queue, log.TWO_SPACES)
        }
    }

    cmakePath := projectDirectory + io.Sep + ".wio" + io.Sep + "build" + io.Sep + "dependencies.cmake"

    if platform == create.AVR {
        avrCmake := generateAvrDependencyCMakeString(cmakeTargets, cmakeTargetsLink)

        return io.NormalIO.WriteFile(cmakePath, []byte(strings.Join(avrCmake, "\n")))
    } else {
        return errors.PlatformNotSupportedError{
            Platform: platform,
            Err:      goerr.New("platform not valid for cmake target creation"),
        }
    }
}
