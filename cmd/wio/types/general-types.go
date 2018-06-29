// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package type contains types for use by other packages
// This file contains all the types that are used throughout the application

package types

import (
    "wio/cmd/wio/constants"
    "bufio"
    "strings"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/errors"
    "gopkg.in/yaml.v2"
    "regexp"
)

// ############################################### Targets ##################################################

// Abstraction of a Target
type Target interface {
    GetSrc() string
    GetName() string
    GetBoard() string
    GetFramework() string
    GetPlatform() string
    GetFlags() TargetFlags
    GetDefinitions() TargetDefinitions

    SetName(name string)
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

func (flags *AppTargetFlags) GetGlobalFlags() []string {
    return flags.GlobalFlags
}

func (flags *AppTargetFlags) GetTargetFlags() []string {
    return flags.TargetFlags
}

func (flags *AppTargetFlags) GetPkgFlags() []string {
    return nil
}

type AppTargetDefinitions struct {
    GlobalFlags []string `yaml:"global_definitions"`
    TargetFlags []string `yaml:"target_definitions"`
}

func (definitions *AppTargetDefinitions) GetGlobalDefinitions() []string {
    return definitions.GlobalFlags
}

func (definitions *AppTargetDefinitions) GetTargetDefinitions() []string {
    return definitions.TargetFlags
}

func (definitions *AppTargetDefinitions) GetPkgDefinitions() []string {
    return nil
}

// Structure to handle individual target inside targets for project of app AVR type
type AppTarget struct {
    Src         string
    Platform    string
    Framework   string
    Board       string
    Flags       AppTargetFlags
    Definitions AppTargetDefinitions

    name string
}

func (target *AppTarget) GetSrc() string {
    return target.Src
}

func (target *AppTarget) GetName() string {
    return target.name
}

func (target *AppTarget) GetBoard() string {
    return target.Board
}

func (target *AppTarget) GetFramework() string {
    return target.Framework
}

func (target *AppTarget) GetPlatform() string {
    return target.Platform
}

func (target *AppTarget) GetFlags() TargetFlags {
    return &target.Flags
}

func (target *AppTarget) GetDefinitions() TargetDefinitions {
    return &target.Definitions
}

func (target *AppTarget) SetName(name string) {
    target.name = name
}

// type for the targets tag in the configuration file for project of app AVR type
type AppTargets struct {
    DefaultTarget string               `yaml:"default"`
    Targets       map[string]*AppTarget `yaml:"create"`
}

func (targets *AppTargets) GetDefaultTarget() string {
    return targets.DefaultTarget
}

func (targets *AppTargets) GetTargets() map[string]Target {
    ret := make(map[string]Target)

    for key, val := range targets.Targets {
        ret[key] = val
    }

    return ret
}

// ######################################### PKG Targets #######################################################

type PkgTargetFlags struct {
    GlobalFlags []string `yaml:"global_flags"`
    TargetFlags []string `yaml:"target_flags"`
    PkgFlags    []string `yaml:"pkg_flags"`
}

func (flags *PkgTargetFlags) GetGlobalFlags() []string {
    return flags.GlobalFlags
}

func (flags *PkgTargetFlags) GetTargetFlags() []string {
    return flags.TargetFlags
}

func (flags *PkgTargetFlags) GetPkgFlags() []string {
    return flags.PkgFlags
}

type PkgTargetDefinitions struct {
    GlobalDefinitions []string `yaml:"global_definitions"`
    TargetDefinitions []string `yaml:"target_definitions"`
    PkgDefinitions    []string `yaml:"pkg_definitions"`
}

func (definitions *PkgTargetDefinitions) GetGlobalDefinitions() []string {
    return definitions.GlobalDefinitions
}

func (definitions *PkgTargetDefinitions) GetTargetDefinitions() []string {
    return definitions.TargetDefinitions
}

func (definitions *PkgTargetDefinitions) GetPkgDefinitions() []string {
    return definitions.PkgDefinitions
}

// Structure to handle individual target inside targets for project of pkg type
type PkgTarget struct {
    Src         string
    Platform    string
    Framework   string
    Board       string
    Flags       PkgTargetFlags
    Definitions PkgTargetDefinitions

    name string
}

func (target *PkgTarget) GetSrc() string {
    return target.Src
}

func (target *PkgTarget) GetName() string {
    return target.name
}

func (target *PkgTarget) GetBoard() string {
    return target.Board
}

func (target *PkgTarget) GetFlags() TargetFlags {
    return &target.Flags
}

func (target *PkgTarget) GetPlatform() string {
    return target.Platform
}

func (target *PkgTarget) GetFramework() string {
    return target.Framework
}

func (target *PkgTarget) GetDefinitions() TargetDefinitions {
    return &target.Definitions
}

func (target *PkgTarget) SetName(name string) {
    target.name = name
}

// type for the targets tag in the configuration file for project of pkg type
type PkgAVRTargets struct {
    DefaultTarget string               `yaml:"default"`
    Targets       map[string]*PkgTarget `yaml:"create"`
}

func (targets *PkgAVRTargets) GetDefaultTarget() string {
    return targets.DefaultTarget
}

func (targets *PkgAVRTargets) GetTargets() map[string]Target {
    ret := make(map[string]Target)

    for key, val := range targets.Targets {
        ret[key] = val
    }

    return ret
}

// ##########################################  Dependencies ################################################

// Structure to handle individual library inside libraries
type DependencyTag struct {
    Version               string
    Vendor                bool
    LinkVisibility        string              `yaml:"link_visibility"`
    Flags                 []string            `yaml:"flags"`
    Definitions           []string            `yaml:"definitions"`
    DependencyFlags       map[string][]string `yaml:"dependency_flags"`
    DependencyDefinitions map[string][]string `yaml:"dependency_definitions"`
}

// type for the libraries tag in the main wio.yml file
type DependenciesTag map[string]*DependencyTag

// ############################################### Project ##################################################

type MainTag interface {
    GetName() string
    GetVersion() string
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

func (options *AppCompileOptions) IsHeaderOnly() bool {
    return false
}

func (options *AppCompileOptions) GetPlatform() string {
    return options.Platform
}

func (app *AppTag) GetName() string {
    return app.Name
}

func (app *AppTag) GetVersion() string {
    return "1.0.0"
}

func (app *AppTag) GetConfigurations() Configurations {
    return app.Config
}

func (app *AppTag) GetCompileOptions() CompileOptions {
    return &app.CompileOptions
}

func (app *AppTag) GetIde() string {
    return app.Ide
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
    Platform   string
}

func (options *PkgCompileOptions) IsHeaderOnly() bool {
    return options.HeaderOnly
}

func (options *PkgCompileOptions) GetPlatform() string {
    return options.Platform
}

type Flags struct {
    AllowOnlyGlobalFlags   bool     `yaml:"allow_only_global_flags"`
    AllowOnlyRequiredFlags bool     `yaml:"allow_only_required_flags"`
    GlobalFlags            []string `yaml:"global_flags"`
    RequiredFlags          []string `yaml:"required_flags"`
    IncludedFlags          []string `yaml:"included_flags"`
    Visibility             string
}

type Definitions struct {
    AllowOnlyGlobalDefinitions   bool     `yaml:"allow_only_global_definitions"`
    AllowOnlyRequiredDefinitions bool     `yaml:"allow_only_required_definitions"`
    GlobalDefinitions            []string `yaml:"global_definitions"`
    RequiredDefinitions          []string `yaml:"required_definitions"`
    IncludedDefinitions          []string `yaml:"included_definitions"`
    Visibility                   string
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

func (tag *PkgTag) GetName() string {
    return tag.Meta.Name
}
func (tag *PkgTag) GetVersion() string {
    return tag.Meta.Version
}

func (tag *PkgTag) GetConfigurations() Configurations {
    return tag.Config
}

func (tag *PkgTag) GetIde() string {
    return tag.Ide
}

func (tag *PkgTag) GetCompileOptions() CompileOptions {
    return &tag.CompileOptions
}

type Type int

const (
    App Type = 0
    Pkg Type = 1
)

type Config struct {
    Config IConfig
    Type   Type
}

func (c Config) GetType() string {
    return c.Config.GetType()
}

func (c Config) GetMainTag() MainTag {
    return c.Config.GetMainTag()
}

func (c Config) GetTargets() Targets {
    return c.Config.GetTargets()
}

func (c Config) GetDependencies() DependenciesTag {
    return c.Config.GetDependencies()
}

func (c Config) SetDependencies(tag DependenciesTag) {
    c.Config.SetDependencies(tag)
}

func (c Config) PrettyPrint(path string) error {
    return prettyPrintConfig(c.Config, path)
}

type IConfig interface {
    GetType() string
    GetMainTag() MainTag
    GetTargets() Targets
    GetDependencies() DependenciesTag
    SetDependencies(tag DependenciesTag)
}

type AppConfig struct {
    MainTag         AppTag          `yaml:"app"`
    TargetsTag      AppTargets      `yaml:"targets"`
    DependenciesTag DependenciesTag `yaml:"dependencies"`
}

func (config *AppConfig) GetType() string {
    return constants.APP
}

func (config *AppConfig) GetMainTag() MainTag {
    return &config.MainTag
}

func (config *AppConfig) GetTargets() Targets {
    return &config.TargetsTag
}

func (config *AppConfig) GetDependencies() DependenciesTag {
    return config.DependenciesTag
}

func (config *AppConfig) SetDependencies(tag DependenciesTag) {
    config.DependenciesTag = tag
}

type PkgConfig struct {
    MainTag         PkgTag          `yaml:"pkg"`
    TargetsTag      PkgAVRTargets   `yaml:"targets"`
    DependenciesTag DependenciesTag `yaml:"dependencies"`
}

func (config *PkgConfig) GetType() string {
    return constants.PKG
}

func (config *PkgConfig) GetMainTag() MainTag {
    return &config.MainTag
}

func (config *PkgConfig) GetTargets() Targets {
    return &config.TargetsTag
}

func (config *PkgConfig) GetDependencies() DependenciesTag {
    return config.DependenciesTag
}

func (config *PkgConfig) SetDependencies(tag DependenciesTag) {
    config.DependenciesTag = tag
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

// Pretty print wio.yml
func (config *PkgConfig) PrettyPrint(path string) error {
    return prettyPrintConfig(config, path)
}

func (config *AppConfig) PrettyPrint(path string) error {
    return prettyPrintConfig(config, path)
}

func prettyPrintConfig(config IConfig, path string) error {
    data, err := yaml.Marshal(config)
    if err != nil {
        return err
    }
    return io.NormalIO.WriteFile(path, data)
}

// Write configuration with nice spacing and information
func prettyPrintHelp(config IConfig, filePath string) error {
    appInfoPath := "templates" + io.Sep + "config" + io.Sep + "app-helper.txt"
    pkgInfoPath := "templates" + io.Sep + "config" + io.Sep + "pkg-helper.txt"
    targetsInfoPath := "templates" + io.Sep + "config" + io.Sep + "targets-helper.txt"
    dependenciesInfoPath := "templates" + io.Sep + "config" + io.Sep + "dependencies-helper.txt"

    var ymlData []byte
    var appInfoData []byte
    var pkgInfoData []byte
    var targetsInfoData []byte
    var dependenciesInfoData []byte
    var err error

    if appInfoData, err = io.AssetIO.ReadFile(appInfoPath); err != nil {
        return errors.ReadFileError{
            FileName: appInfoPath,
            Err:      err,
        }
    }
    if pkgInfoData, err = io.AssetIO.ReadFile(pkgInfoPath); err != nil {
        return errors.ReadFileError{
            FileName: pkgInfoPath,
            Err:      err,
        }
    }
    if targetsInfoData, err = io.AssetIO.ReadFile(targetsInfoPath); err != nil {
        return errors.ReadFileError{
            FileName: targetsInfoPath,
            Err:      err,
        }
    }
    if dependenciesInfoData, err = io.AssetIO.ReadFile(dependenciesInfoPath); err != nil {
        return errors.ReadFileError{
            FileName: dependenciesInfoPath,
            Err:      err,
        }
    }

    // marshall yml data
    if ymlData, err = yaml.Marshal(config); err != nil {
        marshallError := errors.YamlMarshallError{
            Err: err,
        }
        return marshallError
    }

    finalStr := ""

    // configuration tags
    appTagPat := regexp.MustCompile(`(^app:)|((\s| |^\w)app:(\s+|))`)
    pkgTagPat := regexp.MustCompile(`(^pkg:)|((\s| |^\w)pkg:(\s+|))`)
    targetsTagPat := regexp.MustCompile(`(^targets:)|((\s| |^\w)targets:(\s+|))`)
    dependenciesTagPat := regexp.MustCompile(`(^dependencies:)|((\s| |^\w)dependencies:(\s+|))`)
    configTagPat := regexp.MustCompile(`(^config:)|((\s| |^\w)config:(\s+|))`)
    compileOptionsTagPat := regexp.MustCompile(`(^compile_options:)|((\s| |^\w)compile_options:(\s+|))`)
    metaTagPat := regexp.MustCompile(`(^meta:)|((\s| |^\w)meta:(\s+|))`)

    // empty array
    emptyArrayPat := regexp.MustCompile(`:\s+\[]`)
    // empty object
    emptyMapPat := regexp.MustCompile(`:\s+{}`)
    // empty tag
    emptyTagPat := regexp.MustCompile(`:\s+\n+|:\s+"\s+"|:\s+""|:"\s+"|:""`)

    scanner := bufio.NewScanner(strings.NewReader(string(ymlData)))
    for scanner.Scan() {
        line := scanner.Text()

        // ignore empty arrays, objects and tags
        if emptyArrayPat.MatchString(line) || emptyMapPat.MatchString(line) || emptyTagPat.MatchString(line) {
            if !(strings.Contains(line, "global_flags: []") ||
                strings.Contains(line, "target_flags: []") ||
                strings.Contains(line, "pkg_flags: []") ||
                strings.Contains(line, "global_definitions: []") ||
                strings.Contains(line, "target_definitions: []") ||
                strings.Contains(line, "pkg_definitions: []")) {
                continue
            }
        }

        if appTagPat.MatchString(line) {
            finalStr += string(appInfoData) + "\n"
            finalStr += line
        } else if pkgTagPat.MatchString(line) {
            finalStr += string(pkgInfoData) + "\n"
            finalStr += line
        } else if targetsTagPat.MatchString(line) {
            finalStr += "\n"
            finalStr += string(targetsInfoData) + "\n"
            finalStr += line
        } else if dependenciesTagPat.MatchString(line) {
            finalStr += "\n"
            finalStr += string(dependenciesInfoData) + "\n"
            finalStr += line
        } else if configTagPat.MatchString(line) || compileOptionsTagPat.MatchString(line) ||
            metaTagPat.MatchString(line) {
            finalStr += "\n"
            finalStr += line
        } else {
            finalStr += line
        }

        finalStr += "\n"
    }

    if err = io.NormalIO.WriteFile(filePath, []byte(finalStr)); err != nil {
        return errors.WriteFileError{
            FileName: filePath,
            Err:      err,
        }
    }

    return nil
}
