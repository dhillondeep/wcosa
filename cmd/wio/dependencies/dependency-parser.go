package dependencies

import (
    "wio/cmd/wio/utils/io"
    "io/ioutil"
    "wio/cmd/wio/utils"
    "github.com/go-errors/errors"
    "wio/cmd/wio/types"
    "path/filepath"
)

const remoteName = "node_modules"
const vendorName = "vendor"
const wioYmlName = "wio.yml"

/*
1) We need to scan all the node_modules and make a map of all the targets to create
2) Each remote target will be name-version
3) Each vendor target will be name-vendor
 */

type DependencyPackage struct {
    Name       string
    Directory string
    Version string
    FromVendor bool
    MainTag types.PkgTag
}

func createDependencyPackage(depPath string, fromVendor bool) (*DependencyPackage, error) {
    wioPath := depPath + io.Sep + wioYmlName
    wioObject := types.PkgConfig{}
    dependencyPackage := &DependencyPackage{}

    if err := io.NormalIO.ParseYml(wioPath, &wioObject); err != nil {
        return nil, err
    } else {
        dependencyPackage.Directory = depPath
        dependencyPackage.Name = wioObject.MainTag.Name
        dependencyPackage.FromVendor = fromVendor
        dependencyPackage.Version = wioObject.MainTag.Version
        dependencyPackage.MainTag = wioObject.MainTag

        return dependencyPackage, nil
    }
}

func recursiveRemoteScan(currDirectory string, dependencies map[string]*DependencyPackage) (error) {
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

            // path for the directory we are looking at
            dirPath := currDirectory + io.Sep + dir.Name()

            // throw and error if wio.yml file does not exist
            if !utils.PathExists(dirPath + io.Sep + wioYmlName) {
                return errors.New(dir.Name() + " is not a valid wio package")
            }

            var fromVendor = false

            // check if the current directory is for remote or vendor
            if filepath.Base(currDirectory) == vendorName {
                fromVendor = true
            }

            // create DependencyPackage
            if dependencyPackage, err := createDependencyPackage(dirPath, fromVendor);
                err != nil {
                return nil
            } else {
                if dependencyPackage.FromVendor {
                    dependencies[dependencyPackage.Name+"-vendor"] = dependencyPackage
                } else {
                    dependencies[dependencyPackage.Name+"-"+dependencyPackage.Version] = dependencyPackage
                }
            }

            if utils.PathExists(dirPath + io.Sep + remoteName) {
                // if remote directory exists
                if err := recursiveRemoteScan(dirPath + io.Sep + remoteName, dependencies); err != nil {
                    return err
                }
            } else if utils.PathExists(dirPath + io.Sep + vendorName) {
                // if vendor directory exists
                // if remote directory exists
                if err := recursiveRemoteScan(dirPath + io.Sep + vendorName, dependencies); err != nil {
                    return err
                }
            }
        }
    }

    return nil
}

func ScanDependencyPackages(directory string) (error) {
    remotePackagesPath := directory + io.Sep + ".wio" + io.Sep + remoteName
    vendorPackagesPath := directory + io.Sep + vendorName
    dependencyPackages := map[string]*DependencyPackage{}

    // scan vendor folder
    if err := recursiveRemoteScan(vendorPackagesPath, dependencyPackages); err != nil {
        return err
    }

    // scan remote folder
    if err := recursiveRemoteScan(remotePackagesPath, dependencyPackages); err != nil {
        return err
    }

    return nil
}
