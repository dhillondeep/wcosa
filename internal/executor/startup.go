package executor

import (
    "wio/internal/config/root"
)

func ExecuteStartup() error {
    // make sure wio root folder exists
    if err := root.CreateWioRoot(); err != nil {
        return err
    }

    // load up environment for wio
    if err := root.LoadEnv(); err != nil {
        return err
    }

    return nil
}
