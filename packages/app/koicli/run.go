//nolint:wrapcheck
package koicli

import (
	"fmt"

	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"gopkg.ilharper.com/koi/app/ui/tray"
	"gopkg.ilharper.com/koi/core/god"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/logger"
)

const (
	serviceCommandRun = "gopkg.ilharper.com/koi/app/koicli/command.Run"

	serviceActionRun       = "gopkg.ilharper.com/koi/app/koicli/action.Run"
	serviceActionRunDaemon = "gopkg.ilharper.com/koi/app/koicli/action.RunDaemon"
	serviceActionRunUI     = "gopkg.ilharper.com/koi/app/koicli/action.RunUi"
)

func newRunCommand(i *do.Injector) (*cli.Command, error) {
	do.ProvideNamed(i, serviceActionRun, newRunAction)
	do.ProvideNamed(i, serviceActionRunDaemon, newRunDaemonAction)
	do.ProvideNamed(i, serviceActionRunUI, newRunUIAction)

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
			{
				Name:   "ui",
				Usage:  "Run UI",
				Action: do.MustInvokeNamed[cli.ActionFunc](i, serviceActionRunUI),
			},
		},
	}, nil
}

func newRunAction(i *do.Injector) (cli.ActionFunc, error) {
	l := do.MustInvoke[*logger.Logger](i)

	return func(c *cli.Context) error {
		var err error

		l.Debug("Trigger action: run")

		cfg, err := do.Invoke[*koiconfig.Config](i)
		if err != nil {
			return err
		}

		switch cfg.Data.Mode {
		case "cli":
			return do.MustInvokeNamed[cli.ActionFunc](i, serviceActionRunDaemon)(c)
		case "ui":
			return do.MustInvokeNamed[cli.ActionFunc](i, serviceActionRunUI)(c)
		default:
			return fmt.Errorf("unknown mode: %s", cfg.Data.Mode)
		}
	}, nil
}

func newRunDaemonAction(i *do.Injector) (cli.ActionFunc, error) {
	l := do.MustInvoke[*logger.Logger](i)

	return func(c *cli.Context) error {
		l.Debug("Trigger action: run daemon")

		return god.Daemon(i)
	}, nil
}

func newRunUIAction(i *do.Injector) (cli.ActionFunc, error) {
	l := do.MustInvoke[*logger.Logger](i)

	return func(c *cli.Context) error {
		l.Debug("Trigger action: run ui")

		do.Provide(i, tray.NewTrayDaemon)

		return do.MustInvoke[*tray.TrayDaemon](i).Run()
	}, nil
}
