package types

type PropertiesImpl struct {
    Global  []string `yaml:"global,omitempty"`
    Target  []string `yaml:"target,omitempty"`
    Package []string `yaml:"package,omitempty"`
}

func (p *PropertiesImpl) GetGlobal() []string {
    if p == nil {
        return []string{}
    }
    return p.Global
}

func (p *PropertiesImpl) GetTarget() []string {
    if p == nil {
        return []string{}
    }
    return p.Target
}

func (p *PropertiesImpl) GetPackage() []string {
    if p == nil {
        return []string{}
    }
    return p.Package
}

type TargetImpl struct {
    Source      string          `yaml:"src"`
    Platform    string          `yaml:"platform,omitempty"`
    Framework   string          `yaml:"framework,omitempty"`
    Board       string          `yaml:"board,omitempty"`
    Flags       *PropertiesImpl `yaml:"flags,omitempty"`
    Definitions *PropertiesImpl `yaml:"definitions,omitempty"`
    LinkerFlags []string        `yaml:"linker_flags,omitempty"`
    name        string
}

func (t *TargetImpl) GetSource() string {
    if t == nil {
        return ""
    }
    return t.Source
}

func (t *TargetImpl) GetPlatform() string {
    if t == nil {
        return ""
    }
    return t.Platform
}

func (t *TargetImpl) GetFramework() string {
    if t == nil {
        return ""
    }
    return t.Framework
}

func (t *TargetImpl) GetBoard() string {
    if t == nil {
        return ""
    }
    return t.Board
}

func (t *TargetImpl) GetFlags() Properties {
    return t.Flags
}

func (t *TargetImpl) GetDefinitions() Properties {
    return t.Definitions
}

func (t *TargetImpl) GetLinkerFlags() []string {
    return t.LinkerFlags
}

func (t *TargetImpl) GetName() string {
    return t.name
}

func (t *TargetImpl) SetName(name string) {
    t.name = name
}

type LibraryImpl struct {
    Package            bool              `yaml:"cmake_package,omitempty"`
    ImportedTargets    bool              `yaml:"use_imported_targets,omitempty"`
    Version            string            `yaml:"version,omitempty"`
    Required           bool              `yaml:"required,omitempty"`
    Variables          map[string]string `yaml:"variables,omitempty"`
    LibrariesTag       string            `yaml:"libraries_tag,omitempty"`
    IncludesTag        string            `yaml:"includes_tag,omitempty"`
    OsSupported        []string          `yaml:"os_supported,omitempty"`
    RequiredComponents []string          `yaml:"required_components,omitempty"`
    OptionalComponents []string          `yaml:"optional_components,omitempty"`
    Path               string            `yaml:"path,omitempty"`
    LibPath            []string          `yaml:"lib_path,omitempty"`
    IncludePath        []string          `yaml:"include_path,omitempty"`
    LinkVisibility     string            `yaml:"link_visibility,omitempty"`
    LinkerFlags        []string          `yaml:"linker_flags,omitempty"`
}

func (l *LibraryImpl) IsCmakePackage() bool {
    return l.Package
}

func (l *LibraryImpl) UseImportedTargets() bool {
    return l.ImportedTargets
}

func (l *LibraryImpl) GetVersion() string {
    return l.Version
}

func (l *LibraryImpl) IsRequired() bool {
    return l.Required
}

func (l *LibraryImpl) GetVariables() map[string]string {
    return l.Variables
}

func (l *LibraryImpl) GetLibrariesTag() string {
    return l.LibrariesTag
}

func (l *LibraryImpl) GetIncludesTag() string {
    return l.IncludesTag
}

func (l *LibraryImpl) GetOsSupported() []string {
    return l.OsSupported
}

func (l *LibraryImpl) GetRequiredComponents() []string {
    return l.RequiredComponents
}

func (l *LibraryImpl) GetOptionalComponents() []string {
    return l.OptionalComponents
}

func (l *LibraryImpl) GetPath() string {
    return l.Path
}

