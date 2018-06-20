package errors

import (
    "fmt"
    "strings"
)

type Error interface {
    error
}

const (
    Spaces = "         "
)

type ProgramArgumentsError struct {
    CommandName  string
    ArgumentName string
    Err          error
}

func (err ProgramArgumentsError) Error() string {
    str := fmt.Sprintf(`"%s" argument is invalid or not provided for "%s" command`, strings.ToLower(err.ArgumentName), err.CommandName)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type ProgrammingArgumentAssumption struct {
    CommandName  string
    ArgumentName string
    Err          error
}

func (err ProgrammingArgumentAssumption) Error() string {
    str := fmt.Sprintf(`"%s" argument is set to default by "%s" command`, strings.ToLower(err.ArgumentName), err.CommandName)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type PathDoesNotExist struct {
    Path string
    Err  error
}

func (err PathDoesNotExist) Error() string {
    str := fmt.Sprintf(`path does not exist: %s`, err.Path)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type ConfigMissing struct {
    Err error
}

func (err ConfigMissing) Error() string {
    str := fmt.Sprintf(`wio.yml file does not exist: Not a valid wio project`)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type ConfigParsingError struct {
    Err error
}

func (err ConfigParsingError) Error() string {
    str := fmt.Sprintf(`wio.yml file could not be parsed successfully`)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type UnsupportedWioConfigVersion struct {
    PackageName string
    Version     string
    Err         error
}

func (err UnsupportedWioConfigVersion) Error() string {
    str := fmt.Sprintf(`current wio does not support config file of %s => version is invalid: %s`,
        err.PackageName, err.Version)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type ProjectTypeMismatchError struct {
    GivenType  string
    ParsedType string
    Err        error
}

func (err ProjectTypeMismatchError) Error() string {
    str := fmt.Sprintf(`project type given is "%s" but parsed type from wio.yml is "%s"`, err.GivenType, err.ParsedType)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type OverridePossibilityError struct {
    Path string
    Err  error
}

func (err OverridePossibilityError) Error() string {
    str := fmt.Sprintf(`path is not empty and may be overwritten: %s`, err.Path)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type PlatformNotSupportedError struct {
    Platform string
    Err      error
}

func (err PlatformNotSupportedError) Error() string {
    str := fmt.Sprintf(`"%s" platform is not supported by wio`, err.Platform)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type FrameworkNotSupportedError struct {
    Framework string
    Platform string
    Err      error
}

func (err FrameworkNotSupportedError) Error() string {
    str := fmt.Sprintf(`"%s" framework is not supported for %s platform by wio`, err.Framework, err.Platform)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type ProjectStructureConstrainError struct {
    Constrain string
    Path      string
    Err       error
}

func (err ProjectStructureConstrainError) Error() string {
    str := fmt.Sprintf(`"%s" constrain not specified for file/dir: %s`, err.Constrain, err.Path)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type YamlMarshallError struct {
    Err error
}

func (err YamlMarshallError) Error() string {
    str := fmt.Sprintf(`yaml data could not be marshalled`)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type ReadFileError struct {
    FileName string
    Err      error
}

func (err ReadFileError) Error() string {
    str := fmt.Sprintf(`"%s" file read failed`, err.FileName)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type WriteFileError struct {
    FileName string
    Err      error
}

func (err WriteFileError) Error() string {
    str := fmt.Sprintf(`"%s" file write failed`, err.FileName)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type DeleteDirectoryError struct {
    DirName string
    Err     error
}

func (err DeleteDirectoryError) Error() string {
    str := fmt.Sprintf(`"%s" directory failed to be deleted`, err.DirName)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type DeleteFileError struct {
    FileName string
    Err      error
}

func (err DeleteFileError) Error() string {
    str := fmt.Sprintf(`"%s" file failed to be deleted`, err.FileName)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type NotValidWioProjectError struct {
    Directory string
    Err       error
}

func (err NotValidWioProjectError) Error() string {
    str := fmt.Sprintf(`"%s" is not a valid wio project`, err.Directory)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type DependencyDoesNotExistError struct {
    DependencyName string
    Vendor         bool
    Err            error
}

func (err DependencyDoesNotExistError) Error() string {
    var str string

    if err.Vendor {
        str = fmt.Sprintf(`vendor dependency named: "%s" does not exist. Check the vendor folder`, err.DependencyName)
    } else {
        str = fmt.Sprintf(`remote dependency named: "%s" does not exist. Pull the dependency first`, err.DependencyName)
    }

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type InvalidPlaceholderReferenceError struct {
    DependencyName string
    Placeholder    string
    Err            error
}

func (err InvalidPlaceholderReferenceError) Error() string {

    str := fmt.Sprintf(`invalid placeholder reference for dependency: %s and placeholder: %s`, err.DependencyName, err.Placeholder)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type GlobalFlagsMissingError struct {
    DependencyName string
    ProvidedFlags  []string
    MissingFlags   []string
    Err            error
}

func (err GlobalFlagsMissingError) Error() string {
    str := fmt.Sprintf("global flag/definition missing\n%sdependency: %s\n%sprovided flags/definitions: %s\n%smissing flags/definitions: %s",
        Spaces, err.DependencyName, Spaces, strings.Join(err.ProvidedFlags, ","), Spaces,
        strings.Join(err.MissingFlags, ","))

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type RequiredFlagsMissingError struct {
    From          string
    To            string
    ProvidedFlags []string
    MissingFlags  []string
    Err           error
}

func (err RequiredFlagsMissingError) Error() string {
    str := fmt.Sprintf("required flag/definition missing\n%sfrom: %s    to: %s\n%sprovided flags/definitions: %s\n%smissing flags/definitions: %s",
        Spaces, err.From, err.To, Spaces, strings.Join(err.ProvidedFlags, ","), Spaces,
        strings.Join(err.MissingFlags, ","))

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type OnlyRequiredFlagsError struct {
    Dependency       string
    FlagCategoryName string
    Err              error
}

func (err OnlyRequiredFlagsError) Error() string {
    str := fmt.Sprintf(`only required flags are accepted but provided %s for dependency: %s`,
        err.FlagCategoryName, err.Dependency)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type OnlyRequiredDefinitionsError struct {
    Dependency       string
    FlagCategoryName string
    Err              error
}

func (err OnlyRequiredDefinitionsError) Error() string {
    str := fmt.Sprintf(`only required definitions are accepted but provided %s for dependency: %s`,
        err.FlagCategoryName, err.Dependency)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type OnlyGlobalFlagsError struct {
    Dependency       string
    FlagCategoryName string
    Err              error
}

func (err OnlyGlobalFlagsError) Error() string {
    str := fmt.Sprintf(`only global flags are accepted but provided %s for dependency: %s`,
        err.FlagCategoryName, err.Dependency)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type OnlyGlobalDefinitionsError struct {
    Dependency       string
    FlagCategoryName string
    Err              error
}

func (err OnlyGlobalDefinitionsError) Error() string {
    str := fmt.Sprintf(`only global definitions are accepted but provided %s for dependency: %s`,
        err.FlagCategoryName, err.Dependency)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type LinkerVisibilityError struct {
    From            string
    To              string
    GivenVisibility string
    Err             error
}

func (err LinkerVisibilityError) Error() string {
    str := fmt.Sprintf(`"%s" => "%s" :linker visbility error. Provided visbility: %s`,
        err.From, err.To, err.GivenVisibility)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type FlagsDefinitionsVisibilityError struct {
    PackageName     string
    GivenVisibility string
    Err             error
}

func (err FlagsDefinitionsVisibilityError) Error() string {
    str := fmt.Sprintf(`"%s" :flags/definitions visbility error. Provided visbility: %s`,
        err.PackageName, err.GivenVisibility)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type TargetDoesNotExistError struct {
    TargetName string
    Err        error
}

func (err TargetDoesNotExistError) Error() string {
    str := fmt.Sprintf(`"%s" target does not exist. Skipping the build`, err.TargetName)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type CommandStartError struct {
    CommandName string
    Err         error
}

func (err CommandStartError) Error() string {
    str := fmt.Sprintf(`error starting command named: "%s" `, err.CommandName)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type CommandWaitError struct {
    CommandName string
    Err         error
}

func (err CommandWaitError) Error() string {
    str := fmt.Sprintf(`error occured while "%s" command was running `, err.CommandName)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type AutomaticPortNotDetectedError struct {
    CommandName string
    Err         error
}

func (err AutomaticPortNotDetectedError) Error() string {
    str := fmt.Sprintf("port could not be detected automatically. Provide the port using --port flag")

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type ActionNotSupportedByPlatform struct {
    Platform    string
    CommandName string
    Err         error
}

func (err ActionNotSupportedByPlatform) Error() string {
    str := fmt.Sprintf("%s platform does not support %s", err.Platform, err.CommandName)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type FatalError struct {
    Log         interface{}
    Err         error
}

func (err FatalError) Error() string {
    str := fmt.Sprintf("a fatal error occured. Contact developers for a fix")
    str += "\n" + Spaces + fmt.Sprintln(err.Log)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}
