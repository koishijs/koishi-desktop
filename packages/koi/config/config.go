package config

import log "github.com/sirupsen/logrus"

var (
	// Log
	l = log.WithField("package", "config")
)

type KoiConfig struct {
	Mode string `yaml:"mode"`

	// Internal
	ConfigDir string `yaml:"configDir,omitempty"`
}
