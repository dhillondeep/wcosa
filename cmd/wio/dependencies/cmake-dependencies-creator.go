package dependencies

import (
    "sort"
    "strings"
    "wio/cmd/wio/utils/io/log"
    "wio/cmd/wio/utils"
    "strconv"
    "os"
    "errors"
    "wio/cmd/wio/types"
)

// recursively goes through dependencies and creates CMake target and CMake Link
func recursivelyGoThroughTransDependencies(parentName string, parentHeaderOnly bool, dependencyPackages map[string]*DependencyScanPackage, dependencies types.DependenciesTag, globalFlags []string, requiredFlags []string) (error) {
    for transDependencyName, transDependencyPackage := range dependencies {
        var err error

        // based on required flags, fill placeholder values
        transDependencyPackage.DependencyFlags, err = fillPlaceholderFlags(requiredFlags, transDependencyPackage.DependencyFlags, transDependencyName)
        if err != nil {
            return err
        }

        var dependencyTargetName string
        var dependencyTarget *DependencyScanPackage
        dependencyName := transDependencyName + "@" + transDependencyPackage.Version

        if transDependencyPackage.Vendor {
            dependencyTargetName = transDependencyName + "__vendor"
        } else {
            dependencyTargetName = transDependencyName + "__" + packageVersions[transDependencyName]
        }

        if dependencyTarget = dependencyPackages[dependencyTargetName]; dependencyTarget == nil {
            return errors.New(dependencyName + " does not exist. Pull the dependency or check vendor folder")
        }

        requiredFlags, err := createCMakeTargets(parentName, parentHeaderOnly, dependencyName, dependencyTargetName, dependencyTarget,
            globalFlags, transDependencyPackage.DependencyFlags)
        if err != nil {
            os.Exit(2)
        }

        if err := recursivelyGoThroughTransDependencies(dependencyTargetName, dependencyTarget.MainTag.HeaderOnly, dependencyPackages, dependencyTarget.Dependencies, globalFlags, requiredFlags); err != nil {
            return err
        }
    }

    return nil
}

// Recursively calls recursivelyGoThroughTransDependencies function and creates CMake targets
func createCMakeTargets(parentTargetName string, parentTargetHeaderOnly bool, dependencyName string, dependencyTargetName string,
    dependencyTarget *DependencyScanPackage, globalFlags []string, otherFlags []string) ([]string, error) {
    var allFlags []string
    var dependencyTargetRequiredFlags []string
    var optionalFlags []string

    // match global flags with global flags requested and error out if they do not exist
    dependencyTargetGlobalFlags, err := fillGlobalFlags(globalFlags, dependencyTarget.MainTag.GlobalFlags,
        dependencyName)
    if err != nil {
        return nil, err
    }

    if dependencyTarget.MainTag.AllowOnlyRequiredFlags && !dependencyTarget.MainTag.AllowOnlyGlobalFlags &&
        len(dependencyTargetGlobalFlags) > 0 {
        log.Norm.Red(true, "Cannot require global flags when only required flags are allowed")
        log.Norm.Cyan(true, "  Dependency: "+dependencyName)
        return nil, errors.New("invalid requirements for " + dependencyName + " package")
    }

    allFlags = utils.AppendIfMissing(allFlags, dependencyTargetGlobalFlags)

    if dependencyTarget.MainTag.AllowOnlyGlobalFlags && len(otherFlags) > 0 {
        log.Norm.Yellow(true, "Warning: "+dependencyName+" only accepts global flags")
        log.Norm.Write(true, "  Provided flags will be ignored in the build process")
    } else {
        // match required flags with required flags requested and error out if they do not exist
        dependencyTargetRequiredFlags, optionalFlags, err = fillRequiredFlags(otherFlags,
            dependencyTarget.MainTag.RequiredFlags, dependencyName, parentTargetName, true)
        if err != nil {
            return nil, err
        }

        allFlags = utils.AppendIfMissing(allFlags, dependencyTargetRequiredFlags)

        if dependencyTarget.MainTag.AllowOnlyRequiredFlags && len(optionalFlags) > 0 {
            log.Norm.Yellow(true, "Warning: "+dependencyName+" only accepts required flags")
            log.Norm.Write(true, "  Provided optional flags will be ignored in the build process")
        } else {
            allFlags = utils.AppendIfMissing(allFlags, optionalFlags)
        }
    }

    sort.Strings(allFlags)
    hash := dependencyTargetName + strings.Join(allFlags, "")

    if val, exists := cmakeTargets[hash]; exists {
        cmakeTargetsLink = append(cmakeTargetsLink, CMakeTargetLink{From: parentTargetName, To: val.TargetName, HeaderOnly: parentTargetHeaderOnly})
    } else {
        dependencyNameToUse := dependencyTargetName
        counter := 2

        for {
            if _, exists := cmakeTargetNames[dependencyNameToUse]; exists {
                dependencyNameToUse += "__" + strconv.Itoa(counter)
                counter++
            } else {
                break
            }
        }

        // add include flags
        allFlags = utils.AppendIfMissing(allFlags, dependencyTarget.MainTag.IncludedFlags)

        cmakeTargets[hash] = &CMakeTarget{TargetName: dependencyNameToUse,
            Path: dependencyTarget.Directory, Flags: allFlags, HeaderOnly: dependencyTarget.MainTag.HeaderOnly}
        cmakeTargetsLink = append(cmakeTargetsLink, CMakeTargetLink{From: parentTargetName, To: dependencyNameToUse, HeaderOnly: parentTargetHeaderOnly})
        cmakeTargetNames[dependencyNameToUse] = true
    }

    return dependencyTargetRequiredFlags, nil
}
