package config

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hil"
)

type hilString struct {
	Value string
}

func (hilString hilString) Get(config *hil.EvalConfig) (string, error) {
	return applyHilString(hilString.Value, config)
}

// //////////////////////

type variableImpl struct {
	Name  string `mapstructure:"name"`
	Value string `mapstructure:"value"`
}

func (variableImpl variableImpl) GetName() string {
	return variableImpl.Name
}

func (variableImpl variableImpl) GetValue() string {
	return variableImpl.Value
}

// //////////////////////

type argumentImpl struct {
	Name  string `mapstructure:"name"`
	Value string `mapstructure:"value"`
}

func (argumentImpl argumentImpl) GetName() string {
	return argumentImpl.Name
}

func (argumentImpl argumentImpl) GetValue(config *hil.EvalConfig) (string, error) {
	return applyHilString(argumentImpl.Value, config)
}

// //////////////////////

type toolchainImpl struct {
	Name string `mapstructure:"name"`
	Ref  string `mapstructure:"ref"`
}

func (toolchainImpl toolchainImpl) GetName() string {
	return toolchainImpl.Name
}

func (toolchainImpl toolchainImpl) GetRef(config *hil.EvalConfig) (string, error) {
	return applyHilString(toolchainImpl.Name, config)
}

// //////////////////////

type linkerOptionsImpl struct {
	Flags      []string `mapstructure:"flags"`
	Visibility string   `mapstructure:"visibility"`
}

func (linkerOptionsImpl linkerOptionsImpl) GetFlags() Flags {
	var flags Flags
	for _, flag := range linkerOptionsImpl.Flags {
		flags = append(flags, hilString{Value: flag})
	}
	return flags
}

func (linkerOptionsImpl linkerOptionsImpl) GetVisibility(config *hil.EvalConfig) (string, error) {
	return applyHilString(linkerOptionsImpl.Visibility, config)
}

// //////////////////////

type compileOptionsImpl struct {
	Flags       []string `mapstructure:"flags"`
	Definitions []string `mapstructure:"definitions"`
	CXXStandard string   `mapstructure:"cxx_standard"`
	CStandard   string   `mapstructure:"c_standard"`
}

func (compileOptionsImpl compileOptionsImpl) GetFlags() Flags {
	var flags Flags
	for _, flag := range compileOptionsImpl.Flags {
		flags = append(flags, hilString{Value: flag})
	}
	return flags
}

func (compileOptionsImpl compileOptionsImpl) GetDefinitions() Definitions {
	var definitions Definitions
	for _, definition := range compileOptionsImpl.Definitions {
		definitions = append(definitions, hilString{Value: definition})
	}
	return definitions
}

func (compileOptionsImpl compileOptionsImpl) GetCXXStandard(config *hil.EvalConfig) (string, error) {
	return applyHilString(compileOptionsImpl.CXXStandard, config)
}

func (compileOptionsImpl compileOptionsImpl) GetCStandard(config *hil.EvalConfig) (string, error) {
	return applyHilString(compileOptionsImpl.CXXStandard, config)
}

// //////////////////////

type dependencyImpl struct {
	Ref           string            `mapstructure:"ref"`
	Arguments     []argumentImpl    `mapstructure:"arguments"`
	LinkerOptions linkerOptionsImpl `mapstructure:"linker_options"`
}

func (dependencyImpl dependencyImpl) GetRef(config *hil.EvalConfig) (string, error) {
	return applyHilString(dependencyImpl.Ref, config)
}

func (dependencyImpl dependencyImpl) GetArguments() Arguments {
	var arguments Arguments
	for _, argumentImpl := range dependencyImpl.Arguments {
		arguments = append(arguments, argumentImpl)
	}
	return arguments
}

func (dependencyImpl dependencyImpl) GetLinkerOptions() LinkerOptions {
	return dependencyImpl.LinkerOptions
}

// //////////////////////

type packageOptionsImpl struct {
	HeaderOnly bool   `mapstructure:"header_only"`
	Type       string `mapstructure:"type"`
}

func (packageOptionsImpl packageOptionsImpl) IsHeaderOnly() bool {
	return packageOptionsImpl.HeaderOnly
}

