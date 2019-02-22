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

	"github.com/inconshreveable/go-update"
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

const (
	currWioReleaseName = "wio{{extension}}"
	currWioReleaseUrl  = "https://github.com/wio/wio/releases/download/v{{version}}/wio_{{version}}_{{platform}}_{{arch}}.{{format}}"
)

var preArchMapping = map[string]string{
	"386":   "i386",
	"amd64": "x86_64",
	"arm":   "arm7",
	"arm64": "arm7",
}

var preOsMapping = map[string]string{
	"darwin":  "darwin",
	"windows": "windows",
	"linux":   "linux",
}

const (
	preWioReleaseName = "wio_{{platform}}_{{arch}}{{extension}}"
	preWioReleaseUrl  = "https://github.com/wio/wio/releases/download/v{{version}}/wio_{{platform}}_{{arch}}.{{format}}"
)

// Runs the build command when cli build option is provided
func (upgrade Upgrade) Execute() error {
	forceFlag := upgrade.Context.Bool("force")

	var version string
	var err error
	info := resolve.NewInfo(root.GetUpdatePath())
	if upgrade.Context.NArg() > 0 {
		version = upgrade.Context.Args()[0]
	} else {
		if version, err = info.GetLatest(constants.Wio); err != nil {
			return util.Error("no latest version found for wio")
		}
	}

	versionToUpgradeSem := semver.Parse(version)
	if versionToUpgradeSem == nil {
		return util.Error("wio version %s is invalid", version)
	}

	if !forceFlag && versionToUpgradeSem.LT(*semver.Parse("0.7.0")) {
		return util.Error("wio can only be upgraded/downgraded to versions >= 0.7.0")
	}

	version = versionToUpgradeSem.String()

	wioReleaseName := currWioReleaseName
	wioReleaseUrl := currWioReleaseUrl
	osMapping := currOsMapping
	archMapping := currArchMapping

	if version < "0.9.0" {
		wioReleaseName = preWioReleaseName
		wioReleaseUrl = preWioReleaseUrl
		osMapping = preOsMapping
		archMapping = preArchMapping
	}

	releaseName := template.Replace(wioReleaseName, map[string]string{
		"platform":  strings.ToLower(env.GetOS()),
		"arch":      strings.ToLower(archMapping[env.GetArch()]),
		"extension": extensionMapping[env.GetOS()],
	})

	releaseUrl := template.Replace(wioReleaseUrl, map[string]string{
		"version":  version,
		"platform": strings.ToLower(osMapping[env.GetOS()]),
		"arch":     strings.ToLower(archMapping[env.GetArch()]),
		"format":   strings.ToLower(formatMapping[env.GetOS()]),
	})

	if err := os.MkdirAll(sys.Path(root.GetUpdatePath(), sys.Download), os.ModePerm); err != nil {
		return util.Error("wio@%s error creating cache folder", version)
	}

	wioTarPath := sys.Path(root.GetUpdatePath(), sys.Download, fmt.Sprintf("%s__%s.%s", constants.Wio, version,
		formatMapping[env.GetOS()]))

	wioFolderPath := sys.Path(root.GetUpdatePath(), fmt.Sprintf("%s_%s", constants.Wio, version))

	if sys.Exists(wioTarPath) && forceFlag {
		os.RemoveAll(wioTarPath)
	}

	if sys.Exists(wioFolderPath) && forceFlag {
		os.RemoveAll(wioFolderPath)
	}

	log.Info(log.Cyan, "downloading wio version %s... ", version)
	if !sys.Exists(wioTarPath) {
		out, err := os.Create(wioTarPath)
		if err != nil {
			log.WriteFailure()
			return util.Error("wio@%s error creating %s file stream", version, formatMapping[env.GetOS()])
		}
		defer out.Close()

		resp, err := http.Get(releaseUrl)
		if err != nil || resp.StatusCode == 404 {
			log.WriteFailure()
			return util.Error("wio@%s version does not exist", version)
		}
		defer resp.Body.Close()

		if _, err := io.Copy(out, resp.Body); err != nil {
			log.WriteFailure()
			return util.Error("wio@%s error writing data to %s file", version, formatMapping[env.GetOS()])
		}

		log.WriteSuccess()
		log.Info(log.Cyan, "extracting wio version %s %s file... ", version, formatMapping[env.GetOS()])

		if err := archiver.Unarchive(wioTarPath, wioFolderPath); err != nil {
			log.WriteFailure()
			return util.Error("wio@%s error while extracting %s file", version, formatMapping[env.GetOS()])
		}
		log.WriteSuccess()
	} else {
		log.Infoln(log.Green, "already exists!")
	}

	log.Info(log.Cyan, "Updating ")
	log.Info(log.Green, "wio@%s ", meta.Version)
	log.Info(log.Cyan, "-> ")
	log.Info(log.Green, "wio@%s", version)
	log.Info(log.Cyan, "... ")

	newWioExec, err := os.Open(sys.Path(wioFolderPath, releaseName))
	if err != nil {
		log.WriteFailure()
		log.Write(err.Error())
		return util.Error("wio@%s executable not downloaded or corrupted", version)
	}

	if err = update.Apply(newWioExec, update.Options{}); err != nil {
		log.WriteFailure()
		log.Write(err.Error())
		return util.Error("failed to upgrade wio to version %s", version)
	} else {
		wioRootConfig := &root.WioRootConfig{}
		if err := sys.NormalIO.ParseJson(root.GetConfigFilePath(), wioRootConfig); err != nil {
			log.Write(err.Error())
			return util.Error("wio upgraded to %s but an error occurred and it's not complete", version)
		}
		wioRootConfig.Updated = true
		if err := sys.NormalIO.WriteJson(root.GetConfigFilePath(), wioRootConfig); err != nil {
			log.Write(err.Error())
			return util.Error("wio upgraded to %s but an error occurred and it's not complete", version)
		}
	}

	log.WriteSuccess()

	return nil
}
