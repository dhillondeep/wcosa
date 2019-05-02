package utils

import "strings"

func IsStringEmpty(value string) bool {
	return strings.TrimSpace(value) == ""
}