func (packageOptionsImpl packageOptionsImpl) GetPackageType(config *hil.EvalConfig) (string, error) {
	return applyHilString(packageOptionsImpl.Type, config)
}

// //////////////////////

type projectImpl struct {
	Name           string              `mapstructure:"name"`
	Version        string              `mapstructure:"version"`
	Author         string              `mapstructure:"author"`
	Contributors   []string            `mapstructure:"contributors"`
	Description    string              `mapstructure:"description"`
	Homepage       string              `mapstructure:"homepage"`
	Repository     []string            `mapstructure:"repository"`
	CompileOptions compileOptionsImpl  `mapstructure:"compile_options"`
	PackageOptions *packageOptionsImpl `mapstructure:"package_options"` // pkg only
}

func (projectImpl projectImpl) GetName(config *hil.EvalConfig) (string, error) {
	return applyHilString(projectImpl.Name, config)
}

func (projectImpl projectImpl) GetVersion(config *hil.EvalConfig) (*version.Version, error) {
	ver, err := applyHilString(projectImpl.Version, config)
	if err != nil {
		return nil, err
	}
	return version.NewVersion(ver)
}

func (projectImpl projectImpl) GetAuthor(config *hil.EvalConfig) (string, error) {
	return applyHilString(projectImpl.Author, config)
}

func (projectImpl projectImpl) GetContributors() Contributors {
	var contributors Contributors
	for _, contributor := range projectImpl.Contributors {
		contributors = append(contributors, hilString{Value: contributor})
	}

	return contributors
}

func (projectImpl projectImpl) GetDescription(config *hil.EvalConfig) (string, error) {
	return applyHilString(projectImpl.Description, config)
}

func (projectImpl projectImpl) GetRepository(config *hil.EvalConfig) (string, error) {
	return applyHilString(projectImpl.Homepage, config)
}

func (projectImpl projectImpl) GetHomepage(config *hil.EvalConfig) (string, error) {
	return applyHilString(projectImpl.Homepage, config)
}

func (projectImpl projectImpl) GetCompileOptions() CompileOptions {
	return projectImpl.CompileOptions
}

func (projectImpl projectImpl) GetPackageOptions() PackageOptions {
	return projectImpl.PackageOptions
}

// //////////////////////

type executableOptionsImpl struct {
	Source    []string      `mapstructure:"source"`
	MainFile  string        `mapstructure:"main_file"` // only for targets and not for tests
	Platform  string        `mapstructure:"platform"`
	Toolchain toolchainImpl `mapstructure:"toolchain"`
}

func (executableOptionsImpl executableOptionsImpl) GetSource() Sources {
	var sources Sources
	for _, source := range executableOptionsImpl.Source {
		sources = append(sources, hilString{Value: source})
	}

	return sources
}

func (executableOptionsImpl executableOptionsImpl) GetMainFile(config *hil.EvalConfig) (string, error) {
	return applyHilString(executableOptionsImpl.MainFile, config)
}

func (executableOptionsImpl executableOptionsImpl) GetPlatform(config *hil.EvalConfig) (string, error) {
	return applyHilString(executableOptionsImpl.Platform, config)
}

func (executableOptionsImpl executableOptionsImpl) GetToolchain() Toolchain {
	return executableOptionsImpl.Toolchain
}

// //////////////////////

type targetImpl struct {
	ExecutableOptions *executableOptionsImpl `mapstructure:"executable_options"` // app only
	PackageOptions    *packageOptionsImpl    `mapstructure:"package_options"`    // pkg only
	Arguments         []argumentImpl         `mapstructure:"arguments"`
	CompileOptions    compileOptionsImpl     `mapstructure:"compile_options"`
	LinkerOptions     linkerOptionsImpl      `mapstructure:"linker_options"`
}

func (targetImpl targetImpl) GetExecutableOptions() ExecutableOptions {
	return targetImpl.ExecutableOptions
}

func (targetImpl targetImpl) GetPackageOptions() PackageOptions {
	return targetImpl.PackageOptions
}

func (targetImpl targetImpl) GetArguments() Arguments {
	var arguments Arguments
	for _, argumentImpl := range targetImpl.Arguments {
		arguments = append(arguments, argumentImpl)
	}
	return arguments
}

