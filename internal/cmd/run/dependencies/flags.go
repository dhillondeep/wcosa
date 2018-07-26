package dependencies

import (
    "fmt"
    "regexp"
    "strings"
    "wio/pkg/util"
)

var placeholderMatch = regexp.MustCompile(`^\$\([a-zA-Z_-][a-zA-Z0-9_]*\)$`)

// Verifies the placeholder syntax
func IsPlaceholder(flag string) bool {
    return placeholderMatch.MatchString(strings.Trim(flag, " "))
}

// matches a flag by the requested flag
func TryMatch(key, given string) (string, bool) {
    pat := regexp.MustCompile(`^` + key + `(=|->).*$`)
    if !pat.MatchString(given) {
        return "", false
    }
    if strings.Contains(given, "->") {
        return strings.Split(given, "->")[1], true
    }
    return given, true
}

// fill placeholder flags and error if some are left unfilled
func fillPlaceholders(givenFlags, requiredFlags []string) ([]string, error) {
    var ret []string
    for _, required := range requiredFlags {
        if !IsPlaceholder(required) {
            ret = append(ret, required)
            continue
        }
        // look for a match
        for _, given := range givenFlags {
            key := required[2 : len(required)-1]
            if res, match := TryMatch(key, given); match {
                ret = append(ret, res)
                goto Continue
            }
        }
        return nil, util.Error(fmt.Sprintf("placeholder flag/definition \"%s\" unfilled in ", required) + "%s")

    Continue:
        continue
    }
    return ret, nil
}

// this fills global flags if they are requested
func fillDefinition(givenFlags, requiredFlags []string) ([]string, error) {
    var ret []string
    for _, required := range requiredFlags {
        for _, given := range givenFlags {
            if res, match := TryMatch(required, given); match {
                ret = append(ret, res)
                goto Continue
            }
        }
        return nil, util.Error("%s" + fmt.Sprintf(" definition \"%s\" unfilled in ", required) + "%s")

    Continue:
        continue
    }

    return ret, nil
}
