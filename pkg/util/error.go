package util

import (
	"errors"
	"fmt"
)

func Error(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return errors.New(msg)
}
