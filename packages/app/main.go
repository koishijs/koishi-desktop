package main

import (
	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"gopkg.ilharper.com/koi/app/koicli"
	"gopkg.ilharper.com/koi/app/util"
	"gopkg.ilharper.com/koi/core/logger"
	"os"
)

const (
	defaultCommand = "run"
)

func main() {
	i := do.New()
	do.Provide(i, logger.NewConsoleTarget)
	do.Provide(i, logger.BuildNewLogger(0))
	do.Provide(i, koicli.NewCli)

	l := do.MustInvoke[*logger.Logger](i)
	l.Register(do.MustInvoke[*logger.ConsoleTarget](i))

	l.Infof("Koishi Desktop v%s", util.AppVersion)

	args := os.Args
	if len(args) <= 1 {
		args = append(args, defaultCommand)
	}
	err := do.MustInvoke[*cli.App](i).Run(args)
	if err != nil {
		l.Error(err)
		os.Exit(1)
	}
}
