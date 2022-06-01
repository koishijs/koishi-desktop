package config

import (
	l "koi/util/logger"
)

var (
	Version = "v0.0.0"

	Config *KoiConfig

	defaultConfig = KoiConfig{
		Mode: "cli",

		Open: false,

		Strict: false,

		Env: nil,

		UseDataHome: true,
		UseDataTemp: true,
	}
)

type KoiConfig struct {
	Mode   string `yaml:"mode"`
	Target string `yaml:"target"`

	Open bool `yaml:"open"`

	Strict bool `yaml:"strict"`

	Env []string `yaml:"env"`

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
