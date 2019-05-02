package config

import (
	"bytes"
	"github.com/dhillondeep/viper"
	"github.com/mitchellh/mapstructure"
	"wio/internal/constants"
	"wio/pkg/sys"
	"wio/pkg/utils"
	"wio/templates"
)

// ReadConfig parses config file from project path and provides ProjectConfig
func ReadConfig(projectPath string) (ProjectConfig, []string, error) {
	projectConfig := projectConfigImpl{}

	var warnings []string
	updateWarnings := func(warning string) {
		warnings = append(warnings, warning)
	}

	viper.SetConfigFile(sys.JoinPaths(projectPath, constants.WioConfigFile))
	viper.SetKeyDelim(":")
	viper.SetFs(sys.GetFileSystem())
	if err := viper.ReadInConfig(); err != nil {
		return nil, nil, err
	}

	// get project type
	projectConfig.Type = viper.GetString("type")

	if err := viper.Unmarshal(&projectConfig, func(config *mapstructure.DecoderConfig) {
		config.ErrorUnused = true
		config.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			config.DecodeHook,
			dependencyShortFunc(),
			splitKeyValToMapFunc(),
			warningHookFunc(projectConfig.Type, updateWarnings),
		)
	}); err != nil {
		return nil, nil, err
	}

	useGlobalValues(projectConfig)

	var g interface{} = &projectConfig
	return g.(ProjectConfig), warnings, nil
}

// useGlobalValues overrides scope specific values with global if not provided
func useGlobalValues(projectConfig projectConfigImpl) {
	for _, target := range projectConfig.Targets {
		// override if not provided
		if target.CompileOptions == nil {
			target.CompileOptions = projectConfig.Project.CompileOptions
		} else if projectConfig.Project.CompileOptions != nil {
			// override if not provided otherwise append
			if target.CompileOptions.Flags == nil {
				target.CompileOptions.Flags = projectConfig.Project.CompileOptions.Flags
			} else if projectConfig.Project.CompileOptions.Flags != nil {
				target.CompileOptions.Flags = append(target.CompileOptions.Flags,
					projectConfig.Project.CompileOptions.Flags...)
			}

			// override if not provided otherwise append
			if target.CompileOptions.Definitions == nil {
				target.CompileOptions.Definitions = projectConfig.Project.CompileOptions.Definitions
			} else if projectConfig.Project.CompileOptions.Definitions != nil {
				target.CompileOptions.Definitions = append(target.CompileOptions.Definitions,
					projectConfig.Project.CompileOptions.Definitions...)
			}

			// override if not provided
			if utils.IsStringEmpty(target.CompileOptions.CXXStandard) {
				target.CompileOptions.CXXStandard = projectConfig.Project.CompileOptions.CXXStandard
			}

			// override if not provided
			if utils.IsStringEmpty(target.CompileOptions.CStandard) {
				target.CompileOptions.CStandard = projectConfig.Project.CompileOptions.CStandard
			}
		}

		// override if not provided
		if projectConfig.Type == constants.PKG && target.PackageOptions == nil {
			target.PackageOptions = projectConfig.Project.PackageOptions
		} else if projectConfig.Type == constants.PKG && projectConfig.Project.PackageOptions != nil {
			// override if not provided
			if utils.IsStringEmpty(target.PackageOptions.Type) {
				target.PackageOptions.Type = projectConfig.Project.PackageOptions.Type
			}
			// override if not provided
			if utils.IsStringEmpty(target.PackageOptions.HeaderOnly) {
				target.PackageOptions.HeaderOnly = projectConfig.Project.PackageOptions.HeaderOnly
			}
		}
	}
}

// CreateConfig creates initial config file
func CreateConfig(creation templates.ProjectCreation) error {
	var buf bytes.Buffer
	if creation.Type == constants.APP {
		templates.WriteCreateAppConfig(&buf, creation)
	} else {
		templates.WriteCreatePkgConfig(&buf, creation)
	}

	return sys.WriteFile(sys.JoinPaths(creation.ProjectPath, constants.WioConfigFile), buf.Bytes())
}
