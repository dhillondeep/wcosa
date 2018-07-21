package dependencies

import (
    "path/filepath"
    "strings"
    "wio/cmd/wio/commands/run/cmake"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/toolchain/npm/resolve"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/utils/template"
)

const (
    MainTarget = "${TARGET_NAME}"
)

var libraryStrings = map[string]map[bool]string{
    "avr":    {false: cmake.AvrLibrary, true: cmake.AvrHeader},
    "native": {false: cmake.DesktopLibrary, true: cmake.DesktopHeader},
}

// This creates CMake dependency string using build targets that will be used to link dependencies
func GenerateCMakeDependencies(cmakePath string, platform string, targets *TargetSet) error {
    cmakeStrings := make([]string, 0, 256)

    for target := range targets.TargetIterator() {
        finalString := libraryStrings[platform][target.HeaderOnly]

        private := types.Private
        public := types.Public

        if target.HeaderOnly {
            private = types.Interface
            public = types.Interface
        }

        finalString = template.Replace(finalString, map[string]string{
            "DEPENDENCY_PATH":                 filepath.ToSlash(target.Path),
            "DEPENDENCY_NAME":                 target.Name + "__" + target.Version,
            "DEPENDENCY_FLAGS":                strings.Join(target.Flags, " "),
            "PRIVATE_DEFINITIONS_DEFINITIONS": strings.Join(target.Definitions[types.Private], " "),
            "PRIVATE_DEFINITIONS_VISIBILITY":  private,
            "PUBLIC_DEFINITIONS_DEFINITIONS":  strings.Join(target.Definitions[types.Public], " "),
            "PUBLIC_DEFINITIONS_VISIBILITY":   public,
        })
        cmakeStrings = append(cmakeStrings, finalString+"\n")
    }

    for link := range targets.LinkIterator() {
        if link.From.HeaderOnly && link.LinkInfo.Visibility != "INTERFACE" {
            return errors.Stringf("header library %s@%s must have INTERFACE link visibility", link.From.Name,
                link.From.Version)
        } else if strings.Trim(link.LinkInfo.Visibility, " ") == "" {
            link.LinkInfo.Visibility = "PRIVATE"
        }

        linkerName := link.From.Name + "__" + link.From.Version
        if link.From.Name == MainTarget {
            linkerName = MainTarget
        }

        finalString := template.Replace(cmake.LinkString, map[string]string{
            "LINKER_NAME":     linkerName,
            "DEPENDENCY_NAME": link.To.Name + "__" + link.To.Version,
            "LINK_VISIBILITY": link.LinkInfo.Visibility,
        })
        cmakeStrings = append(cmakeStrings, finalString)
    }

    return io.NormalIO.WriteFile(cmakePath, []byte(strings.Join(cmakeStrings, "\n")))
}

// Scans the dependency tree and creates build targets that will be converted into CMake targets
func CreateBuildTargets(projectDir string, target types.Target) (*TargetSet, error) {
    targetSet := NewTargetSet()

    i := resolve.NewInfo(projectDir)
    config, err := utils.ReadWioConfig(projectDir)
    if err != nil {
        return nil, err
    }

    err = i.ResolveRemote(config)
    if err != nil {
        return nil, err
    }

    if config.GetType() == constants.APP {
        for _, dep := range i.GetRoot().Dependencies {
            var configDependency types.Dependency
            var exists bool

            if configDependency, exists = config.GetDependencies()[dep.Name]; !exists {
                return nil, errors.Stringf("%s@%s dependency is invalid and information is wrong in wio.yml",
                    dep.Name, dep.ResolvedVersion.Str())
            }

            parentInfo := &parentGivenInfo{
                flags:          configDependency.GetFlags(),
                definitions:    configDependency.GetDefinitions(),
                linkVisibility: configDependency.GetVisibility(),
            }

            // all direct dependencies will link to the main target
            err := resolveTree(i, dep, &Target{
                Name: MainTarget,
            }, targetSet, target.GetFlags().GetGlobal(),
                target.GetDefinitions().GetGlobal(), parentInfo)
            if err != nil {
                return nil, err
            }
        }
    } else {
        parentInfo := &parentGivenInfo{
            flags:          target.GetFlags().GetPackage(),
            definitions:    target.GetDefinitions().GetPackage(),
            linkVisibility: "PRIVATE",
        }

        // this package will link to the main target
        err := resolveTree(i, i.GetRoot(), &Target{
            Name: MainTarget,
        }, targetSet, target.GetFlags().GetGlobal(),
            target.GetDefinitions().GetGlobal(), parentInfo)
        if err != nil {
            return nil, err
        }
    }

    return targetSet, nil
}