func (targetImpl targetImpl) GetCompileOptions() CompileOptions {
	return targetImpl.CompileOptions
}

func (targetImpl targetImpl) GetLinkerOptions() LinkerOptions {
	return targetImpl.LinkerOptions
}

// //////////////////////

type testImpl struct {
	ExecutableOptions executableOptionsImpl `mapstructure:"executable_options"`
	Arguments         []argumentImpl        `mapstructure:"arguments"`
	TargetName        string                `mapstructure:"target_name"`
	TargetArguments   []argumentImpl        `mapstructure:"target_arguments"`
	CompileOptions    compileOptionsImpl    `mapstructure:"compile_options"`
	LinkerOptions     linkerOptionsImpl     `mapstructure:"linker_options"`
}

func (testImpl testImpl) GetExecutableOptions() ExecutableOptions {
	return testImpl.ExecutableOptions
}

func (testImpl testImpl) GetArguments() Arguments {
	var arguments Arguments
	for _, argumentImpl := range testImpl.Arguments {
		arguments = append(arguments, argumentImpl)
	}
	return arguments
}

func (testImpl testImpl) GetTargetName(config *hil.EvalConfig) (string, error) {
	return applyHilString(testImpl.TargetName, config)
}

func (testImpl testImpl) GetTargetArguments() Arguments {
	var arguments Arguments
	for _, argumentImpl := range testImpl.TargetArguments {
		arguments = append(arguments, argumentImpl)
	}
	return arguments
}

func (testImpl testImpl) GetCompileOptions() CompileOptions {
	return testImpl.CompileOptions
}

func (testImpl testImpl) GetLinkerOptions() LinkerOptions {
	return testImpl.LinkerOptions
}

// //////////////////////

type projectConfigImpl struct {
	Type             string                     `mapstructure:"type"`
	Project          projectImpl                `mapstructure:"project"`
	Variables        []variableImpl             `mapstructure:"variables"`
	Arguments        []argumentImpl             `mapstructure:"arguments"`
	Scripts          []string                   `mapstructure:"scripts"`
	Targets          map[string]*targetImpl     `mapstructure:"targets"`
	Tests            map[string]*testImpl       `mapstructure:"tests"`
	Dependencies     map[string]*dependencyImpl `mapstructure:"dependencies"`
	TestDependencies map[string]*dependencyImpl `mapstructure:"test_dependencies"`
}

func (projectConfigImpl *projectConfigImpl) GetType() string {
	return projectConfigImpl.Type
}

func (projectConfigImpl *projectConfigImpl) GetProject() Project {
	return projectConfigImpl.Project
}

func (projectConfigImpl *projectConfigImpl) GetVariables() Variables {
	var variables Variables
	for _, variableImpl := range projectConfigImpl.Variables {
		variables = append(variables, variableImpl)
	}
	return variables
}

func (projectConfigImpl *projectConfigImpl) GetArguments() Arguments {
	var arguments Arguments
	for _, argumentImpl := range projectConfigImpl.Arguments {
		arguments = append(arguments, argumentImpl)
	}
	return arguments
}

func (projectConfigImpl *projectConfigImpl) GetScripts() Scripts {
	var scripts Scripts
	for _, script := range projectConfigImpl.Scripts {
		scripts = append(scripts, hilString{Value: script})
	}
	return scripts
}

func (projectConfigImpl *projectConfigImpl) GetTargets() Targets {
	targets := Targets{}
	for name, value := range projectConfigImpl.Targets {
		targets[name] = value
	}
	return targets
}

func (projectConfigImpl *projectConfigImpl) GetTests() Tests {
	tests := Tests{}
	for name, value := range projectConfigImpl.Tests {
		tests[name] = value
	}
	return tests
}

func (projectConfigImpl *projectConfigImpl) GetDependencies() Dependencies {
	dependencies := Dependencies{}
	for name, value := range projectConfigImpl.Dependencies {
		dependencies[name] = value
	}
	return dependencies
}

func (projectConfigImpl *projectConfigImpl) GetTestDependencies() Dependencies {
	dependencies := Dependencies{}
	for name, value := range projectConfigImpl.TestDependencies {
		dependencies[name] = value
	}
	return dependencies
}
