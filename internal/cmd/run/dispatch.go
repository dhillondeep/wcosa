package run

import (
	"wio/internal/cmd/generate"
	"wio/internal/constants"
	"wio/internal/types"
	"wio/pkg/util"
	"wio/pkg/util/sys"
)

func dispatchRunTarget(info *runInfo, target types.Target) error {
	binDir := binaryPath(info, target)
	platform := target.GetPlatform()

	genHardwareFile := func(port string) error {
		// generate hardware files
		infoGen := &generate.InfoGenerate{
			Config:      info.config,
			Directory:   info.directory,
			ProjectType: info.projectType,
			Port:        port,
		}

		if err := generate.HardwareFile(infoGen, target); err != nil {
			return err
		}

		return nil
	}

	switch platform {
	case constants.Avr:
		// auto detect the port or use the provided one
		port, err := generate.GetPort(info.port)
		if err != nil {
			return err
		}
		if err := genHardwareFile(port); err != nil {
			return err
		}

		if err := configTarget(binDir); err != nil {
			return err
		}

		return uploadTarget(binDir)
	case constants.Native:
		if err := genHardwareFile(info.port); err != nil {
			return err
		}

		args := info.context.String("args")
		return runTarget(info.directory, sys.Path(binDir, target.GetName()), args)
	default:
		return util.Error("platform [%s] is not supported", platform)
	}
}

func dispatchCanRunTarget(info *runInfo, target types.Target) bool {
	binDir := binaryPath(info, target)
	platform := target.GetPlatform()
	file := sys.Path(binDir, target.GetName()+platformExtension(platform))
	return sys.Exists(file)
}
