package dependencies

import (
    "wio/cmd/wio/utils/io"
    "io/ioutil"
    "wio/cmd/wio/utils"
    "github.com/go-errors/errors"
    "wio/cmd/wio/types"
    "path/filepath"
    "wio/cmd/wio/utils/io/log"
    "strings"
    "os"
)

const (
    remoteName = "node_modules"
    vendorName = "vendor"
    wioYmlName = "wio.yml"
    targetName = "${TARGET_NAME}"
)

var packageVersions = map[string]string{}    /* Keeps track of versions for the packages */
var cmakeTargets = map[string]*CMakeTarget{} /* CMake Target that will be built */
var cmakeTargetsLink []CMakeTargetLink       /* CMake Target to Link to and from */
var cmakeTargetNames = map[string]bool{}     /* CMake Target Names. Used to check for unique names */

// CMake Target information
type CMakeTarget struct {
    TargetName string
    Path       string
    Flags      []string
    HeaderOnly bool
}

// CMake Target Link information
type CMakeTargetLink struct {
    From       string
    To         string
    HeaderOnly bool
}

// Stores information about every package that is scanned
type DependencyScanPackage struct {
    Name         string
    Directory    string
    Version      string
    FromVendor   bool
    MainTag      types.PkgTag
    Dependencies types.DependenciesTag
}

// Creates Scan structures for all the scanned packages
func createDependencyScanPackage(depPath string, fromVendor bool) (*DependencyScanPackage, error) {
    wioPath := depPath + io.Sep + wioYmlName
    wioObject := types.PkgConfig{}
    dependencyPackage := &DependencyScanPackage{}

    if err := io.NormalIO.ParseYml(wioPath, &wioObject); err != nil {
        return nil, err
    } else {
        dependencyPackage.Directory = depPath
        dependencyPackage.Name = wioObject.MainTag.Name
        dependencyPackage.FromVendor = fromVendor
        dependencyPackage.Version = wioObject.MainTag.Version
        dependencyPackage.MainTag = wioObject.MainTag
        dependencyPackage.Dependencies = wioObject.DependenciesTag

        packageVersions[dependencyPackage.Name] = dependencyPackage.Version

        return dependencyPackage, nil
    }
}

// Go through all the dependency packages and get information about them
func recursiveDependencyScan(currDirectory string, dependencies map[string]*DependencyScanPackage,
    providedFlags map[string][]string) (error) {
    // if directory does not exist, do not do anything
    if !utils.PathExists(currDirectory) {
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

            if !utils.PathExists(dirPath + io.Sep + wioYmlName) {
                return errors.New(dir.Name() + " is not a valid wio package")
            }

            var fromVendor = false

            // check if the current directory is for remote or vendor
            if filepath.Base(currDirectory) == vendorName {
                fromVendor = true
            }

            // create DependencyPackage
            if dependencyPackage, err := createDependencyScanPackage(dirPath, fromVendor); err != nil {
                return nil
            } else {
                if dependencyPackage.FromVendor {
                    dependencies[dependencyPackage.Name+"__vendor"] = dependencyPackage
                } else {
                    dependencies[dependencyPackage.Name+"__"+dependencyPackage.Version] = dependencyPackage
                }
            }

            if utils.PathExists(dirPath + io.Sep + remoteName) {
                // if remote directory exists
                if err := recursiveDependencyScan(dirPath+io.Sep+remoteName, dependencies, providedFlags); err != nil {
                    return err
                }
            } else if utils.PathExists(dirPath + io.Sep + vendorName) {
                // if vendor directory exists
                if err := recursiveDependencyScan(dirPath+io.Sep+vendorName, dependencies, providedFlags); err != nil {
                    return err
                }
            }
        }
    }

    return nil
}

