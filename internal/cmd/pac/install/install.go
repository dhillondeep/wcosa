package install

import (
	"github.com/urfave/cli"
	"strings"
	"wio/internal/cmd"
	"wio/internal/types"
	"wio/pkg/log"
	"wio/pkg/npm/resolve"
	"wio/pkg/util"
)

type Cmd struct {
	Context *cli.Context

	dir    string
	info   *resolve.Info
	config types.Config
}

func (c Cmd) GetContext() *cli.Context {
	return c.Context
}

func (c Cmd) Execute() error {
	var err error
	c.dir, err = cmd.GetDirectory(c)
	if err != nil {
		return err
	}
	c.config, err = types.ReadWioConfig(c.dir, false)
	if err != nil {
		return err
	}
	c.info = resolve.NewInfo(c.dir)

	if len(c.Context.Args()) > 0 {
		if err := c.AddDependency(); err != nil {
			return err
		}
	}

	if err := c.info.ResolveRemote(c.config, true); err != nil {
		return err
	}
	return c.info.InstallResolved()
}

func (c Cmd) AddDependency() error {
	urlFlagVal := c.Context.String("url")
	urlDirVal := c.Context.String("dir")
	urlOptionsVal := c.Context.String("options")

	urlProvided := !util.IsEmptyString(urlFlagVal)

	name, ver, err := c.getArgs(c.info, !urlProvided)
	if err != nil {
		return err
	}

	log.Info(log.Cyan, "Adding dependency: ")
	log.Infoln(log.Green, "%s@%s", name, ver)

	deps := c.config.GetDependencies()

	if prev, exists := deps[name]; exists && prev.GetVersion() != ver {
		log.Warnln("Replacing previous version %s", prev.GetVersion())
	} else if exists {
		log.Warnln("Same version already exists")
	}

	newDependency := &types.DependencyImpl{
		Version: ver,
		Vendor:  false,
	}

	if urlProvided {
		newDependency.Url = &types.DependencyUrlImpl{
			Name: urlFlagVal,
			Dir:  urlDirVal,
		}

		options := strings.Split(urlOptionsVal, ",")
		if len(options) > 0 {
			newDependency.Url.Options = map[string]string{}
		}
		for _, option := range options {
			optionList := strings.Split(option, "=")

			if len(optionList) > 1 {
				newDependency.Url.Options[optionList[0]] = optionList[1]
			}
		}
	}

	c.config.AddDependency(name, newDependency)
	return types.WriteWioConfig(c.dir, c.config)
}
