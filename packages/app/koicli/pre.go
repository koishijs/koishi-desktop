package koicli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"golang.org/x/text/message"
	"gopkg.ilharper.com/koi/app/config"
	"gopkg.ilharper.com/koi/core/koishell"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/core/util"
	"gopkg.ilharper.com/x/rpl"
)

const (
	serviceActionPre = "gopkg.ilharper.com/koi/app/koicli/action.Pre"
)

func newPreAction(i *do.Injector) (cli.BeforeFunc, error) {
	l := do.MustInvoke[*logger.Logger](i)
	p := do.MustInvoke[*message.Printer](i)
	consoleTarget := do.MustInvokeNamed[*logger.KoiFileTarget](i, logger.ServiceConsoleTarget)

	return func(c *cli.Context) error {
		l.Debug(p.Sprintf("Trigger pseudo action: pre"))
		l.Debug(p.Sprintf("You're seeing debug output because this RPL target is running in debug mode (which is normal for syslog). This will not be controlled by --debug flag."))

		if c.Bool("debug") {
			consoleTarget.Level = rpl.LevelDebug
			l.Debug(p.Sprintf("--debug flag detected - debug mode enabled."))
		}

		l.Debug(p.Sprintf("PID: %s", fmt.Sprintf("%d", os.Getpid())))
		exe, err := os.Executable()
		if err != nil {
			return errors.New(p.Sprintf("failed to get executable: %v", err))
		}
		do.ProvideNamedValue(i, util.ServiceExecutable, exe)
		l.Debug(p.Sprintf("Executable: %s", exe))
		l.Debug(p.Sprintf("Command line arguments:\n%#+v", os.Args))

		do.Provide(i, config.BuildLoadConfig("koi.yml"))

		fileLogger, err := logger.BuildNewFileLogger()(i)
		if err != nil {
			l.Error(p.Sprintf("failed to build file logger: %v", err))
		} else {
			localReceiver := do.MustInvoke[*rpl.Receiver](i)
			localReceiver.Register(fileLogger)
		}

		var shellName string
		if runtime.GOOS == "windows" {
			shellName = "koishell.exe"
		} else {
			shellName = "koishell"
		}
		do.Provide(i, koishell.BuildKoiShell(filepath.Join(filepath.Dir(exe), shellName)))

		return nil
	}, nil
}
