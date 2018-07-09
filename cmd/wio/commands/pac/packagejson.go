package pac

import (
    goerr "errors"
    "os"
    "regexp"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
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
        Repository:  "repo",
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

func updateNpmConfig(directory string, strict bool) error {
    config, err := utils.ReadWioConfig(directory)
    if err != nil {
        return err
    }
    npmConfig := createNpmConfig(config)
    if err := validateNpmConfig(npmConfig); strict && err != nil {
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
    dotWioPath := io.Path(directory, io.Folder)
    if err := os.MkdirAll(dotWioPath, os.ModePerm); err != nil {
        return err
    }
    return io.NormalIO.WriteJson(io.Path(dotWioPath, "package.json"), npmConfig)
}

func validateNpmConfig(npmConfig *types.NpmConfig) error {
    versionPat := regexp.MustCompile(`[0-9]+.[0-9]+.[0-9]+`)
    stringPat := regexp.MustCompile(`[\w"]+`)
    if !stringPat.MatchString(npmConfig.Author) {
        return goerr.New("author must be specified for a package")
    }
    if !stringPat.MatchString(npmConfig.Description) {
        return goerr.New("description must be specified for a package")
    }
    if !versionPat.MatchString(npmConfig.Version) {
        return goerr.New("package does not have a valid version")
    }
    if !stringPat.MatchString(npmConfig.License) {
        npmConfig.License = "MIT"
    }
    return nil
}
