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
    depPkgs map[string]*DependencyScanStructure, dependencies types.DependenciesTag,
    globalFlags []string, requiredFlags []string, gblDefs []string, requiredDefinitions []string,
    projectDependency *types.DependencyTag) error {
    for depName, depPkg := range dependencies {
        depVer := packageVersions[depName]

        var err error
        var depTgtName string
        var depTgt *DependencyScanStructure

        if depPkg.Vendor {
            depTgtName = depName + "__vendor"
        } else {
            depTgtName = depName + "__" + depVer
        }

        // based on required flags and global flags, fill placeholder values
        depPkg.Flags, err = fillPlaceholders(append(requiredFlags, globalFlags...),
            depPkg.Flags)
        if err != nil {
            return err
        }

        // based on required definitions and global definitions, fill placeholder values
        depPkg.Definitions, err = fillPlaceholders(append(requiredDefinitions, gblDefs...),
            depPkg.Definitions)
        if err != nil {
            return err
        }

        if depTgt = depPkgs[depTgtName]; depTgt == nil {
            return errors.Stringf("missing dependency %s@%s", depName, depVer)
        }
        log.Verb(queue, "creating cmake target for %s@%s ...", depName, depVer)

        subQueue := log.GetQueue()
        requiredFlags, requiredDefinitions, err := CreateCMakeTargets(subQueue, parentName, parentHeaderOnly,
            depName, depTgtName, depTgt, globalFlags,
            gblDefs, projectDependency, depPkg)
        if err != nil {
            log.WriteFailure(queue, log.VERB)
            log.CopyQueue(subQueue, queue, log.NO_SPACES)
            return err
        } else {
            log.WriteSuccess(queue, log.VERB)
            log.CopyQueue(subQueue, queue, log.NO_SPACES)
        }

        log.Verb(queue, "creating cmake targets for %s@%s", depName, depVer)
        subQueue = log.GetQueue()
        if err := traverseDependencies(queue, depTgtName,
            depTgt.MainTag.GetCompileOptions().IsHeaderOnly(), depPkgs,
            depTgt.Dependencies, globalFlags, gblDefs, requiredFlags, requiredDefinitions,
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
    depName string, depTgtName string, depTgt *DependencyScanStructure,
    globalFlags []string, gblDefs []string, config *types.DependencyTag,
    target *types.DependencyTag) ([]string, []string, error) {

    var allFlags []string
    var allDefinitions []string
    otherFlags := target.Flags
    otherDefinitions := target.Definitions

    // get all the manual flags provided
    if val, exists := config.DependencyFlags[depName]; exists {
        otherFlags = append(otherFlags, val...)
    }

    // get all the manual definitions provided
    if val, exists := config.DependencyDefinitions[depName]; exists {
        otherDefinitions = append(otherDefinitions, val...)
    }

    var depReqFlags []string
    var depReqDefn []string
    var optionalFlags []string
    var optionalDefinitions []string
    mainTag := depTgt.MainTag

    // match global flags
    dependencyTargetGlobalFlags, err := fillGlobal(globalFlags, mainTag.Flags.GlobalFlags)
    if err != nil {
        log.WriteFailure(queue, log.VERB)
        return nil, nil, err
    } else {
        log.WriteSuccess(queue, log.VERB)
    }

    // match global definitions
    dependencyTargetGlobalDefinitions, err := fillGlobal(gblDefs, mainTag.Definitions.GlobalDefinitions)
    if err != nil {
        log.WriteFailure(queue, log.VERB)
        return nil, nil, err
    } else {
        log.WriteSuccess(queue, log.VERB)
    }

    // when only required flags are allowed
    if mainTag.Flags.AllowOnlyRequiredFlags &&
        !mainTag.Flags.AllowOnlyGlobalFlags &&
        len(dependencyTargetGlobalFlags) > 0 {
        log.Warnln(queue, "Ignoring global flags")
        dependencyTargetGlobalFlags = []string{}
    }

    // when only required definitions are allowed
    if mainTag.Definitions.AllowOnlyRequiredDefinitions &&
        !mainTag.Definitions.AllowOnlyGlobalDefinitions &&
        len(dependencyTargetGlobalDefinitions) > 0 {
        log.Warnln(queue, "Ignoring global definitions")
        dependencyTargetGlobalDefinitions = []string{}
    }

    allFlags = utils.AppendIfMissing(allFlags, dependencyTargetGlobalFlags)
    allDefinitions = utils.AppendIfMissing(allDefinitions, dependencyTargetGlobalDefinitions)

    ////////////////////////////////// Handle other flags besides global ////////////////////////////////////////
    if mainTag.Flags.AllowOnlyGlobalFlags && len(otherFlags) > 0 {
        log.Warnln(queue, "Ignoring non-global flags")
    } else {
        // match required flags with required flags requested and error out if they do not exist
        depReqFlags, optionalFlags, err = fillRequired(otherFlags,
            mainTag.Flags.RequiredFlags)
        if err != nil {
            return nil, nil, err
        }

        allFlags = utils.AppendIfMissing(allFlags, depReqFlags)
        if mainTag.Flags.AllowOnlyRequiredFlags && len(optionalFlags) > 0 {
            log.Warnln(queue, "Ignoring optional flags")
        } else {
            allFlags = utils.AppendIfMissing(allFlags, optionalFlags)
        }
    }

    //////////////////////////////////// Handle other definitions besides global ////////////////////////////////////////
    if mainTag.Definitions.AllowOnlyGlobalDefinitions && len(otherDefinitions) > 0 {
        log.Warnln(queue, "Ignoring non-global definitions")
    } else {
        // match required flags with required flags requested and error out if they do not exist
        depReqDefn, optionalDefinitions, err = fillRequired(otherDefinitions,
            mainTag.Definitions.RequiredDefinitions)
        if err != nil {
            return nil, nil, err
        }

        allDefinitions = utils.AppendIfMissing(allDefinitions, depReqDefn)

        if mainTag.Definitions.AllowOnlyRequiredDefinitions && len(optionalFlags) > 0 {
            log.Warnln(queue, "Ignoring optional flags")
        } else {
            allDefinitions = utils.AppendIfMissing(allDefinitions, optionalDefinitions)
        }
    }

    // this is to make sure when a flag or a definition is unique, we create a new target
    sort.Strings(allFlags)
    sort.Strings(allDefinitions)
    hash := depTgtName + strings.Join(allFlags, "") + strings.Join(allDefinitions, "")

    linkVisibility := strings.ToUpper(target.LinkVisibility)

    if val, exists := cmakeTargets[hash]; exists {
        linkVisibility = linkVisibilityVerify(linkVisibility, parentTargetHeaderOnly)
        cmakeTargetsLink = append(cmakeTargetsLink, cmake.TargetLink{From: parentTargetName, To: val.TargetName, Visibility: linkVisibility})
    } else {
        dependencyNameToUse := depTgtName
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
        allFlags = utils.AppendIfMissing(allFlags, mainTag.Flags.IncludedFlags)

        // add include definitions
        allDefinitions = utils.AppendIfMissing(allDefinitions, mainTag.Definitions.IncludedDefinitions)

        // verify linker visibility
        linkVisibility = linkVisibilityVerify(linkVisibility, parentTargetHeaderOnly)

        cmakeTargets[hash] = &cmake.Target{TargetName: dependencyNameToUse,
            Path: depTgt.Directory, Flags: allFlags, Definitions: allDefinitions,
            HeaderOnly: mainTag.CompileOptions.HeaderOnly}

        cmakeTargetsLink = append(cmakeTargetsLink, cmake.TargetLink{From: parentTargetName, To: dependencyNameToUse,
            Visibility: linkVisibility})

        cmakeTargetNames[dependencyNameToUse] = true
    }

    cmakeTargets[hash].FlagsVisibility = flagsDefinitionsVisibility(
        mainTag.Flags.Visibility, mainTag.CompileOptions.HeaderOnly)

    cmakeTargets[hash].DefinitionsVisibility = flagsDefinitionsVisibility(
        mainTag.Definitions.Visibility, mainTag.CompileOptions.HeaderOnly)

    return depReqFlags, depReqDefn, nil
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
