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

	do.ProvideNamed(i, "gopkg.ilharper.com/koi/app/koicli/command.Daemon", newDaemonCommand)

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
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Use configuration from `FILE`",
				EnvVars: []string{"KOI_CONFIG"},
			},

			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug mode",
			},

			cli.HelpFlag,
			cli.VersionFlag,
			cli.BashCompletionFlag,
		},

		Commands: []*cli.Command{
			do.MustInvokeNamed[*cli.Command](i, "gopkg.ilharper.com/koi/app/koicli/command.Daemon"),
		},

		Before: buildPreAction(i),
		CommandNotFound: func(context *cli.Context, s string) {
			l.Errorf("Command not found: %s", s)
		},
		ExitErrHandler: func(context *cli.Context, err error) {
			l.Error(err)
		},
	}, nil
}
