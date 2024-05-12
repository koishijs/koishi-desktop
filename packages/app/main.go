package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/cloudfoundry/jibber_jabber"
	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.ilharper.com/koi/app/koicli"
	"gopkg.ilharper.com/koi/app/util"
	"gopkg.ilharper.com/koi/core/logger"
	coreUtil "gopkg.ilharper.com/koi/core/util"
	"gopkg.ilharper.com/x/hideconsole"
	"gopkg.ilharper.com/x/rpl"
	"gopkg.ilharper.com/x/setconsoleutf8"
)

func main() {
	lang, _ := jibber_jabber.DetectIETF()
	langTag, langTagErr := language.Parse(lang)
	if langTagErr != nil {
		langTag = language.MustParse("en-US")
	}
	p := message.NewPrinter(langTag)

	i := do.New()

	do.Provide(i, logger.BuildNewLogger(0))
	l := do.MustInvoke[*logger.Logger](i)

	do.ProvideNamedValue(i, coreUtil.ServiceAppVersion, util.AppVersion)
	do.ProvideNamedValue(i, coreUtil.ServiceAppBuildNumber, util.AppBuildNumber)

	do.ProvideValue(i, langTag)
	do.ProvideValue(i, p)

	wg := &sync.WaitGroup{}
	do.ProvideValue(i, wg)

	setupLogger(i)

	do.Provide(i, koicli.NewCli)

	l.Info(p.Sprintf("Cordis Desktop v%s", util.AppVersion))

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
			l.Warn(p.Sprintf("Failed to hide console: %v", hideConsoleErr))
		}
	}

	consoleUTF8Err := setconsoleutf8.SetConsoleUTF8()
	if consoleUTF8Err != nil {
		l.Warn(p.Sprintf("Failed to set console codepage to UTF-8: %v", consoleUTF8Err))
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
					l.Debug(p.Sprintf("Received signal %s. Gracefully shutting down", sig))
					_ = i.Shutdown()
					l.Close()
					wg.Wait()
					os.Exit(0)
				})
			}(s)
		}
	}()

	err := do.MustInvoke[*cli.App](i).Run(args)
	_ = i.Shutdown()
	l.Close()
	wg.Wait()
	if err != nil {
		os.Exit(1)
	}
}

func setupLogger(i *do.Injector) {
	p := do.MustInvoke[*message.Printer](i)

	// Get the local logger.
	l := do.MustInvoke[*logger.Logger](i)

	// Setup local receiver.
	//
	// Local receiver is a receiver hub that will receive logs from 2 sources:
	// - The local logger
	// - The KoiProc subprocess
	localReceiver := rpl.NewReceiver()

	// Register local receiver unnamed.
	//
	// This way, these RPL sources can be registered:
	//
	// - The local logger
	// - The KoiProc subprocess
	// - The remote logger (like the logger in daemonserv.scopedI)
	//
	// Using:
	//
	// source.Register(do.MustInvoke[*rpl.Receiver](i))
	do.ProvideValue(i, localReceiver)

	// Register local receiver to local logger.
	l.Register(localReceiver)

	// Setup remote receiver named.
	remoteReceiver := rpl.NewReceiver()
	// Set ChOffset to 100.
	remoteReceiver.ChOffset = 100
	do.ProvideNamedValue(i, logger.ServiceRemoteReceiver, remoteReceiver)

	// Setup console target named.
	do.ProvideNamed(i, logger.ServiceConsoleTarget, logger.BuildNewKoiFileTarget(os.Stderr))
	consoleTarget := do.MustInvokeNamed[*logger.KoiFileTarget](i, logger.ServiceConsoleTarget)
	// Register console target to receivers.
	localReceiver.Register(consoleTarget)
	remoteReceiver.Register(consoleTarget)

	// Register SysLogger to local receiver.
	sysLogger, err := logger.BuildNewSysLogger()(i)
	if err != nil {
		l.Error(p.Sprintf("Failed to create system logger: %v", err))
	} else {
		localReceiver.Register(sysLogger)
	}

	// FileLogger will be registered later, in the pseudo action pre.
}
