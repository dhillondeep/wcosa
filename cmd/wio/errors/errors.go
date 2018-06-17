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
	str := fmt.Sprintf(`"yaml data could not be marshalled`)

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
	str := fmt.Sprintf(`"%s file read failed`, err.FileName)

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
	str := fmt.Sprintf(`"%s file write failed`, err.FileName)

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
	str := fmt.Sprintf(`"%s directory failed to be deleted`, err.DirName)

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
	str := fmt.Sprintf(`"%s file failed to be deleted`, err.FileName)

	if err.Err != nil {
		str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
	}

	return str
}

// An error for exceptions that are intended to be seen by the user.
//
// These exceptions won't have any debugging information printed when they're
// thrown.
type ApplicationError struct {
	error
}

// An exception class for exceptions that are intended to be seen by the user
// and are associated with a problem in a file at some path.
type FileError struct {
	error
	path string
}

type IoError struct {
	error
}

func (ioError IoError) GetType() string {
	return "ioError"
}

func main() {
	err := IoError{}
	err.GetType()
}
