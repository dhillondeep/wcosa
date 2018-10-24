package dependencies

import (
    "wio/internal/cmd/run/cmake"
    "wio/internal/types"
    "wio/pkg/npm/resolve"
    "wio/pkg/util"
)

type definitionsInfo struct {
    name         string
    globalsGiven []string
    otherGiven   []string
    globals      map[string][]string
    required     map[string][]string
    optional     map[string][]string
    ingest       map[string][]string
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

    globalPrivate, err := fillDefinition(info.globalsGiven, info.globals[types.Private], false)
    if err != nil {
        return nil, util.Error(err.Error(), "global", info.name)
    }
    globalPublic, err := fillDefinition(info.globalsGiven, info.globals[types.Public], false)
    if err != nil {
        return nil, util.Error(err.Error(), "global", info.name)
    }

    all[types.Private] = util.AppendIfMissing(all[types.Private], globalPrivate)
    all[types.Public] = util.AppendIfMissing(all[types.Public], globalPublic)

    if !info.singleton {
        if len(info.required[types.Private]) > 0 || len(info.required[types.Public]) > 0 {
            requiredPrivate, err := fillDefinition(info.otherGiven, info.required[types.Private], false)
            if err != nil {
                return nil, util.Error(err.Error(), "required", info.name)
            }
            requiredPublic, err := fillDefinition(info.otherGiven, info.required[types.Public], false)
            if err != nil {
                return nil, util.Error(err.Error(), "required", info.name)
            }
            all[types.Private] = util.AppendIfMissing(all[types.Private], requiredPrivate)
            all[types.Public] = util.AppendIfMissing(all[types.Public], requiredPublic)
        }

        if len(info.optional[types.Private]) > 0 || len(info.optional[types.Public]) > 0 {
            optionalPrivate, _ := fillDefinition(info.otherGiven, info.optional[types.Private], true)
            optionalPublic, _ := fillDefinition(info.otherGiven, info.optional[types.Public], true)

            all[types.Private] = util.AppendIfMissing(all[types.Private], optionalPrivate)
            all[types.Public] = util.AppendIfMissing(all[types.Public], optionalPublic)
        }
    }

    all[types.Private] = util.AppendIfMissing(all[types.Private], info.ingest[types.Private])
    all[types.Public] = util.AppendIfMissing(all[types.Public], info.ingest[types.Public])

    return all, nil
}

func resolveTree(i *resolve.Info, currNode *resolve.Node, parentTarget *Target, targetSet *TargetSet,
    librarySet *TargetSet, globalFlags, globalDefinitions []string, parentGiven *parentGivenInfo) error {
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

    defIngest := map[string][]string{
        types.Private: definitions.GetIngest().GetPrivate(),
        types.Public:  definitions.GetIngest().GetPublic(),
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
            ingest:       defIngest,
            singleton:    definitions.IsSingleton(),
        }); err != nil {
        return err
    }

    // flags
    currTarget.Flags = util.AppendIfMissing(pkg.Config.GetInfo().GetOptions().GetFlags(), parentGiven.flags)

    targetSet.Add(currTarget, false)
    targetSet.Link(parentTarget, currTarget, &TargetLinkInfo{
        Visibility: parentGiven.linkVisibility,
        Flags:      parentGiven.linkFlags,
    })

    // go over all the shared libraries for this package
    pkgConfig, err := types.ReadWioConfig(pkg.Path)
    if err != nil {
        return err
    }

    for name, library := range pkgConfig.GetLibraries() {
        libraryTarget := &Target{
            Name:       name,
            ParentPath: currTarget.Path,
            Library:    library,
        }

        librarySet.Add(libraryTarget, true)
        librarySet.Link(currTarget, libraryTarget, &TargetLinkInfo{
            Visibility: library.GetLinkerVisibility(),
            Flags:      library.GetLinkerFlags(),
        })
    }

    for _, dep := range currNode.Dependencies {
        if configDependency, exists := pkg.Config.GetDependencies()[dep.Name]; !exists {
            return util.Error("%s@%s dependency's information is wrong in wio.yml", dep.Name,
                dep.ResolvedVersion.Str())
        } else {
            // resolve placeholders
            tDef := util.AppendIfMissing(currTarget.Definitions[types.Private], currTarget.Definitions[types.Public])
            parentDefinitions, err := fillPlaceholders(tDef, configDependency.GetDefinitions())
            if err != nil {
                return util.Error(err.Error(), currTarget.Name)
            }

            parentInfo := &parentGivenInfo{
                flags:          configDependency.GetCompileFlags(),
                definitions:    parentDefinitions,
                linkVisibility: configDependency.GetVisibility(),
                linkFlags:      configDependency.GetLinkerFlags(),
            }
            if err := resolveTree(i, dep, currTarget, targetSet, librarySet, globalFlags, globalDefinitions, parentInfo); err != nil {
                return err
            }
        }

    }

    return nil
}
