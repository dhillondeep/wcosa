// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create and update command and sub commands provided by the tool.
// Creates, updates and initializes a wio project.
package create

import (
	"path/filepath"
	"wio/internal/config/defaults"
	"wio/internal/config/meta"
	"wio/internal/constants"
	"wio/internal/types"
	"wio/pkg/log"
	"wio/pkg/util/sys"

	"github.com/fatih/color"
)

// Creation of AVR projects
func (create Create) createProject(dir string) (*createInfo, error) {
	info := &createInfo{
		context:     create.Context,
		directory:   dir,
		projectType: create.Type,
		name:        filepath.Base(dir),
		platform:    create.Context.String("platform"),
		framework:   create.Context.String("framework"),
		board:       create.Context.String("board"),
		ide:         create.Context.String("ide"),
		configOnly:  create.Context.Bool("only-config"),
		headerOnly:  create.Context.Bool("header-only"),
		updateOnly:  create.Update,
	}
	info.toLowerCase()

	// Generate project structure
	queue := log.GetQueue()
	if !info.configOnly {
		log.Info(log.Cyan, "Creating project structure... ")
		if err := createStructure(queue, info); err != nil {
			log.WriteFailure()
			return nil, err
		} else {
			log.WriteSuccess()
		}
		log.PrintQueue(queue, log.TWO_SPACES)
	}

	// Fill configuration file
	queue = log.GetQueue()
	log.Info(log.Cyan, "Configuring project files... ")
	if err := fillConfig(info); err != nil {
		log.WriteFailure()
		return nil, err
	}

	log.PrintQueue(queue, log.TWO_SPACES)

	// print structure summary
	info.printPackageCreateSummary()
	return info, nil
}

// Copy and generate files for a package project
func createStructure(queue *log.Queue, info *createInfo) error {
	log.Verb(queue, "Reading paths.json file... ")
	structureData := &StructureConfigData{}

	// read configurationsFile
	if err := sys.AssetIO.ParseJson("configurations/structure.json", structureData); err != nil {
		log.WriteFailure(queue, log.VERB)
		return err
	} else {
		log.WriteSuccess(queue, log.VERB)
	}

	log.Verb(queue, "copying asset files ... ")
	subQueue := log.GetQueue()

	dataType := &structureData.Pkg
	if info.projectType == constants.App {
		dataType = &structureData.App
	}
	dataType.Paths = append(dataType.Paths, structureData.Shared.Paths...)

	if err := copyProjectAssets(subQueue, info, dataType, false); err != nil {
		log.WriteFailure(queue, log.VERB)
		log.CopyQueue(subQueue, queue, log.FOUR_SPACES)
		return err
	} else {
		log.WriteSuccess(queue, log.VERB)
		log.CopyQueue(subQueue, queue, log.FOUR_SPACES)
	}

	readmeFile := sys.Path(info.directory, "README.md")
	err := info.fillReadMe(queue, readmeFile)

	return err
}

// Generate wio.yml for app project
func fillConfig(info *createInfo) error {
	switch info.projectType {
	case constants.App:
		return fillAppConfig(info)
	case constants.Pkg:
		return fillPackageConfig(info)
	}
	return nil
}

func fillAppConfig(info *createInfo) error {
	cfg := &types.ConfigImpl{
		Type: constants.App,
		Info: &types.InfoImpl{
			Name:     info.name,
			Version:  defaults.Version,
			Keywords: defaults.App.Keywords,
			Options: &types.OptionsImpl{
				Default: defaults.App.Target,
				Version: meta.Version,
			},
		},
		Targets: map[string]*types.TargetImpl{
			defaults.App.Target: {
				Source:    defaults.App.Source,
				Platform:  getPlatform(info.platform),
				Framework: getFramework(info.framework),
				Board:     getBoard(info.board),
			},
		},
	}
	return types.WriteWioConfig(info.directory, cfg)
}

// Generate wio.yml for package project
func fillPackageConfig(info *createInfo) error {
	cfg := &types.ConfigImpl{
		Type: constants.Pkg,
		Info: &types.InfoImpl{
			Name:     info.name,
			Version:  defaults.Version,
			Keywords: defaults.Pkg.Keywords,
			Options: &types.OptionsImpl{
				Default: defaults.Pkg.Target,
				Version: meta.Version,
			},
		},
		Targets: map[string]*types.TargetImpl{
			defaults.Pkg.Target: {
				Source:    defaults.Pkg.Source,
				Platform:  getPlatform(info.platform),
				Framework: getFramework(info.framework),
				Board:     getBoard(info.board),
			},
		},
	}
	return types.WriteWioConfig(info.directory, cfg)
}

// Print package creation summary
func (info createInfo) printPackageCreateSummary() {
	log.Infoln()
	log.Infoln(log.Yellow.Add(color.Underline), "Project structure summary")
	if !info.headerOnly {
		log.Info(log.Cyan, "src              ")
		log.Infoln("source/non client files")
	}

	if info.projectType == constants.Pkg {
		log.Info(log.Cyan, "tests            ")
		log.Infoln("source files for test target")
		log.Info(log.Cyan, "include          ")
		log.Infoln("public headers for the package")
	}

	// print project summary
	log.Writeln()
	log.Infoln(log.Yellow.Add(color.Underline), "Project creation summary")
	log.Info(log.Cyan, "path             ")
	log.Writeln(info.directory)
	log.Info(log.Cyan, "project type     ")
	log.Writeln(info.projectType)
	log.Info(log.Cyan, "platform         ")
	log.Writeln(info.platform)
	log.Info(log.Cyan, "framework        ")
	log.Writeln(info.framework)

	if info.framework == constants.Avr {
		log.Info(log.Cyan, "board            ")
		log.Writeln(info.board)
	}
}
