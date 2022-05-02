package daemon

import (
	log "github.com/sirupsen/logrus"
	"koi/config"
	"koi/env"
	"koi/util"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	// Log
	l = log.WithField("package", "daemon")

	// Working dir
	dir string
	// Current Koishi process
	process *NodeCmd
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

// Daemon is to start and daemon (IntelliDaemon) the Koishi.
//
// The only thing to emphasize is that the Daemon will always be
// the main goroutine, and will last for the whole lifecycle for Koi
// once called.
//
// If you want to do something else, start a new goroutine.
func Daemon() error {
	l.Debug("Starting daemon.")

	resolvedDir, err := util.Resolve(config.Config.InternalInstanceDir, config.Config.Target)
	if err != nil {
		l.Fatalf("Failed to resolve target: %s", config.Config.Target)
	}
	dir = resolvedDir

	daemonHandleExit()

	return daemonMain()
}

func daemonHandleExit() {
	c := make(chan os.Signal)
	l.Debug("Setting up signal.Notify.")
	signal.Notify(
		c,
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		syscall.SIGINT,  // Ctrl-C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGKILL, // May not be caught
		syscall.SIGHUP,  // Terminal disconnected. SIGHUP also needs gracefully terminating
	)
	go func() {
		l.Debug("Starting daemonHandleExit.")
		for {
			// Waiting signal from c.
			// This will block the goroutine.
			s := <-c

			// Once received signal,
			// start another goroutine immediately and restore signal watching.
			// This can prevent the second signal terminating Koi.
			go func() {
				sig := s
				l.Debugf("Received signal %s", sig)
				l.Info("Terminating Koishi.")

				if process == nil {
					l.Debug("No active Koishi process.")
				} else {
					l.Debugf("Terminating process %d.", process.Cmd.Process.Pid)
					err := process.Cmd.Process.Signal(syscall.SIGINT)
					if err != nil {
						l.Debugf("Failed to send SIGINT to %d. Has exited?", process.Cmd.Process.Pid)
					} else {
						err = process.Cmd.Wait()
						if err != nil {
							l.Debugf("Koishi exited with %s", err)
						}
					}
				}

				l.Info("Exiting Koi.")
				os.Exit(0)
			}()
		}
	}()
}

func daemonMain() error {
	var t *time.Timer
	defer func() {
		t.Stop()
	}()

	currentInterval := initialInterval

	for {
		l.Info("Starting Koishi.")
		code := daemonRunCmd()

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

		l.Infof("Wait for %s.", currentInterval.String())
		if t == nil {
			t = time.NewTimer(currentInterval)
		} else {
			t.Reset(currentInterval)
		}
		<-t.C
	}
}

func daemonRunCmd() int {
	yarnPath, err := ResolveYarn()
	if err != nil {
		l.Fatal(err)
	}
	cmd := CreateNodeCmd(
		"node",
		[]string{yarnPath, "start"},
		dir,
	)

	l.Debug("Now start Koishi process.")
	err = cmd.Start()
	if err != nil {
		l.Error("Cannot start Koishi process.")
		l.Error(err)
		return StatusErrorRestart
	}

	l.Debug("Koishi process started.")
	l.Debugf("PID: %d", cmd.Cmd.Process.Pid)
	process = cmd

	l.Debug("Waiting process.")
	err = cmd.Wait()

	l.Debug("Cleaning process.")
	process = nil

	if err == nil {
		return StatusErrorRestart
	} else {
		return StatusErrorRestart
	}
}
