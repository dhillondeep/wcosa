package dependencies

import (
    "sort"
    "strconv"
    "strings"
    "wio/cmd/wio/commands/run/cmake"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/log"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
)

// recursively goes through dependencies and creates CMake target and CMake Link
func traverseDependencies(queue *log.Queue, parentName string, parentHeaderOnly bool,
    dependencyPackages map[string]*DependencyScanStructure, dependencies types.DependenciesTag,
    globalFlags []string, requiredFlags []string, globalDefinitions []string, requiredDefinitions []string,
    projectDependency *types.DependencyTag) error {
    for transDependencyName, transDependencyPackage := range dependencies {
        var err error

        var dependencyTargetName string
        var dependencyTarget *DependencyScanStructure
        dependencyNameToUseForLogs := transDependencyName + "@" + packageVersions[transDependencyName]

        if transDependencyPackage.Vendor {
            dependencyTargetName = transDependencyName + "__vendor"
        } else {
            dependencyTargetName = transDependencyName + "__" + packageVersions[transDependencyName]
        }

        // based on required flags and global flags, fill placeholder values
        transDependencyPackage.Flags, err = fillPlaceholderFlags(append(requiredFlags, globalFlags...),
            transDependencyPackage.Flags, dependencyNameToUseForLogs)
        if err != nil {
            return err
        }

        // based on required definitions and global definitions, fill placeholder values
        transDependencyPackage.Definitions, err = fillPlaceholderFlags(append(requiredDefinitions, globalDefinitions...),
            transDependencyPackage.Definitions, dependencyNameToUseForLogs)
        if err != nil {
            return err
        }

        if dependencyTarget = dependencyPackages[dependencyTargetName]; dependencyTarget == nil {
            return errors.DependencyDoesNotExistError{
                DependencyName: dependencyNameToUseForLogs,
                Vendor:         transDependencyPackage.Vendor,
            }
        }

        log.Verb(queue, "creating cmake target for %s ...", dependencyNameToUseForLogs)

        subQueue := log.GetQueue()

        requiredFlags, requiredDefinitions, err := CreateCMakeTargets(subQueue, parentName, parentHeaderOnly,
            dependencyNameToUseForLogs, dependencyTargetName, dependencyTarget, globalFlags,
            globalDefinitions, projectDependency, transDependencyPackage)
        if err != nil {
            log.WriteFailure(queue, log.VERB)
            log.CopyQueue(subQueue, queue, log.NO_SPACES)
            return err
        } else {
            log.WriteSuccess(queue, log.VERB)
            log.CopyQueue(subQueue, queue, log.NO_SPACES)
        }

        log.Verb(queue, "recursively creating cmake targets for %s dependencies ...", dependencyNameToUseForLogs)
        subQueue = log.GetQueue()

        if err := traverseDependencies(queue, dependencyTargetName,
            dependencyTarget.MainTag.GetCompileOptions().IsHeaderOnly(), dependencyPackages,
            dependencyTarget.Dependencies, globalFlags, globalDefinitions, requiredFlags, requiredDefinitions,
            projectDependency); err != nil {
            log.WriteFailure(queue, log.VERB)
            log.CopyQueue(subQueue, queue, log.NO_SPACES)
            return err
        } else {
            log.WriteSuccess(queue, log.VERB)
            log.CopyQueue(subQueue, queue, log.NO_SPACES)
        }
    }

    return nil
}

