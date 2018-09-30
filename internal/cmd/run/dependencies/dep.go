package dependencies

import (
    "io/ioutil"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "wio/internal/cmd/run/cmake"
    "wio/internal/constants"
    "wio/internal/types"
    "wio/pkg/npm/resolve"
    "wio/pkg/util"
    "wio/pkg/util/template"
)

const (
    MainTarget = "${TARGET_NAME}"
)

var libraryStrings = map[string]map[bool]string{
    "avr":    {false: cmake.AvrLibrary, true: cmake.AvrHeader},
    "native": {false: cmake.DesktopLibrary, true: cmake.DesktopHeader},
}

// This creates CMake dependency string using build targets that will be used to link dependencies
func GenerateCMakeDependencies(directory string, cmakePath string, platform string, targets *TargetSet,
    shared *TargetSet) error {
    cmakeStrings := make([]string, 0, 256)

    for shared := range shared.TargetIterator() {
        finalString := cmake.SharedLibrary

        finalString = template.Replace(finalString, map[string]string{
            "SHARED_LIB_NAME":         shared.Name,
            "SHARED_LIB_NAME_ORG":     strings.Split(shared.Name, "__")[0],
            "SHARED_LIB_PATH":         strings.Replace(shared.Path, "$(PROJECT_PATH)", directory, -1),
            "SHARED_LIB_INCLUDE_PATH": strings.Replace(shared.SharedIncludePath, "$(PROJECT_PATH)", directory, -1),
        })

        cmakeStrings = append(cmakeStrings, finalString+"\n")
    }

    for target := range targets.TargetIterator() {
        finalString := libraryStrings[platform][target.HeaderOnly]

        finalString = template.Replace(finalString, map[string]string{
            "DEPENDENCY_PATH":     filepath.ToSlash(target.Path),
            "DEPENDENCY_NAME":     target.Name + "__" + target.Version,
            "DEPENDENCY_FLAGS":    strings.Join(target.Flags, " "),
            "PRIVATE_DEFINITIONS": strings.Join(target.Definitions[types.Private], " "),
            "PUBLIC_DEFINITIONS":  strings.Join(target.Definitions[types.Public], " "),
            "CXX_STANDARD":        target.CXXStandard,
            "C_STANDARD":          target.CStandard,
        })
        cmakeStrings = append(cmakeStrings, finalString+"\n")
    }

    for sharedLink := range shared.LinkIterator() {
        if sharedLink.From.HeaderOnly {
            sharedLink.LinkInfo.Visibility = types.Interface
        } else if strings.Trim(sharedLink.LinkInfo.Visibility, " ") == "" {
            sharedLink.LinkInfo.Visibility = types.Private
        }

        linkerName := sharedLink.From.Name + "__" + sharedLink.From.Version
        if sharedLink.From.Name == MainTarget {
            linkerName = MainTarget
        }

        finalString := template.Replace(cmake.LinkString, map[string]string{
            "LINKER_NAME":     linkerName,
            "DEPENDENCY_NAME": "${LIB_" + sharedLink.To.Name + "}",
            "LINK_VISIBILITY": sharedLink.LinkInfo.Visibility,
            "LINKER_FLAGS":    strings.Join(sharedLink.LinkInfo.Flags, " "),
        })
        cmakeStrings = append(cmakeStrings, finalString)
    }

    for link := range targets.LinkIterator() {
        if link.From.HeaderOnly {
            link.LinkInfo.Visibility = types.Interface
        } else if strings.Trim(link.LinkInfo.Visibility, " ") == "" {
            link.LinkInfo.Visibility = types.Private
        }

        linkerName := link.From.Name + "__" + link.From.Version
        if link.From.Name == MainTarget {
            linkerName = MainTarget
        }

        finalString := template.Replace(cmake.LinkString, map[string]string{
            "LINKER_NAME":     linkerName,
            "DEPENDENCY_NAME": link.To.Name + "__" + link.To.Version,
            "LINK_VISIBILITY": link.LinkInfo.Visibility,
            "LINKER_FLAGS":    strings.Join(link.LinkInfo.Flags, " "),
        })
        cmakeStrings = append(cmakeStrings, finalString)
    }
    fileContents := []byte(strings.Join(cmakeStrings, "\n"))
    return ioutil.WriteFile(cmakePath, fileContents, os.ModePerm)
}

// Scans the dependency tree and creates build targets that will be converted into CMake targets
func CreateBuildTargets(projectDir string, target types.Target) (*TargetSet, *TargetSet, error) {
    targetSet := NewTargetSet()
    sharedTargetSet := NewTargetSet()

    i := resolve.NewInfo(projectDir)
    config, err := types.ReadWioConfig(projectDir)
    if err != nil {
        return nil, nil, err
    }

    err = i.ResolveRemote(config)
    if err != nil {
        return nil, nil, err
    }

    if config.GetType() == constants.App {
        parentTarget := &Target{
            Name: MainTarget,
        }

        // link all the libraries for the application
        for name, shared := range config.GetLibraries() {
            currTarget := &Target{
                Name:              name,
                Path:              shared.GetPath(),
                SharedIncludePath: shared.GetIncludePath(),
            }

            sharedTargetSet.Add(currTarget)
            sharedTargetSet.Link(parentTarget, currTarget, &TargetLinkInfo{
                Visibility: shared.GetLinkerVisibility(),
                Flags:      shared.GetLinkerFlags(),
            })
        }

        for _, dep := range i.GetRoot().Dependencies {
            var configDependency types.Dependency
            var exists bool

            if configDependency, exists = config.GetDependencies()[dep.Name]; !exists {
                return nil, nil, util.Error("%s@%s dependency is invalid and information is wrong in wio.yml",
                    dep.Name, dep.ResolvedVersion.Str())
            }

            parentInfo := &parentGivenInfo{
                flags:          configDependency.GetCompileFlags(),
                definitions:    configDependency.GetDefinitions(),
                linkVisibility: configDependency.GetVisibility(),
                linkFlags:      configDependency.GetLinkerFlags(),
            }

            // all direct dependencies will link to the main target
            err := resolveTree(i, dep, parentTarget, targetSet, sharedTargetSet, target.GetFlags().GetGlobal(),
                target.GetDefinitions().GetGlobal(), parentInfo)
            if err != nil {
                return nil, nil, err
            }
        }
    } else {
        parentInfo := &parentGivenInfo{
            flags:          target.GetFlags().GetPackage(),
            definitions:    target.GetDefinitions().GetPackage(),
            linkVisibility: "PRIVATE",
        }

        // separate normal flags with linker flags
        linkerRegex := regexp.MustCompile(`-l((\s+[A-Za-z]+)|([A-Za-z]+))`)

        var compileFlags []string
        var linkerFlags []string

        for _, flag := range parentInfo.flags {
            if linkerRegex.MatchString(flag) {
                flag = strings.Trim(strings.Replace(flag, "-l", "", 1), " ")
                linkerFlags = append(linkerFlags, flag)
            } else {
                compileFlags = append(compileFlags, flag)
            }
        }

        parentInfo.flags = compileFlags
        parentInfo.linkFlags = linkerFlags

        // this package will link to the main target
        err := resolveTree(i, i.GetRoot(), &Target{
            Name: MainTarget,
        }, targetSet, sharedTargetSet, target.GetFlags().GetGlobal(),
            target.GetDefinitions().GetGlobal(), parentInfo)
        if err != nil {
            return nil, nil, err
        }
    }

    return targetSet, sharedTargetSet, nil
}
