package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsStringEmpty(t *testing.T) {
	// empty string
	require.Equal(t, true, IsStringEmpty(""))

	// space string
	require.Equal(t, true, IsStringEmpty("    "))

	// normal string without spaces
	require.Equal(t, false, IsStringEmpty("kashdjk"))

	// normal string with spaces
	require.Equal(t, false, IsStringEmpty(" kashdjk   "))
}
