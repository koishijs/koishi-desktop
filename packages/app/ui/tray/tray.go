package tray

import (
	"runtime"
	"time"

	"fyne.io/systray"
	"github.com/mitchellh/mapstructure"
	"github.com/samber/do"
	"gopkg.ilharper.com/koi/app/ui/icon"
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/koicmd"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/sdk/client"
	"gopkg.ilharper.com/koi/sdk/manage"
)

const refreshDuration = 3 * time.Second

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
		addHide(i)

		cfg, err := do.Invoke[*koiconfig.Config](i)
		if err != nil {
			l.Error(err)
			systray.Quit()
		}

		manager := manage.NewKoiManager(cfg.Computed.Exe, cfg.Computed.DirLock)
		// Ensure() only once. No need to Ensure() every tick.
		conn, err := manager.Ensure()
		if err != nil {
			l.Error(err)
			systray.Quit()
		}

		go func() {
			for {
				respC, logC, err := client.Ps(conn, true)
				if err != nil {
					l.Error(err)

					continue
				}

				logger.LogChannel(i, logC)

				var result proto.Result
				response := <-respC
				if response == nil {
					l.Error("failed to get result, response is nil")

					continue
				}

				if response.Type != proto.TypeResponseResult {
					l.Errorf("failed to parse result %#+v: response type '%s' is not '%s': %v", response, response.Type, proto.TypeResponseResult, err)

					continue
				}

				err = mapstructure.Decode(response.Data, &result)
				if err != nil {
					l.Errorf("failed to parse result %#+v: %v", response, err)

					continue
				}

				if result.Code != 0 {
					s, ok := result.Data.(string)
					if !ok {
						l.Errorf("result data %#+v is not string", result.Data)

						continue
					}
					l.Error(s)

					continue
				}

				var resultPs koicmd.ResultPs
				err = mapstructure.Decode(result.Data, &resultPs)
				if err != nil {
					l.Errorf("failed to parse result %#+v: %w", result, err)

					continue
				}

				instances := resultPs.Instances

				systray.ResetMenu()

				for _, instance := range instances {
					systray.AddMenuItem(instance.Name, instance.Name)
				}

				addHide(i)

				<-time.NewTimer(3 * time.Second).C
			}
		}()
	}
}

func addHide(i *do.Injector) {
	l := do.MustInvoke[*logger.Logger](i)

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Hide", "Hide Tray Button")

	go func() {
		_, ok := <-mQuit.ClickedCh
		if ok {
			l.Debugf("Exiting systray")
			systray.Quit()
		}
	}()
}
