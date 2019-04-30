package config

import (
	"github.com/d5/tengo/script"
	"github.com/hashicorp/hil"
	"reflect"
	"regexp"
	"strings"
	"wio/internal/evaluators/tengolang"
)

var reg = regexp.MustCompile(`(?s)^\s*\$exec\s*{.*}\s*$`)
var beginReg = regexp.MustCompile(`(?s)^\s*\$exec\s*{`)
var endReg = regexp.MustCompile(`(?s)}\s*$`)

// applyHilGeneric applies Hil language parser on string and returns an interface
func applyHilGeneric(val string, config *hil.EvalConfig) (*hil.EvaluationResult, error) {
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
	result, err := applyHilGeneric(val, config)
	if err != nil {
		return "", err
	}
	return result.Value.(string), err
}

// applyEvaluator applies Hil or script execution based on regex matching
func applyEvaluator(val string, config *hil.EvalConfig) (string, error) {
	if reg.Match([]byte(val)) {
		content := endReg.ReplaceAll(beginReg.ReplaceAll([]byte(val), []byte("")), []byte(""))

		// evaluate hil
		result, err := applyHilString(string(content), config)
		if err != nil {
			return "", err
		}

		s := script.New([]byte(result))
		s.SetImports(tengolang.GetModuleMap("os", "text", "math", "times", "rand", "json", "enum", "wstrings"))

		if err := s.Add("out", ""); err != nil {
			return "", err
		}

		if c, err := s.Run(); err != nil {
			return "", err
		} else {
			return c.Get("out").String(), nil
		}
	} else {
		result, err := applyHilString(val, config)
		if err != nil {
			return "", err
		}
		return result, err
	}
}

// stringToStringSlice convert string reflect value to a slice of string based on the seperator
func stringToStringSlice(val reflect.Value, sep string) []string {
	newSlice := strings.Split(val.String(), sep)
	if len(newSlice) < 2 {
		newSlice = append(newSlice, "")
	}

	return newSlice
}
