package tray

import (
	"runtime"

	"fyne.io/systray"
	"github.com/samber/do"
	"gopkg.ilharper.com/koi/app/ui/icon"
	"gopkg.ilharper.com/koi/core/logger"
)

func Run(i *do.Injector) error {
	systray.Run(buildOnReady(i), nil)

	return nil
}

func buildOnReady(i *do.Injector) func() {
	l := do.MustInvoke[*logger.Logger](i)

	return func() {
		l.Debug("Tray ready.")

		if runtime.GOOS != "darwin" {
			systray.SetTitle("Koishi")
		}
		systray.SetTooltip("Koishi")
		systray.SetTemplateIcon(icon.Data, icon.Data)

		mStarting := systray.AddMenuItem("Starting...", "Starting...")
		mStarting.Disable()
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Hide", "Hide Tray Button")

		_, ok := <-mQuit.ClickedCh
		if ok {
			l.Debugf("Exiting systray")
			systray.Quit()
		}
	}
}
