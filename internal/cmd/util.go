package cmd

import "os"

func GetDirectory(cmd Command) (string, error) {
    ctx := cmd.GetContext()
    if ctx.IsSet("dir") {
        return ctx.String("dir"), nil
    }
    return os.Getwd()
}
