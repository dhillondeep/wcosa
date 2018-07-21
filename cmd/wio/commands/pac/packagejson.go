package pac

import (
    goerr "errors"
    "os"
    "regexp"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
)

func createNpmConfig(config types.Config) *types.NpmConfig {
    info := config.GetInfo()
    return &types.NpmConfig{
        Name:         info.GetName(),
        Version:      info.GetVersion(),
        Description:  "Wio package",
        Repository:   "",
        Main:         ".wio.js",
        Keywords:     utils.AppendIfMissing(info.GetKeywords(), []string{"wio", "pkg"}),
        Author:       "",
        License:      info.GetLicense(),
        Contributors: []string{},
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
        if !value.IsVendor() {
            if err := dependencyCheck(directory, name, value.GetVersion()); err != nil {
                return err
            }
            npmConfig.Dependencies[name] = value.GetVersion()
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
