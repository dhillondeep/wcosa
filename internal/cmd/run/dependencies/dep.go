package dependencies

import (
	"fmt"
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
func GenerateCMakeDependencies(cmakePath string, platform string, dependencies *TargetSet, libraries *TargetSet) error {
	cmakeStrings := make([]string, 0, 256)
	cmakeStrings = append(cmakeStrings, cmake.ImportedTargetsMacro+"\n")

	// joins the slice or if empty puts a cmake comment
	cmakeSliceJoin := func(slice []string, message string) string {
		if len(slice) > 0 {
			return strings.Join(slice, " ")
		} else {
			return "# " + message
		}
	}

	// inserts the string or if empty puts a cmake comment
	cmakeString := func(str string, message string) string {
		if str == "" {
			return str
		} else {
			return "# " + message
		}
	}

	// create cmake targets for libraries
	for library := range libraries.TargetIterator() {
		var finalString string
		libOriginalName := GetOriginalName(library, true)
		configLibrary := library.Library

		for variableName, variableValue := range configLibrary.GetVariables() {
			finalString += fmt.Sprintf("set(%s %s)\n", variableName, variableValue)
		}

		// Find<LIB_NAME>.cmake file exists for the library
		if configLibrary.IsCmakePackage() {
			if len(configLibrary.GetPath()) > 0 {
				finalString += fmt.Sprintf("list(APPEND CMAKE_MODULE_PATH %s CACHE FORCE)", configLibrary.GetPath())
				finalString += "\n" + cmake.LibraryPackageFind
			} else {
				finalString += cmake.LibraryPackageFind
			}

			finalString = template.Replace(finalString, map[string]string{
				"LIB_VERSION": configLibrary.GetVersion(),
				"LIB_REQUIRED_COMPONENTS": cmakeSliceJoin(configLibrary.GetRequiredComponents(),
					"no required components"),
				"LIB_OPTIONAL_COMPONENTS": cmakeSliceJoin(configLibrary.GetOptionalComponents(),
					"no optional components"),
			})
		} else {
			finalString += cmake.LibraryFind
		}

		pathHintsCMake := func(prefix string) string {
			if len(configLibrary.GetLibPath()) > 0 {
				return prefix + " " + cmakeSliceJoin(configLibrary.GetLibPath(), "no path/hint provided")
			} else {
				return "# no path/hint provided"
			}
		}

		finalString = template.Replace(finalString, map[string]string{
			"LIB_NAME_VAR": library.Name,
			"LIB_NAME":     libOriginalName,
			"LIB_PATHS":    pathHintsCMake("PATHS"),
			"LIB_HINTS":    pathHintsCMake("HINTS"),
			"LIB_REQUIRED": func() string {
				if configLibrary.IsRequired() {
					return "REQUIRED"
				} else {
					return cmakeString("", "not required")
				}
			}(),
		})

		cmakeStrings = append(cmakeStrings, finalString+"\n")
	}

	// create cmake targets for dependencies
	for dependency := range dependencies.TargetIterator() {
		finalString := libraryStrings[platform][dependency.HeaderOnly]

		finalString = template.Replace(finalString, map[string]string{
			"DEPENDENCY_PATH":  filepath.ToSlash(dependency.Path),
			"DEPENDENCY_NAME":  dependency.Name,
			"DEPENDENCY_FLAGS": cmakeSliceJoin(dependency.Flags, "no flags provided"),
			"PRIVATE_DEFINITIONS": cmakeSliceJoin(dependency.Definitions[types.Private],
				"no private definitions provided"),
			"PUBLIC_DEFINITIONS": cmakeSliceJoin(dependency.Definitions[types.Public],
				"no public definitions provided"),
			"CXX_STANDARD": dependency.CXXStandard,
			"C_STANDARD":   dependency.CStandard,
		})
		cmakeStrings = append(cmakeStrings, finalString+"\n")
	}

	for libraryLink := range libraries.LinkIterator() {
		var finalString string
		configLibrary := libraryLink.To.Library
		toOriginalName := GetOriginalName(libraryLink.To, true)

		if configLibrary.IsCmakePackage() && configLibrary.UseImportedTargets() {
			finalString = template.Replace(cmake.LibraryLinkImportedTargets, map[string]string{
				"LIB_NAME": toOriginalName,
				"LIB_REQUIRED_COMPONENTS": func() string {
					var requiredComps []string
					if len(configLibrary.GetRequiredComponents()) <= 0 {
						requiredComps = append(requiredComps, toOriginalName, strings.ToLower(toOriginalName),
							strings.ToUpper(toOriginalName), strings.Title(toOriginalName))
					} else {
						requiredComps = configLibrary.GetRequiredComponents()
					}

					return strings.Join(requiredComps, " ")
				}(),
				"LIB_OPTIONAL_COMPONENTS": strings.Join(configLibrary.GetOptionalComponents(), " "),
			})
		} else {
			finalString = cmake.LibraryLink
		}

		if libraryLink.From.HeaderOnly {
			libraryLink.LinkInfo.Visibility = types.Interface
		} else if strings.Trim(libraryLink.LinkInfo.Visibility, " ") == "" {
			libraryLink.LinkInfo.Visibility = types.Private
		}

		finalString = template.Replace(finalString, map[string]string{
			"LINK_FROM":          libraryLink.From.Name,
			"INCLUDE_VISIBILITY": libraryLink.LinkInfo.Visibility,
			"LIB_INCLUDE_PATHS": func() string {
				if configLibrary.IsCmakePackage() {
					if strings.Trim(configLibrary.GetIncludesTag(), " ") != "" {
						return fmt.Sprintf("${%s}", configLibrary.GetIncludesTag())
					} else {
						return fmt.Sprintf("${%s_INCLUDE_DIR} ${%s_INCLUDE_DIR} ${%s_INCLUDE_DIR} ${%s_INCLUDE_DIR}",
							toOriginalName, strings.ToLower(toOriginalName),
							strings.ToUpper(toOriginalName), strings.Title(toOriginalName))
					}
				} else {
					return cmakeSliceJoin(configLibrary.GetIncludePath(), "no include paths provided")
				}
			}(),
			"LINK_VISIBILITY": libraryLink.LinkInfo.Visibility,
			"LINK_TO": func() string {
				if configLibrary.IsCmakePackage() {
					if strings.Trim(configLibrary.GetLibrariesTag(), " ") != "" {
						return fmt.Sprintf("${%s}", configLibrary.GetLibrariesTag())
					} else {
						return fmt.Sprintf("${%s_LIBRARIES} ${%s_LIBRARIES} ${%s_LIBRARIES} ${%s_LIBRARIES}",
							toOriginalName, strings.ToLower(toOriginalName),
							strings.ToUpper(toOriginalName), strings.Title(toOriginalName))
					}
				} else {
					return fmt.Sprintf("${%s}", libraryLink.To.Name)
				}
			}(),
			"LINKER_FLAGS": cmakeSliceJoin(libraryLink.LinkInfo.Flags, "no linker flags provided"),
		})

		cmakeStrings = append(cmakeStrings, finalString+"\n")
	}

	for dependencyLink := range dependencies.LinkIterator() {
		if dependencyLink.From.HeaderOnly {
			dependencyLink.LinkInfo.Visibility = types.Interface
		} else if strings.Trim(dependencyLink.LinkInfo.Visibility, " ") == "" {
			dependencyLink.LinkInfo.Visibility = types.Private
		}

		finalString := template.Replace(cmake.LinkString, map[string]string{
			"LINK_FROM":       dependencyLink.From.Name,
			"LINK_VISIBILITY": dependencyLink.LinkInfo.Visibility,
			"LINK_TO":         dependencyLink.To.Name,
			"LINKER_FLAGS":    cmakeSliceJoin(dependencyLink.LinkInfo.Flags, "no linker flags provided"),
		})
		cmakeStrings = append(cmakeStrings, finalString+"\n")
	}

	fileContents := []byte(strings.Join(cmakeStrings, "\n"))
	return ioutil.WriteFile(cmakePath, fileContents, os.ModePerm)
}

// Scans the dependency tree and creates build targets that will be converted into CMake targets
func CreateBuildTargets(projectDir string, target types.Target) (*TargetSet, *TargetSet, error) {
	targetSet := NewTargetSet()
	libraryTargetSet := NewTargetSet()

	i := resolve.NewInfo(projectDir)
	config, err := types.ReadWioConfig(projectDir, true)
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
		for name, library := range config.GetLibraries() {
			libraryTarget := &Target{
				Name:       name,
				ParentPath: projectDir,
				Library:    library,
			}

			libraryTargetSet.Add(libraryTarget, true)
			libraryTargetSet.Link(parentTarget, libraryTarget, &TargetLinkInfo{
				Visibility: library.GetLinkVisibility(),
				Flags:      library.GetLinkerFlags(),
			})
		}

		for _, dep := range i.GetRoot().Dependencies {
			var configDependency types.Dependency
			var exists bool

			if configDependency, exists = config.GetDependencies()[dep.Name]; !exists {
				return nil, nil, util.Error("%s@%s dependency is invalid and information is wrong in wio.yml",
					dep.Name, dep.ResolvedVersion.String())
			}

			parentInfo := &parentGivenInfo{
				flags:          configDependency.GetCompileFlags(),
				definitions:    configDependency.GetDefinitions(),
				linkVisibility: configDependency.GetVisibility(),
				linkFlags:      configDependency.GetLinkerFlags(),
				OsSupported:    configDependency.GetOsSupported(),
			}

			// all direct dependencies will link to the main target
			err := resolveTree(i, dep, parentTarget, targetSet, libraryTargetSet, target.GetFlags().GetGlobal(),
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
		}, targetSet, libraryTargetSet, target.GetFlags().GetGlobal(),
			target.GetDefinitions().GetGlobal(), parentInfo)
		if err != nil {
			return nil, nil, err
		}
	}

	return targetSet, libraryTargetSet, nil
}
