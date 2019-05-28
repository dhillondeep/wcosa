package upgrade

const (
	wioExecName             = "wio{{extension}}"
	preWioExecName          = "wio_{{platform}}_{{arch}}{{extension}}"
	wioAssetName            = "wio_.*_{{platform}}_{{arch}}.{{format}}"
	preWioAssetName         = "wio_{{platform}}_{{arch}}.{{format}}"
	checksumFile            = "asset:wio_.*_checksums.txt"
	preChecksumFile         = "asset:checksums.txt"
	compatibilityLowerBound = "0.7.0"
	assetNameChangeVersion  = "0.9.0"
)

var currArchMapping = map[string]string{
	"386":   "32bit",
	"amd64": "64bit",
	"arm":   "arm",
	"arm64": "arm64",
}

var currOsMapping = map[string]string{
	"darwin":  "macOS",
	"windows": "windows",
	"linux":   "linux",
}

var formatMapping = map[string]string{
	"windows": "zip",
	"linux":   "tar.gz",
	"darwin":  "tar.gz",
}

var extensionMapping = map[string]string{
	"windows": ".exe",
	"linux":   "",
	"darwin":  "",
}

// preArchMapping holds arch mappings for versions < 0.9.0
var preArchMapping = map[string]string{
	"386":   "i386",
	"amd64": "x86_64",
	"arm":   "arm7",
	"arm64": "arm7",
}

// preOsMapping holds os mappings for versions < 0.9.0
var preOsMapping = map[string]string{
	"darwin":  "darwin",
	"windows": "windows",
	"linux":   "linux",
}
