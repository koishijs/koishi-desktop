package koicli

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"gopkg.ilharper.com/koi/app/util"
)

type KoiApp struct {
	cliApp  *cli.App
	actions map[string]func(c *cli.Context) error
}

func newApp(kcli *KoiCli) *KoiApp {
	app := &KoiApp{}

	builders := []func(kcli *KoiCli) (map[string]func(c *cli.Context) error, *cli.Command){
		buildDaemonCommand,
	}
	var commands []*cli.Command
	for _, builder := range builders {
		commandActions, c := builder(kcli)
		for actionName, action := range commandActions {
			app.actions[actionName] = action
		}
		commands = append(commands, c)
	}

	app.cliApp = &cli.App{
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

		Commands: commands,

		Before: buildPreAction(kcli),
		CommandNotFound: func(context *cli.Context, s string) {
			kcli.l.Errorf("Command not found: %s", s)
		},
		ExitErrHandler: func(context *cli.Context, err error) {
			kcli.l.Error(err)
		},
	}

	return app
}
