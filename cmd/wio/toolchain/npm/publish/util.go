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
    "time"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/toolchain/npm"
    "wio/cmd/wio/toolchain/npm/login"
    "wio/cmd/wio/toolchain/npm/semver"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"

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
    src := io.Path(dir, io.Folder, "package")
    content := []string{
        io.Path(src, "include"),
        io.Path(src, "src"),
        io.Path(src, "package.json"),
        io.Path(src, "wio.yml"),
        io.Path(src, "README.md"),
        io.Path(src, ".wio.js"),
    }
    return archiver.TarGz.Make(dst, content)
}

func GeneratePackage(dir string, data *npm.Version) error {
    pkg := io.Path(dir, io.Folder, "package")
    if err := os.RemoveAll(pkg); err != nil {
        return err
    }
    if err := os.MkdirAll(pkg, os.ModePerm); err != nil {
        return err
    }
    dst := io.Path(pkg, "package.json")
    if !io.Exists(dst) {
        data, _ := json.MarshalIndent(data, "", login.Indent)
        if err := ioutil.WriteFile(dst, data, os.ModePerm); err != nil {
            return err
        }
    }
    dst = io.Path(pkg, ".wio.js")
    if !io.Exists(dst) {
        str := "console.log('Hi!!! Welcome to Wio world!')\n"
        if err := ioutil.WriteFile(dst, []byte(str), os.ModePerm); err != nil {
            return err
        }
    }
    copies := []string{"include", "src", "wio.yml", "README.md"}
    for _, copy := range copies {
        if err := io.Copy(io.Path(dir, copy), io.Path(pkg, copy)); err != nil {
            return err
        }
    }
    return nil
}

// Generate version data (package.json) based on project config (wio.yml)
// and parses the README.md file. Also verifies package and dependencies.
func VersionData(dir string, cfg types.Config) (*npm.Version, error) {
    if cfg.GetType() != constants.PKG {
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
    readme, err := ioutil.ReadFile(io.Path(dir, "README.md"))
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
    }, nil
}
