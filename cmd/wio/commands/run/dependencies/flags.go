package dependencies

import (
    "regexp"
    "strings"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/errors"
)

// Verifies the placeholder syntax
func placeholderSyntaxValid(flag string) bool {
    pat := regexp.MustCompile(`\$\([a-zA-Z0-9=_-]+\)`)
    s := pat.FindString(flag)

    return s != ""
}

// matches a flag by the requested flag
func matchFlag(providedFlag string, requestedFlag string) string {
    pat := regexp.MustCompile(`^` + requestedFlag + `\b`)
    s := pat.FindString(providedFlag)

    return s
}

// fills placeholder flags and errors out if flag is not valid
func fillPlaceholderFlags(providedFlags []string, desiredFlags []string, dependencyName string) ([]string, error) {
    var newFlags []string
    for _, desiredFlag := range desiredFlags {
        if desiredFlag[0] != '$' || !placeholderSyntaxValid(desiredFlag) {
            newFlags = append(newFlags, desiredFlag)
            continue
        }

        oldLength := len(newFlags)

        for _, providedFlag := range providedFlags {
            newFlag := strings.Replace(desiredFlag, "$", "", 1)
            newFlag = strings.Replace(newFlag, "(", "", 1)
            newFlag = strings.Replace(newFlag, ")", "", 1)

            s := matchFlag(providedFlag, newFlag)

            if s != "" {
                newFlags = append(newFlags, providedFlag)
                break
            }
        }

        if len(newFlags) == oldLength {
            err := errors.InvalidPlaceholderReferenceError{
                DependencyName: dependencyName,
                Placeholder: desiredFlag,
            }

            return nil, err
        }
    }

    return newFlags, nil
}

// this fills global flags if they are requested
func fillGlobalFlags(globalFlags []string, dependencyGlobalFlagsRequired []string, dependencyName string) ([]string, error) {
    var filledFlags []string
    var notFilledFlags []string

    if len(globalFlags) == 0 {
        notFilledFlags = dependencyGlobalFlagsRequired
    } else {
        for _, requiredGlobalFlag := range dependencyGlobalFlagsRequired {
            numFilledFlags := len(filledFlags)

            for _, givenGlobalFlag := range globalFlags {
                s := matchFlag(givenGlobalFlag, requiredGlobalFlag)

                if s != "" {
                    filledFlags = append(filledFlags, givenGlobalFlag)
                    break
                }
            }

            // this means any of the global flag did not match
            if len(filledFlags) == numFilledFlags {
                notFilledFlags = append(notFilledFlags, requiredGlobalFlag)
            }
        }
    }

    // print errors when global flags are not provided
    if len(dependencyGlobalFlagsRequired) != len(filledFlags) {
        err := errors.GlobalFlagsMissingError{
            DependencyName: dependencyName,
            ProvidedFlags: filledFlags,
            MissingFlags: notFilledFlags,
        }

        return nil, err
    }

    return filledFlags, nil
}

// this fills required flags if they are requested
func fillRequiredFlags(providedFlags []string, dependencyFlagsRequired []string,
    dependencyName string, fromName string, isTarget bool) ([]string, []string, error) {
    var filledFlags []string
    var nonFilledRequiredFlags []string

    if len(providedFlags) == 0 {
        nonFilledRequiredFlags = dependencyFlagsRequired
    } else {
        for _, requiredFlag := range dependencyFlagsRequired {
            numFilledFlags := len(filledFlags)

            for _, givenFlag := range providedFlags {
                s := matchFlag(givenFlag, requiredFlag)

                if s != "" {
                    filledFlags = append(filledFlags, givenFlag)
                }
            }

            if len(filledFlags) == numFilledFlags {
                nonFilledRequiredFlags = append(nonFilledRequiredFlags, requiredFlag)
            }
        }
    }

    // print errors when global flags are not provided
    if len(dependencyFlagsRequired) != len(filledFlags) {
        var from string
        var to string

        if isTarget {
            from = "Target::" + fromName
            to = "Package::" + dependencyName
        } else {
            from = "Package::" + fromName
            to = "Package::" + dependencyName
        }


        err := errors.RequiredFlagsMissingError{
            From: from,
            To: to,
            ProvidedFlags: filledFlags,
            MissingFlags: nonFilledRequiredFlags,
        }

        return nil, nil, err
    }

    return filledFlags, utils.Difference(providedFlags, filledFlags), nil
}
