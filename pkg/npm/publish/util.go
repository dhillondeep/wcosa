package publish

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"wio/internal/constants"
	"wio/internal/types"
	"wio/pkg/npm"
	"wio/pkg/npm/login"
	"wio/pkg/npm/semver"
	"wio/pkg/util/sys"

	"github.com/mholt/archiver"
)

const ByteMax = 1 << 8
const ByteLen = 8

var Encoder = base64.StdEncoding.WithPadding(base64.NoPadding)

// Sessions are identified by a random 8-byte base64
// encoded string. Generate one seeded with current time.
func GenerateSession() string {
	value := make([]byte, 0, ByteLen)
	ret := make([]byte, 0, 2*ByteLen)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len(value); i++ {
		value = append(value, byte(rand.Int()%ByteMax))
	}
	Encoder.Encode(ret, value)
	return string(ret)
}

// Perform SHA1 checksum on the package tarball and return
// in base64 encoded form.
func Shasum(data []byte) string {
	ret := sha1.Sum(data)
	return hex.EncodeToString(ret[:])
}

func TarEncode(data []byte) string {
	ret := make([]byte, Encoder.EncodedLen(len(data)))
	Encoder.Encode(ret, data)
	return string(ret)
}

func MakeTar(dir, dst string) error {
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}
	content := sys.Path(dir, sys.WioFolder, "package")
	return archiver.Archive([]string{content}, dst)
}

func GeneratePackage(dir string, data *npm.Version) error {
	pkg := sys.Path(dir, sys.WioFolder, "package")
	if err := os.RemoveAll(pkg); err != nil {
		return err
	}
	if err := os.MkdirAll(pkg, os.ModePerm); err != nil {
		return err
	}
	dst := sys.Path(pkg, "package.json")
	if !sys.Exists(dst) {
		data, _ := json.MarshalIndent(data, "", login.Indent)
		if err := ioutil.WriteFile(dst, data, os.ModePerm); err != nil {
			return err
		}
	}
	dst = sys.Path(pkg, ".wio.js")
	if !sys.Exists(dst) {
		str := "console.log('Hi!!! Welcome to Wio world!')\n"
		if err := ioutil.WriteFile(dst, []byte(str), os.ModePerm); err != nil {
			return err
		}
	}

	var ignorePathsReg []*regexp.Regexp
	for _, ignorePath := range data.IgnorePaths {
		isDir, err := sys.IsDir(ignorePath)
		if err != nil {
			return err
		}

		if isDir {
			ignorePath = sys.Path(ignorePath, "*")
		}

		ignorePathsReg = append(ignorePathsReg, regexp.MustCompile(
			strings.Replace(ignorePath, "/", `\/`, -1)))
	}

	ignorePathsReg = append(ignorePathsReg, regexp.MustCompile(sys.Path(dir, `\.+.+`)))

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if path == dir {
			return nil
		}

		// ignore all the paths specified
		for _, ignorePathReg := range ignorePathsReg {
			if ignorePathReg.MatchString(path) {
				return nil
			}
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		if err := sys.Copy(path, sys.Path(pkg, relPath)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// Generate version data (package.json) based on project config (wio.yml)
// and parses the README.md file. Also verifies package and dependencies.
func VersionData(dir string, cfg types.Config) (*npm.Version, error) {
	if cfg.GetType() != constants.Pkg {
		return nil, InvalidProjectType{}
	}
	info := cfg.GetInfo()
	if semver.Parse(info.GetVersion()) == nil {
		return nil, InvalidProjectVersion{info.GetVersion()}
	}
	deps := cfg.DependencyMap()
	for name, ver := range deps {
		if semver.MakeQuery(ver) == nil {
			return nil, InvalidDependencyVersion{name, ver}
		}
	}
	readme, err := ioutil.ReadFile(sys.Path(dir, "README.md"))
	if err != nil {
		return nil, err
	}
	return &npm.Version{
		Name:        info.GetName(),
		Description: info.GetDescription(),
		Keywords:    info.GetKeywords(),
		Readme:      string(readme),
		ReadmeFile:  "README.md",

		Version: info.GetVersion(),
		Main:    ".wio.js",

		Dependencies: deps,
		Contributors: info.GetContributors(),
		Bugs:         info.GetBugs(),
		Author:       info.GetAuthor(),
		License:      info.GetLicense(),
		Homepage:     info.GetHomepage(),
		Repository:   info.GetRepository(),

		IgnorePaths: info.GetIgnoreFiles(),
	}, nil
}
