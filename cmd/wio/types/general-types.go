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
    GetFramework() string
    GetBoard() string
    GetFlags() TargetFlags
    GetDefinitions() TargetDefinitions
}

// Abstraction of targets that have been created
type Targets interface {
    GetDefaultTarget() string
    GetTargets() map[string]Target
}

// Abstraction of targets flags
type TargetFlags interface {
    GetGlobalFlags() []string
    GetTargetFlags() []string
    GetPkgFlags() []string
}

// Abstraction of targets definitions
type TargetDefinitions interface {
    GetGlobalDefinitions() []string
    GetTargetDefinitions() []string
    GetPkgDefinitions() []string
}

// ############################################# APP Targets ###############################################
type AppTargetFlags struct {
    GlobalFlags []string `yaml:"global_flags"`
    TargetFlags []string `yaml:"target_flags"`
}

func (appTargetFlags AppTargetFlags) GetGlobalFlags() ([]string) {
    return appTargetFlags.GlobalFlags
}

func (appTargetFlags AppTargetFlags) GetTargetFlags() ([]string) {
    return appTargetFlags.TargetFlags
}

func (appTargetFlags AppTargetFlags) GetPkgFlags() ([]string) {
    return nil
}

type AppTargetDefinitions struct {
    GlobalFlags []string `yaml:"global_definitions"`
    TargetFlags []string `yaml:"target_definitions"`
}

func (appTargetDefinitions AppTargetDefinitions) GetGlobalDefinitions() ([]string) {
    return appTargetDefinitions.GlobalFlags
}

func (appTargetDefinitions AppTargetDefinitions) GetTargetDefinitions() ([]string) {
    return appTargetDefinitions.TargetFlags
}

func (appTargetDefinitions AppTargetDefinitions) GetPkgDefinitions() ([]string) {
    return nil
}

// Structure to handle individual target inside targets for project of app AVR type
type AppAVRTarget struct {
    Framework   string
    Board       string
    Flags       AppTargetFlags
    Definitions AppTargetDefinitions
}

func (appTargetTag AppAVRTarget) GetBoard() string {
    return appTargetTag.Board
}

func (appTargetTag AppAVRTarget) GetFramework() string {
    return appTargetTag.Framework
}

func (appTargetTag AppAVRTarget) GetFlags() TargetFlags {
    return appTargetTag.Flags
}

func (appTargetTag AppAVRTarget) GetDefinitions() TargetDefinitions {
    return appTargetTag.Definitions
}

// type for the targets tag in the configuration file for project of app AVR type
type AppAVRTargets struct {
    DefaultTarget string            `yaml:"default"`
    Targets       map[string]AppAVRTarget `yaml:"create"`
}

func (appTargetsTag AppAVRTargets) GetDefaultTarget() string {
    return appTargetsTag.DefaultTarget
}

func (appTargetsTag AppAVRTargets) GetTargets() map[string]Target {
    targets := make(map[string]Target)

    for key, val := range appTargetsTag.Targets {
        targets[key] = val
    }

    return targets
}

// ######################################### PKG Targets #######################################################

type PkgTargetFlags struct {
    GlobalFlags []string `yaml:"global_flags"`
    TargetFlags []string `yaml:"target_flags"`
    PkgFlags    []string `yaml:"pkg_flags"`
}

func (pkgTargetFlags PkgTargetFlags) GetGlobalFlags() ([]string) {
    return pkgTargetFlags.GlobalFlags
}

func (pkgTargetFlags PkgTargetFlags) GetTargetFlags() ([]string) {
    return pkgTargetFlags.TargetFlags
}

func (pkgTargetFlags PkgTargetFlags) GetPkgFlags() ([]string) {
    return pkgTargetFlags.PkgFlags
}

type PkgTargetDefinitions struct {
    GlobalDefinitions []string `yaml:"global_definitions"`
    TargetDefinitions []string `yaml:"target_definitions"`
    PkgDefinitions    []string `yaml:"pkg_definitions"`
}

func (pkgTargetDefinitions PkgTargetDefinitions) GetGlobalDefinitions() ([]string) {
    return pkgTargetDefinitions.GlobalDefinitions
}

func (pkgTargetDefinitions PkgTargetDefinitions) GetTargetDefinitions() ([]string) {
    return pkgTargetDefinitions.TargetDefinitions
}

func (pkgTargetDefinitions PkgTargetDefinitions) GetPkgDefinitions() ([]string) {
    return pkgTargetDefinitions.PkgDefinitions
}

// Structure to handle individual target inside targets for project of pkg type
type PkgAVRTarget struct {
    Framework string
    Board       string
    Flags       PkgTargetFlags
    Definitions PkgTargetDefinitions
}

func (pkgAVRTarget PkgAVRTarget) GetBoard() string {
    return pkgAVRTarget.Board
}

func (pkgAVRTarget PkgAVRTarget) GetFlags() TargetFlags {
    return pkgAVRTarget.Flags
}

func (pkgAVRTarget PkgAVRTarget) GetFramework() string{
    return pkgAVRTarget.Framework
}

func (pkgAVRTarget PkgAVRTarget) GetDefinitions() TargetDefinitions {
    return pkgAVRTarget.Definitions
}

