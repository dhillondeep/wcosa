package npm

type Data struct {
    Time           map[string]string  `json:"time"`
    Name           string             `json:"name"`
    DistTags       map[string]string  `json:"dist-tags"`
    Versions       map[string]Version `json:"versions"`
    Maintainers    []Author           `json:"maintainers"`
    Keywords       []string           `json:"keywords"`
    Bugs           Bug                `json:"bugs"`
    Readme         string             `json:"readme"`
    ReadmeFilename string             `json:"readmeFilename"`
    Error          string             `json:"error"`

    //Repository   Repository `json:"repository"`
    //Contributors []Author   `json:"contributors"`
    //Author       Author     `json:"author"`
    //License      string     `json:"license"`
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

type Bug struct {
    Url string `json:"url"`
}

type Dist struct {
    Integrity    string `json:"integrity"`
    ShaSum       string `json:"shasum"`
    Tarball      string `json:"tarball"`
    FileCount    int    `json:"fileCount"`
    UnpackedSize int    `json:"unpackedSize"`
    NpmSignature string `json:"npm-signature"`
}

type Version struct {
    Name         string            `json:"name"`
    Version      string            `json:"version"`
    Description  string            `json:"description"`
    Main         string            `json:"main"`
    Keywords     []string          `json:"keywords"`
    Dependencies map[string]string `json:"dependencies"`
    Bugs         Bug               `json:"bugs"`
    Dist         Dist              `json:"dist"`
    Maintainers  []Author          `json:"maintainers"`

    //Repository   Repository `json:"repository"`
    //Author       Author     `json:"author"`
    //License      string     `json:"license"`
    //Contributors []Author   `json:"contributors"`
    //Homepage     string     `json:"homepage"`
}
