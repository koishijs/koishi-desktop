package webview

import (
	"strings"

	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/koishell"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/x/browser"
)

func run(
	i *do.Injector,
	name string,
	listen string,
) {
	l := do.MustInvoke[*logger.Logger](i)
	shell := do.MustInvoke[*koishell.KoiShell](i)

	err := shell.WebView(name, listen)
	if err != nil {
		l.Errorf("WebView error: %v", err)
	}
}

func Open(
	i *do.Injector,
	name string,
	listen string,
) {
	l := do.MustInvoke[*logger.Logger](i)
	cfg := do.MustInvoke[*koiconfig.Config](i)

	listen = strings.ReplaceAll(listen, "0.0.0.0", "localhost")

	switch cfg.Data.Open {
	case "integrated":
		l.Debugf(
			"Running webview for instance %s: %s",
			name,
			listen,
		)
		go func() {
			run(i, name, listen)
		}()
	case "external":
		l.Debugf(
			"Running webview for instance %s: %s",
			name,
			listen,
		)
		err := browser.OpenURL(listen)
		if err != nil {
			l.Warnf("cannot open browser: %s", err.Error())
		}
	}
}
