package types

const (
    Private   = "PRIVATE"
    Public    = "PUBLIC"
    Interface = "PUBLIC"
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

    GetName() string
    SetName(name string)
}

type Dependency interface {
    GetVersion() string
    GetVisibility() string
    GetFlags() []string
    GetDefinitions() []string
    IsVendor() bool
}

type Options interface {
    GetWioVersion() string
    GetIsHeaderOnly() bool
    GetStandard() string
    GetDefault() string
    GetFlags() []string
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
}

type Info interface {
    GetName() string
    GetVersion() string
    GetKeywords() []string
    GetLicense() string
    GetOptions() Options
    GetDefinitions() Definitions
}

type Config interface {
    GetType() string
    GetName() string
    GetVersion() string

    GetInfo() Info
    GetTargets() map[string]Target
    GetDependencies() map[string]Dependency

    DependencyMap() map[string]string
}
