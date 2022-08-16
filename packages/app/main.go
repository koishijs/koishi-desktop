package main

import (
	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"gopkg.ilharper.com/koi/app/koicli"
	"gopkg.ilharper.com/koi/app/util"
	"gopkg.ilharper.com/koi/core/logger"
	coreUtil "gopkg.ilharper.com/koi/core/util"
	"gopkg.ilharper.com/x/rpl"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	defaultCommand = "run"
)

func main() {
	i := do.New()

	do.ProvideNamedValue(i, coreUtil.ServiceAppVersion, util.AppVersion)

	wg := &sync.WaitGroup{}
	do.ProvideValue(i, wg)

	do.Provide(i, logger.BuildNewKoiFileTarget(os.Stderr))
	do.Provide(i, logger.BuildNewLogger(0))
	receiver := rpl.NewReceiver()
	receiver.ChOffset = 100
	// Use ProvideValue() here because x/rpl didn't provide a do ctor
	do.ProvideValue(i, receiver)
	do.Provide(i, koicli.NewCli)

	l := do.MustInvoke[*logger.Logger](i)
	consoleTarget := do.MustInvoke[*logger.KoiFileTarget](i)
	receiver.Register(consoleTarget)
	l.Register(consoleTarget)

	l.Infof("Koishi Desktop v%s", util.AppVersion)

	args := os.Args
	if len(args) <= 1 {
		args = append(args, defaultCommand)
	}

	c := make(chan os.Signal)
	signal.Notify(
		c,
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		syscall.SIGINT,  // Ctrl-C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGKILL, // May not be caught
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
						l.Errorf("failed to gracefully shutdown: %w", err)
					}
					l.Close()
					wg.Wait()
					os.Exit(0)
				})
			}(s)
		}
	}()

	err := do.MustInvoke[*cli.App](i).Run(args)
	l.Close()
	wg.Wait()
	if err != nil {
		os.Exit(1)
	}
}
