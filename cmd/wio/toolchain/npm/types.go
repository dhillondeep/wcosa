package npm

type packageData struct {
    Time           map[string]string         `json:"time"`
    Name           string                    `json:"name"`
    DistTags       map[string]string         `json:"dist-tags"`
    Versions       map[string]packageVersion `json:"versions"`
    Maintainers    []packageAuthor           `json:"maintainers"`
    Keywords       []string                  `json:"keywords"`
    Repository     packageRepository         `json:"repository"`
    Contributors   []packageAuthor           `json:"contributors"`
    Author         packageAuthor             `json:"author"`
    Bugs           packageBug                `json:"bugs"`
    License        string                    `json:"license"`
    Readme         string                    `json:"readme"`
    ReadmeFilename string                    `json:"readmeFilename"`
    Error          string                    `json:"error"`
}

type packageRepository struct {
    PType string `json:"type"`
    Url   string `json:"url"`
}

type packageAuthor struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Url   string `json:"url"`
}

type packageBug struct {
    Url string `json:"url"`
}

type packageDist struct {
    Integrity    string `json:"integrity"`
    Shasum       string `json:"shasum"`
    Tarball      string `json:"tarball"`
    FileCount    int    `json:"fileCount"`
    UnpackedSize int    `json:"unpackedSize"`
    NpmSignature string `json:"npm-signature"`
}

type packageVersion struct {
    Name         string            `json:"name"`
    Version      string            `json:"version"`
    Description  string            `json:"description"`
    Repository   packageRepository `json:"repository"`
    Main         string            `json:"main"`
    Keywords     []string          `json:"keywords"`
    Author       packageAuthor     `json:"author"`
    License      string            `json:"license"`
    Contributors []packageAuthor   `json:"contributors"`
    Dependencies map[string]string `json:"dependencies"`
    Bugs         packageBug        `json:"bugs"`
    Homepage     string            `json:"homepage"`
    Dist         packageDist       `json:"dist"`
    Maintainers  []packageAuthor   `json:"maintainers"`
}
