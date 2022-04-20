package cli

import (
	"github.com/urfave/cli/v2"
	"koi/daemon"
)

func runAction(c *cli.Context) error {
	l.Debug("Trigger action: run")
	daemon.Daemon()
	return nil
}
