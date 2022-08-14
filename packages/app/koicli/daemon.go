package koicli

import (
	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/sdk/manage"
)

const (
	serviceCommandDaemon = "gopkg.ilharper.com/koi/app/koicli/command.Daemon"

	serviceActionDaemonPing = "gopkg.ilharper.com/koi/app/koicli/action.DaemonPing"
)

func newDaemonCommand(i *do.Injector) (*cli.Command, error) {
	do.ProvideNamed(i, serviceActionDaemonPing, newDaemonPingAction)

	return &cli.Command{
		Name:  "daemon",
		Usage: "Manage daemon",
		Subcommands: []*cli.Command{
			{
				Name:   "ping",
				Usage:  "Ping current daemon",
				Action: do.MustInvokeNamed[cli.ActionFunc](i, serviceActionDaemonPing),
			},
		},
	}, nil
}

func newDaemonPingAction(i *do.Injector) (cli.ActionFunc, error) {
	l := do.MustInvoke[*logger.Logger](i)

	return func(c *cli.Context) (err error) {
		cfg, err := do.Invoke[*koiconfig.Config](i)
		if err != nil {
			return
		}

		manager := manage.NewKoiManager(cfg.Computed.DirExe)
		conn, err := manager.Available()
		if err != nil {
			return
		}

		l.Success("PONG at:\n%#+v", conn)

		return
	}, nil
}
