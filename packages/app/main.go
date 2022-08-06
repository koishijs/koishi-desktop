package main

import (
	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"gopkg.ilharper.com/koi/app/koicli"
	"gopkg.ilharper.com/koi/app/util"
	"gopkg.ilharper.com/koi/core/logger"
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
	do.Provide(i, logger.NewConsoleTarget)
	do.Provide(i, logger.BuildNewLogger(0))
	do.Provide(i, koicli.NewCli)

	l := do.MustInvoke[*logger.Logger](i)
	l.Register(do.MustInvoke[*logger.ConsoleTarget](i))

	l.Infof("Koishi Desktop v%s", util.AppVersion)

	args := os.Args
	if len(args) <= 1 {
		args = append(args, defaultCommand)
	}

	var c chan os.Signal
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
					os.Exit(0)
				})
			}(s)
		}
	}()

	err := do.MustInvoke[*cli.App](i).Run(args)
	if err != nil {
		l.Error(err)
		os.Exit(1)
	}
}