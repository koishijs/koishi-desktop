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
	InternalConfigDir   string `yaml:"internalConfigDir"`
	InternalDataDir     string `yaml:"internalDataDir"`
	InternalHomeDir     string `yaml:"internalHomeDir"`
	InternalNodeDir     string `yaml:"internalNodeDir"`
	InternalTempDir     string `yaml:"internalTempDir"`
	InternalInstanceDir string `yaml:"internalInstanceDir"`
}

func LoadConfig(configPath string) {
	l.Debug("Now loading koi config.")
	config, err := ReadConfig(configPath)
	if err != nil {
		l.Fatal(err)
	}
	Config = config
}
