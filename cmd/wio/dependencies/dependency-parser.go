package dependencies

import (
    "wio/cmd/wio/utils/io"
    "io/ioutil"
    "wio/cmd/wio/utils"
    "github.com/go-errors/errors"
    "wio/cmd/wio/types"
    "path/filepath"
    "wio/cmd/wio/utils/io/log"
    "regexp"
    "strings"
    "os"
)

const remoteName = "node_modules"
const vendorName = "vendor"
const wioYmlName = "wio.yml"

var buildStatus = true

/*
1) We need to scan all the node_modules and make a map of all the targets to create
2) Each remote target will be name-version
3) Each vendor target will be name-vendor
 */

type DependencyPackage struct {
    Name                   string
    Directory              string
    Version                string
    FromVendor             bool
    MainTag                types.PkgTag
    TransitiveDependencies types.DependenciesTag
}

func matchFlag(providedFlag string, requestedFlag string) (string) {
    pat := regexp.MustCompile(`^` + strings.ToLower(requestedFlag) + `\b`)
    s := pat.FindString(strings.ToLower(providedFlag))

    return s
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

        // add transitive dependencies
        dependencyPackage.TransitiveDependencies = wioObject.DependenciesTag

        return dependencyPackage, nil
    }
}

func recursiveRemoteScan(currDirectory string, dependencies map[string]*DependencyPackage,
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
                    dependencies[dependencyPackage.Name+"__vendor"] = dependencyPackage
                } else {
                    dependencies[dependencyPackage.Name+"__"+dependencyPackage.Version] = dependencyPackage
                }
            }

            if utils.PathExists(dirPath + io.Sep + remoteName) {
                // if remote directory exists
                if err := recursiveRemoteScan(dirPath+io.Sep+remoteName, dependencies, providedFlags); err != nil {
                    return err
                }
            } else if utils.PathExists(dirPath + io.Sep + vendorName) {
                // if vendor directory exists
                // if remote directory exists
                if err := recursiveRemoteScan(dirPath+io.Sep+vendorName, dependencies, providedFlags); err != nil {
                    return err
                }
            }
        }
    }

    return nil
}

func fillGlobalFlags(globalFlags []string, dependencyGlobalFlagsRequired []string, dependencyName string) ([]string) {
    var filledFlags []string
    var notFilledFlags [] string

    if len(globalFlags) == 0 {
        notFilledFlags = dependencyGlobalFlagsRequired
    } else {
        for _, requiredGlobalFlag := range dependencyGlobalFlagsRequired {
            for _, giveGlobalFlag := range globalFlags {
                s := matchFlag(giveGlobalFlag, requiredGlobalFlag)

                if s != "" {
                    filledFlags = append(filledFlags, giveGlobalFlag)
                } else {
                    notFilledFlags = append(notFilledFlags, requiredGlobalFlag)
                }
            }
        }
    }

    // print errors when global flags are not provided
    if len(dependencyGlobalFlagsRequired) != len(filledFlags) {
        if buildStatus {
            log.Norm.Red(true, "Global flags missing")
        }
        buildStatus = false

        log.Norm.Cyan(true, "  Dependency: "+dependencyName)

        if len(filledFlags) != 0 {
            log.Norm.Write(true, "    Provided Global Flags: "+strings.Join(filledFlags, ","))
        }
        log.Norm.Write(true, "    Missing Global Flags: "+strings.Join(notFilledFlags, ","))
    }

    return filledFlags
}

func fillOtherFlags(providedFlags []string, dependencyPackage *DependencyPackage, fromDep string, toDep string) ([]string) {
    var filledFlags []string
    provideRequiredFlag := make(map[string]bool)

    for _, givenFlag := range providedFlags {
        flagUsed := false

        for _, requiredFlag := range dependencyPackage.MainTag.RequiredFlags {
            if s := matchFlag(givenFlag, requiredFlag); s != "" {
                flagUsed = true
                filledFlags = append(filledFlags, givenFlag)
                provideRequiredFlag[requiredFlag] = true
            }
        }

        if !flagUsed {
            filledFlags = append(filledFlags, givenFlag)
        }
    }

    if len(dependencyPackage.MainTag.RequiredFlags) < len(providedFlags) {
        if buildStatus {
            log.Norm.Red(true, "Required flags missing")
        }
        buildStatus = false

        log.Norm.Cyan(true, "  From : " + fromDep + "\t" + toDep)

        // case where no flag is provided
        if len(providedFlags) == 0 {
            log.Norm.Write(true, "    Missing Flags: "+strings.Join(dependencyPackage.MainTag.RequiredFlags, ","))
        } else {
            providedString := "    Provided Flags: "
            missingString := "    Missing Flags: "

            for pkgName, pkgStatus := range provideRequiredFlag {
                if pkgStatus {
                    providedString += pkgName + ", "
                } else {
                    missingString += pkgName + ", "
                }
            }
        }
    }

    return filledFlags
}

func ScanDependencyPackages(directory string, providedFlags map[string][]string) (error) {
    remotePackagesPath := directory + io.Sep + ".wio" + io.Sep + remoteName
    vendorPackagesPath := directory + io.Sep + vendorName
    dependencyPackages := map[string]*DependencyPackage{}

    // scan vendor folder
    if err := recursiveRemoteScan(vendorPackagesPath, dependencyPackages, providedFlags); err != nil {
        return err
    }

    // scan remote folder
    if err := recursiveRemoteScan(remotePackagesPath, dependencyPackages, providedFlags); err != nil {
        return err
    }

    // fill in global flags and error out if global flag is missing
    for depName, depVal := range dependencyPackages {
        // fill in global flags
        depVal.MainTag.GlobalFlags = fillGlobalFlags(providedFlags["global_flags"], depVal.MainTag.GlobalFlags, depName)
    }

    // if error occurs for global flags, exit
    if !buildStatus {
        os.Exit(2)
    }

    /*
    // fill in required flags for all the direct dependencies


    // fill in required flags and create cmake string
    for depName, depVal := range dependencyPackages {
        // fill in global flags
        depVal.MainTag.GlobalFlags = fillGlobalFlags(providedFlags["global_flags"], depVal.MainTag.GlobalFlags, depName)


        // go through dependencies
        for transitiveDepName, transitiveDepValue := range depVal.TransitiveDependencies {
            fillOtherFlags(transitiveDepValue.DependencyFlags,
                dependencyPackages[transitiveDepName + "__" + transitiveDepValue.Version], depName, transitiveDepName)
        }

        //log.Norm.Write(true, dependencyPackageToCMakeString(depName, depVal))
        //log.Norm.Write(true, "")
    }

    // if error occurs exit
    if !buildStatus {
        os.Exit(2)
    }


    if !buildStatus {
        os.Exit(2)
    }
    */

    return nil
}
