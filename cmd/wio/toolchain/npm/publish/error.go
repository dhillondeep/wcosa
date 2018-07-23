package publish

import "fmt"

type InvalidProjectType struct{}

func (e InvalidProjectType) Error() string {
    return "wio can only publish package projects"
}

type InvalidProjectVersion struct {
    ver string
}

func (e InvalidProjectVersion) Error() string {
    format := "package project has invalid version: %s"
    return fmt.Sprintf(format, e.ver)
}

type InvalidDependencyVersion struct {
    name string
    ver  string
}

func (e InvalidDependencyVersion) Error() string {
    format := "dependency %s has invalid version: %s"
    return fmt.Sprintf(format, e.name, e.ver)
}

type HttpFailed struct {
    status int
}

func (e HttpFailed) Error() string {
    format := "PUT request failed with %d"
    return fmt.Sprintf(format, e.status)
}

type UnknownError struct{}

func (e UnknownError) Error() string {
    return "publish request failed with an unknown error"
}