// Recursively calls traverseDependencies function and creates CMake targets
func CreateCMakeTargets(queue *log.Queue, parentTargetName string, parentTargetHeaderOnly bool,
    dependencyNameToUseForLogs string, dependencyTargetName string, dependencyTarget *DependencyScanStructure,
    globalFlags []string, globalDefinitions []string, configDependency *types.DependencyTag,
    targetDependency *types.DependencyTag) ([]string, []string, error) {

    var allFlags []string
    var allDefinitions []string
    otherFlags := targetDependency.Flags
    otherDefinitions := targetDependency.Definitions

    // get all the manual flags provided
    if val, exists := configDependency.DependencyFlags[strings.Split(dependencyNameToUseForLogs, "@")[0]]; exists {
        otherFlags = append(otherFlags, val...)
    }

    // get all the manual definitions provided
    if val, exists := configDependency.DependencyDefinitions[strings.Split(dependencyNameToUseForLogs, "@")[0]]; exists {
        otherDefinitions = append(otherDefinitions, val...)
    }

    var dependencyTargetRequiredFlags []string
    var dependencyTargetRequiredDefinitions []string
    var optionalFlags []string
    var optionalDefinitions []string

    log.Verb(queue, "matching global flags for dependency: %s ... ", dependencyNameToUseForLogs)

    // match global flags
    dependencyTargetGlobalFlags, err := fillGlobalFlags(globalFlags, dependencyTarget.MainTag.Flags.GlobalFlags,
        dependencyNameToUseForLogs)
    if err != nil {
        log.WriteFailure(queue, log.VERB)
        return nil, nil, err
    } else {
        log.WriteSuccess(queue, log.VERB)
    }

    log.Verb(queue, "matching global definitions for dependency: %s ... ", dependencyNameToUseForLogs)

    // match global definitions
    dependencyTargetGlobalDefinitions, err := fillGlobalFlags(globalDefinitions, dependencyTarget.MainTag.Definitions.GlobalDefinitions,
        dependencyNameToUseForLogs)
    if err != nil {
        log.WriteFailure(queue, log.VERB)
        return nil, nil, err
    } else {
        log.WriteSuccess(queue, log.VERB)
    }

    // when only required flags are allowed
    if dependencyTarget.MainTag.Flags.AllowOnlyRequiredFlags &&
        !dependencyTarget.MainTag.Flags.AllowOnlyGlobalFlags &&
        len(dependencyTargetGlobalFlags) > 0 {
        log.Warnln(queue, "Ignoring global flags")
        dependencyTargetGlobalFlags = []string{}
    }

    // when only required definitions are allowed
    if dependencyTarget.MainTag.Definitions.AllowOnlyRequiredDefinitions &&
        !dependencyTarget.MainTag.Definitions.AllowOnlyGlobalDefinitions &&
        len(dependencyTargetGlobalDefinitions) > 0 {
        log.Warnln(queue, "Ignoring global definitions")
        dependencyTargetGlobalDefinitions = []string{}
    }

    allFlags = utils.AppendIfMissing(allFlags, dependencyTargetGlobalFlags)
    allDefinitions = utils.AppendIfMissing(allDefinitions, dependencyTargetGlobalDefinitions)

    ////////////////////////////////// Handle other flags besides global ////////////////////////////////////////
    if dependencyTarget.MainTag.Flags.AllowOnlyGlobalFlags && len(otherFlags) > 0 {
        log.Warnln(queue, "Ignoring non-global flags")
    } else {
        // match required flags with required flags requested and error out if they do not exist
        dependencyTargetRequiredFlags, optionalFlags, err = fillRequiredFlags(otherFlags,
            dependencyTarget.MainTag.Flags.RequiredFlags, dependencyNameToUseForLogs, parentTargetName, true)
        if err != nil {
            return nil, nil, err
        }

        allFlags = utils.AppendIfMissing(allFlags, dependencyTargetRequiredFlags)
        if dependencyTarget.MainTag.Flags.AllowOnlyRequiredFlags && len(optionalFlags) > 0 {
            log.Warnln(queue, "Ignoring optional flags")
        } else {
            allFlags = utils.AppendIfMissing(allFlags, optionalFlags)
        }
    }

    //////////////////////////////////// Handle other definitions besides global ////////////////////////////////////////
    if dependencyTarget.MainTag.Definitions.AllowOnlyGlobalDefinitions && len(otherDefinitions) > 0 {
        log.Warnln(queue, "Ignoring non-global definitions")
    } else {
        // match required flags with required flags requested and error out if they do not exist
        dependencyTargetRequiredDefinitions, optionalDefinitions, err = fillRequiredFlags(otherDefinitions,
            dependencyTarget.MainTag.Definitions.RequiredDefinitions, dependencyNameToUseForLogs, parentTargetName, true)
        if err != nil {
            return nil, nil, err
        }

        allDefinitions = utils.AppendIfMissing(allDefinitions, dependencyTargetRequiredDefinitions)

        if dependencyTarget.MainTag.Definitions.AllowOnlyRequiredDefinitions && len(optionalFlags) > 0 {
            log.Warnln(queue, "Ignoring optional flags")
        } else {
            allDefinitions = utils.AppendIfMissing(allDefinitions, optionalDefinitions)
        }
    }

    // this is to make sure when a flag or a definition is unique, we create a new target
    sort.Strings(allFlags)
    sort.Strings(allDefinitions)
    hash := dependencyTargetName + strings.Join(allFlags, "") + strings.Join(allDefinitions, "")

    linkVisibility := strings.ToUpper(configDependency.LinkVisibility)

    if val, exists := cmakeTargets[hash]; exists {
        linkVisibility = linkVisibilityVerify(linkVisibility, parentTargetHeaderOnly)
        cmakeTargetsLink = append(cmakeTargetsLink, cmake.TargetLink{From: parentTargetName, To: val.TargetName, Visibility: linkVisibility})
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
        allFlags = utils.AppendIfMissing(allFlags, dependencyTarget.MainTag.Flags.IncludedFlags)

        // add include definitions
        allDefinitions = utils.AppendIfMissing(allDefinitions, dependencyTarget.MainTag.Definitions.IncludedDefinitions)

        // verify linker visibility
        linkVisibility = linkVisibilityVerify(linkVisibility, parentTargetHeaderOnly)

        cmakeTargets[hash] = &cmake.Target{TargetName: dependencyNameToUse,
            Path: dependencyTarget.Directory, Flags: allFlags, Definitions: allDefinitions,
            HeaderOnly: dependencyTarget.MainTag.CompileOptions.HeaderOnly}

        cmakeTargetsLink = append(cmakeTargetsLink, cmake.TargetLink{From: parentTargetName, To: dependencyNameToUse,
            Visibility: linkVisibility})

        cmakeTargetNames[dependencyNameToUse] = true
    }

    cmakeTargets[hash].FlagsVisibility = flagsDefinitionsVisibility(
        dependencyTarget.MainTag.Flags.Visibility, dependencyTarget.MainTag.CompileOptions.HeaderOnly)

    cmakeTargets[hash].DefinitionsVisibility = flagsDefinitionsVisibility(
        dependencyTarget.MainTag.Definitions.Visibility, dependencyTarget.MainTag.CompileOptions.HeaderOnly)

    return dependencyTargetRequiredFlags, dependencyTargetRequiredDefinitions, nil
}

// This checks and verifies flag/definition visibility to make sure it is valid
func flagsDefinitionsVisibility(givenVisibility string, headerOnly bool) string {
    if givenVisibility != "PRIVATE" && givenVisibility != "PUBLIC" && givenVisibility != "INTERFACE" {
        if headerOnly {
            givenVisibility = "INTERFACE"
        } else {
            givenVisibility = "PRIVATE"
        }
    } else {
        if headerOnly {
            givenVisibility = "INTERFACE"
        }
    }

    return givenVisibility
}

// This checks and verifies link visibility to make sure it is valid
func linkVisibilityVerify(givenVisibility string, headerOnly bool) string {
    givenVisibility = strings.ToUpper(givenVisibility)

    if givenVisibility != "PRIVATE" && givenVisibility != "PUBLIC" && givenVisibility != "INTERFACE" {
        if headerOnly {
            givenVisibility = "INTERFACE"
        } else {
            givenVisibility = "PRIVATE"
        }
    } else {
        if headerOnly && givenVisibility != "INTERFACE" {
            givenVisibility = "INTERFACE"
        }
    }

    return givenVisibility
}
