package types

type PropertiesImpl struct {
	Global  []string `yaml:"global,omitempty"`
	Target  []string `yaml:"target,omitempty"`
	Package []string `yaml:"package,omitempty"`
}

func (p *PropertiesImpl) GetGlobal() []string {
	return p.Global
}

func (p *PropertiesImpl) GetTarget() []string {
	return p.Target
}

func (p *PropertiesImpl) GetPackage() []string {
	return p.Package
}

type TargetImpl struct {
	Source      string          `yaml:"src"`
	Platform    string          `yaml:"platform"`
	Framework   string          `yaml:"framework,omitempty"`
	Board       string          `yaml:"board,omitempty"`
	Flags       *PropertiesImpl `yaml:"flags,omitempty"`
	Definitions *PropertiesImpl `yaml:"definitions,omitempty"`
}

func (t *TargetImpl) GetSource() string {
	return t.Source
}

func (t *TargetImpl) GetPlatform() string {
	return t.Platform
}

func (t *TargetImpl) GetFramework() string {
	return t.Framework
}

func (t *TargetImpl) GetBoard() string {
	return t.Board
}

func (t *TargetImpl) GetFlags() Properties {
	return t.Flags
}

func (t *TargetImpl) GetDefinitions() Properties {
	return t.Definitions
}

type DependencyImpl struct {
	Vendor      bool     `yaml:"vendor,omitempty"`
	Version     string   `yaml:"version"`
	Visibility  string   `yaml:"link_visibility,omitempty"`
	Flags       []string `yaml:"compile_flags,omitempty"`
	Definitions []string `yaml:"definitions"`
}

func (d *DependencyImpl) GetVersion() string {
	return d.Version
}

func (d *DependencyImpl) GetVisibility() string {
	return d.Visibility
}

func (d *DependencyImpl) GetFlags() []string {
	return d.Flags
}

func (d *DependencyImpl) GetDefinitions() []string {
	return d.Definitions
}

func (d *DependencyImpl) GetIsVendor() bool {
	return d.Vendor
}

type OptionsImpl struct {
	Version  string   `yaml:"wio_version"`
	Header   bool     `yaml:"header_only,omitempty"`
	Standard string   `yaml:"standard,omitempty"`
	Default  string   `yaml:"default_target,omitempty"`
	Flags    []string `yaml:"flags,omitempty"`
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

type DefinitionSetImpl struct {
	Public  []string `yaml:"public,omitempty"`
	Private []string `yaml:"private,omitempty"`
}

func (d *DefinitionSetImpl) GetPublic() []string {
	return d.Public
}

func (d *DefinitionSetImpl) GetPrivate() []string {
	return d.Private
}

type DefinitionsImpl struct {
	Singleton bool               `yaml:"singleton"`
	Global    *DefinitionSetImpl `yaml"global"`
	Required  *DefinitionSetImpl `yaml:"required"`
	Optional  *DefinitionSetImpl `yaml:"optional"`
}

func (d *DefinitionsImpl) GetIsSingleton() bool {
	return d.Singleton
}

func (d *DefinitionsImpl) GetGlobal() DefinitionSet {
	return d.Global
}

func (d *DefinitionsImpl) GetRequired() DefinitionSet {
	return d.Required
}

func (d *DefinitionsImpl) GetOptional() DefinitionSet {
	return d.Optional
}

type InfoImpl struct {
	Name        string           `yaml:"name"`
	Version     string           `yaml:"version"`
	Keywords    []string         `yaml:"keywords,omitempty"`
	License     string           `yaml:"license,omitempty"`
	Options     *OptionsImpl     `yaml:"compile_options"`
	Definitions *DefinitionsImpl `yaml:"definitions"`
}

func (i *InfoImpl) GetName() string {
	return i.Name
}

func (i *InfoImpl) GetVersion() string {
	return i.Version
}

func (i *InfoImpl) GetKeywords() []string {
	return i.Keywords
}

func (i *InfoImpl) GetLicense() string {
	return i.License
}

func (i *InfoImpl) GetOptions() Options {
	return i.Options
}

func (i *InfoImpl) GetDefinitions() Definitions {
	return i.Definitions
}

type ConfigImpl struct {
	Type         string                `yaml:"type"`
	Info         *InfoImpl             `yaml:"project"`
	Targets      map[string]Target     `yaml:"targets"`
	Dependencies map[string]Dependency `yaml:"dependencies"`
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
	return c.Targets
}

func (c *ConfigImpl) GetDependencies() map[string]Dependency {
	return c.Dependencies
}

func (c *ConfigImpl) DependencyMap() map[string]string {
	ret := map[string]string{}
	for name, dep := range c.GetDependencies() {
		ret[name] = dep.GetVersion()
	}
	return ret
}
