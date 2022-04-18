package main

import (
	"koi/cli"
	"koi/config"
	"os"
	"runtime"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

var (
	// Log
	l = log.WithField("package", "main")
)

func main() {
	// Initialize log
	log.SetFormatter(&formatter.Formatter{
		FieldsOrder: []string{"package", "component"},
		HideKeys:    true,
	})

	// Initialize environment
	l.Infof("Koi %s", config.Version)
	l.Infof("Go: %s", runtime.Version())

	for true {
		l.Debug("Start spin")
		err := cli.Run(os.Args)
		if err == nil {
			l.Debug("err == nil. Breaking.")
			break
		}
		l.Debug("Err: ", err)
		l.Debug("Spin")
	}
}
