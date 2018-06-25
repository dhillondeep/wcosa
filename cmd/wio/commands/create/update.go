package create

import (
    "wio/cmd/wio/types"
    "wio/cmd/wio/constants"
    "wio/cmd/wio/log"
    "github.com/fatih/color"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
    "strings"
    "wio/cmd/wio/errors"
    "path/filepath"
)

func (create Create) updateApp(directory string, config *types.AppConfig) error {
    return nil
}

func (create Create) updatePackage(directory string, config *types.PkgConfig) error {
    info := &createInfo{
        Directory: directory,
        Type:      constants.PKG,
        Name:      config.MainTag.GetName(),
    }

    log.Info(log.Cyan, "updating package files ... ")
    queue := log.GetQueue()
    if err := create.updateProjectFiles(queue, info); err != nil {
        log.WriteFailure()
        return err
    }
    log.WriteSuccess()
    log.PrintQueue(queue, log.TWO_SPACES)

    log.Info(log.Cyan, "updating wio.yml ... ")
    queue = log.GetQueue()
    if err := updatePackageConfig(queue, config, info); err != nil {
        log.WriteFailure()
        log.Errln(err)
        return err
    }
    log.WriteSuccess()
    log.PrintQueue(queue, log.TWO_SPACES)

    log.Writeln()
    log.Infoln(log.Yellow.Add(color.Underline), "Project update summary")
    log.Info(log.Cyan, "path             ")
    log.Writeln(directory)
    log.Info(log.Cyan, "project name     ")
    log.Writeln(info.Name)
    log.Info(log.Cyan, "wio version      ")
    log.Writeln(config.GetMainTag().GetConfigurations().WioVersion)
    log.Info(log.Cyan, "project type     ")
    log.Writeln("pkg")

    return nil
}

// Update Wio project
func (create Create) handleUpdate(directory string) error {
    cfg, err := utils.ReadWioConfig(directory + io.Sep + "wio.yml")
    if err != nil {
        return err
    }
    switch cfg.Type {
    case types.App:
        err = create.updateApp(directory, cfg.Config.(*types.AppConfig))
        break
    default:
        err = create.updatePackage(directory, cfg.Config.(*types.PkgConfig))
        break
    }
    return err
}

func updatePackageConfig(queue *log.Queue, config *types.PkgConfig, info *createInfo) error {
    // Ensure a minimim wio version is specified
    if strings.Trim(config.MainTag.Config.WioVersion, " ") == "" {
        return errors.String("wio.yml missing `minimum_wio_version`")
    }
    if config.MainTag.GetName() != filepath.Base(info.Directory) {
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
    configPath := info.Directory + io.Sep + "wio.yml"
    return config.PrettyPrint(configPath)
}

/*// Update wio.yml file
func (create Create) updateConfig(queue *log.Queue, projectConfig types.Config, directory string) error {
    isApp, err := utils.IsAppType(directory + io.Sep + "wio.yml")
    if err != nil {
        return errors.ConfigParsingError{
            Err: err,
        }
    }

    if isApp {
        appConfig := projectConfig.(*types.AppConfig)

        appConfig.MainTag.Name = filepath.Base(directory)

        updateAVRAppTargets(&appConfig.TargetsTag, directory)

        // make current wio version as the default version if no version is provided
        if strings.Trim(appConfig.MainTag.Config.WioVersion, " ") == "" {
            appConfig.MainTag.Config.WioVersion = config.ProjectMeta.Version
        }
    } else {
    pkgConfig := projectConfig.(*types.PkgConfig)

    pkgConfig.MainTag.Meta.Name = filepath.Base(directory)

    // make sure framework is AVR
    pkgConfig.MainTag.CompileOptions.Platform = constants.AVR

    // make current wio version as the default version if no version is provided
    if strings.Trim(pkgConfig.MainTag.Config.WioVersion, " ") == "" {
        pkgConfig.MainTag.Config.WioVersion = config.ProjectMeta.Version
    }

    // keywords
    pkgConfig.MainTag.Meta.Keywords = utils.AppendIfMissing(pkgConfig.MainTag.Meta.Keywords,
        []string{"wio", "c", "c++"})
    pkgConfig.MainTag.Meta.Keywords = utils.AppendIfMissing(pkgConfig.MainTag.Meta.Keywords,
        pkgConfig.MainTag.Config.SupportedFrameworks)
    pkgConfig.MainTag.Meta.Keywords = utils.AppendIfMissing(pkgConfig.MainTag.Meta.Keywords,
        pkgConfig.MainTag.Config.SupportedPlatforms)
    pkgConfig.MainTag.Meta.Keywords = utils.AppendIfMissing(pkgConfig.MainTag.Meta.Keywords,
        pkgConfig.MainTag.Config.SupportedBoards)

    // flags
    if strings.Trim(pkgConfig.MainTag.Flags.Visibility, " ") == "" {
        if pkgConfig.MainTag.CompileOptions.HeaderOnly {
            pkgConfig.MainTag.Flags.Visibility = "INTERFACE"
        } else {
            pkgConfig.MainTag.Flags.Visibility = "PRIVATE"
        }
    }

    // definitions
    if strings.Trim(pkgConfig.MainTag.Definitions.Visibility, " ") == "" {
        if pkgConfig.MainTag.CompileOptions.HeaderOnly {
            pkgConfig.MainTag.Definitions.Visibility = "INTERFACE"
        } else {
            pkgConfig.MainTag.Definitions.Visibility = "PRIVATE"
        }
    }

    // version
    if strings.Trim(pkgConfig.MainTag.Meta.Version, " ") == "" {
        pkgConfig.MainTag.Meta.Version = "0.0.1"
    }

    if err := utils.PrettyPrintConfig(projectConfig, directory+io.Sep+"wio.yml", false); err != nil {
        return errors.WriteFileError{
            FileName: directory + io.Sep + "wio.yml",
            Err:      err,
        }
    }

    return nil
}*/

// Update project `files
func (create Create) updateProjectFiles(queue *log.Queue, info *createInfo) error {
    log.Verb(queue, "reading paths.json file ... ")
    structureData := &StructureConfigData{}
    if err := io.AssetIO.ParseJson("configurations/structure-avr.json", structureData); err != nil {
        log.WriteFailure(queue, log.VERB)
        return err
    } else {
        log.WriteSuccess(queue, log.VERB)
    }
    create.copyProjectAssets(queue, info, structureData.Pkg)
    return nil
}
