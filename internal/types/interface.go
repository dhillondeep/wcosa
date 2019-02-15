package types

const (
	Private   = "PRIVATE"
	Public    = "PUBLIC"
	Interface = "INTERFACE"
)

type Properties interface {
	GetGlobal() []string
	GetTarget() []string
	GetPackage() []string
}

type Target interface {
	GetSource() string
	GetPlatform() string
	GetFramework() string
	GetBoard() string
	GetFlags() Properties
	GetDefinitions() Properties
	GetLinkerFlags() []string

	GetName() string
	SetName(name string)
}

type Library interface {
	IsCmakePackage() bool
	UseImportedTargets() bool
	GetVersion() string
	IsRequired() bool
	GetVariables() map[string]string
	GetLibrariesTag() string
	GetIncludesTag() string
	GetOsSupported() []string
	GetRequiredComponents() []string
	GetOptionalComponents() []string
	GetPath() string
	GetLibPath() []string
	GetIncludePath() []string
	GetLinkVisibility() string
	GetLinkerFlags() []string
}

type Dependency interface {
	IsVendor() bool
	GetVersion() string
	GetOsSupported() []string
	GetVisibility() string
	GetLinkerFlags() []string
	GetCompileFlags() []string
	GetDefinitions() []string
}

type Options interface {
	GetWioVersion() string
	GetIsHeaderOnly() bool
	GetStandard() string
	GetDefault() string
	GetFlags() []string
	GetLinkerFlags() []string
	GetLinkVisibility() string
}

type DefinitionSet interface {
	GetPublic() []string
	GetPrivate() []string
}

type Definitions interface {
	IsSingleton() bool
	GetGlobal() DefinitionSet
	GetRequired() DefinitionSet
	GetOptional() DefinitionSet
	GetIngest() DefinitionSet
}

type Info interface {
	GetName() string
	GetVersion() string

	GetDescription() string
	GetRepository() string
	GetHomepage() string
	GetLicense() string
	GetAuthor() string
	GetBugs() string
	GetContributors() []string
	GetKeywords() []string
	GetIgnoreFiles() []string

	GetOptions() Options
	GetDefinitions() Definitions
}

type Config interface {
	GetType() string
	GetName() string
	GetVersion() string
	SetVersion(version string)

	GetInfo() Info
	GetTargets() map[string]Target
	GetDependencies() map[string]Dependency
	GetLibraries() map[string]Library

	AddDependency(name string, dep Dependency)

	DependencyMap() map[string]string
}
