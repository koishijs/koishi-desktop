package tray

import (
	"runtime"
	"time"

	"fyne.io/systray"
	"github.com/mitchellh/mapstructure"
	"github.com/samber/do"
	"gopkg.ilharper.com/koi/app/util"
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/koicmd"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/core/ui/icon"
	"gopkg.ilharper.com/koi/sdk/client"
	"gopkg.ilharper.com/koi/sdk/manage"
)

const refreshDuration = 3 * time.Second

func Run(i *do.Injector) error {
	do.ProvideNamed(i, serviceTrayChannelRegistry, NewChannelRegistry)

	systray.Run(buildOnReady(i), nil)

	return nil
}

func buildOnReady(i *do.Injector) func() {
	l := do.MustInvoke[*logger.Logger](i)

	return func() {
		var err error

		l.Debug("Tray ready.")

		if runtime.GOOS != "darwin" {
			systray.SetTitle("Koishi")
		}
		if runtime.GOOS != "windows" {
			systray.SetTooltip("Koishi")
		}
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
		go trayDaemon(i, manager)
	}
}

func trayDaemon(i *do.Injector, manager *manage.KoiManager) {
	l := do.MustInvoke[*logger.Logger](i)
	channelRegistry := do.MustInvokeNamed[*ChannelRegistry](i, serviceTrayChannelRegistry)

	for {
		var err error

		conn, err := manager.Ensure()
		if err != nil {
			l.Error(err)

			continue
		}

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

		// Clear all menu items.
		systray.ResetMenu()

		addInfo(i)

		// Iterate all instances.
		for _, instance := range instances {
			// Add menu items for each instance.
			m := systray.AddMenuItem(instance.Name, instance.Name)
			mOpen := m.AddSubMenuItem("Open", "Open")
			mStart := m.AddSubMenuItem("Start", "Start")
			mRestart := m.AddSubMenuItem("Restart", "Restart")
			mStop := m.AddSubMenuItem("Stop", "Stop")
			if instance.Running {
				mStart.Disable()
			} else {
				mRestart.Disable()
				mStop.Disable()
			}

			channelRegistry.Insert(mOpen.ClickedCh)
			channelRegistry.Insert(mStart.ClickedCh)
			channelRegistry.Insert(mRestart.ClickedCh)
			channelRegistry.Insert(mStop.ClickedCh)
		}

		addHide(i)

		<-time.NewTimer(refreshDuration).C
	}
}

func addInfo(i *do.Injector) {
	channelRegistry := do.MustInvokeNamed[*ChannelRegistry](i, serviceTrayChannelRegistry)

	mTitle := systray.AddMenuItem("Koishi Desktop", "Koishi Desktop")
	mTitle.Disable()
	version := "v" + util.AppVersion
	mVersion := systray.AddMenuItem(version, version)
	mVersion.Disable()
	systray.AddSeparator()

	channelRegistry.Insert(mTitle.ClickedCh)
	channelRegistry.Insert(mVersion.ClickedCh)
}

func addHide(i *do.Injector) {
	l := do.MustInvoke[*logger.Logger](i)
	channelRegistry := do.MustInvokeNamed[*ChannelRegistry](i, serviceTrayChannelRegistry)

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Hide", "Hide Tray Button")

	channelRegistry.Insert(mQuit.ClickedCh)

	go func() {
		_, ok := <-mQuit.ClickedCh
		if ok {
			l.Debugf("Exiting systray")
			systray.Quit()
		}
	}()
}
