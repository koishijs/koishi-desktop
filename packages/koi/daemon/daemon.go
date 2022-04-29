package daemon

import (
	log "github.com/sirupsen/logrus"
	"koi/config"
	"koi/util"
)

var (
	// Log
	l = log.WithField("package", "daemon")
)

func Daemon() {
	l.Debug("Starting daemon.")
	for {
		l.Info("Starting Koishi.")

		dir, err := util.Resolve(config.Config.InternalInstanceDir, config.Config.Target, true)
		if err != nil {
			l.Fatalf("Failed to resolve target: %s", config.Config.Target)
		}

		err = RunNodeCmd(
			"npm",
			[]string{"run", "start"},
			dir,
		)

		if err == nil {
			l.Info("Koishi exited successfully.")
		} else {
			l.Error("Koishi exited with error:")
			l.Error(err)
		}
	}
}
