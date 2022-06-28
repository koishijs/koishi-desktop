//go:generate goversioninfo

package main

import (
	"koi/cli"
	"koi/config"
	l "koi/util/logger"
	"os"
)

func main() {
	l.Infof("Koi %s", config.Version)

	err := cli.Run(os.Args)
	if err != nil {
		l.Fatal(err)
	}
}
