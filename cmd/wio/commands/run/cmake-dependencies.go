package run

import (
    goerr "errors"
    "sort"
    "strconv"
    "strings"
    "wio/cmd/wio/log"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "github.com/fatih/color"
    "wio/cmd/wio/errors"
)

// recursively goes through dependencies and creates CMake target and CMake Link
func recursivelyGoThroughTransDependencies(queue *log.Queue, parentName string, parentHeaderOnly bool,
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
                Vendor: transDependencyPackage.Vendor,
            }
        }

        log.QueueWrite(queue, log.VERB, nil,"creating cmake target for %s ...", dependencyNameToUseForLogs)

        subQueue := log.GetQueue()

        requiredFlags, requiredDefinitions, err := createCMakeTargets(subQueue, parentName, parentHeaderOnly,
            dependencyNameToUseForLogs, dependencyTargetName, dependencyTarget, globalFlags,
            globalDefinitions, projectDependency, transDependencyPackage)
        if err != nil {
            log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
            log.CopyQueue(subQueue, queue, log.NO_SPACES)
            return err
        } else {
            log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
            log.CopyQueue(subQueue, queue, log.NO_SPACES)
        }

        log.QueueWrite(queue, log.VERB, nil,"recursively creating cmake targets for %s dependencies ...", dependencyNameToUseForLogs)
        subQueue = log.GetQueue()

        if err := recursivelyGoThroughTransDependencies(queue, dependencyTargetName,
            dependencyTarget.MainTag.GetCompileOptions().IsHeaderOnly(), dependencyPackages,
            dependencyTarget.Dependencies, globalFlags, globalDefinitions, requiredFlags, requiredDefinitions,
            projectDependency); err != nil {
            log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
            log.CopyQueue(subQueue, queue, log.NO_SPACES)
            return err
        } else {
            log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
            log.CopyQueue(subQueue, queue, log.NO_SPACES)
        }
    }

    return nil
}

