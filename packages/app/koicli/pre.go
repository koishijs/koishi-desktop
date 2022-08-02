package koicli

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"gopkg.ilharper.com/koi/app/config"
	"gopkg.ilharper.com/x/rpl"
)

func buildPreAction(kcli *KoiCli) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		var err error

		kcli.l.Debug("Trigger pseudo action: pre")

		kcli.l.Debug("Checking flag debug...")
		if c.Bool("debug") {
			kcli.consoleTarget.Level = rpl.LevelDebug
		}

		kcli.l.Debug("Checking config file...")
		configPath := c.String("config")
		if configPath != "" {
			kcli.l.Debugf("Using flag provided config path: %s", configPath)
		} else {
			configPath = "koi.yml"
		}
		kcli.l.Infof("Using config file: %s", configPath)
		kcli.config, err = config.LoadConfig(kcli.l, configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		return nil
	}
}
