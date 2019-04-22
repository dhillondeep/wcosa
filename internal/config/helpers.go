package config

import (
	"github.com/hashicorp/hil"
	"reflect"
	"strings"
)

// applyHilGeneric applies Hil language parser on string and returns an interface
func applyHilGeneric(val string, config *hil.EvalConfig, def interface{}) (*hil.EvaluationResult, error) {
	tree, err := hil.Parse(val)
	if err != nil {
		return nil, err
	}

	result, err := hil.Eval(tree, config)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// applyHilString applies Hil language parser on string and returns a string
func applyHilString(val string, config *hil.EvalConfig) (string, error) {
	result, err := applyHilGeneric(val, config, "")
	if err != nil {
		return "", err
	}
	return result.Value.(string), err
}

// stringToStringSlice convert string reflect value to a slice of string based on the seperator
func stringToStringSlice(val reflect.Value, sep string) []string {
	newSlice := strings.Split(val.String(), sep)
	if len(newSlice) < 2 {
		newSlice = append(newSlice, "")
	}

	return newSlice
}
