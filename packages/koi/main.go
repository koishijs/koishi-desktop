//go:generate goversioninfo

package main

import (
	"koi/cli"
	"koi/config"
	l "koi/util/logger"
	"os"
	"runtime"
)

func main() {
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