// type for the targets tag in the configuration file for project of pkg type
type PkgAVRTargets struct {
    DefaultTarget string                  `yaml:"default"`
    Targets       map[string]PkgAVRTarget `yaml:"create"`
}

func (pkgAVRTargets PkgAVRTargets) GetDefaultTarget() string {
    return pkgAVRTargets.DefaultTarget
}

func (pkgAVRTargets PkgAVRTargets) GetTargets() map[string]Target {
    targets := make(map[string]Target)

    for key, val := range pkgAVRTargets.Targets {
        targets[key] = val
    }

    return targets
}

// ##########################################  Dependencies ################################################

// Structure to handle individual library inside libraries
type DependencyTag struct {
    Version         string
    Vendor          bool
    LinkVisibility  string   `yaml:"link_visibility"`
    DependencyFlags []string `yaml:"dependency_flags"`
    DependencyDefinitions []string `yaml:"dependency_definitions"`
}

// type for the libraries tag in the main wio.yml file
type DependenciesTag map[string]*DependencyTag

// ############################################### Project ##################################################

type MainTag interface {
    GetName() string
    GetConfigurations() Configurations
    GetCompileOptions() CompileOptions
    GetIde() string
}

type CompileOptions interface {
    IsHeaderOnly() bool
    GetPlatform() string
}

type Configurations struct {
    WioVersion            string   `yaml:"minimum_wio_version"`
    SupportedPlatforms    []string `yaml:"supported_platforms"`
    UnSupportedPlatforms  []string `yaml:"unsupported_platforms"`
    SupportedFrameworks   []string `yaml:"supported_frameworks"`
    UnSupportedFrameworks []string `yaml:"unsupported_frameworks"`
    SupportedBoards       []string `yaml:"supported_boards"`
    UnSupportedBoards     []string `yaml:"unsupported_boards"`
}

// ############################################# APP Project ###############################################

// Structure to hold information about project type: app
type AppTag struct {
    Name           string
    Ide            string
    Config         Configurations
    CompileOptions AppCompileOptions `yaml:"compile_options"`
}

type AppCompileOptions struct {
    Platform string
}

func (appCompileOptions AppCompileOptions) IsHeaderOnly() bool {
    return false
}

func (appCompileOptions AppCompileOptions) GetPlatform() string {
    return appCompileOptions.Platform
}

func (appTag AppTag) GetName() string {
    return appTag.Name
}

func (appTag AppTag) GetConfigurations() Configurations {
    return appTag.Config
}

func (appTag AppTag) GetCompileOptions() CompileOptions {
    return appTag.CompileOptions
}

func (appTag AppTag) GetIde() string {
    return appTag.Ide
}

// ############################################# PKG Project ###############################################

type PackageMeta struct {
    Name         string
    Description  string
    Repository   string
    Version      string
    Author       string
    Contributors []string
    Organization string
    Keywords     []string
    License      string
}

type PkgCompileOptions struct {
    HeaderOnly bool `yaml:"header_only"`
    Platform string
}

func (pkgCompileOptions PkgCompileOptions) IsHeaderOnly() bool {
    return pkgCompileOptions.HeaderOnly
}

func (pkgCompileOptions PkgCompileOptions) GetPlatform() string {
    return pkgCompileOptions.Platform
}

type Flags struct {
    AllowOnlyGlobalFlags   bool     `yaml:"allow_only_global_flags"`
    AllowOnlyRequiredFlags bool     `yaml:"allow_only_required_flags"`
    GlobalFlags            []string `yaml:"global_flags"`
    RequiredFlags          []string `yaml:"required_flags"`
    IncludedFlags          []string `yaml:"included_flags"`
}

type Definitions struct {
    AllowOnlyGlobalDefinitions   bool     `yaml:"allow_only_global_definitions"`
    AllowOnlyRequiredDefinitions bool     `yaml:"allow_only_required_definitions"`
    GlobalDefinitions            []string `yaml:"global_definitions"`
    RequiredDefinitions          []string `yaml:"required_definitions"`
    IncludedDefinitions          []string `yaml:"included_definitions"`
}

// Structure to hold information about project type: lib
type PkgTag struct {
    Ide            string
    Meta           PackageMeta
    Config         Configurations
    CompileOptions PkgCompileOptions `yaml:"compile_options"`
    Flags          Flags
    Definitions    Definitions
}

func (pkgTag PkgTag) GetName() string {
    return pkgTag.Meta.Name
}

func (pkgTag PkgTag) GetConfigurations() Configurations {
    return pkgTag.Config
}

func (pkgTag PkgTag) GetIde() string {
    return pkgTag.Ide
}

func (pkgTag PkgTag) GetCompileOptions() CompileOptions {
    return pkgTag.CompileOptions
}

type Config interface {
    GetMainTag() MainTag
    GetTargets() Targets
    GetDependencies() DependenciesTag
}

type AppConfig struct {
    MainTag         AppTag          `yaml:"app"`
    TargetsTag      AppAVRTargets   `yaml:"targets"`
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
    TargetsTag      PkgAVRTargets   `yaml:"targets"`
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
