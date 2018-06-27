package run

import (
    "os"
)

func readDirectory(args []string) (string, error) {
    if len(args) <= 0 {
        return os.Getwd()
    }
    return args[0], nil
}
