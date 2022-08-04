package koicli

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"gopkg.ilharper.com/koi/app/config"
	"gopkg.ilharper.com/koi/core/god"
	"gopkg.ilharper.com/koi/core/logger"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

func newDaemonCommand(i *do.Injector) (*cli.Command, error) {
	do.ProvideNamed(i, "gopkg.ilharper.com/koi/app/koicli/action.DaemonRun", newDaemonRunAction)

	return &cli.Command{
		Name:  "daemon",
		Usage: "Manage daemon",
		Subcommands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "Run daemon",
				Action: do.MustInvokeNamed[cli.ActionFunc](i, "gopkg.ilharper.com/koi/app/koicli/action.DaemonRun"),
			},
		},
	}, nil
}

func newDaemonRunAction(i *do.Injector) (cli.ActionFunc, error) {
	l := do.MustInvoke[*logger.Logger](i)
	cfg := do.MustInvoke[*config.Config](i)

	return func(c *cli.Context) (err error) {
		// Construct TCP listener
		listener, err := net.Listen("tcp4", "localhost:")
		if err != nil {
			return fmt.Errorf("failed to start daemon: %w", err)
		}
		addr := listener.Addr().String()

		l.Debug("Writing daemon.lock...")
		lock, err := os.OpenFile(
			filepath.Join(cfg.Computed.DirLock, "daemon.lock"),
			os.O_WRONLY|os.O_CREATE|os.O_EXCL, // Must create new file and write only
			0444,                              // -r--r--r--
		)
		// 【管理员】昵称什么的能吃吗 22:39:06
		// 死了也无所谓了 下次启动 check 一下

		daemonLock := &god.DaemonLock{
			Pid:  os.Getpid(),
			Addr: addr,
		}
		daemonLockJson, err := json.Marshal(daemonLock)
		if err != nil {
			return fmt.Errorf("failed to generate daemon lock data: %w", err)
		}
		_, err = lock.Write(daemonLockJson)
		if err != nil {
			return fmt.Errorf("failed to write daemon lock data: %w", err)
		}

		// Construct Daemon
		daemon := god.NewDaemon(i)

		mux := http.NewServeMux()
		mux.Handle("/api", daemon.Handler)

		server := &http.Server{Addr: addr, Handler: mux}
		l.Debug("Serving daemon...")
		err = server.Serve(listener)
		if err != nil {
			return fmt.Errorf("daemon closed: %w", err)
		}

		return
	}, nil
}
