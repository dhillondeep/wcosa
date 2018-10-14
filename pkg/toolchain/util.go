package toolchain

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "os"
    "time"
    "wio/pkg/log"
    "wio/pkg/util"
    "wio/pkg/util/sys"

    "github.com/google/go-github/github"
    "github.com/mholt/archiver"
    "github.com/tj/go-spin"
    "gopkg.in/cheggaaa/pb.v1"
)

const (
    Org                 = "wio"
    FrameworkConfigName = "framework.json"
)

type AssetRelease struct {
    RemoteName  string
    LocalName   string
    Version     string
    Size        int
    DownloadURL string
}

type githubFetcher struct {
    client *github.Client
}

var githubFetcherImpl = githubFetcher{}

func GetClient() *github.Client {
    if githubFetcherImpl.client == nil {
        githubFetcherImpl.client = github.NewClient(nil)
    }

    return githubFetcherImpl.client
}

func GetRemoteAsset(repoName string, localName string, version string) (*AssetRelease, error) {
    client := GetClient()

    releases, resp, err := client.Repositories.ListReleases(context.Background(), Org, repoName, nil)
    if err != nil {
        return nil, util.Error("Error occurred while downloading toolchain/framework")
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
    if !util.IsEmptyString(version) {
        for _, release := range releases {
            if *release.Name == "v"+version {
                return &AssetRelease{
                    RemoteName:  *release.Assets[0].Name,
                    LocalName:   fmt.Sprintf("%s@%s", localName, version),
                    Version:     version,
                    Size:        *release.Assets[0].Size,
                    DownloadURL: release.Assets[0].GetBrowserDownloadURL(),
                }, nil
            }
        }

        return nil, nil
    }

    frameworkVersion := (*releases[0].Name)[1:]

    return &AssetRelease{
        RemoteName:  *releases[0].Assets[0].Name,
        LocalName:   fmt.Sprintf("%s@%s", localName, frameworkVersion),
        Version:     frameworkVersion,
        Size:        *releases[0].Assets[0].Size,
        DownloadURL: releases[0].Assets[0].GetBrowserDownloadURL(),
    }, nil

    return nil, nil
}

func DownloadAsset(releaseAsset *AssetRelease, path string) error {
    if releaseAsset == nil || util.IsEmptyString(releaseAsset.DownloadURL) {
        return util.Error("toolchain/framework does not exist")
    }

    bar := pb.New(releaseAsset.Size).SetUnits(pb.U_BYTES)
    bar.Start()

    // Create the file
    out, err := os.Create(sys.Path(path))
    if err != nil {
        return err
    }
    defer out.Close()

    // Get the data
    resp, err := http.Get(releaseAsset.DownloadURL)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // create proxy reader
    reader := bar.NewProxyReader(resp.Body)

    // and copy
    if _, err := io.Copy(out, reader); err != nil {
        return err
    }

    bar.Finish()

    return nil
}

func ExtractTarball(srcPath, destPath, description string) error {
    // a channel to tell it to stop
    stopchan := make(chan struct{})
    // a channel to signal that it's stopped
    stoppedchan := make(chan struct{})

    s := spin.New()
    s.Set(spin.Spin9)
    go func() {
        defer close(stoppedchan)
        for {
            select {
            default:
                log.Write(log.Cyan, "\r"+description)
                log.Write(log.Default, "%s ", s.Next())
                time.Sleep(100 * time.Millisecond)
            case <-stopchan:
                return
            }

        }
    }()

    err := archiver.TarGz.Open(srcPath, destPath)
    if err != nil {
        return err
    }

    close(stopchan)

    return nil
}