// Recursively calls recursivelyGoThroughTransDependencies function and creates CMake targets
func createCMakeTargets(queue *log.Queue, parentTargetName string, parentTargetHeaderOnly bool,
    dependencyNameToUseForLogs string, dependencyTargetName string, dependencyTarget *DependencyScanStructure,
    globalFlags []string, globalDefinitions []string, configDependency *types.DependencyTag,
        targetDependency *types.DependencyTag)([]string, []string, error) {

    var allFlags []string
    var allDefinitions []string
    otherFlags := targetDependency.Flags
    otherDefinitions := targetDependency.Definitions

    // get all the manual flags provided
    if val, exists := configDependency.DependencyFlags[strings.Split(dependencyNameToUseForLogs, "@")[0]]; exists {
        otherFlags =  append(otherFlags, val...)
    }

    // get all the manual definitions provided
    if val, exists := configDependency.DependencyDefinitions[strings.Split(dependencyNameToUseForLogs, "@")[0]]; exists {
        otherDefinitions =  append(otherDefinitions, val...)
    }

    var dependencyTargetRequiredFlags []string
    var dependencyTargetRequiredDefinitions []string
    var optionalFlags []string
    var optionalDefinitions []string

    log.QueueWrite(queue, log.VERB, nil, "matching global flags for dependency: %s ...", dependencyNameToUseForLogs)

    // match global flags
    dependencyTargetGlobalFlags, err := fillGlobalFlags(globalFlags, dependencyTarget.MainTag.Flags.GlobalFlags,
        dependencyNameToUseForLogs)
    if err != nil {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
        return nil, nil, err
    } else {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
    }

    log.QueueWrite(queue, log.VERB, nil, "matching global definitions for dependency: %s ...", dependencyNameToUseForLogs)

    // match global definitions
    dependencyTargetGlobalDefinitions, err := fillGlobalFlags(globalDefinitions, dependencyTarget.MainTag.Definitions.GlobalDefinitions,
        dependencyNameToUseForLogs)
    if err != nil {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgRed), "failure")
        return nil, nil, err
    } else {
        log.QueueWriteln(queue, log.VERB_NONE, color.New(color.FgGreen), "success")
    }

    // when only required flags are allowed
    if dependencyTarget.MainTag.Flags.AllowOnlyRequiredFlags &&
        !dependencyTarget.MainTag.Flags.AllowOnlyGlobalFlags &&
        len(dependencyTargetGlobalFlags) > 0 {
        err = errors.OnlyRequiredFlagsError{
            Dependency: dependencyNameToUseForLogs,
            FlagCategoryName: "global flags",
            Err: goerr.New("global flags will be ignored"),
        }

        log.QueueWriteln(queue, log.WARN, nil, err.Error())
        dependencyTargetGlobalFlags = []string{}
    }

    // when only required definitions are allowed
    if dependencyTarget.MainTag.Definitions.AllowOnlyRequiredDefinitions &&
        !dependencyTarget.MainTag.Definitions.AllowOnlyGlobalDefinitions &&
        len(dependencyTargetGlobalDefinitions) > 0 {
        err = errors.OnlyRequiredFlagsError{
            Dependency: dependencyNameToUseForLogs,
            FlagCategoryName: "global definitions",
            Err: goerr.New("global definitions will be ignored"),
        }

        log.QueueWriteln(queue, log.WARN, nil, err.Error())
        dependencyTargetGlobalDefinitions = []string{}
    }

    allFlags = utils.AppendIfMissing(allFlags, dependencyTargetGlobalFlags)
    allDefinitions = utils.AppendIfMissing(allDefinitions, dependencyTargetGlobalDefinitions)

    ////////////////////////////////// Handle other flags besides global ////////////////////////////////////////
    if dependencyTarget.MainTag.Flags.AllowOnlyGlobalFlags && len(otherFlags) > 0 {
        log.QueueWriteln(queue, log.WARN, nil, errors.OnlyGlobalFlagsError{
            Dependency: dependencyNameToUseForLogs,
            FlagCategoryName: "other flags",
            Err: goerr.New("other flags will be ignored in the build process"),
        }.Error())
    } else {
        // match required flags with required flags requested and error out if they do not exist
        dependencyTargetRequiredFlags, optionalFlags, err = fillRequiredFlags(otherFlags,
            dependencyTarget.MainTag.Flags.RequiredFlags, dependencyNameToUseForLogs, parentTargetName, true)
        if err != nil {
            return nil, nil, err
        }

        allFlags = utils.AppendIfMissing(allFlags, dependencyTargetRequiredFlags)

        if dependencyTarget.MainTag.Flags.AllowOnlyRequiredFlags && len(optionalFlags) > 0 {
            log.QueueWriteln(queue, log.WARN, nil, errors.OnlyRequiredFlagsError{
                Dependency: dependencyNameToUseForLogs,
                FlagCategoryName: "optional flags",
                Err: goerr.New("optional flags will be ignored in the build process"),
            }.Error())
        } else {
            allFlags = utils.AppendIfMissing(allFlags, optionalFlags)
        }
    }

    //////////////////////////////////// Handle other definitions besides global ////////////////////////////////////////
    if dependencyTarget.MainTag.Definitions.AllowOnlyGlobalDefinitions && len(otherDefinitions) > 0 {
        log.QueueWriteln(queue, log.WARN, nil, errors.OnlyGlobalDefinitionsError{
            Dependency: dependencyNameToUseForLogs,
            FlagCategoryName: "other definitions",
            Err: goerr.New("other defintions will be ignored in the build process"),
        }.Error())
    } else {
        // match required flags with required flags requested and error out if they do not exist
        dependencyTargetRequiredDefinitions, optionalDefinitions, err = fillRequiredFlags(otherDefinitions,
            dependencyTarget.MainTag.Definitions.RequiredDefinitions, dependencyNameToUseForLogs, parentTargetName, true)
        if err != nil {
            return nil, nil, err
        }

        allDefinitions = utils.AppendIfMissing(allDefinitions, dependencyTargetRequiredDefinitions)

        if dependencyTarget.MainTag.Definitions.AllowOnlyRequiredDefinitions && len(optionalFlags) > 0 {
            log.QueueWriteln(queue, log.WARN, nil, errors.OnlyRequiredDefinitionsError{
                Dependency: dependencyNameToUseForLogs,
                FlagCategoryName: "optional definitions",
                Err: goerr.New("optional definitions will be ignored in the build process"),
            }.Error())
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
        linkVisibility = linkVisibilityVerify(queue, parentTargetName, val.TargetName, linkVisibility, parentTargetHeaderOnly)
        cmakeTargetsLink = append(cmakeTargetsLink, CMakeTargetLink{From: parentTargetName, To: val.TargetName, LinkVisibility: linkVisibility})
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
        linkVisibility = linkVisibilityVerify(queue, parentTargetName, dependencyNameToUse, linkVisibility, parentTargetHeaderOnly)

        cmakeTargets[hash] = &CMakeTarget{TargetName: dependencyNameToUse,
            Path: dependencyTarget.Directory, Flags: allFlags, Definitions: allDefinitions,
            HeaderOnly: dependencyTarget.MainTag.CompileOptions.HeaderOnly}


        cmakeTargetsLink = append(cmakeTargetsLink, CMakeTargetLink{From: parentTargetName, To: dependencyNameToUse,
            LinkVisibility: linkVisibility})

        cmakeTargetNames[dependencyNameToUse] = true
    }

    cmakeTargets[hash].FlagsVisibility = flagsDefinitionsVisibility(queue, dependencyNameToUseForLogs,
        dependencyTarget.MainTag.Flags.Visibility, dependencyTarget.MainTag.CompileOptions.HeaderOnly)

    cmakeTargets[hash].DefinitionsVisibility = flagsDefinitionsVisibility(queue, dependencyNameToUseForLogs,
        dependencyTarget.MainTag.Definitions.Visibility, dependencyTarget.MainTag.CompileOptions.HeaderOnly)

    return dependencyTargetRequiredFlags, dependencyTargetRequiredDefinitions, nil
}