// When we are building for pkg type, we will copy the files into the remote directory
// This will be picked up while scanning and hence the rest of build process stays the same
func convertPkgToDependency(remotePackagesPath string, projectName string, projectDirectory string) (error) {
    if !utils.PathExists(remotePackagesPath) {
        if err := os.MkdirAll(remotePackagesPath, os.ModePerm); err != nil {
            return err
        }
    }

    if utils.PathExists(remotePackagesPath + io.Sep + projectName) {
        if err := os.RemoveAll(remotePackagesPath + io.Sep + projectName); err != nil {
            return err
        }
    }

    // copy src directory
    if err := utils.CopyDir(projectDirectory+io.Sep+"src", remotePackagesPath+io.Sep+projectName+io.Sep+"src"); err != nil {
        return err
    }

    // copy include directory
    if err := utils.CopyDir(projectDirectory+io.Sep+"include", remotePackagesPath+io.Sep+projectName+io.Sep+"include"); err != nil {
        return err
    }

    // copy wio.yml file
    if err := utils.CopyFile(projectDirectory+io.Sep+"wio.yml", remotePackagesPath+io.Sep+projectName+io.Sep+"wio.yml"); err != nil {
        return err
    }

    return nil
}

// parses dependencies and creates a dependencies.cmake file
func CreateCMakeDependencies(projectName string, directory string, providedFlags map[string][]string,
    projectDependencies types.DependenciesTag, isAPP bool, headerOnly bool) (error) {

    remotePackagesPath := directory + io.Sep + ".wio" + io.Sep + remoteName
    vendorPackagesPath := directory + io.Sep + vendorName
    dependencyPackages := map[string]*DependencyScanPackage{}

    if !isAPP {
        if err := convertPkgToDependency(remotePackagesPath, projectName, directory); err != nil {
            return err
        }
    }

    // scan vendor folder
    if err := recursiveDependencyScan(vendorPackagesPath, dependencyPackages, providedFlags); err != nil {
        return err
    }

    // scan remote folder
    if err := recursiveDependencyScan(remotePackagesPath, dependencyPackages, providedFlags); err != nil {
        return err
    }

    globalFlags := providedFlags["global_flags"]
    projectTarget := targetName

    // go through all the direct dependencies and create a cmake targets
    for projectDependencyName, projectDependency := range projectDependencies {
        var dependencyTargetName string
        var dependencyTarget *DependencyScanPackage
        dependencyName := projectDependencyName + "@" + projectDependency.Version

        if projectDependency.Vendor {
            dependencyTargetName = projectDependencyName + "__vendor"
        } else {
            dependencyTargetName = projectDependencyName + "__" + packageVersions[projectDependencyName]
        }

        if dependencyTarget = dependencyPackages[dependencyTargetName]; dependencyTarget == nil {
            return errors.New(dependencyName + " does not exist. Pull the dependency or check vendor folder")
        }

        requiredFlags, err := createCMakeTargets(projectTarget, headerOnly, dependencyName, dependencyTargetName, dependencyTarget,
            globalFlags, projectDependency.DependencyFlags)
        if err != nil {
            return err
        }

        if err := recursivelyGoThroughTransDependencies(dependencyTargetName, dependencyTarget.MainTag.HeaderOnly, dependencyPackages, dependencyTarget.Dependencies, globalFlags, requiredFlags); err != nil {
            return err
        }
    }

    return io.NormalIO.WriteFile(directory+io.Sep+".wio"+io.Sep+"build"+io.Sep+"dependencies.cmake", []byte(strings.Join(generateAvrDependencyCMakeString(cmakeTargets, cmakeTargetsLink), "\n")))
}

// Creates main cmake file that will build the project
func CreateMainCMake(projectName string, directory string, board string, port string, framework string, targetName string, projectFlags map[string][]string) (error) {
    // create cmake for App type
    if err := generateAvrMainCMakeLists(projectName, directory, board, port, framework, targetName, projectFlags); err != nil {
        log.Verb.Verbose(true, "failure")
        return err
    }

    return nil
}
