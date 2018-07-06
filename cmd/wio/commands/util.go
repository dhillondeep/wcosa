package commands

import "os"

func getDirectory(cmd Command) (string, error) {
    ctx := cmd.GetContext()
    if ctx.IsSet("dir") {
        return ctx.String("dir"), nil
    }
    return os.Getwd()
}
