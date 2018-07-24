package publish

import (
	"wio/internal/cmd"
	"wio/internal/types"
	"wio/pkg/npm/publish"

	"github.com/urfave/cli"
)

type Cmd struct {
	Context *cli.Context
}

func (c Cmd) GetContext() *cli.Context {
	return c.Context
}

func (c Cmd) Execute() error {
	dir, err := cmd.GetDirectory(c)
	if err != nil {
		return err
	}
	cfg, err := types.ReadWioConfig(dir)
	if err != nil {
		return err
	}
	return publish.Do(dir, cfg)
}
