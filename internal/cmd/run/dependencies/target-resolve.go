package dependencies

import (
    "wio/internal/cmd/run/cmake"
    "wio/internal/types"
    "wio/pkg/npm/resolve"
    "wio/pkg/npm/semver"
    "wio/pkg/util"
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
        return nil, util.Error(err.Error(), "global", info.name)
    }
    globalPublic, err := fillDefinition(info.globalsGiven, info.globals[types.Public])
    if err != nil {
        return nil, util.Error(err.Error(), "global", info.name)
    }

    all[types.Private] = util.AppendIfMissing(all[types.Private], globalPrivate)
    all[types.Public] = util.AppendIfMissing(all[types.Public], globalPublic)

    if !info.singleton {
        if len(info.required[types.Private]) > 0 || len(info.required[types.Public]) > 0 {
            requiredPrivate, err := fillDefinition(info.otherGiven, info.required[types.Private])
            if err != nil {
                return nil, util.Error(err.Error(), "required", info.name)
            }
            requiredPublic, err := fillDefinition(info.otherGiven, info.required[types.Public])
            if err != nil {
                return nil, util.Error(err.Error(), "required", info.name)
            }
            all[types.Private] = util.AppendIfMissing(all[types.Private], requiredPrivate)
            all[types.Public] = util.AppendIfMissing(all[types.Public], requiredPublic)
        } else if len(info.optional[types.Private]) > 0 || len(info.optional[types.Public]) > 0 {
            optionalPrivate := fillOptionalDefinition(info.otherGiven, info.optional[types.Private])
            optionalPublic := fillOptionalDefinition(info.otherGiven, info.optional[types.Public])

            all[types.Private] = util.AppendIfMissing(all[types.Private], optionalPrivate)
            all[types.Public] = util.AppendIfMissing(all[types.Public], optionalPublic)
        }
    } else {
        all[types.Private] = util.AppendIfMissing(all[types.Private], info.optional[types.Private])
        all[types.Public] = util.AppendIfMissing(all[types.Public], info.optional[types.Public])
    }

    return all, nil
}

func resolveTree(i *resolve.Info, currNode *resolve.Node, parentTarget *Target, targetSet *TargetSet,
    sharedSet *TargetSet, globalFlags, globalDefinitions []string, parentGiven *parentGivenInfo) error {
    var err error

    pkg, err := i.GetPkg(currNode.Name, currNode.ResolvedVersion.Str())
    if err != nil {
        return err
    } else if pkg == nil {
        return util.Error("%s dependency does not exist. Check vendor or wio install", currNode.Name)
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
    currTarget.Flags = util.AppendIfMissing(pkg.Config.GetInfo().GetOptions().GetFlags(), parentGiven.flags)

    targetSet.Add(currTarget)
    targetSet.Link(parentTarget, currTarget, &TargetLinkInfo{
        Visibility: parentGiven.linkVisibility,
        Flags:      parentGiven.linkFlags,
    })

    // go over all the shared libraries for this package
    pkgConfig, err := types.ReadWioConfig(pkg.Path)
    if err != nil {
        return err
    }

    wioVerStr := pkgConfig.GetInfo().GetOptions().GetWioVersion()
    wioVer := semver.Parse(wioVerStr)
    if wioVer == nil {
        return util.Error("Invalid wio version in wio.yml: %s", wioVerStr)
    }
    if wioVer.Ge(semver.Parse("0.5.0")) {
        for name, shared := range pkgConfig.GetLibraries() {
            sharedTarget := &Target{
                Name:              name,
                Path:              shared.GetPath(),
                SharedIncludePath: shared.GetIncludePath(),
                ParentPath:        currTarget.Path,
            }

            sharedSet.Add(sharedTarget)
            sharedSet.Link(currTarget, sharedTarget, &TargetLinkInfo{
                Visibility: shared.GetLinkerVisibility(),
                Flags:      shared.GetLinkerFlags(),
            })
        }
    }

    for _, dep := range currNode.Dependencies {
        if configDependency, exists := pkg.Config.GetDependencies()[dep.Name]; !exists {
            return util.Error("%s@%s dependency's information is wrong in wio.yml", dep.Name,
                dep.ResolvedVersion.Str())
        } else {
            // resolve placeholders
            parentFlags, err := fillPlaceholders(currTarget.Flags, configDependency.GetCompileFlags())
            if err != nil {
                return util.Error(err.Error(), currTarget.Name)
            }

            tDef := util.AppendIfMissing(currTarget.Definitions[types.Private], currTarget.Definitions[types.Public])
            parentDefinitions, err := fillPlaceholders(tDef, configDependency.GetDefinitions())
            if err != nil {
                return util.Error(err.Error(), currTarget.Name)
            }

            parentInfo := &parentGivenInfo{
                flags:          parentFlags,
                definitions:    parentDefinitions,
                linkVisibility: configDependency.GetVisibility(),
            }
            if err := resolveTree(i, dep, currTarget, targetSet, sharedSet, globalFlags, globalDefinitions, parentInfo); err != nil {
                return err
            }
        }

    }

    return nil
}
