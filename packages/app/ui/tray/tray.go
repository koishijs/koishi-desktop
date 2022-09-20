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

var refreshWaitDuration = 2 * time.Second

type TrayDaemon struct { //nolint:golint
	i       *do.Injector
	chanReg []chan struct{}
	manager *manage.KoiManager
}

func NewTrayDaemon(i *do.Injector) (*TrayDaemon, error) {
	cfg := do.MustInvoke[*koiconfig.Config](i)

	return &TrayDaemon{
		i:       i,
		manager: manage.NewKoiManager(cfg.Computed.Exe, cfg.Computed.DirLock),
	}, nil
}

func (tray *TrayDaemon) Run() error {
	systray.Run(tray.onReady, nil)

	return nil
}

func (tray *TrayDaemon) onReady() {
	l := do.MustInvoke[*logger.Logger](tray.i)

	l.Debug("Tray ready.")

	// systray.SetTitle("Koishi")
	if runtime.GOOS != "windows" {
		systray.SetTooltip("Koishi")
	}
	systray.SetTemplateIcon(icon.Koishi, icon.Koishi)

	mStarting := systray.AddMenuItem("Starting...", "")
	mStarting.Disable()
	tray.chanReg = append(tray.chanReg, mStarting.ClickedCh)

	tray.addItemsAfter()

	if _, err := tray.manager.Ensure(); err != nil {
		l.Error(err)
	}

	tray.rebuild()
}

