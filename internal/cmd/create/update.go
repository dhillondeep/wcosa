package create

import (
	"path/filepath"
	"strings"
	"wio/internal/constants"
	"wio/internal/types"
	"wio/pkg/log"
	"wio/pkg/util"
	"wio/pkg/util/sys"

	"github.com/fatih/color"
)

func (create Create) updateApp(info *createInfo, directory string, config types.Config) error {
	info.projectType = constants.App

	log.Info(log.Cyan, "updating app files ... ")
	return info.update(config)
}

func (create Create) updatePackage(info *createInfo, directory string, config types.Config) error {
	info.projectType = constants.Pkg

	log.Info(log.Cyan, "updating package files ... ")
	return info.update(config)
}

func (info *createInfo) update(config types.Config) error {
	queue := log.GetQueue()

	if !info.configOnly {
		if err := updateProjectFiles(queue, info); err != nil {
			log.WriteFailure()
			log.PrintQueue(queue, log.TWO_SPACES)
			return err
		}
		log.WriteSuccess()
		log.PrintQueue(queue, log.TWO_SPACES)
	}

	log.Info(log.Cyan, "updating wio.yml ... ")
	queue = log.GetQueue()
	if err := updateConfig(queue, config, info); err != nil {
		log.WriteFailure()
		log.PrintQueue(queue, log.TWO_SPACES)
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
	log.Writeln(config.GetInfo().GetOptions().GetWioVersion())
	log.Info(log.Cyan, "project type     ")
	log.Writeln(info.projectType)

	return nil
}

// Update Wio project
func (create Create) handleUpdate(directory string) (*createInfo, error) {
	cfg, err := types.ReadWioConfig(directory, false)
	if err != nil {
		return nil, err
	}

	info := &createInfo{
		context:    create.Context,
		directory:  directory,
		name:       cfg.GetName(),
		ide:        create.Context.String("ide"),
		fullUpdate: create.Context.Bool("full"),
		updateOnly: create.Update,
		configOnly: create.Context.Bool("only-config"),
	}

	newUpdateVer := create.Context.String("version")
	if !util.IsEmptyString(newUpdateVer) {
		cfg.SetVersion(newUpdateVer)
	}

	switch cfg.GetType() {
	case constants.App:
		err = create.updateApp(info, directory, cfg)
		break
	case constants.Pkg:
		err = create.updatePackage(info, directory, cfg)
		break
	}
	return info, err
}

// Update configurations
func updateConfig(queue *log.Queue, config types.Config, info *createInfo) error {
	switch info.projectType {
	case constants.App:
		return updateAppConfig(queue, config, info)
	case constants.Pkg:
		return updatePackageConfig(queue, config, info)
	}
	return nil
}

func updatePackageConfig(queue *log.Queue, config types.Config, info *createInfo) error {
	// Ensure a minimum wio version is specified
	if strings.Trim(config.GetInfo().GetOptions().GetWioVersion(), " ") == "" {
		return util.Error("wio.yml missing `minimum_wio_version`")
	}
	if config.GetName() != filepath.Base(info.directory) {
		log.Warnln(queue, "Base directory different from project name")
	}
	return types.WriteWioConfig(info.directory, config)
}

func updateAppConfig(queue *log.Queue, config types.Config, info *createInfo) error {
	// Ensure a minimum wio version is specified
	if strings.Trim(config.GetInfo().GetOptions().GetWioVersion(), " ") == "" {
		return util.Error("wio.yml missing `minimum_wio_version`")
	}
	if config.GetName() != filepath.Base(info.directory) {
		log.Warnln(queue, "Base directory different from project name")
	}
	return types.WriteWioConfig(info.directory, config)
}

// Update project files
func updateProjectFiles(queue *log.Queue, info *createInfo) error {
	log.Verb(queue, "reading paths.json file ... ")
	structureData := &StructureConfigData{}
	if err := sys.AssetIO.ParseJson("configurations/structure.json", structureData); err != nil {
		log.WriteFailure(queue, log.VERB)
		return err
	} else {
		log.WriteSuccess(queue, log.VERB)
	}
	dataType := &structureData.Pkg
	if info.projectType == constants.App {
		dataType = &structureData.App
	}
	dataType.Paths = append(dataType.Paths, structureData.Shared.Paths...)

	copyProjectAssets(queue, info, dataType, info.fullUpdate)
	return nil
}
