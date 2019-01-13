package npm

type Data struct {
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Keywords    []string `json:"keywords"`
    Readme      string   `json:"readme"`

    Time     map[string]string  `json:"time"`
    DistTags map[string]string  `json:"dist-tags"`
    Versions map[string]Version `json:"versions"`

    Maintainers  interface{} `json:"maintainers"`
    Contributors interface{} `json:"contributors"`
    Bugs         interface{} `json:"bugs"`
    Author       interface{} `json:"author"`
    License      interface{} `json:"license"`
    Homepage     interface{} `json:"homepage"`
    Repository   interface{} `json:"repository"`

    Error string `json:"error"`
}

type Dist struct {
    Integrity    string `json:"integrity"`
    Shasum       string `json:"shasum"`
    Tarball      string `json:"tarball"`
    FileCount    int    `json:"fileCount"`
    UnpackedSize int    `json:"unpackedSize"`
    NpmSignature string `json:"npm-signature"`
}

type Version struct {
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Keywords    []string `json:"keywords"`
    Readme      string   `json:"readme"`
    ReadmeFile  string   `json:"readmeFile"`

    Version string `json:"version"`
    Main    string `json:"main"`
    Dist    Dist   `json:"dist"`

    Scripts      map[string]string `json:"scripts"`
    Dependencies map[string]string `json:"dependencies"`

    Maintainers  interface{} `json:"maintainers"`
    Contributors interface{} `json:"contributors"`
    Bugs         interface{} `json:"bugs"`
    Author       interface{} `json:"author"`
    License      interface{} `json:"license"`
    Homepage     interface{} `json:"homepage"`
    Repository   interface{} `json:"repository"`

    IgnorePaths []string `json:"ignore-files"`
}

type Repository struct {
    Type string `json:"type"`
    Url  string `json:"url"`
}

type Author struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Url   string `json:"url"`
}

type Bugs struct {
    Url   string `json:"url"`
    Email string `json:"email"`
}

type License struct {
    License string `json:"license"`
    Type    string `json:"type"`
    Url     string `json:"url"`
}
