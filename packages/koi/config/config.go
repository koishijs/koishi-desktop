package config

import log "github.com/sirupsen/logrus"

var (
	// Log
	l = log.WithField("package", "config")

	Version = "INTERNAL"
)

type KoiConfig struct {
	Mode   string `yaml:"mode"`
	Target string `yaml:"target"`

	// Internal
	ConfigDir string `yaml:"configDir,omitempty"`
}
