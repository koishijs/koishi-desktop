package cli

import (
	"github.com/urfave/cli/v2"
	"koi/config"
	"koi/daemon"
	"koi/util"
	l "koi/util/logger"
	"koi/util/strutil"
)

var (
	yarnCommand = &cli.Command{
		Name:  "yarn",
		Usage: "Run yarn command",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "target",
				Usage: "Target to run",
			},
		},

		Action: yarnAction,
	}
)

func yarnAction(c *cli.Context) error {
	l.Info("Run yarn command:")

	args := c.Args().Slice()
	l.Info(args)

	target := strutil.Trim(c.String("target"))
	if target == "" {
		target = config.Config.Target
	}
	dir, err := util.Resolve(config.Config.InternalInstanceDir, target)
	if err != nil {
		l.Fatalf("Failed to resolve target: %s", target)
	}

	l.Infof("In instance %s", target)

	err = daemon.RunYarn(args, dir)
	if err != nil {
		l.Error("Yarn exited with error.")
		l.Fatal(err)
	}

	return nil
}
