package config

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hil"
)

type Dependencies map[string]Dependency
type Tests map[string]Test
type Targets map[string]Target
type Scripts []HilString
type Arguments []Argument
type Variables []Variable
type Definitions []HilString
type Flags []HilString
type Contributors []HilString
type Sources []HilString

type HilString interface {
	Get(config *hil.EvalConfig) (string, error)
}

type Toolchain interface {
	GetName() string
	GetRef(config *hil.EvalConfig) (string, error)
}

type ExecutableOptions interface {
	GetSource() Sources
	GetMainFile(config *hil.EvalConfig) (string, error)
	GetPlatform(config *hil.EvalConfig) (string, error)
	GetToolchain() Toolchain
}

type LinkerOptions interface {
	GetFlags() Flags
	GetVisibility(config *hil.EvalConfig) (string, error)
}

type CompileOptions interface {
	GetFlags() Flags
	GetDefinitions() Definitions
	GetCXXStandard(config *hil.EvalConfig) (string, error)
	GetCStandard(config *hil.EvalConfig) (string, error)
}

type PackageOptions interface {
	IsHeaderOnly() bool
	GetPackageType(config *hil.EvalConfig) (string, error)
}

type Dependency interface {
	GetRef(config *hil.EvalConfig) (string, error)
	GetArguments() Arguments
	GetLinkerOptions() LinkerOptions
}

type Test interface {
	GetExecutableOptions() ExecutableOptions
	GetArguments() Arguments
	GetTargetName(config *hil.EvalConfig) (string, error)
	GetTargetArguments() Arguments
	GetCompileOptions() CompileOptions
	GetLinkerOptions() LinkerOptions
}

type Target interface {
	GetExecutableOptions() ExecutableOptions
	GetPackageOptions() PackageOptions
	GetArguments() Arguments
	GetCompileOptions() CompileOptions
	GetLinkerOptions() LinkerOptions
}

type Argument interface {
	GetName() string
	GetValue(config *hil.EvalConfig) (string, error)
}

type Variable interface {
	GetName() string
	GetValue() string
}

type Project interface {
	GetName(config *hil.EvalConfig) (string, error)
	GetVersion(config *hil.EvalConfig) (*version.Version, error)
	GetAuthor(config *hil.EvalConfig) (string, error)
	GetDescription(config *hil.EvalConfig) (string, error)
	GetContributors() Contributors
	GetHomepage(config *hil.EvalConfig) (string, error)
	GetRepository(config *hil.EvalConfig) (string, error)

	GetCompileOptions() CompileOptions
	GetPackageOptions() PackageOptions
}

type ProjectConfig interface {
	GetType() string
	GetProject() Project
	GetVariables() Variables
	GetArguments() Arguments
	GetScripts() Scripts
	GetTargets() Targets
	GetTests() Tests
	GetDependencies() Dependencies
	GetTestDependencies() Dependencies
}
