package cli

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"koi/config"
)

func preAction(c *cli.Context) error {
	l.Debug("Trigger pseudo action: pre")

	l.Debug("Checking flag debug...")
	if c.Bool("debug") {
		log.SetLevel(log.TraceLevel)
	}

	l.Debug("Checking config file...")
	configPath := c.String("config")
	if configPath != "" {
		l.Debugf("Using flag provided config path: %s", configPath)
	} else {
		configPath = "koi.yml"
	}
	l.Infof("Using config file: %s", configPath)
	config.LoadConfig(configPath)

	return nil
}
