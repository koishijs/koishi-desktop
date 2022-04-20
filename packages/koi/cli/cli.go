package cli

import (
	"github.com/urfave/cli/v2"
	"koi/config"

	log "github.com/sirupsen/logrus"
)

const (
	KoiErrSpin = "KOI_ERR_SPIN"
)

var (
	// Log
	l = log.WithField("package", "cli")

	app = &cli.App{
		Name:    "Koi",
		Usage:   "The Koishi Launcher.",
		Version: config.Version,
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
			{
				Name:   "run",
				Usage:  "Run Koishi",
				Action: runAction,
			},

			instanceCommand,
		},

		Before: preAction,
		CommandNotFound: func(context *cli.Context, s string) {
			l.Fatalf("Command not found: %s", s)
		},
	}

	defaultCommand = "run"
)

func Run(args []string) error {
	if len(args) <= 1 {
		args = append(args, defaultCommand)
		return runIntl(args)
	}

	var commands []string
	for _, cmd := range app.Commands {
		commands = append(commands, cmd.Names()...)
	}

	for _, x := range args[1:] {
		for _, y := range commands {
			if x == y {
				return runIntl(args)
			}
		}
	}

	/* koi [global options] command [command options]
	 * Kinda hAcK for now !cuz "run" don't have options
	 */
	// newArgs := []string{args[0], defaultCommand}
	newArgs := []string{args[0]}
	newArgs = append(newArgs, args[1:]...)
	newArgs = append(newArgs, defaultCommand)

	return runIntl(newArgs)
}

func runIntl(args []string) error {
	l.Debug("Constructing cli app")

	cli.VersionPrinter = func(c *cli.Context) {
		l.Info(config.Version)
	}

	l.Debug("Running cli app")
	err := app.Run(args)
	if err != nil && err.Error() != KoiErrSpin {
		l.Fatal(err)
	}
	return err
}
