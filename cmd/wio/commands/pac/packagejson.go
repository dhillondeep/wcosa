package pac

import (
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/types"
    "wio/cmd/wio/log"
    "regexp"
    goerr "errors"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/utils"
)

func createPkgNpmConfig(pkgConfig *types.PkgConfig) *types.NpmConfig {
    meta := pkgConfig.MainTag.Meta
    return &types.NpmConfig{
        Name:         pkgConfig.GetMainTag().GetName(),
        Version:      pkgConfig.GetMainTag().GetVersion(),
        Description:  meta.Description,
        Repository:   meta.Repository,
        Main:         ".wio.js",
        Keywords:     meta.Keywords,
        Author:       meta.Author,
        License:      meta.License,
        Contributors: meta.Contributors,
    }
}

func createAppNpmConfig(appConfig *types.AppConfig) *types.NpmConfig {
    return &types.NpmConfig{
        Name:        appConfig.GetMainTag().GetName(),
        Version:     appConfig.GetMainTag().GetVersion(),
        Description: "A wio application",
        Main:        ".wio.js",
        Keywords:    []string{"wio", "app"},
        Author:      "wio",
        License:     "MIT",
    }
}

func createNpmConfig(config types.IConfig) *types.NpmConfig {
    if config.GetType() == constants.APP {
        return createAppNpmConfig(config.(*types.AppConfig))
    } else {
        return createPkgNpmConfig(config.(*types.PkgConfig))
    }
}

func updateNpmConfig(directory string) error {
    config, err := utils.ReadWioConfig(directory)
    if err != nil {
        return err
    }
    npmConfig := createNpmConfig(config)
    if err := validateNpmConfig(npmConfig); err != nil {
        return err
    }
    npmConfig.Dependencies = make(types.NpmDependencyTag)
    for name, value := range config.GetDependencies() {
        if !value.Vendor {
            if err := dependencyCheck(directory, name, value.Version); err != nil {
                return err
            }
            npmConfig.Dependencies[name] = value.Version
        }
    }

    return io.NormalIO.WriteJson(directory+io.Sep+"package.json", npmConfig)
}

func validateNpmConfig(npmConfig *types.NpmConfig) error {
    versionPat := regexp.MustCompile(`[0-9]+.[0-9]+.[0-9]+`)
    stringPat := regexp.MustCompile(`[\w"]+`)
    if !stringPat.MatchString(npmConfig.Author) {
        log.WriteFailure()
        return goerr.New("author must be specified for a package")
    }
    if !stringPat.MatchString(npmConfig.Description) {
        log.WriteFailure()
        return goerr.New("description must be specified for a package")
    }
    if !versionPat.MatchString(npmConfig.Version) {
        log.WriteFailure()
        return goerr.New("package does not have a valid version")
    }
    if !stringPat.MatchString(npmConfig.License) {
        npmConfig.License = "MIT"
    }
    return nil
}
