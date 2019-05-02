package config

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hil"
	"strconv"
	"strings"
)

type ExpressionImpl struct {
	Value string
}

func (expressionImpl ExpressionImpl) Eval(config *hil.EvalConfig) (string, error) {
	return applyEvaluator(expressionImpl.Value, config)
}

// //////////////////////

type VariableImpl struct {
	Name  string `mapstructure:"name"`
	Value string `mapstructure:"value"`
}

func (variableImpl VariableImpl) GetName() string {
	return variableImpl.Name
}

func (variableImpl VariableImpl) GetValue(config *hil.EvalConfig) (string, error) {
	return applyEvaluator(variableImpl.Value, config)
}

// //////////////////////

type ArgumentImpl struct {
	Name  string `mapstructure:"name"`
	Value string `mapstructure:"value"`
}

func (argumentImpl ArgumentImpl) GetName() string {
	return argumentImpl.Name
}

func (argumentImpl ArgumentImpl) GetValue(config *hil.EvalConfig) (string, error) {
	return applyEvaluator(argumentImpl.Value, config)
}

// //////////////////////

type ToolchainImpl struct {
	Name string `mapstructure:"name"`
	Ref  string `mapstructure:"ref"`
}

func (toolchainImpl ToolchainImpl) GetName(config *hil.EvalConfig) (string, error) {
	return applyEvaluator(toolchainImpl.Name, config)
}

func (toolchainImpl ToolchainImpl) GetRef(config *hil.EvalConfig) (string, error) {
	return applyEvaluator(toolchainImpl.Ref, config)
}

// //////////////////////

type LinkerOptionsImpl struct {
	Flags      []string `mapstructure:"flags"`
	Visibility string   `mapstructure:"visibility"`
}

func (linkerOptionsImpl LinkerOptionsImpl) GetFlags() Flags {
	var flags Flags
	for _, flag := range linkerOptionsImpl.Flags {
		flags = append(flags, ExpressionImpl{Value: strings.TrimSpace(flag)})
	}
	return flags
}

func (linkerOptionsImpl LinkerOptionsImpl) GetVisibility(config *hil.EvalConfig) (string, error) {
	return applyEvaluator(linkerOptionsImpl.Visibility, config)
}

// //////////////////////

type CompileOptionsImpl struct {
	Flags       []string `mapstructure:"flags"`
	Definitions []string `mapstructure:"definitions"`
	CXXStandard string   `mapstructure:"cxx_standard"`
	CStandard   string   `mapstructure:"c_standard"`
}

func (compileOptionsImpl CompileOptionsImpl) GetFlags() Flags {
	var flags Flags
	for _, flag := range compileOptionsImpl.Flags {
		flags = append(flags, ExpressionImpl{Value: strings.TrimSpace(flag)})
	}
	return flags
}

func (compileOptionsImpl CompileOptionsImpl) GetDefinitions() Definitions {
	var definitions Definitions
	for _, definition := range compileOptionsImpl.Definitions {
		definitions = append(definitions, ExpressionImpl{Value: strings.TrimSpace(definition)})
	}
	return definitions
}

func (compileOptionsImpl CompileOptionsImpl) GetCXXStandard(config *hil.EvalConfig) (string, error) {
	return applyEvaluator(compileOptionsImpl.CXXStandard, config)
}

func (compileOptionsImpl CompileOptionsImpl) GetCStandard(config *hil.EvalConfig) (string, error) {
	return applyEvaluator(compileOptionsImpl.CStandard, config)
}

// //////////////////////

type DependencyImpl struct {
	Ref           string            `mapstructure:"ref"`
	Arguments     []ArgumentImpl    `mapstructure:"arguments"`
	LinkerOptions LinkerOptionsImpl `mapstructure:"linker_options"`
}

func (dependencyImpl DependencyImpl) GetRef(config *hil.EvalConfig) (string, error) {
	return applyEvaluator(dependencyImpl.Ref, config)
}

func (dependencyImpl DependencyImpl) GetArguments() Arguments {
	var arguments Arguments
	for _, argumentImpl := range dependencyImpl.Arguments {
		arguments = append(arguments, argumentImpl)
	}
	return arguments
}

func (dependencyImpl DependencyImpl) GetLinkerOptions() LinkerOptions {
	return dependencyImpl.LinkerOptions
}

// //////////////////////

