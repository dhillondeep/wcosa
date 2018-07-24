package user

import (
	"os"
	"wio/internal/cmd"
	"wio/pkg/log"
	"wio/pkg/util"
	"wio/pkg/util/sys"
)

func (c Logout) Execute() error {
	dir, err := cmd.GetDirectory(c)
	if err != nil {
		return nil
	}
	path := sys.Path(dir, sys.Folder, "token.json")
	if !sys.Exists(path) {
		return util.Error("not logged in")
	}
	log.Info(log.Cyan, "Logging out ... ")
	if err := os.RemoveAll(path); err != nil {
		log.WriteFailure()
		return err
	}
	log.WriteSuccess()
	return nil
}
