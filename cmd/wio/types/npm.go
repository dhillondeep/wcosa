package types

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
