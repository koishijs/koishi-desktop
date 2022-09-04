package tray

import (
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
		systray.SetTitle("Koishi")
		systray.SetTooltip("Koishi")
		systray.SetTemplateIcon(icon.Data, icon.Data)

		mQuit := systray.AddMenuItem("Hide", "Hide Tray Button")

		go func() {
			for {
				select {
				case <-mQuit.ClickedCh:
					l.Debugf("Exiting systray")
					systray.Quit()
				}
			}
		}()
	}
}
