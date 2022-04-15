package cli

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

func PreAction(c *cli.Context) {
	l.Debug("Trigger pseudo action: pre")

	l.Debug("Checking flag debug...")
	if c.Bool("debug") {
		log.SetLevel(log.TraceLevel)
	}

	l.Debug("Checking config file...")
	cConfigPath := c.String("config")
	if cConfigPath != "" {
		l.Debugf("Using flag provided config path: %s", cConfigPath)
		configPath = cConfigPath
	}
	l.Infof("Using config file: %s", configPath)
	configFileRaw, err := os.ReadFile(configPath)
	if err != nil {
		l.Errorf("Err when reading config file: %s", configPath)
		l.Fatal(err)
	}
	configFile = string(configFileRaw)
}
