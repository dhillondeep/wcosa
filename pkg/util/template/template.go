package template

import (
	"strings"
	"wio/pkg/util/sys"
)

func IOReplace(path string, values map[string]string) error {
	data, err := sys.NormalIO.ReadFile(path)
	if err != nil {
		return err
	}
	result := Replace(string(data), values)
	err = sys.NormalIO.WriteFile(path, []byte(result))
	if err != nil {
		return err
	}
	return nil
}

func Replace(template string, values map[string]string) string {
	for match, replace := range values {
		template = strings.Replace(template, "{{"+match+"}}", replace, -1)
	}
	return template
}
