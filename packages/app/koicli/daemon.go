package koicli

import (
	"fmt"

	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/sdk/manage"
)

const (
	serviceCommandDaemon = "gopkg.ilharper.com/koi/app/koicli/command.Daemon"

	serviceActionDaemonPing  = "gopkg.ilharper.com/koi/app/koicli/action.DaemonPing"
	serviceActionDaemonStart = "gopkg.ilharper.com/koi/app/koicli/action.DaemonStart"
	serviceActionDaemonStop  = "gopkg.ilharper.com/koi/app/koicli/action.DaemonStop"
	serviceActionDaemonKill  = "gopkg.ilharper.com/koi/app/koicli/action.DaemonKill"
)

func newDaemonCommand(i *do.Injector) (*cli.Command, error) {
	do.ProvideNamed(i, serviceActionDaemonPing, newDaemonPingAction)
	do.ProvideNamed(i, serviceActionDaemonStart, newDaemonStartAction)
	do.ProvideNamed(i, serviceActionDaemonStop, newDaemonStopAction)
	do.ProvideNamed(i, serviceActionDaemonKill, newDaemonKillAction)

	return &cli.Command{
		Name:  "daemon",
		Usage: "Manage daemon",
		Subcommands: []*cli.Command{
			{
				Name:   "ping",
				Usage:  "Ping current daemon",
				Action: do.MustInvokeNamed[cli.ActionFunc](i, serviceActionDaemonPing),
			},
			{
				Name:   "start",
				Usage:  "Start daemon",
				Action: do.MustInvokeNamed[cli.ActionFunc](i, serviceActionDaemonStart),
			},
			{
				Name:   "stop",
				Usage:  "Stop all daemons",
				Action: do.MustInvokeNamed[cli.ActionFunc](i, serviceActionDaemonStop),
			},
			{
				Name:   "kill",
				Usage:  "Kill all daemons",
				Action: do.MustInvokeNamed[cli.ActionFunc](i, serviceActionDaemonKill),
			},
		},
	}, nil
}

func newDaemonPingAction(i *do.Injector) (cli.ActionFunc, error) {
	l := do.MustInvoke[*logger.Logger](i)

	return func(c *cli.Context) error {
		var err error

		l.Debug("Trigger action: daemon ping")

		cfg, err := do.Invoke[*koiconfig.Config](i)
		if err != nil {
			return err
		}

		manager := manage.NewKoiManager(cfg.Computed.Exe, cfg.Computed.DirLock)
		conn, err := manager.Available()
		if err != nil {
			return fmt.Errorf("failed to get daemon status: %w", err)
		}

		l.Successf("PONG at:\n%#+v", conn)

		return nil
	}, nil
}

func newDaemonStartAction(i *do.Injector) (cli.ActionFunc, error) {
	l := do.MustInvoke[*logger.Logger](i)

	return func(c *cli.Context) error {
		var err error

		l.Debug("Trigger action: daemon start")

		cfg, err := do.Invoke[*koiconfig.Config](i)
		if err != nil {
			return err
		}

		manager := manage.NewKoiManager(cfg.Computed.Exe, cfg.Computed.DirLock)
		conn, err := manager.Ensure(true)
		if err != nil {
			return fmt.Errorf("failed to start daemon: %w", err)
		}

		l.Successf("Daemon started at:\n%#+v", conn)

		return nil
	}, nil
}

func newDaemonStopAction(i *do.Injector) (cli.ActionFunc, error) {
	l := do.MustInvoke[*logger.Logger](i)

	return func(c *cli.Context) error {
		var err error

		l.Debug("Trigger action: daemon stop")

		cfg, err := do.Invoke[*koiconfig.Config](i)
		if err != nil {
			return err
		}

		manager := manage.NewKoiManager(cfg.Computed.Exe, cfg.Computed.DirLock)
		manager.Stop()

		l.Success("All daemon stopped.")

		return nil
	}, nil
}

func newDaemonKillAction(i *do.Injector) (cli.ActionFunc, error) {
	l := do.MustInvoke[*logger.Logger](i)

	return func(c *cli.Context) error {
		l.Debug("Trigger action: daemon kill")

		cfg, err := do.Invoke[*koiconfig.Config](i)
		if err != nil {
			return err
		}

		killed := manage.NewKoiManager(cfg.Computed.Exe, cfg.Computed.DirLock).Kill()

		l.Successf("%d Daemon killed.", killed)

		return nil
	}, nil
}
