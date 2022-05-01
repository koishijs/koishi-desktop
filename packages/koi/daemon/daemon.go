package daemon

import (
	log "github.com/sirupsen/logrus"
	"koi/config"
	"koi/env"
	"koi/util"
	"os"
	"time"
)

var (
	// Log
	l = log.WithField("package", "daemon")
)

// Koishi exit status.
const (
	// StatusNormalRestart is for Koishi normal restart,
	// like user triggered config reload.
	StatusNormalRestart = iota

	// StatusErrorRestart is for Koishi error restart,
	// like port binding errors.
	// Here we use exponential backoff.
	StatusErrorRestart

	// StatusErrorExit is for Koishi error exit,
	// like koishi.yml config error.
	StatusErrorExit

	// StatusKoiStop is for Koishi request.
	StatusKoiStop

	// StatusKoiRestart is for Koishi request.
	// Here we use env.KoiErrSpin.
	StatusKoiRestart
)

const (
	initialInterval = 1 * time.Second
	maxInterval     = 60 * time.Second
	multiplier      = 2
)

func Daemon() error {
	l.Debug("Starting daemon.")

	dir, err := util.Resolve(config.Config.InternalInstanceDir, config.Config.Target, true)
	if err != nil {
		l.Fatalf("Failed to resolve target: %s", config.Config.Target)
	}

	var t *time.Timer
	defer func() {
		t.Stop()
	}()

	currentInterval := initialInterval

	for {
		l.Info("Starting Koishi.")
		code := daemonIntl(dir)

		if code == StatusErrorExit {
			l.Fatal("Exit due to Koishi error.")
		}

		if code == StatusKoiStop {
			l.Info("Exit due to request of Koishi.")
			os.Exit(0)
		}

		if code == StatusKoiRestart {
			l.Info("Restart due to request of Koishi.")
			return env.KoiErrSpin
		}

		if code == StatusNormalRestart {
			l.Info("Koishi exited successfully.")

			// Reset interval
			currentInterval = initialInterval
		} else {
			// StatusErrorRestart
			l.Error("Koishi exited with error.")

			next := currentInterval * multiplier
			if next > maxInterval {
				currentInterval = maxInterval
			} else {
				currentInterval = next
			}
		}

		l.Debugf("Wait for %s.", currentInterval.String())
		if t == nil {
			t = time.NewTimer(currentInterval)
		} else {
			t.Reset(currentInterval)
		}
		<-t.C
	}
}

func daemonIntl(dir string) int {
	err := RunYarn(
		[]string{"run", "start"},
		dir,
	)

	if err == nil {
		return StatusErrorRestart
	} else {
		return StatusErrorRestart
	}
}