// This checks and verifies flag/definition visibility to make sure it is valid
func flagsDefinitionsVisibility(queue *log.Queue, packageName string, givenVisibility string, headerOnly bool) string {
    if givenVisibility != "PRIVATE" && givenVisibility != "PUBLIC" && givenVisibility != "INTERFACE" {
        var err error
        if headerOnly {
            err = errors.FlagsDefinitionsVisibilityError{
                PackageName: packageName,
                GivenVisibility: givenVisibility,
                Err:             goerr.New("header only => INTERFACE visibility will be used"),
            }
            givenVisibility = "INTERFACE"
        } else {
            err = errors.FlagsDefinitionsVisibilityError{
                PackageName: packageName,
                GivenVisibility: givenVisibility,
                Err:             goerr.New("PRIVATE visibility will be used"),
            }
            givenVisibility = "PRIVATE"
        }

        log.QueueWriteln(queue, log.WARN, nil, err.Error())
    } else {
        if headerOnly && givenVisibility != "INTERFACE" {
            err := errors.FlagsDefinitionsVisibilityError{
                PackageName: packageName,
                GivenVisibility: givenVisibility,
                Err:             goerr.New("package is header only => INTERFACE visibility will be used"),
            }

            givenVisibility = "INTERFACE"
            log.QueueWriteln(queue, log.WARN, nil, err.Error())
        }
    }

    return givenVisibility
}

// This checks and verifies link visibility to make sure it is valid
func linkVisibilityVerify(queue *log.Queue, from string, to string, givenVisibility string, headerOnly bool) string {
    givenVisibility = strings.ToUpper(givenVisibility)

    if givenVisibility != "PRIVATE" && givenVisibility != "PUBLIC" && givenVisibility != "INTERFACE" {
        var err error
        if headerOnly {
            err = errors.LinkerVisibilityError{
                From:            from,
                To:              to,
                GivenVisibility: givenVisibility,
                Err:             goerr.New("header only => INTERFACE visibility will be used"),
            }
            givenVisibility = "INTERFACE"
        } else {
            err = errors.LinkerVisibilityError{
                From:            from,
                To:              to,
                GivenVisibility: givenVisibility,
                Err:             goerr.New("PRIVATE visibility will be used"),
            }
            givenVisibility = "PRIVATE"
        }

        log.QueueWriteln(queue, log.WARN, nil, err.Error())
    } else {
        if headerOnly && givenVisibility != "INTERFACE" {
            err := errors.LinkerVisibilityError{
                From:            from,
                To:              to,
                GivenVisibility: givenVisibility,
                Err:             goerr.New("target is header only => INTERFACE visibility will be used"),
            }

            givenVisibility = "INTERFACE"
            log.QueueWriteln(queue, log.WARN, nil, err.Error())
        }
    }

    return givenVisibility
}
