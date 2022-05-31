package config

import log "github.com/sirupsen/logrus"

var (
	// Log
	l = log.WithField("package", "config")

	Version = "v0.0.0"

	Config *KoiConfig

	defaultConfig = KoiConfig{
		Mode: "portable",

		Open: false,

		Strict: false,

		UseDataHome: true,
		UseDataTemp: true,
	}
)

type KoiConfig struct {
	Mode   string `yaml:"mode"`
	Target string `yaml:"target"`

	Open bool `yaml:"open"`

	Strict bool `yaml:"strict"`

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
