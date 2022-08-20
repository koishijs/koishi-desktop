package koicli

import (
	"os"

	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"gopkg.ilharper.com/koi/app/config"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/x/rpl"
)

const (
	serviceActionPre = "gopkg.ilharper.com/koi/app/koicli/action.Pre"
)

func newPreAction(i *do.Injector) (cli.BeforeFunc, error) {
	l := do.MustInvoke[*logger.Logger](i)
	consoleTarget := do.MustInvoke[*logger.KoiFileTarget](i)

	return func(c *cli.Context) (err error) {
		l.Debug("Trigger pseudo action: pre")
		l.Debug("You're seeing debug output because you have a RPL target running in debug mode. This will not be controlled by --debug flag.")

		if c.Bool("debug") {
			consoleTarget.Level = rpl.LevelDebug
			l.Debug("--debug flag detected - debug mode enabled.")
		}

		l.Debugf("Command line arguments:\n%#+v", os.Args)

		do.Provide(i, config.BuildLoadConfig("koi.yml"))

		return
	}, nil
}