func (l *LibraryImpl) GetLibPath() []string {
    return l.LibPath
}

func (l *LibraryImpl) GetIncludePath() []string {
    return l.IncludePath
}

func (l *LibraryImpl) GetLinkVisibility() string {
    return l.LinkVisibility
}

func (l *LibraryImpl) GetLinkerFlags() []string {
    return l.LinkerFlags
}

type DependencyImpl struct {
    Vendor       bool     `yaml:"vendor,omitempty"`
    Version      string   `yaml:"version"`
    OsSupported  []string `yaml:"os_supported,omitempty"`
    Visibility   string   `yaml:"link_visibility,omitempty"`
    LinkerFlags  []string `yaml:"linker_flags,omitempty"`
    CompileFlags []string `yaml:"compile_flags,omitempty"`
    Definitions  []string `yaml:"definitions,omitempty"`
}

func (d *DependencyImpl) GetVersion() string {
    return d.Version
}

func (d *DependencyImpl) GetVisibility() string {
    return d.Visibility
}

func (d *DependencyImpl) GetOsSupported() []string {
    return d.OsSupported
}

func (d *DependencyImpl) GetCompileFlags() []string {
    return d.CompileFlags
}

func (d *DependencyImpl) GetLinkerFlags() []string {
    return d.LinkerFlags
}

func (d *DependencyImpl) GetDefinitions() []string {
    return d.Definitions
}

func (d *DependencyImpl) IsVendor() bool {
    return d.Vendor
}

type OptionsImpl struct {
    Version        string   `yaml:"wio_version"`
    Header         bool     `yaml:"header_only,omitempty"`
    Standard       string   `yaml:"standard,omitempty"`
    Default        string   `yaml:"default_target,omitempty"`
    Flags          []string `yaml:"flags,omitempty"`
    LinkerFlags    []string `yaml:"linker_flags,omitempty"`
    LinkVisibility string   `yaml:"link_visibility,omitempty"`
}

func (o *OptionsImpl) GetWioVersion() string {
    return o.Version
}

func (o *OptionsImpl) GetIsHeaderOnly() bool {
    return o.Header
}

func (o *OptionsImpl) GetStandard() string {
    return o.Standard
}

func (o *OptionsImpl) GetDefault() string {
    return o.Default
}

func (o *OptionsImpl) GetFlags() []string {
    return o.Flags
}

func (o *OptionsImpl) GetLinkerFlags() []string {
    return o.LinkerFlags
}

func (o *OptionsImpl) GetLinkVisibility() string {
    return o.LinkVisibility
}

type DefinitionSetImpl struct {
    Public  []string `yaml:"public,omitempty"`
    Private []string `yaml:"private,omitempty"`
}

func (d *DefinitionSetImpl) GetPublic() []string {
    if d == nil {
        return []string{}
    }
    return d.Public
}

func (d *DefinitionSetImpl) GetPrivate() []string {
    if d == nil {
        return []string{}
    }
    return d.Private
}

type DefinitionsImpl struct {
    Singleton bool               `yaml:"singleton,omitempty"`
    Global    *DefinitionSetImpl `yaml:"global,omitempty"`
    Required  *DefinitionSetImpl `yaml:"required,omitempty"`
    Optional  *DefinitionSetImpl `yaml:"optional,omitempty"`
    Ingest    *DefinitionSetImpl `yaml:"ingest,omitempty"`
}

func (d *DefinitionsImpl) IsSingleton() bool {
    if d == nil {
        return false
    }
    return d.Singleton
}

func (d *DefinitionsImpl) GetIngest() DefinitionSet {
    if d == nil {
        return &DefinitionSetImpl{}
    }
    return d.Ingest
}

func (d *DefinitionsImpl) GetGlobal() DefinitionSet {
    if d == nil {
        return &DefinitionSetImpl{}
    }
    return d.Global
}

func (d *DefinitionsImpl) GetRequired() DefinitionSet {
    if d == nil {
        return &DefinitionSetImpl{}
    }
    return d.Required
}

