package koicli

import (
	"github.com/urfave/cli/v2"
	"gopkg.ilharper.com/koi/app/config"
	"gopkg.ilharper.com/koi/core/logger"
)

const (
	defaultCommand = "run"
)

type KoiCli struct {
	l             *logger.Logger
	consoleTarget *logger.ConsoleTarget

	app    *cli.App
	config *config.Config
}

func NewCli(
	l *logger.Logger,
	consoleTarget *logger.ConsoleTarget,
) *KoiCli {
	kcli := &KoiCli{
		l:             l,
		consoleTarget: consoleTarget,
	}

	kcli.app = newApp(kcli)
	return kcli
}

func (kcli *KoiCli) Run(args []string) error {
	return kcli.app.Run(args)
}

func Run(
	args []string,
	l *logger.Logger,
	consoleTarget *logger.ConsoleTarget,
) error {
	if len(args) <= 1 {
		args = append(args, defaultCommand)
	}

	return NewCli(l, consoleTarget).Run(args)
}
