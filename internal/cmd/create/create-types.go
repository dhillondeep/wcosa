package create

import (
	"wio/internal/cmd/generate"
	"wio/internal/types"
	"wio/pkg/log"

	"github.com/urfave/cli"
)

type Create struct {
	Context *cli.Context
	Update  bool
	Type    string
}

type createInfo struct {
	context *cli.Context

	directory   string
	projectType string
	name        string

	platform   string
	framework  string
	board      string
	ide        string
	fullUpdate bool

	configOnly bool
	headerOnly bool
	updateOnly bool
}

// get context for the command
func (create Create) GetContext() *cli.Context {
	return create.Context
}

// Executes the create command
func (create Create) Execute() error {
	var err error
	directory, err := performDirectoryCheck(create.Context)
	if err != nil {
		return err
	}

	var info *createInfo
	if create.Update {
		// this checks if wio.yml file exists for it to update
		if err := performWioExistsCheck(directory); err != nil {
			return err
		}
		if info, err = create.handleUpdate(directory); err != nil {
			return err
		}
	} else {
		// this checks if directory is empty before create can be triggered
		onlyConfig := create.Context.Bool("only-config")
		if err := performPreCreateCheck(directory, onlyConfig); err != nil {
			return err
		}
		if info, err = create.createProject(directory); err != nil {
			return err
		}
	}

	if create.Update && info.configOnly {
		return nil
	}

	log.Writeln()

	// read config file
	config, err := types.ReadWioConfig(directory, true)
	if err != nil {
		return err
	}

	// handle ide
	target, err := create.generateIdeFiles(info.ide, directory, config)
	if err != nil {
		return err
	}

	infoGen := &generate.InfoGenerate{
		Config:      config,
		Directory:   directory,
		ProjectType: config.GetType(),
		Port:        "none",
		Retool:      false,
		NoDepCreate: true,
	}

	// generate cmake file for default target
	if err := generate.CMakeListsFile(infoGen, target); err != nil {
		return err
	}

	// generate cmake dependencies file for default target
	if err := generate.DependenciesFile(infoGen, target); err != nil {
		return err
	}

	// generate cmake hardware file for default target
	if err := generate.HardwareFile(infoGen, target); err != nil {
		return err
	}

	return nil
}
