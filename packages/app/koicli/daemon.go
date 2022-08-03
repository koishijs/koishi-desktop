package koicli

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/urfave/cli/v2"
	"gopkg.ilharper.com/koi/core/god"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

func buildDaemonCommand(kcli *KoiCli) (map[string]func(c *cli.Context) error, *cli.Command) {
	actions := map[string]func(c *cli.Context) error{
		"daemon run": buildDaemonRunAction(kcli),
	}

	cmd := &cli.Command{
		Name:  "daemon",
		Usage: "Manage daemon",
		Subcommands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "Run daemon",
				Action: actions["daemon run"],
			},
		},
	}

	return actions, cmd
}

func buildDaemonRunAction(kcli *KoiCli) func(c *cli.Context) error {
	return func(c *cli.Context) (err error) {
		// Construct TCP listener
		listener, err := net.Listen("tcp4", "localhost:")
		if err != nil {
			return fmt.Errorf("failed to start daemon: %w", err)
		}
		addr := listener.Addr().String()

		kcli.l.Debug("Writing daemon.lock...")
		lock, err := os.OpenFile(
			filepath.Join(kcli.config.Computed.DirLock, "daemon.lock"),
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
		daemon := god.NewDaemon(kcli.l)

		mux := http.NewServeMux()
		mux.Handle("/api", daemon.Handler)

		server := &http.Server{Addr: addr, Handler: mux}
		kcli.l.Debug("Serving daemon...")
		err = server.Serve(listener)
		if err != nil {
			return fmt.Errorf("daemon closed: %w", err)
		}

		return
	}
}
