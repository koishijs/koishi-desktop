package config

import log "github.com/sirupsen/logrus"

var (
	// Log
	l = log.WithField("package", "config")

	Version = "INTERNAL"

	Config *KoiConfig

	defaultConfig = KoiConfig{
		Mode: "portable",

		UseDataHome: true,
		UseDataTemp: true,
	}
)

type KoiConfig struct {
	Mode   string `yaml:"mode"`
	Target string `yaml:"target"`

	// Env override
	UseDataHome bool `yaml:"useDataHome"`
	UseDataTemp bool `yaml:"useDataTemp"`

	// Internal
	InternalConfigDir   string
	InternalDataDir     string
	InternalHomeDir     string
	InternalNodeDir     string
	InternalNodeExeDir  string
	InternalTempDir     string
	InternalInstanceDir string
}

func LoadConfig(configPath string) {
	l.Debug("Now loading koi config.")
	config, err := ReadConfig(configPath)
	if err != nil {
		l.Fatal(err)
	}
	Config = config
}
