package koicli

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"gopkg.ilharper.com/koi/core/god"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/logger"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

const (
	serviceActionRun       = "gopkg.ilharper.com/koi/app/koicli/action.Run"
	serviceActionRunDaemon = "gopkg.ilharper.com/koi/app/koicli/action.RunDaemon"
)

func newRunCommand(i *do.Injector) (*cli.Command, error) {
	do.ProvideNamed(i, serviceActionRun, newRunAction)
	do.ProvideNamed(i, serviceActionRunDaemon, newRunDaemonAction)

	return &cli.Command{
		Name:   "run",
		Usage:  "Run Koishi Desktop",
		Action: do.MustInvokeNamed[cli.ActionFunc](i, serviceActionRun),
		Subcommands: []*cli.Command{
			{
				Name:   "daemon",
				Usage:  "Run daemon",
				Action: do.MustInvokeNamed[cli.ActionFunc](i, serviceActionRunDaemon),
			},
		},
	}, nil
}

func newRunAction(i *do.Injector) (cli.ActionFunc, error) {
	return func(c *cli.Context) (err error) {
		cfg, err := do.Invoke[*koiconfig.Config](i)
		if err != nil {
			return
		}

		switch cfg.Data.Mode {
		case "cli":
			err = do.MustInvokeNamed[cli.ActionFunc](i, serviceActionRunDaemon)(c)
			return
		default:
			err = fmt.Errorf("unknown mode: %s", cfg.Data.Mode)
			return
		}
	}, nil
}

func newRunDaemonAction(i *do.Injector) (cli.ActionFunc, error) {
	l := do.MustInvoke[*logger.Logger](i)

	return func(c *cli.Context) (err error) {
		do.Provide(i, newDaemonUnlocker)

		cfg, err := do.Invoke[*koiconfig.Config](i)
		if err != nil {
			return
		}

		// Construct TCP listener
		listener, err := net.Listen("tcp4", "localhost:")
		if err != nil {
			return fmt.Errorf("failed to start daemon: %w", err)
		}
		addr := listener.Addr().String()

		l.Debug("Writing daemon.lock...")
		lock, err := os.OpenFile(
			filepath.Join(cfg.Computed.DirLock, "daemon.lock"),
			os.O_WRONLY|os.O_CREATE|os.O_EXCL, // Must create new file and write only
			0444,                              // -r--r--r--
		)

		daemonLock := &god.DaemonLock{
			Pid:  os.Getpid(),
			Addr: addr,
		}
		daemonLockJson, err := json.Marshal(daemonLock)
		if err != nil {
			return fmt.Errorf("failed to generate daemon lock data: %w", err)
		}
		_, err = lock.Write(daemonLockJson)
		if err != nil {
			return fmt.Errorf("failed to write daemon lock data: %w", err)
		}
		err = lock.Close()
		if err != nil {
			return fmt.Errorf("failed to close daemon lock: %w", err)
		}

		// Construct Daemon
		daemon := god.NewDaemon(i)

		mux := http.NewServeMux()
		mux.Handle("/api", daemon.Handler)

		server := &http.Server{Addr: addr, Handler: mux}
		l.Debug("Serving daemon...")
		err = server.Serve(listener)
		if err != nil {
			return fmt.Errorf("daemon closed: %w", err)
		}

		return
	}, nil
}

type daemonUnlocker struct {
	l      *logger.Logger
	config *koiconfig.Config
}

func newDaemonUnlocker(i *do.Injector) (*daemonUnlocker, error) {
	cfg, err := do.Invoke[*koiconfig.Config](i)
	if err != nil {
		return nil, err
	}

	return &daemonUnlocker{
		l:      do.MustInvoke[*logger.Logger](i),
		config: cfg,
	}, nil
}

func (unlocker *daemonUnlocker) Shutdown() error {
	err := os.Remove(filepath.Join(unlocker.config.Computed.DirLock, "daemon.lock"))
	if err != nil {
		unlocker.l.Errorf("failed to delete daemon lock: %w", err)
	}
	// Do not short other do.Shutdownable
	return nil
}
