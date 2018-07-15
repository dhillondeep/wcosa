package create

import (
    "path/filepath"
    "strings"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/errors"
    "wio/cmd/wio/log"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"

    "github.com/fatih/color"
)

func (create Create) updateApp(directory string, config *types.AppConfig) error {
    info := &createInfo{
        context:     create.Context,
        directory:   directory,
        projectType: constants.APP,
        name:        config.MainTag.GetName(),
    }
    log.Info(log.Cyan, "updating app files ... ")
    return info.update(config)
}

func (create Create) updatePackage(directory string, config *types.PkgConfig) error {
    info := &createInfo{
        context:     create.Context,
        directory:   directory,
        projectType: constants.PKG,
        name:        config.MainTag.GetName(),
    }

    log.Info(log.Cyan, "updating package files ... ")
    return info.update(config)
}

func (info *createInfo) update(config types.IConfig) error {
    queue := log.GetQueue()
    if err := updateProjectFiles(queue, info); err != nil {
        log.WriteFailure()
        return err
    }
    log.WriteSuccess()
    log.PrintQueue(queue, log.TWO_SPACES)

    log.Info(log.Cyan, "updating wio.yml ... ")
    queue = log.GetQueue()
    if err := updateConfig(queue, config, info); err != nil {
        log.WriteFailure()
        log.Errln(err)
        return err
    }
    log.WriteSuccess()
    log.PrintQueue(queue, log.TWO_SPACES)

    log.Writeln()
    log.Infoln(log.Yellow.Add(color.Underline), "Project update summary")
    log.Info(log.Cyan, "path             ")
    log.Writeln(info.directory)
    log.Info(log.Cyan, "project name     ")
    log.Writeln(info.name)
    log.Info(log.Cyan, "wio version      ")
    log.Writeln(config.GetMainTag().GetConfigurations().WioVersion)
    log.Info(log.Cyan, "project type     ")
    log.Writeln(info.projectType)

    return nil
}

// Update Wio project
func (create Create) handleUpdate(directory string) error {
    cfg, err := utils.ReadWioConfig(directory)
    if err != nil {
        return err
    }
    switch cfg.GetType() {
    case constants.APP:
        err = create.updateApp(directory, cfg.(*types.AppConfig))
        break
    case constants.PKG:
        err = create.updatePackage(directory, cfg.(*types.PkgConfig))
        break
    }
    return err
}

// Update configurations
func updateConfig(queue *log.Queue, config types.IConfig, info *createInfo) error {
    switch info.projectType {
    case constants.APP:
        return updateAppConfig(queue, config.(*types.AppConfig), info)
    case constants.PKG:
        return updatePackageConfig(queue, config.(*types.PkgConfig), info)
    }
    return nil
}

func updatePackageConfig(queue *log.Queue, config *types.PkgConfig, info *createInfo) error {
    // Ensure a minimum wio version is specified
    if strings.Trim(config.MainTag.Config.WioVersion, " ") == "" {
        return errors.String("wio.yml missing `minimum_wio_version`")
    }
    if config.MainTag.GetName() != filepath.Base(info.directory) {
        log.Warnln(queue, "Base directory different from project name")
    }
    if config.MainTag.CompileOptions.HeaderOnly {
        config.MainTag.Flags.Visibility = "INTERFACE"
        config.MainTag.Definitions.Visibility = "INTERFACE"
    } else {
        config.MainTag.Flags.Visibility = "PRIVATE"
        config.MainTag.Definitions.Visibility = "PRIVATE"
    }
    if strings.Trim(config.MainTag.Meta.Version, " ") == "" {
        config.MainTag.Meta.Version = "0.0.1"
    }
    configPath := info.directory + io.Sep + io.Config
    return config.PrettyPrint(configPath)
}

func updateAppConfig(queue *log.Queue, config *types.AppConfig, info *createInfo) error {
    // Ensure a minimum wio version is specified
    if strings.Trim(config.MainTag.Config.WioVersion, " ") == "" {
        return errors.String("wio.yml missing `minimum_wio_version`")
    }
    if config.MainTag.GetName() != filepath.Base(info.directory) {
        log.Warnln(queue, "Base directory different from project name")
    }
    configPath := info.directory + io.Sep + io.Config
    return config.PrettyPrint(configPath)
}

// Update project files
func updateProjectFiles(queue *log.Queue, info *createInfo) error {
    log.Verb(queue, "reading paths.json file ... ")
    structureData := &StructureConfigData{}
    if err := io.AssetIO.ParseJson("configurations/structure-avr.json", structureData); err != nil {
        log.WriteFailure(queue, log.VERB)
        return err
    } else {
        log.WriteSuccess(queue, log.VERB)
    }
    dataType := &structureData.Pkg
    if info.projectType == constants.APP {
        dataType = &structureData.App
    }
    copyProjectAssets(queue, info, dataType)
    return nil
}
