package koicli

import (
	"fmt"

	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"gopkg.ilharper.com/koi/app/util"
	"gopkg.ilharper.com/koi/core/logger"
)

func NewCli(i *do.Injector) (*cli.App, error) {
	l := do.MustInvoke[*logger.Logger](i)

	do.ProvideNamed(i, serviceActionPre, newPreAction)
	do.ProvideNamed(i, serviceCommandRun, newRunCommand)
	do.ProvideNamed(i, serviceCommandImport, newImportCommand)
	do.ProvideNamed(i, serviceCommandDaemon, newDaemonCommand)

	return &cli.App{
		Name:    "Koishi Desktop",
		Usage:   "Launch Koishi from your desktop.",
		Version: fmt.Sprintf("v%s", util.AppVersion),
		Authors: []*cli.Author{
			{
				Name:  "Il Harper",
				Email: "hi@ilharper.com",
			},
		},

		UseShortOptionHandling: true,
		EnableBashCompletion:   true,

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug mode",
			},

			cli.HelpFlag,
			cli.VersionFlag,
			cli.BashCompletionFlag,
		},

		Commands: []*cli.Command{
			do.MustInvokeNamed[*cli.Command](i, serviceCommandRun),
			do.MustInvokeNamed[*cli.Command](i, serviceCommandImport),
			do.MustInvokeNamed[*cli.Command](i, serviceCommandDaemon),
		},

		Before: do.MustInvokeNamed[cli.BeforeFunc](i, serviceActionPre),
		CommandNotFound: func(context *cli.Context, s string) {
			l.Errorf("Command not found: %s", s)
		},
		ExitErrHandler: func(context *cli.Context, err error) {
			l.Error(err)
		},
	}, nil
}