type PackageOptionsImpl struct {
	HeaderOnly string `mapstructure:"header_only"`
	Type       string `mapstructure:"type"`
}

func (packageOptionsImpl PackageOptionsImpl) IsHeaderOnly(config *hil.EvalConfig) (bool, error) {
	result, err := applyEvaluator(packageOptionsImpl.HeaderOnly, config)
	if err != nil {
		return false, err
	}

	return strconv.ParseBool(result)
}

func (packageOptionsImpl PackageOptionsImpl) GetPackageType(config *hil.EvalConfig) (string, error) {
	return applyEvaluator(packageOptionsImpl.Type, config)
}

// //////////////////////

type ProjectImpl struct {
	Name           string              `mapstructure:"name"`
	Version        string              `mapstructure:"version"`
	Author         string              `mapstructure:"author"`
	Contributors   []string            `mapstructure:"contributors"`
	Description    string              `mapstructure:"description"`
	Homepage       string              `mapstructure:"homepage"`
	Repository     []string            `mapstructure:"repository"`
	CompileOptions *CompileOptionsImpl `mapstructure:"compile_options"`
	PackageOptions *PackageOptionsImpl `mapstructure:"package_options"` // pkg only
}

func (projectImpl ProjectImpl) GetName(config *hil.EvalConfig) (string, error) {
	return applyEvaluator(projectImpl.Name, config)
}

func (projectImpl ProjectImpl) GetVersion(config *hil.EvalConfig) (*version.Version, error) {
	ver, err := applyEvaluator(projectImpl.Version, config)
	if err != nil {
		return nil, err
	}
	return version.NewVersion(ver)
}

func (projectImpl ProjectImpl) GetAuthor(config *hil.EvalConfig) (string, error) {
	return applyEvaluator(projectImpl.Author, config)
}

func (projectImpl ProjectImpl) GetContributors() Contributors {
	var contributors Contributors
	for _, contributor := range projectImpl.Contributors {
		contributors = append(contributors, ExpressionImpl{Value: strings.TrimSpace(contributor)})
	}

	return contributors
}

func (projectImpl ProjectImpl) GetDescription(config *hil.EvalConfig) (string, error) {
	return applyEvaluator(projectImpl.Description, config)
}

func (projectImpl ProjectImpl) GetRepository() Repositories {
	var repositories Repositories
	for _, repository := range projectImpl.Repository {
		repositories = append(repositories, ExpressionImpl{Value: strings.TrimSpace(repository)})
	}

	return repositories
}

func (projectImpl ProjectImpl) GetHomepage(config *hil.EvalConfig) (string, error) {
	return applyEvaluator(projectImpl.Homepage, config)
}

func (projectImpl ProjectImpl) GetCompileOptions() CompileOptions {
	if projectImpl.CompileOptions == nil {
		return nil
	}
	return projectImpl.CompileOptions
}

func (projectImpl ProjectImpl) GetPackageOptions() PackageOptions {
	if projectImpl.PackageOptions == nil {
		return nil
	}
	return projectImpl.PackageOptions
}

// //////////////////////

type ExecutableOptionsImpl struct {
	Source    []string      `mapstructure:"source"`
	MainFile  string        `mapstructure:"main_file"` // only for targets and not for tests
	Platform  string        `mapstructure:"platform"`
	Toolchain ToolchainImpl `mapstructure:"toolchain"`
}

func (executableOptionsImpl ExecutableOptionsImpl) GetSource() Sources {
	var sources Sources
	for _, source := range executableOptionsImpl.Source {
		sources = append(sources, ExpressionImpl{Value: strings.TrimSpace(source)})
	}

	return sources
}

func (executableOptionsImpl ExecutableOptionsImpl) GetMainFile(config *hil.EvalConfig) (string, error) {
	return applyEvaluator(executableOptionsImpl.MainFile, config)
}

func (executableOptionsImpl ExecutableOptionsImpl) GetPlatform(config *hil.EvalConfig) (string, error) {
	return applyEvaluator(executableOptionsImpl.Platform, config)
}

func (executableOptionsImpl ExecutableOptionsImpl) GetToolchain() Toolchain {
	return executableOptionsImpl.Toolchain
}

// //////////////////////

