package cmake

import (
	"path/filepath"
    "strings"
    "wio/cmd/wio/utils/template"
)

type Target struct {
    TargetName            string
    Path                  string
    Flags                 []string
    Definitions           []string
    FlagsVisibility       string
    DefinitionsVisibility string
    HeaderOnly            bool
}

// CMake Target Link information
type TargetLink struct {
    From       string
    To         string
    Visibility string
}

var libraryStrings = map[string]map[bool]string{
    "avr":    {false: avrLibrary, true: avrHeader},
    "native": {false: desktopLibrary, true: desktopHeader},
}

// This creates CMake library string that will be used to link libraries
func GenerateDependencies(platform string, targets map[string]*Target, links []TargetLink) []string {
    cmakeStrings := make([]string, 0, 256)

    for _, target := range targets {
        finalString := libraryStrings[platform][target.HeaderOnly]

        finalString = template.Replace(finalString, map[string]string{
            "DEPENDENCY_PATH":        filepath.ToSlash(target.Path),
            "DEPENDENCY_NAME":        target.TargetName,
            "DEPENDENCY_FLAGS":       strings.Join(target.Flags, " "),
            "DEPENDENCY_DEFINITIONS": strings.Join(target.Definitions, " "),
            "FLAGS_VISIBILITY":       target.FlagsVisibility,
            "DEFINITIONS_VISIBILITY": target.DefinitionsVisibility,
        })
        cmakeStrings = append(cmakeStrings, finalString+"\n")
    }

    for _, link := range links {
        finalString := template.Replace(linkString, map[string]string{
            "LINKER_NAME":     link.From,
            "DEPENDENCY_NAME": link.To,
            "LINK_VISIBILITY": link.Visibility,
        })
        cmakeStrings = append(cmakeStrings, finalString)
    }

    return cmakeStrings
}
