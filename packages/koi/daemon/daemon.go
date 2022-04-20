package daemon

import (
	log "github.com/sirupsen/logrus"
	"koi/config"
	"path"
)

var (
	// Log
	l = log.WithField("package", "daemon")
)

func Daemon() {
	l.Debug("Starting daemon.")
	for {
		l.Info("Starting Koishi.")

		err := RunNodeCmd(
			"yarn",
			[]string{"start"},
			path.Join(config.Config.InternalInstanceDir, config.Config.Target),
		)

		if err == nil {
			l.Info("Koishi exited successfully.")
		} else {
			l.Error("Koishi exited with error:")
			l.Error(err)
		}
	}
}
