package main

import (
	"gopkg.ilharper.com/koi/app/koicli"
	"gopkg.ilharper.com/koi/app/util"
	"gopkg.ilharper.com/koi/core/logger"
	"os"
)

func main() {
	l := logger.NewLogger(0)
	l.Register(logger.NewConsoleTarget())

	l.Infof("Koishi Desktop v%s", util.AppVersion)

	koicli.Run(os.Args, l)
}
