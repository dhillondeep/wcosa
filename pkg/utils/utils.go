package utils

import "strings"

// IsStringEmpty returns true or false based on if the string is empty (spaces are ignored)
func IsStringEmpty(value string) bool {
	return strings.TrimSpace(value) == ""
}
