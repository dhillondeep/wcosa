package dependencies

import (
    "regexp"
    "strings"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/utils"
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
func fillPlaceholders(givenFlags, requiredFlags []string, name string) ([]string, error) {
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
        return nil, errors.Stringf("placeholder flag/definition %s unfilled in %s", required, name)

    Continue:
        continue
    }
    return ret, nil
}

// this fills global flags if they are requested
func fillGlobal(givenFlags, requiredFlags []string, name string) ([]string, error) {
    var ret []string
    for _, required := range requiredFlags {
        for _, given := range givenFlags {
            if res, match := TryMatch(required, given); match {
                ret = append(ret, res)
                goto Continue
            }
        }
        return nil, errors.Stringf("global flag/definition %s unfilled in %s", required, name)

    Continue:
        continue
    }

    if len(givenFlags) < len(requiredFlags) {
        return nil, errors.Stringf("global flag(s)/definition(s) %s unfilled in %s",
            strings.Join(utils.Difference(requiredFlags, ret), ", "), name)
    }

    return ret, nil
}

// this fills required flags if they are requested
func fillRequired(givenFlags []string, requiredFlags []string, name string) ([]string, []string, error) {
    var ret []string
    for _, required := range requiredFlags {
        for _, given := range givenFlags {
            if res, match := TryMatch(required, given); match {
                ret = append(ret, res)
                goto Continue
            }
            return nil, nil, errors.Stringf("required flag/definition %s unfilled in %s", required, name)

        Continue:
            continue
        }
    }

    if len(givenFlags) < len(requiredFlags) {
        return nil, nil, errors.Stringf("required flag(s)/definition(s) %s unfilled in %s",
            strings.Join(utils.Difference(requiredFlags, ret), ", "), name)
    }

    return ret, utils.Difference(givenFlags, ret), nil
}