func (tray *TrayDaemon) rebuild() {
	var err error

	l := do.MustInvoke[*logger.Logger](tray.i)

	conn, err := tray.manager.Available()
	if err != nil {
		systray.ResetMenu()
		for _, c := range tray.chanReg {
			close(c)
		}
		tray.chanReg = []chan struct{}{}
		tray.addItemsBefore()
		systray.AddSeparator()
		tray.addItemsAfter()

		return
	}

	respC, logC, err := client.Ps(conn, true)
	if err != nil {
		l.Error(err)

		return
	}

	logger.LogChannel(tray.i, logC)

	var result proto.Result
	response := <-respC
	if response == nil {
		l.Error("failed to get result, response is nil")

		return
	}

	if response.Type != proto.TypeResponseResult {
		l.Errorf("failed to parse result %#+v: response type '%s' is not '%s': %v", response, response.Type, proto.TypeResponseResult, err)

		return
	}

	err = mapstructure.Decode(response.Data, &result)
	if err != nil {
		l.Errorf("failed to parse result %#+v: %v", response, err)

		return
	}

	if result.Code != 0 {
		s, ok := result.Data.(string)
		if !ok {
			l.Errorf("result data %#+v is not string", result.Data)

			return
		}
		l.Error(s)

		return
	}

	var resultPs koicmd.ResultPs
	err = mapstructure.Decode(result.Data, &resultPs)
	if err != nil {
		l.Errorf("failed to parse result %#+v: %v", result, err)

		return
	}

	instances := resultPs.Instances

	// Clear all menu items.
	systray.ResetMenu()
	for _, c := range tray.chanReg {
		close(c)
	}
	tray.chanReg = []chan struct{}{}

	tray.addItemsBefore()

	if len(instances) > 0 {
		systray.AddSeparator()
	}

	// Iterate all instances.
	for _, instance := range instances {
		// Add menu items for each instance.
		m := systray.AddMenuItem(instance.Name, instance.Name)
		mOpen := m.AddSubMenuItem("Open", "Open")
		mOpen.SetTemplateIcon(icon.Open, icon.Open)
		mStart := m.AddSubMenuItem("Start", "Start")
		mStart.SetTemplateIcon(icon.Start, icon.Start)
		mRestart := m.AddSubMenuItem("Restart", "Restart")
		mRestart.SetTemplateIcon(icon.Restart, icon.Restart)
		mStop := m.AddSubMenuItem("Stop", "Stop")
		mStop.SetTemplateIcon(icon.Stop, icon.Stop)
		if instance.Running {
			mStart.Disable()
		} else {
			mOpen.Disable()
			mRestart.Disable()
			mStop.Disable()
		}

		tray.chanReg = append(tray.chanReg, mOpen.ClickedCh)
		tray.chanReg = append(tray.chanReg, mStart.ClickedCh)
		tray.chanReg = append(tray.chanReg, mRestart.ClickedCh)
		tray.chanReg = append(tray.chanReg, mStop.ClickedCh)

		go func(name string) {
			for {
				_, ok := <-mOpen.ClickedCh
				if !ok {
					break
				}

				l.Debugf("Opening instance %s", name)

				conn, err := tray.manager.Available()
				if err != nil {
					l.Error(err)

					continue
				}

				respC, logC, err := client.Open(
					conn,
					[]string{name},
				)
				if err != nil {
					l.Error(err)

					continue
				}

				logger.LogChannel(tray.i, logC)

				var result proto.Result
				for {
					response := <-respC
					if response == nil {
						l.Error("failed to get result, response is nil")

						break
					}
					if response.Type == proto.TypeResponseResult {
						err = mapstructure.Decode(response.Data, &result)
						if err != nil {
							l.Error("failed to parse result %#+v: %w", response, err)

							break
						}

						break
					}
					// Ignore other type of responses
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

				<-time.After(refreshWaitDuration)
				l.Debug("Rebuilding tray")
				tray.rebuild()
			}
		}(instance.Name)

		go func(name string) {
			for {
				_, ok := <-mStart.ClickedCh
				if !ok {
					break
				}

				l.Debugf("Starting instance %s", name)

				conn, err := tray.manager.Available()
				if err != nil {
					l.Error(err)

					continue
				}

				respC, logC, err := client.Start(
					conn,
					[]string{name},
				)
				if err != nil {
					l.Error(err)

					continue
				}

				logger.LogChannel(tray.i, logC)

				var result proto.Result
				for {
					response := <-respC
					if response == nil {
						l.Error("failed to get result, response is nil")

						break
					}
					if response.Type == proto.TypeResponseResult {
						err = mapstructure.Decode(response.Data, &result)
						if err != nil {
							l.Error("failed to parse result %#+v: %w", response, err)

							break
						}

						break
					}
					// Ignore other type of responses
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

				<-time.After(refreshWaitDuration)
				l.Debug("Rebuilding tray")
				tray.rebuild()
			}
		}(instance.Name)

		go func(name string) {
			for {
				_, ok := <-mStop.ClickedCh
				if !ok {
					break
				}

				l.Debugf("Stopping instance %s", name)

				conn, err := tray.manager.Available()
				if err != nil {
					l.Error(err)

					continue
				}

				respC, logC, err := client.Stop(
					conn,
					[]string{name},
				)
				if err != nil {
					l.Error(err)

					continue
				}

				logger.LogChannel(tray.i, logC)

				var result proto.Result
				for {
					response := <-respC
					if response == nil {
						l.Error("failed to get result, response is nil")

						break
					}
					if response.Type == proto.TypeResponseResult {
						err = mapstructure.Decode(response.Data, &result)
						if err != nil {
							l.Error("failed to parse result %#+v: %w", response, err)

							break
						}

						break
					}
					// Ignore other type of responses
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

				<-time.After(refreshWaitDuration)
				l.Debug("Rebuilding tray")
				tray.rebuild()
			}
		}(instance.Name)

		go func(name string) {
			for {
				_, ok := <-mRestart.ClickedCh
				if !ok {
					break
				}

				l.Debugf("Restarting instance %s", name)

				conn, err := tray.manager.Available()
				if err != nil {
					l.Error(err)

					continue
				}

				respC, logC, err := client.Restart(
					conn,
					[]string{name},
				)
				if err != nil {
					l.Error(err)

					continue
				}

				logger.LogChannel(tray.i, logC)

				var result proto.Result
				for {
					response := <-respC
					if response == nil {
						l.Error("failed to get result, response is nil")

						break
					}
					if response.Type == proto.TypeResponseResult {
						err = mapstructure.Decode(response.Data, &result)
						if err != nil {
							l.Error("failed to parse result %#+v: %w", response, err)

							break
						}

						break
					}
					// Ignore other type of responses
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

				<-time.After(refreshWaitDuration)
				l.Debug("Rebuilding tray")
				tray.rebuild()
			}
		}(instance.Name)
	}

	systray.AddSeparator()
	tray.addItemsAfter()
}

func (tray *TrayDaemon) addItemsBefore() {
	mTitle := systray.AddMenuItem("Koishi Desktop", "")
	mTitle.Disable()
	mTitle.SetTemplateIcon(icon.Koishi, icon.Koishi)
	version := "v" + util.AppVersion
	mVersion := systray.AddMenuItem(version, "")
	mVersion.Disable()

	tray.chanReg = append(tray.chanReg, mTitle.ClickedCh)
	tray.chanReg = append(tray.chanReg, mVersion.ClickedCh)
}

func (tray *TrayDaemon) addItemsAfter() {
	l := do.MustInvoke[*logger.Logger](tray.i)

	mAdvanced := systray.AddMenuItem("Advanced", "")
	mRefresh := mAdvanced.AddSubMenuItem("Refresh", "")
	mRefresh.SetTemplateIcon(icon.Restart, icon.Restart)
	mStartDaemon := mAdvanced.AddSubMenuItem("Start Daemon", "")
	mStartDaemon.SetTemplateIcon(icon.Start, icon.Start)
	mStopDaemon := mAdvanced.AddSubMenuItem("Stop Daemon", "")
	mStopDaemon.SetTemplateIcon(icon.Stop, icon.Stop)
	mKillDaemon := mAdvanced.AddSubMenuItem("Kill Daemon", "")
	mKillDaemon.SetTemplateIcon(icon.Kill, icon.Kill)
	mExit := mAdvanced.AddSubMenuItem("Stop and Exit", "")
	mExit.SetTemplateIcon(icon.Exit, icon.Exit)
	mQuit := systray.AddMenuItem("Hide", "")
	mQuit.SetTemplateIcon(icon.Hide, icon.Hide)

	tray.chanReg = append(tray.chanReg, mAdvanced.ClickedCh)
	tray.chanReg = append(tray.chanReg, mRefresh.ClickedCh)
	tray.chanReg = append(tray.chanReg, mStartDaemon.ClickedCh)
	tray.chanReg = append(tray.chanReg, mStopDaemon.ClickedCh)
	tray.chanReg = append(tray.chanReg, mKillDaemon.ClickedCh)
	tray.chanReg = append(tray.chanReg, mExit.ClickedCh)
	tray.chanReg = append(tray.chanReg, mQuit.ClickedCh)

	go func() {
		for {
			_, ok := <-mRefresh.ClickedCh
			if !ok {
				break
			}
			l.Debug("Rebuilding tray")
			tray.rebuild()
		}
	}()

	go func() {
		for {
			_, ok := <-mStartDaemon.ClickedCh
			if !ok {
				break
			}
			l.Debug("Starting daemon")
			err := tray.manager.Start()
			if err != nil {
				l.Error(err)

				continue
			}
			l.Debug("Rebuilding tray")
			tray.rebuild()
		}
	}()

	go func() {
		for {
			_, ok := <-mStopDaemon.ClickedCh
			if !ok {
				break
			}
			l.Debug("Stopping daemon")
			tray.manager.Stop()
			l.Debug("Rebuilding tray")
			tray.rebuild()
		}
	}()

	go func() {
		for {
			_, ok := <-mKillDaemon.ClickedCh
			if !ok {
				break
			}
			l.Debug("Killing daemon")
			tray.manager.Kill()
			l.Debug("Rebuilding tray")
			tray.rebuild()
		}
	}()

	go func() {
		for {
			_, ok := <-mExit.ClickedCh
			if !ok {
				break
			}
			l.Debug("Stopping daemon")
			tray.manager.Stop()
			l.Debug("Exiting systray")
			systray.Quit()
		}
	}()

	go func() {
		_, ok := <-mQuit.ClickedCh
		if ok {
			l.Debug("Exiting systray")
			systray.Quit()
		}
	}()
}