func (d *DefinitionsImpl) GetOptional() DefinitionSet {
    if d == nil {
        return &DefinitionSetImpl{}
    }
    return d.Optional
}

type InfoImpl struct {
    Name    string `yaml:"name"`
    Version string `yaml:"version"`

    Organization string   `yaml:"organization,omitempty"`
    Description  string   `yaml:"description,omitempty"`
    Repository   string   `yaml:"repository,omitempty"`
    Homepage     string   `yaml:"homepage,omitempty"`
    License      string   `yaml:"license,omitempty"`
    Author       string   `yaml:"author,omitempty"`
    Bugs         string   `yaml:"bugs,omitempty"`
    Contributors []string `yaml:"contributors,omitempty"`
    Keywords     []string `yaml:"keywords,omitempty"`
    IgnoreFiles  []string `yaml:"ignore_files,omitempty"`

    Options     *OptionsImpl     `yaml:"compile_options"`
    Definitions *DefinitionsImpl `yaml:"definitions,omitempty"`
}

func (i *InfoImpl) GetName() string {
    return i.Name
}

func (i *InfoImpl) GetVersion() string {
    return i.Version
}

func (i *InfoImpl) GetDescription() string {
    return i.Description
}

func (i *InfoImpl) GetRepository() string {
    return i.Repository
}

func (i *InfoImpl) GetHomepage() string {
    return i.Homepage
}

func (i *InfoImpl) GetLicense() string {
    return i.License
}

func (i *InfoImpl) GetAuthor() string {
    return i.Author
}

func (i *InfoImpl) GetBugs() string {
    return i.Bugs
}

func (i *InfoImpl) GetContributors() []string {
    return i.Contributors
}

func (i *InfoImpl) GetKeywords() []string {
    return i.Keywords
}

func (i *InfoImpl) GetIgnoreFiles() []string {
    return i.IgnoreFiles
}

func (i *InfoImpl) GetOptions() Options {
    return i.Options
}

func (i *InfoImpl) GetDefinitions() Definitions {
    return i.Definitions
}

type ConfigImpl struct {
    Type         string                     `yaml:"type"`
    Info         *InfoImpl                  `yaml:"project"`
    Targets      map[string]*TargetImpl     `yaml:"targets"`
    Dependencies map[string]*DependencyImpl `yaml:"dependencies,omitempty"`
    Libraries    map[string]*LibraryImpl    `yaml:"libraries,omitempty"`
}

func (c *ConfigImpl) GetType() string {
    return c.Type
}

func (c *ConfigImpl) GetName() string {
    return c.GetInfo().GetName()
}

func (c *ConfigImpl) GetVersion() string {
    return c.GetInfo().GetVersion()
}

func (c *ConfigImpl) GetInfo() Info {
    return c.Info
}

func (c *ConfigImpl) GetTargets() map[string]Target {
    if c.Targets == nil {
        c.Targets = map[string]*TargetImpl{}
    }
    s := map[string]Target{}
    for name, value := range c.Targets {
        s[name] = value
    }
    return s
}

func (c *ConfigImpl) GetDependencies() map[string]Dependency {
    if c.Dependencies == nil {
        c.Dependencies = map[string]*DependencyImpl{}
    }
    s := map[string]Dependency{}
    for name, value := range c.Dependencies {
        s[name] = value
    }
    return s
}

func (c *ConfigImpl) GetLibraries() map[string]Library {
    if c.Libraries == nil {
        c.Libraries = map[string]*LibraryImpl{}
    }
    s := map[string]Library{}
    for name, value := range c.Libraries {
        s[name] = value
    }
    return s
}

func (c *ConfigImpl) AddDependency(name string, dep Dependency) {
    if c.Dependencies == nil {
        c.Dependencies = map[string]*DependencyImpl{}
    }
    c.Dependencies[name] = dep.(*DependencyImpl)
}

func (c *ConfigImpl) DependencyMap() map[string]string {
    ret := map[string]string{}
    for name, dep := range c.GetDependencies() {
        ret[name] = dep.GetVersion()
    }
    return ret
}
