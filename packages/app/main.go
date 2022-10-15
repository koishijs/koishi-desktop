package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"gopkg.ilharper.com/koi/app/koicli"
	"gopkg.ilharper.com/koi/app/util"
	"gopkg.ilharper.com/koi/core/logger"
	coreUtil "gopkg.ilharper.com/koi/core/util"
	"gopkg.ilharper.com/koi/core/util/hideconsole"
	"gopkg.ilharper.com/x/rpl"
)

func main() {
	l, _ := logger.BuildNewLogger(0)(nil)

	i := do.NewWithOpts(&do.InjectorOpts{
		Logf: func(format string, args ...any) {
			l.Debugf(format, args...)
		},
	})

	do.ProvideNamedValue(i, coreUtil.ServiceAppVersion, util.AppVersion)

	wg := &sync.WaitGroup{}
	do.ProvideValue(i, wg)

	do.Provide(i, logger.BuildNewKoiFileTarget(os.Stderr))
	do.ProvideValue(i, l)
	receiver := rpl.NewReceiver()
	receiver.ChOffset = 100
	// Use ProvideValue() here because x/rpl didn't provide a do ctor
	do.ProvideValue(i, receiver)
	do.Provide(i, koicli.NewCli)

	consoleTarget := do.MustInvoke[*logger.KoiFileTarget](i)
	receiver.Register(consoleTarget)
	l.Register(consoleTarget)

	l.Infof("Koishi Desktop v%s", util.AppVersion)

	noConsole := false

	args := os.Args
	if len(args) <= 1 {
		args = append(args, "--no-console", "run")
		noConsole = true
	} else {
		for _, arg := range args[1:] {
			if arg == "--no-console" {
				noConsole = true
				break
			}
		}
	}

	if noConsole {
		hideConsoleErr := hideconsole.HideConsole()
		if hideConsoleErr != nil {
			l.Warnf("Failed to hide console: %v", hideConsoleErr)
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(
		c,
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		syscall.SIGINT,  // Ctrl-C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGHUP,  // Terminal disconnected. SIGHUP also needs gracefully terminating
	)
	go func() {
		once := sync.Once{}
		for {
			s := <-c

			// Once received signal,
			// start another goroutine immediately and restore signal watching.
			// This can prevent the second signal terminating.
			go func(s1 os.Signal) {
				once.Do(func() {
					sig := s1
					l.Debugf("Received signal %s. Gracefully shutting down", sig)
					err := i.Shutdown()
					if err != nil {
						l.Errorf("failed to gracefully shutdown: %s", err)
					}
					l.Close()
					wg.Wait()
					os.Exit(0)
				})
			}(s)
		}
	}()

	runErr := do.MustInvoke[*cli.App](i).Run(args)
	if shutdownErr := i.Shutdown(); shutdownErr != nil {
		l.Errorf("failed to gracefully shutdown: %s", runErr)
	}
	l.Close()
	wg.Wait()
	if runErr != nil {
		os.Exit(1)
	}
}