type TargetImpl struct {
	ExecutableOptions *ExecutableOptionsImpl `mapstructure:"executable_options"` // app only
	PackageOptions    *PackageOptionsImpl    `mapstructure:"package_options"`    // pkg only
	Arguments         []ArgumentImpl         `mapstructure:"arguments"`
	CompileOptions    *CompileOptionsImpl    `mapstructure:"compile_options"`
	LinkerOptions     LinkerOptionsImpl      `mapstructure:"linker_options"`
}

func (targetImpl TargetImpl) GetExecutableOptions() ExecutableOptions {
	if targetImpl.ExecutableOptions == nil {
		return nil
	}
	return targetImpl.ExecutableOptions
}

func (targetImpl TargetImpl) GetPackageOptions() PackageOptions {
	if targetImpl.PackageOptions == nil {
		return nil
	}
	return targetImpl.PackageOptions
}

func (targetImpl TargetImpl) GetArguments() Arguments {
	var arguments Arguments
	for _, argumentImpl := range targetImpl.Arguments {
		arguments = append(arguments, argumentImpl)
	}
	return arguments
}

func (targetImpl TargetImpl) GetCompileOptions() CompileOptions {
	if targetImpl.CompileOptions == nil {
		return nil
	}
	return targetImpl.CompileOptions
}

func (targetImpl TargetImpl) GetLinkerOptions() LinkerOptions {
	return targetImpl.LinkerOptions
}

// //////////////////////

type TestImpl struct {
	ExecutableOptions ExecutableOptionsImpl `mapstructure:"executable_options"`
	Arguments         []ArgumentImpl        `mapstructure:"arguments"`
	TargetName        string                `mapstructure:"target_name"`
	TargetArguments   []ArgumentImpl        `mapstructure:"target_arguments"`
	CompileOptions    CompileOptionsImpl    `mapstructure:"compile_options"`
	LinkerOptions     LinkerOptionsImpl     `mapstructure:"linker_options"`
}

func (testImpl TestImpl) GetExecutableOptions() ExecutableOptions {
	return testImpl.ExecutableOptions
}

func (testImpl TestImpl) GetArguments() Arguments {
	var arguments Arguments
	for _, argumentImpl := range testImpl.Arguments {
		arguments = append(arguments, argumentImpl)
	}
	return arguments
}

func (testImpl TestImpl) GetTargetName(config *hil.EvalConfig) (string, error) {
	return applyEvaluator(testImpl.TargetName, config)
}

func (testImpl TestImpl) GetTargetArguments() Arguments {
	var arguments Arguments
	for _, argumentImpl := range testImpl.TargetArguments {
		arguments = append(arguments, argumentImpl)
	}
	return arguments
}

func (testImpl TestImpl) GetCompileOptions() CompileOptions {
	return testImpl.CompileOptions
}

func (testImpl TestImpl) GetLinkerOptions() LinkerOptions {
	return testImpl.LinkerOptions
}

// //////////////////////

type projectConfigImpl struct {
	Type             string                     `mapstructure:"type"`
	Project          ProjectImpl                `mapstructure:"project"`
	Variables        []VariableImpl             `mapstructure:"variables"`
	Arguments        []ArgumentImpl             `mapstructure:"arguments"`
	Scripts          map[string]string          `mapstructure:"scripts"`
	Targets          map[string]*TargetImpl     `mapstructure:"targets"`
	Tests            map[string]*TestImpl       `mapstructure:"tests"`
	Dependencies     map[string]*DependencyImpl `mapstructure:"dependencies"`
	TestDependencies map[string]*DependencyImpl `mapstructure:"test_dependencies"`
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
		variableImpl.Name = strings.TrimSpace(variableImpl.Name)
		variableImpl.Value = strings.TrimSpace(variableImpl.Value)

		variables = append(variables, variableImpl)
	}
	return variables
}

func (projectConfigImpl *projectConfigImpl) GetArguments() Arguments {
	var arguments Arguments
	for _, argumentImpl := range projectConfigImpl.Arguments {
		argumentImpl.Name = strings.TrimSpace(argumentImpl.Name)
		argumentImpl.Value = strings.TrimSpace(argumentImpl.Value)

		arguments = append(arguments, argumentImpl)
	}
	return arguments
}

func (projectConfigImpl *projectConfigImpl) GetScripts() Scripts {
	scripts := Scripts{}
	for name, script := range projectConfigImpl.Scripts {
		scripts[name] = ExpressionImpl{Value: strings.TrimSpace(script)}
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
