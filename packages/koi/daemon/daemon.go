package daemon

import (
	"github.com/ifrstr/browser"
	"koi/config"
	"koi/util"
	l "koi/util/logger"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// Daemon is to start and daemon the Koishi process.
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
		l.Errorf("Failed to resolve target: %s", config.Config.Target)
		l.Error(err)
		return err
	}

	cmd, err := CreateYarnCmd(
		[]string{"start"},
		resolvedDir,
	)
	if err != nil {
		l.Error("Err constructing KoiCmd:")
		l.Error(err)
		return err
	}

	if config.Config.Open {
		out := make(chan KoiCmdOut)
		cmd.Out = &out
		go func(out chan KoiCmdOut) {
			for {
				msg := <-out
				go func(msg KoiCmdOut) {
					if msg.IsErr && strings.Contains(msg.Text, " server listening at ") {
						l.Debug("Try opening browser. Parsing")
						l.Debug(msg.Text)
						s := msg.Text[strings.Index(msg.Text, "http"):]
						s = s[:strings.Index(s, "\u001b\u005b")]
						l.Debugf("Parsed %s", s)
						err := browser.OpenURL(s)
						if err != nil {
							l.Warn("Cannot open browser:")
							l.Warn(err)
						}
					}
				}(msg)
			}
		}(out)
	}

	l.Debug("Now start Koishi process.")
	err = cmd.Start()
	if err != nil {
		l.Error("Cannot start Koishi process.")
		l.Error(err)
		return err
	}

	l.Debug("Koishi process started.")
	l.Debugf("PID: %d", cmd.Cmd.Process.Pid)

	daemonHandleExit(cmd)

	l.Debug("Waiting process.")
	err = cmd.Wait()

	if err != nil {
		l.Error("Koishi exited with:")
		l.Error(err)
		return err
	}

	return nil
}

func daemonHandleExit(process *KoiCmd) {
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
							l.Debug("Koishi exited with:")
							l.Debug(err)
						}
					}
				}

				l.Info("Exiting Koi.")
				os.Exit(0)
			}()
		}
	}()
}
