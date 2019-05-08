package sys

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
)

// ParseJson parses json from the file on filesystem
func ParseJson(fileName string, out interface{}) (err error) {
	text, err := ReadFile(fileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(text), out)
	return err
}

// ParseYaml parses yaml from the file on filesystem
func ParseYaml(fileName string, out interface{}) error {
	text, err := ReadFile(fileName)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(text, out)
}

// WriteJson writes json data to a file on filesystem
func WriteJson(fileName string, in interface{}) error {
	data, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return err
	}

	return WriteFile(fileName, data)
}

// WriteYaml writes yaml data to a file on filesystem
func WriteYaml(fileName string, in interface{}) error {
	data, err := yamlMarshal(in)
	if err != nil {
		return err
	}

	return WriteFile(fileName, data)
}
