package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"wio/internal/constants"
)

// warningHookFunc handles all the unsupported tags and for each unsupported tag, it generates a warning
// and set's the tag value to null so it cannot be used
func warningHookFunc(projectType string, provideWarning func(warning string)) mapstructure.DecodeHookFunc {
	warningTpl := "* %s config does not support %s tag"

	unsupportedTagWarning := func(dataVal reflect.Value, desiredType string,
		tagName string, warnTagName string) reflect.Value {
		if dataVal.MapIndex(reflect.ValueOf(tagName)).Kind() != reflect.Invalid {
			if projectType == desiredType {
				provideWarning(fmt.Sprintf(warnTagName, tagName))
				dataVal.SetMapIndex(reflect.ValueOf(tagName), reflect.ValueOf(nil))
			}
		}

		return dataVal
	}

	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		dataVal := reflect.ValueOf(data)

		if t.ConvertibleTo(reflect.TypeOf(TargetImpl{})) {
			tagName := fmt.Sprintf(warningTpl, projectType, "targets[*].%s")
			dataVal = unsupportedTagWarning(dataVal, constants.PKG, "executable_options", tagName)
			dataVal = unsupportedTagWarning(dataVal, constants.APP, "package_options", tagName)

			return dataVal.Interface(), nil
		} else if t.ConvertibleTo(reflect.TypeOf(ProjectImpl{})) {
			tagName := fmt.Sprintf(warningTpl, projectType, "project.%s")
			dataVal = unsupportedTagWarning(dataVal, constants.APP, "package_options", tagName)

			return dataVal.Interface(), nil
		} else if t.ConvertibleTo(reflect.TypeOf(TestImpl{})) {
			tagName := fmt.Sprintf(warningTpl, projectType, "tests[*].executable_options.%s")
			tagValue := dataVal.MapIndex(reflect.ValueOf("executable_options"))

			if tagValue.Kind() == reflect.Interface {
				tagValInterface := reflect.ValueOf(tagValue.Interface())

				tagValInterface = unsupportedTagWarning(tagValInterface, projectType, "main_file", tagName)
				dataVal.SetMapIndex(reflect.ValueOf("executable_options"), tagValInterface)
				return dataVal.Interface(), nil
			}

		}

		return data, nil
	}
}

func dependencyShortFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		val := reflect.ValueOf(data)

		if t.ConvertibleTo(reflect.TypeOf(DependencyImpl{})) {
			if val.Kind() == reflect.String {
				return reflect.ValueOf(DependencyImpl{
					Ref: val.String(),
				}).Interface(), nil
			}
		}

		return data, nil
	}
}

// splitKeyValToMapFunc is a decode hook function to convert a string to key value pair
func splitKeyValToMapFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		val := reflect.ValueOf(data)

		if val.Kind() == reflect.String {
			if t.ConvertibleTo(reflect.TypeOf(ArgumentImpl{})) || t.ConvertibleTo(reflect.TypeOf(VariableImpl{})) {
				splitVal := stringToStringSlice(val, "=")
				return map[string]string{
					"name":  splitVal[0],
					"value": splitVal[1],
				}, nil
			} else if t.ConvertibleTo(reflect.TypeOf(ToolchainImpl{})) {
				splitVal := stringToStringSlice(val, "::")
				return map[string]string{
					"name": splitVal[0],
					"ref":  splitVal[1],
				}, nil
			}
		}

		return data, nil
	}
}
