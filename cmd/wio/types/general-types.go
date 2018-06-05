// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package type contains types for use by other packages
// This file contains all the types that are used throughout the application

package types

type CliArgs struct {
    AppType   string
    Directory string
    Board     string
    Framework string
    Platform  string
    Ide       string
    Tests     bool
    Update    bool
}

// ############################################### Targets ##################################################

// Abstraction of a Target
type Target interface {
    GetBoard() string
    GetFlags() map[string][]string
}

// Abstraction of targets that have been created
type Targets interface {
    GetDefaultTarget() string
    GetTargets() map[string]Target
}

// ############################################# APP Targets ###############################################

// Structure to handle individual target inside targets for project of app type
type AppTargetTag struct {
    Board              string
    TargetCompileFlags []string `yaml:"compile_flags"`
}

func (appTargetTag AppTargetTag) GetBoard() string {
    return appTargetTag.Board
}

func (appTargetTag AppTargetTag) GetFlags() map[string][]string {
    flags := make(map[string][]string)
    flags["target_compile_flags"] = appTargetTag.TargetCompileFlags
    return flags
}

// type for the targets tag in the configuration file for project of app type
type AppTargetsTag struct {
    DefaultTarget string                  `yaml:"default"`
    Targets       map[string]AppTargetTag `yaml:"create"`
}

func (appTargetsTag AppTargetsTag) GetDefaultTarget() (string) {
    return appTargetsTag.DefaultTarget
}

func (appTargetsTag AppTargetsTag) GetTargets() (map[string]Target) {
    targets := make(map[string]Target)

    for key, val := range appTargetsTag.Targets {
        targets[key] = val
    }

    return targets
}

// ######################################### PKG Targets #######################################################

// Structure to handle individual target inside targets for project of pkg type
type PkgTargetTag struct {
    Board              string
    TargetCompileFlags []string `yaml:"target_compile_flags"`
    PkgCompileFlags    []string `yaml:"pkg_compile_flags"`
}

func (pkgTargetTag PkgTargetTag) GetBoard() string {
    return pkgTargetTag.Board
}

func (pkgTargetTag PkgTargetTag) GetFlags() map[string][]string {
    flags := make(map[string][]string)
    flags["target_compile_flags"] = pkgTargetTag.TargetCompileFlags
    flags["pkg_compile_flags"] = pkgTargetTag.PkgCompileFlags
    return flags
}

// type for the targets tag in the configuration file for project of pkg type
type PkgTargetsTag struct {
    DefaultTarget string                  `yaml:"default"`
    Targets       map[string]PkgTargetTag `yaml:"create"`
}

func (pkgTargetsTag PkgTargetsTag) GetDefaultTarget() (string) {
    return pkgTargetsTag.DefaultTarget
}

func (pkgTargetsTag PkgTargetsTag) GetTargets() (map[string]Target) {
    targets := make(map[string]Target)

    for key, val := range pkgTargetsTag.Targets {
        targets[key] = val
    }

    return targets
}

// ##########################################  Dependencies ################################################

// Structure to handle individual library inside libraries
type DependencyTag struct {
    Version      string
    Vendor       bool
    CompileFlags []string `yaml:"compile_flags"`
}

// type for the libraries tag in the main wio.yml file
type DependenciesTag map[string]*DependencyTag

// ############################################### Project ##################################################

type MainTag interface {
    GetName() string
    GetPlatforms() []string
    GetFrameworks() []string
    GetIde() string
    IsHeaderOnly() bool
}

// ############################################# APP Project ###############################################

// Structure to hold information about project type: app
type AppTag struct {
    Name      string
    Platform  string
    Framework string
    Ide       string
}

func (appTag AppTag) GetName() string {
    return appTag.Name
}

func (appTag AppTag) GetPlatforms() []string {
    return []string{appTag.Platform}
}

func (appTag AppTag) GetFrameworks() []string {
    return []string{appTag.Framework}
}

func (appTag AppTag) GetIde() string {
    return appTag.Ide
}

func (appTag AppTag) IsHeaderOnly() bool {
    return false
}

// ############################################# PKG Project ###############################################

// Structure to hold information about project type: lib
type PkgTag struct {
    Name         string
    Description  string
    Repository   string
    Version      string
    Author       string
    Contributors []string
    Organization string
    Keywords     []string
    License      string
    HeaderOnly   bool     `yaml:"header_only"`
    Platform     string
    Framework    []string
    Board        []string
    CompileFlags []string `yaml:"compile_flags"`
    Ide          string
}

func (pkgTag PkgTag) GetName() string {
    return pkgTag.Name
}

func (pkgTag PkgTag) GetPlatforms() []string {
    return []string{pkgTag.Platform}
}

func (pkgTag PkgTag) GetFrameworks() []string {
    return pkgTag.Framework
}

func (pkgTag PkgTag) GetIde() string {
    return pkgTag.Ide
}

func (pkgTag PkgTag) IsHeaderOnly() bool {
    return pkgTag.HeaderOnly
}

type Config interface {
    GetMainTag() MainTag
    GetTargets() Targets
    GetDependencies() DependenciesTag
}

type AppConfig struct {
    MainTag         AppTag          `yaml:"app"`
    TargetsTag      AppTargetsTag   `yaml:"targets"`
    DependenciesTag DependenciesTag `yaml:"dependencies"`
}

func (appConfig AppConfig) GetMainTag() MainTag {
    return appConfig.MainTag
}

func (appConfig AppConfig) GetTargets() Targets {
    return appConfig.TargetsTag
}

func (appConfig AppConfig) GetDependencies() DependenciesTag {
    return appConfig.DependenciesTag
}

type PkgConfig struct {
    MainTag         PkgTag          `yaml:"pkg"`
    TargetsTag      PkgTargetsTag   `yaml:"targets"`
    DependenciesTag DependenciesTag `yaml:"dependencies"`
}

func (pkgConfig PkgConfig) GetMainTag() MainTag {
    return pkgConfig.MainTag
}

func (pkgConfig PkgConfig) GetTargets() Targets {
    return pkgConfig.TargetsTag
}

func (pkgConfig PkgConfig) GetDependencies() DependenciesTag {
    return pkgConfig.DependenciesTag
}

type NpmDependencyTag map[string]string

type NpmConfig struct {
    Name         string           `json:"name"`
    Version      string           `json:"version"`
    Description  string           `json:"description"`
    Repository   string           `json:"repository"`
    Main         string           `json:"main"`
    Keywords     []string         `json:"keywords"`
    Author       string           `json:"author"`
    License      string           `json:"license"`
    Contributors []string         `json:"contributors"`
    Dependencies NpmDependencyTag `json:"dependencies"`
}

// DConfig contains configurations for default commandline arguments
type DConfig struct {
    Ide       string
    Framework string
    Platform  string
    Port      string
    Version   string
    Board     string
    Btarget   string
    Utarget   string
}
