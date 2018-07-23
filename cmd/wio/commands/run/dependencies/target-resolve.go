package dependencies

import (
    "wio/cmd/wio/commands/run/cmake"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/toolchain/npm/resolve"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
)

type definitionsInfo struct {
    name         string
    globalsGiven []string
    otherGiven   []string
    globals      map[string][]string
    required     map[string][]string
    optional     map[string][]string
    singleton    bool
}

type parentGivenInfo struct {
    flags          []string
    definitions    []string
    linkVisibility string
    linkFlags      []string
}

func fillDefinitions(info definitionsInfo) (map[string][]string, error) {
    var all = map[string][]string{}

    globalPrivate, err := fillDefinition(info.globalsGiven, info.globals[types.Private])
    if err != nil {
        return nil, errors.Stringf(err.Error(), "global", info.name)
    }
    globalPublic, err := fillDefinition(info.globalsGiven, info.globals[types.Public])
    if err != nil {
        return nil, errors.Stringf(err.Error(), "global", info.name)
    }

    all[types.Private] = utils.AppendIfMissing(all[types.Private], globalPrivate)
    all[types.Public] = utils.AppendIfMissing(all[types.Public], globalPublic)

    if !info.singleton && (len(info.required[types.Private]) > 0 || len(info.required[types.Public]) > 0) {
        requiredPrivate, err := fillDefinition(info.otherGiven, info.required[types.Private])
        if err != nil {
            return nil, errors.Stringf(err.Error(), "required", info.name)
        }
        requiredPublic, err := fillDefinition(info.otherGiven, info.required[types.Public])
        if err != nil {
            return nil, errors.Stringf(err.Error(), "required", info.name)
        }
        optionalPrivate, err := fillDefinition(info.otherGiven, info.optional[types.Private])
        if err != nil {
            return nil, errors.Stringf(err.Error(), "other", info.name)
        }
        optionalPublic, err := fillDefinition(info.otherGiven, info.optional[types.Public])
        if err != nil {
            return nil, errors.Stringf(err.Error(), "other", info.name)
        }

        all[types.Private] = utils.AppendIfMissing(all[types.Private], requiredPrivate)
        all[types.Public] = utils.AppendIfMissing(all[types.Public], requiredPublic)
        all[types.Private] = utils.AppendIfMissing(all[types.Private], optionalPrivate)
        all[types.Public] = utils.AppendIfMissing(all[types.Public], optionalPublic)
    }

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

    cxxStandard, cStandard, err := cmake.GetStandard(pkg.Config.GetInfo().GetOptions().GetStandard())
    if err != nil {
        return err
    }

    currTarget := &Target{
        Name:        pkg.Config.GetName(),
        Version:     pkg.Config.GetVersion(),
        Path:        pkg.Path,
        FromVendor:  pkg.Vendor,
        HeaderOnly:  pkg.Config.GetInfo().GetOptions().GetIsHeaderOnly(),
        CXXStandard: cxxStandard,
        CStandard:   cStandard,
    }

    definitions := pkg.Config.GetInfo().GetDefinitions()
    defGlobals := map[string][]string{
        types.Private: definitions.GetGlobal().GetPrivate(),
        types.Public:  definitions.GetGlobal().GetPublic(),
    }

    defRequired := map[string][]string{
        types.Private: definitions.GetRequired().GetPrivate(),
        types.Public:  definitions.GetRequired().GetPublic(),
    }

    defOptional := map[string][]string{
        types.Private: definitions.GetOptional().GetPrivate(),
        types.Public:  definitions.GetOptional().GetPublic(),
    }

    // definitions
    if currTarget.Definitions, err = fillDefinitions(
        definitionsInfo{
            name:         currTarget.Name + "__" + currTarget.Version,
            globalsGiven: globalDefinitions,
            otherGiven:   parentGiven.definitions,
            globals:      defGlobals,
            required:     defRequired,
            optional:     defOptional,
            singleton:    definitions.IsSingleton(),
        }); err != nil {
        return err
    }

    // flags
    currTarget.Flags = utils.AppendIfMissing(pkg.Config.GetInfo().GetOptions().GetFlags(), parentGiven.flags)

    targetSet.Add(currTarget)
    targetSet.Link(parentTarget, currTarget, &TargetLinkInfo{
        Visibility: parentGiven.linkVisibility,
        Flags:      parentGiven.linkFlags,
    })

    for _, dep := range currNode.Dependencies {
        if configDependency, exists := pkg.Config.GetDependencies()[dep.Name]; !exists {
            return errors.Stringf("%s@%s dependency's information is wrong in wio.yml", dep.Name,
                dep.ResolvedVersion.Str())
        } else {
            // resolve placeholders
            parentFlags, err := fillPlaceholders(currTarget.Flags, configDependency.GetCompileFlags())
            if err != nil {
                return errors.Stringf(err.Error(), currTarget.Name)
            }

            tDef := utils.AppendIfMissing(currTarget.Definitions[types.Private], currTarget.Definitions[types.Public])
            parentDefinitions, err := fillPlaceholders(tDef, configDependency.GetDefinitions())
            if err != nil {
                return errors.Stringf(err.Error(), currTarget.Name)
            }

            parentInfo := &parentGivenInfo{
                flags:          parentFlags,
                definitions:    parentDefinitions,
                linkVisibility: configDependency.GetVisibility(),
            }

            return resolveTree(i, dep, currTarget, targetSet, globalFlags, globalDefinitions, parentInfo)
        }

    }

    return nil
}
