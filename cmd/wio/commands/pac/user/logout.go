package user

import (
    "errors"
    "os"
    "wio/cmd/wio/commands"
    "wio/cmd/wio/log"
    "wio/cmd/wio/utils/io"
)

func (cmd Logout) Execute() error {
    dir, err := commands.GetDirectory(cmd)
    if err != nil {
        return nil
    }
    path := io.Path(dir, io.Folder, "token.json")
    if !io.Exists(path) {
        return errors.New("not logged in")
    }
    log.Info(log.Cyan, "Logging out ... ")
    if err := os.RemoveAll(path); err != nil {
        log.WriteFailure()
        return err
    }
    log.WriteSuccess()
    return nil
}
