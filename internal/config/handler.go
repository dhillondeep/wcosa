package config

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"wio/internal/constants"
	"wio/pkg/sys"
	"wio/templates"
)

// ReadConfig parses config file from project path and provides ProjectConfig
func ReadConfig(projectPath string) (ProjectConfig, error) {
	projectConfig := projectConfigImpl{}

	var warnings []string
	updateWarnings := func(warning string) {
		warnings = append(warnings, warning)
	}

	viper.SetConfigFile(sys.JoinPaths(projectPath, constants.WioConfigFile))
	viper.SetFs(sys.GetFileSystem())
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// get project type
	projectConfig.Type = viper.GetString("type")

	if err := viper.Unmarshal(&projectConfig, func(config *mapstructure.DecoderConfig) {
		config.ErrorUnused = true
		config.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			config.DecodeHook,
			mapstructure.StringToSliceHookFunc(" "),
			splitKeyValToMapFunc(),
			warningHookFunc(projectConfig.Type, updateWarnings),
		)
	}); err != nil {
		return nil, err
	}

	if len(warnings) > 0 {
		fmt.Println(strconv.Itoa(len(warnings)) + " warning(s) decoding:\n")
		fmt.Println(strings.Join(warnings, "\n"))
	}

	var g interface{} = &projectConfig
	return g.(ProjectConfig), nil
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
