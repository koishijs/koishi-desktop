package config

import log "github.com/sirupsen/logrus"

var (
	// Log
	l = log.WithField("package", "config")

	Version = "INTERNAL"

	Config *KoiConfig
)

type KoiConfig struct {
	Mode   string `yaml:"mode"`
	Target string `yaml:"target"`

	// Internal
	ConfigDir string `yaml:"configDir,omitempty"`
}

func LoadConfig(configPath string) {
	l.Debug("Now loading koi config.")
	config, err := ReadConfig(configPath)
	if err != nil {
		l.Fatal(err)
	}
	Config = config
}
