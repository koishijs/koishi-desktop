package koicli

import (
	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"gopkg.ilharper.com/koi/app/config"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/x/rpl"
)

func newPreAction(i *do.Injector) (cli.BeforeFunc, error) {
	l := do.MustInvoke[*logger.Logger](i)
	consoleTarget := do.MustInvoke[*logger.ConsoleTarget](i)

	return func(c *cli.Context) (err error) {
		l.Debug("Trigger pseudo action: pre")

		l.Debug("Checking flag debug...")
		if c.Bool("debug") {
			consoleTarget.Level = rpl.LevelDebug
		}

		l.Debug("Checking config file...")
		configPath := c.String("config")
		if configPath != "" {
			l.Debugf("Using flag provided config path: %s", configPath)
		} else {
			configPath = "koi.yml"
		}
		l.Infof("Using config file: %s", configPath)
		do.Provide(i, config.BuildLoadConfig(configPath))

		return
	}, nil
}
