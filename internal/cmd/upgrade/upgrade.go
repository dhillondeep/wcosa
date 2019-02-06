package upgrade

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
    "wio/internal/config/meta"
    "wio/internal/config/root"
    "wio/internal/constants"
    "wio/internal/env"
    "wio/pkg/log"
    "wio/pkg/npm/resolve"
    "wio/pkg/npm/semver"
    "wio/pkg/util"
    "wio/pkg/util/sys"
    "wio/pkg/util/template"

    update "github.com/inconshreveable/go-update"
    "github.com/mholt/archiver"
    "github.com/urfave/cli"
)

type Upgrade struct {
    Context *cli.Context
}

// get context for the command
func (upgrade Upgrade) GetContext() *cli.Context {
    return upgrade.Context
}

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

const (
    wioReleaseName = "wio_{{platform}}_{{arch}}{{format}}"
    wioReleaseUrl = "https://github.com/wio/wio/releases/download/v{{version}}/wio_{{platform}}_{{arch}}.{{format}}"

)

// Runs the build command when cli build option is provided
func (upgrade Upgrade) Execute() error {
    var version string
    var err error
    info := resolve.NewInfo(root.GetUpdatePath())
    if upgrade.Context.NArg() > 0 {
        version = upgrade.Context.Args()[0]
    } else {
        if version, err = info.GetLatest(constants.Wio); err != nil {
            return util.Error("no version found for wio")
        }
    }

    versionToUpgradeSem := semver.Parse(version)
    if versionToUpgradeSem == nil {
        return util.Error("wio version %s is invalid", version)
    }

    if !upgrade.Context.Bool("force") && versionToUpgradeSem.LT(*semver.Parse("0.7.0")) {
        return util.Error("wio can only be upgraded/downgraded to versions >= 0.7.0")
    }

    version = versionToUpgradeSem.String()

    releaseName := template.Replace(wioReleaseName, map[string]string{
        "platform": strings.ToLower(env.GetOS()),
        "arch": strings.ToLower(env.GetOS()),
        "format": formatMapping[env.GetOS()],
    })

    releaseUrl := template.Replace(wioReleaseUrl, map[string]string{
        "version":  version,
        "platform": strings.ToLower(env.GetOS()),
        "arch":     strings.ToLower(env.GetOS()),
        "format":   formatMapping[env.GetOS()],
    })

    resp, err := http.Get(releaseUrl)
    if err != nil {
        log.WriteFailure()
        return util.Error("wio@%s version does not exist", version)
    }
    defer resp.Body.Close()

    wioTarPath := sys.Path(root.GetUpdatePath(), sys.Download, fmt.Sprintf("%s__%s.%s", constants.Wio, version,
        formatMapping[env.GetOS()]))

    wioFolderPath := sys.Path(root.GetUpdatePath(), fmt.Sprintf("%s_%s", constants.Wio, version))

    if err := os.MkdirAll(sys.Path(root.GetUpdatePath(), sys.Download), os.ModePerm); err != nil {
        log.WriteFailure()
        return util.Error("wio@%s executable tar could not extracted", version)
    }

    if !sys.Exists(wioTarPath) {
        out, err := os.Create(wioTarPath)
        if err != nil {
            log.WriteFailure()
            return util.Error("wio@%s executable tar could not extracted", version)
        }
        defer out.Close()

        if _, err := io.Copy(out, resp.Body); err != nil {
            log.WriteFailure()
            return util.Error("wio@%s executable tar could not extracted", version)
        }
    }

    if !sys.Exists(wioFolderPath) {
        if err := archiver.Unarchive(wioTarPath, wioFolderPath); err != nil {
            log.WriteFailure()
            return util.Error("wio@%s executable tar could not extracted", version)
        }
    }

    log.WriteSuccess()
    log.Info(log.Cyan, "Updating ")
    log.Info(log.Green, "wio@%s ", meta.Version)
    log.Info(log.Cyan, "-> ")
    log.Info(log.Green, "wio@%s", version)
    log.Info(log.Cyan, "... ")

    newWioExec, err := os.Open(sys.Path(wioFolderPath, releaseName))
    if err != nil {
        log.WriteFailure()
        return util.Error("wio@%s executable data is corrupted ", version)
    }

    if err = update.Apply(newWioExec, update.Options{}); err != nil {
        log.WriteFailure()
        return util.Error("failed to upgrade wio to version %s", version)
    } else {
        wioRootConfig := &root.WioRootConfig{}
        if err := sys.NormalIO.ParseJson(root.GetConfigFilePath(), wioRootConfig); err != nil {
            return util.Error("wio upgraded to %s but an error occurred and it's not complete", version)
        }
        wioRootConfig.Updated = true
        if err := sys.NormalIO.WriteJson(root.GetConfigFilePath(), wioRootConfig); err != nil {
            return util.Error("wio upgraded to %s but an error occurred and it's not complete", version)
        }
    }

    log.WriteSuccess()

    return nil
}
