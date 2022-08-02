package main

import (
	"gopkg.ilharper.com/koi/app/koicli"
	"gopkg.ilharper.com/koi/app/util"
	"gopkg.ilharper.com/koi/core/logger"
	"os"
)

func main() {
	l := logger.NewLogger(0)
	consoleTarget := logger.NewConsoleTarget()
	l.Register(consoleTarget)

	l.Infof("Koishi Desktop v%s", util.AppVersion)

	err := koicli.Run(os.Args, l, consoleTarget)
	if err != nil {
		l.Error(err)
		os.Exit(1)
	}
}
