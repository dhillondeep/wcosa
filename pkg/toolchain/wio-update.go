package toolchain

import (
    "context"
    "fmt"
    "os"
    "runtime"
    "wio/internal/config"
    "wio/internal/config/meta"
    "wio/pkg/log"
    "wio/pkg/npm/semver"
    "wio/pkg/util"
    "wio/pkg/util/sys"
    "wio/pkg/util/template"

    "github.com/inconshreveable/go-update"
    "github.com/wayneashleyberry/terminal-dimensions"
)

const wioAssetStr = "wio_{{platform}}_{{arch}}.{{format}}"
const repoName = "wio"

var archMapping = map[string]string{
    "386":   "i386",
    "amd64": "x86_64",
    "arm-5": "arm5",
    "arm-6": "arm6",
    "arm-7": "arm7",
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

func getWioAssetString() string {
    platform := sys.GetOS()
    arch := archMapping[runtime.GOARCH]
    format := formatMapping[platform]

    return template.Replace(wioAssetStr, map[string]string{
        "platform": platform,
        "arch":     arch,
        "format":   format,
    })
}

func queryLatestVersion(platformWioStr string) (*assetRelease, error) {
    client := getClient()

    releases, resp, err := client.Repositories.ListReleases(context.Background(), Org, repoName, nil)
    if err != nil {
        return nil, util.Error("Error occurred while querying latest wio version")
    }

    if resp.StatusCode != 200 {
        return nil, nil
    }

    if len(releases) <= 0 {
        return nil, nil
    } else if len(releases[0].Assets) <= 0 {
        return nil, nil
    }

    currVersion := semver.Parse(meta.Version)

    for _, release := range releases {
        if semver.Parse((*release.Name)[1:]).Gt(currVersion) {
            for _, asset := range release.Assets {
                if *asset.Name == platformWioStr {
                    return &assetRelease{
                        Version:     (*release.Name)[1:],
                        Size:        *asset.Size,
                        DownloadURL: asset.GetBrowserDownloadURL(),
                    }, nil
                }
            }
        }
    }

    return nil, nil
}

func getWioReleaseAsset(version *semver.Version, platformWioStr string) (*assetRelease, error) {
    client := getClient()

    releases, resp, err := client.Repositories.ListReleases(context.Background(), Org, repoName, nil)
    if err != nil {
        return nil, util.Error("Error occurred while downloading wio executable for version: %s", version.Str())
    }

    if resp.StatusCode != 200 {
        return nil, nil
    }

    if len(releases) <= 0 {
        return nil, nil
    } else if len(releases[0].Assets) <= 0 {
        return nil, nil
    }

    // look for specific version
    for _, release := range releases {
        if *release.Name == "v"+version.Str() {
            for _, asset := range release.Assets {
                if *asset.Name == platformWioStr {
                    return &assetRelease{
                        Version:     version.Str(),
                        Size:        *asset.Size,
                        DownloadURL: asset.GetBrowserDownloadURL(),
                    }, nil
                }
            }
        }
    }

    return nil, util.Error("No wio executable found for version %s", version.Str())
}

func UpdateWioExecutable(version *semver.Version) error {
    platformWioStr := getWioAssetString()

    var wioAsset *assetRelease
    var err error

    if version == nil {
        wioAsset, err = queryLatestVersion(platformWioStr)
    } else {
        wioAsset, err = getWioReleaseAsset(version, platformWioStr)
    }
    if err != nil {
        return err
    }

    if wioAsset == nil {
        log.Writeln(log.Cyan, "Wio is already up to date: %s", meta.Version)
        return nil
    }

    wioNewVerPath := sys.Path(config.GetVersionsPath(), wioAsset.Version)
    wioNewVerTarPath := sys.Path(wioNewVerPath, platformWioStr)

    if err := os.MkdirAll(wioNewVerPath, os.ModePerm); err != nil {
        return err
    }

    log.Write(log.Cyan, "downloading wio version %s...", wioAsset.Version)

    if !sys.Exists(wioNewVerTarPath) {
        log.Writeln()
        if err := downloadAsset(wioAsset, wioNewVerTarPath); err != nil {
            return err
        }
    } else {
        log.Writeln(log.Green, " already downloaded")
    }

    if err := extractTarball(wioNewVerTarPath, wioNewVerPath, fmt.Sprintf("Extracting %s ", platformWioStr)); err != nil {
        log.WriteFailure()
        return err
    }
    log.WriteSuccess()

    newWioExecPath := sys.Path(wioNewVerPath, "wio"+extensionMapping[sys.GetOS()])

    log.Write(log.Cyan, "updating wio ")
    log.Write(log.Green, "%s", meta.Version)
    log.Write(log.Cyan, " => ")
    log.Write(log.Green, "%s", wioAsset.Version)
    log.Write(log.Cyan, "... ")

    newWioExec, err := os.Open(newWioExecPath)
    if err != nil {
        log.WriteFailure()
        return err
    }

    if err = update.Apply(newWioExec, update.Options{}); err != nil {
        log.WriteFailure()
        return err
    }
    log.WriteSuccess()
    os.Exit(0)

    return nil
}

func PromptWioExecUpdate() error {
    _, isSet := os.LookupEnv("WIO_NO_UPGRADE_PROMPT")
    if isSet {
        return nil
    }

    width, err := terminaldimensions.Width()
    if err != nil {
        return err
    }

    var i uint

    assetRelease, err := queryLatestVersion(getWioAssetString())
    if err != nil {
        return err
    }

    if assetRelease != nil {
        log.Writeln()

        // top
        for i = 0; i < width; i++ {
            log.Write(log.Green, "#")
        }

        log.Writeln(log.Cyan, "newer wio version available %s => %s\n", meta.Version, assetRelease.Version)
        log.Write(log.Green, "'wio upgrade'")
        log.Writeln(log.Cyan, " to upgrade to latest version")

        // bottom
        for i = 0; i < width; i++ {
            log.Write(log.Green, "#")
        }
    }

    return nil
}
