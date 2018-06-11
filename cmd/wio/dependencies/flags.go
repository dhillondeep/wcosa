package dependencies

import (
    "errors"
    "regexp"
    "strings"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io/log"
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
            log.Norm.Red(true, "Invalid placeholder reference")
            log.Norm.Cyan(true, "  Dependency: "+dependencyName+"\t Placeholder: "+desiredFlag)

            return nil, errors.New("invalid placeholder reference in " + dependencyName + " package")
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
        log.Norm.Red(true, "Global flags missing")
        log.Norm.Cyan(true, "  Dependency: "+dependencyName)

        if len(filledFlags) != 0 {
            log.Norm.Write(true, "    Provided Global Flags: "+strings.Join(filledFlags, ","))
        }
        log.Norm.Write(true, "    Missing Global Flags: "+strings.Join(notFilledFlags, ","))

        return nil, errors.New("global flag not provided " + dependencyName + " package")
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
        log.Norm.Red(true, "Required flags missing")

        if isTarget {
            log.Norm.Cyan(true, "  From: Target::"+fromName+"\tTo: Package::"+dependencyName)
        } else {
            log.Norm.Cyan(true, "  From: Package::"+fromName+"\tTo: Package::"+dependencyName)
        }

        if len(filledFlags) != 0 {
            log.Norm.Write(true, "    Provided Flags: "+strings.Join(filledFlags, ","))
        }
        log.Norm.Write(true, "    Missing Flags: "+strings.Join(nonFilledRequiredFlags, ","))

        return nil, nil, errors.New("required flag not provided for " + dependencyName + " package")
    }

    return filledFlags, utils.Difference(providedFlags, filledFlags), nil
}
