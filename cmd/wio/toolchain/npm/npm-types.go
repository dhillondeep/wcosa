package npm

type packageData struct {
	time           map[string]string `json:"time"`
	name           string            `json:"name"`
	distTags       map[string]string `json:"dist-tags"`
	versions       []packageVersion  `json:"versions"`
	maintainers    []packageAuthor   `json:"maintainers"`
	keywords       []string          `json:"keywords"`
	repository     packageRepository `json:"repository"`
	contributors   []packageAuthor   `json:"contributors"`
	author         packageAuthor     `json:"author"`
	bugs           packageBug        `json:"bugs"`
	license        string            `json:"license"`
	readme         string            `json:"readme"`
	readmeFilename string            `json:"readmeFilename"`
}

type packageRepository struct {
	pType string `json:"type"`
	url   string `json:"url"`
}

type packageAuthor struct {
	name  string `json:"name"`
	email string `json:"email"`
	url   string `json:"url"`
}

type packageBug struct {
	url string `json:"url"`
}

type packageDist struct {
	integrity    string `json:"integrity"`
	shasum       string `json:"shasum"`
	tarball      string `json:"tarball"`
	fileCount    int    `json:"fileCount"`
	unpackedSize int    `json:"unpackedSize"`
	npmSignature string `json:"npm-signature"`
}

type packageVersion struct {
	name         string            `json:"name"`
	version      string            `json:"version"`
	description  string            `json:"description"`
	repository   packageRepository `json:"repository"`
	main         string            `json:"main"`
	keywords     []string          `json:"keywords"`
	author       packageAuthor     `json:"author"`
	license      string            `json:"license"`
	contributors []packageAuthor   `json:"contributors"`
	dependencies map[string]string `json:"dependencies"`
	bugs         packageBug        `json:"bugs"`
	homepage     string            `json:"homepage"`
	dist         packageDist       `json:"dist"`
	maintainers  []packageAuthor   `json:"maintainers"`
}
