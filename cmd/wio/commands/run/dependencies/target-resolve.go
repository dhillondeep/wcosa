package dependencies

import (
    "fmt"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/log"
    "wio/cmd/wio/toolchain/npm/resolve"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
)

type flagsDefinitionsInfo struct {
    name         string
    isFlag       bool
    globalsGiven []string
    otherGiven   []string
    globals      []string
    required     []string
    included     []string
    onlyGlobal   bool
    onlyRequired bool
}

type parentGivenInfo struct {
    flags          []string
    definitions    []string
    linkVisibility string
    linkFlags      []string
}

func fillFlagsDefinitions(info flagsDefinitionsInfo) ([]string, error) {
    var all []string

    global, err := fillGlobal(info.globalsGiven, info.globals)
    if err != nil {
        return nil, errors.Stringf(err.Error(), info.name)
    } else {
        all = utils.AppendIfMissing(all, global)
    }

    if !info.onlyGlobal && len(info.required) > 0 {
        required, other, err := fillRequired(info.otherGiven, info.required)
        if err != nil {
            fmt.Println(err.Error())
            return nil, errors.Stringf(err.Error(), info.name)
        } else {
            all = utils.AppendIfMissing(all, required)
            all = utils.AppendIfMissing(all, other)
        }
    } else {
        acceptName := "flags"
        if !info.isFlag {
            acceptName = "definitions"
        }

        if info.onlyGlobal && info.onlyRequired {
            log.Warnln(fmt.Sprintf("%s is set to accept only global %s and only required %s. "+
                "Ignoring required %s...", info.name, acceptName, acceptName, acceptName))
        } else if info.onlyRequired {
            log.Warnln(fmt.Sprintf("%s only accepts global %s but required flags are also requested. "+
                "Ignoring required %s...", info.name, acceptName, acceptName))
        }

        all = utils.AppendIfMissing(all, info.otherGiven)
    }

    all = utils.AppendIfMissing(all, info.included)
    return all, nil
}

func resolveTree(i *resolve.Info, currNode *resolve.Node, parentTarget *Target, targetSet *TargetSet,
    globalFlags, globalDefinitions []string, parentGiven *parentGivenInfo) error {
    var err error

    pkg, err := i.GetPkg(currNode.Name, currNode.ResolvedVersion.Str())
    if err != nil {
        return err
    } else if pkg == nil {
        return errors.Stringf("%s dependency does not exist. Check vendor or wio install", currNode.Name)
    }

    currTarget := &Target{
        Name:                  pkg.Config.Name(),
        Version:               pkg.Config.Version(),
        Path:                  pkg.Path,
        FromVendor:            pkg.Vendor,
        HeaderOnly:            pkg.Config.GetMainTag().GetCompileOptions().IsHeaderOnly(),
        FlagsVisibility:       pkg.Config.MainTag.Flags.Visibility,
        DefinitionsVisibility: pkg.Config.MainTag.Definitions.Visibility,
    }

    // resolve flags
    if currTarget.Flags, err = fillFlagsDefinitions(flagsDefinitionsInfo{
        name:         currTarget.Name + "__" + currTarget.Version,
        globalsGiven: globalFlags,
        otherGiven:   parentGiven.flags,
        globals:      pkg.Config.MainTag.Flags.GlobalFlags,
        required:     pkg.Config.MainTag.Flags.RequiredFlags,
        included:     pkg.Config.MainTag.Flags.IncludedFlags,
        onlyGlobal:   pkg.Config.MainTag.Flags.AllowOnlyGlobalFlags,
        onlyRequired: pkg.Config.MainTag.Flags.AllowOnlyRequiredFlags,
    }); err != nil {
        return err
    }

    // resolve definitions
    if currTarget.Definitions, err = fillFlagsDefinitions(flagsDefinitionsInfo{
        name:         currTarget.Name + "__" + currTarget.Version,
        globalsGiven: globalDefinitions,
        otherGiven:   parentGiven.definitions,
        globals:      pkg.Config.MainTag.Definitions.GlobalDefinitions,
        required:     pkg.Config.MainTag.Definitions.RequiredDefinitions,
        included:     pkg.Config.MainTag.Definitions.IncludedDefinitions,
        onlyGlobal:   pkg.Config.MainTag.Definitions.AllowOnlyGlobalDefinitions,
        onlyRequired: pkg.Config.MainTag.Definitions.AllowOnlyRequiredDefinitions,
    }); err != nil {
        return err
    }

    targetSet.Add(currTarget)

    targetSet.Link(parentTarget, currTarget, &TargetLinkInfo{
        Visibility: parentGiven.linkVisibility,
    })

    for _, dep := range currNode.Dependencies {
        configDependency := &types.DependencyTag{}
        var exists bool

        if configDependency, exists = pkg.Config.GetDependencies()[dep.Name]; !exists {
            return errors.Stringf("%s@%s dependency is invalid and information is wrong in wio.yml",
                dep.Name, dep.ResolvedVersion.Str())
        }

        // resolve placeholders
        parentFlags, err := fillPlaceholders(currTarget.Flags, configDependency.Flags)
        if err != nil {
            return errors.Stringf(err.Error(), currTarget.Name)
        }
        parentDefinitions, err := fillPlaceholders(currTarget.Definitions, configDependency.Definitions)
        if err != nil {
            return errors.Stringf(err.Error(), currTarget.Name)
        }

        parentInfo := &parentGivenInfo{
            flags:          parentFlags,
            definitions:    parentDefinitions,
            linkVisibility: configDependency.LinkVisibility,
        }

        if err = resolveTree(i, dep, currTarget, targetSet, globalFlags, globalDefinitions, parentInfo); err != nil {
            return err
        }
    }

    return nil
}
