package webview

import (
	"os/exec"
	"strings"
	"sync/atomic"
	"time"

	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/koishell"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/x/browser"
)

func run(i *do.Injector, name string, listen string, ) (*exec.Cmd, bool) {
	var success atomic.Bool

	l := do.MustInvoke[*logger.Logger](i)
	shell := do.MustInvoke[*koishell.KoiShell](i)

	go func() {
		<-time.After(3 * time.Second)
		success.Store(true)
	}()

	cmd, err := shell.WebView(name, listen)
	if err != nil {
		l.Errorf("WebView error: %v", err)

		return /* cmd is always nil here */ nil, success.Load()
	}

	return cmd, true
}

func Open(i *do.Injector, name string, listen string, ) *exec.Cmd {
	l := do.MustInvoke[*logger.Logger](i)
	cfg := do.MustInvoke[*koiconfig.Config](i)

	listen = strings.ReplaceAll(listen, "0.0.0.0", "localhost")

	switch cfg.Data.Open {
	case "auto":
		l.Debugf("Running webview for instance %s: %s", name, listen)
		cmd, success := run(i, name, listen)
		if !success {
			go func() {
				l.Debugf("Failed to launch integrated webview. Fallback to external.")
				err := browser.OpenURL(listen)
				if err != nil {
					l.Warnf("cannot open browser: %s", err.Error())
				}
			}()
		}
		return cmd
	case "integrated":
		l.Debugf("Running webview for instance %s: %s", name, listen)
		cmd, _ := run(i, name, listen)
		return cmd
	case "external":
		l.Debugf("Running webview for instance %s: %s", name, listen)
		go func() {
			err := browser.OpenURL(listen)
			if err != nil {
				l.Warnf("cannot open browser: %s", err.Error())
			}
		}()
		return nil
	default:
		return nil
	}
}
