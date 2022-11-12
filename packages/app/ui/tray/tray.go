package tray

import (
	"os"
	"runtime"
	"sync"
	"time"

	"fyne.io/systray"
	"github.com/mitchellh/mapstructure"
	"github.com/samber/do"
	"golang.org/x/text/message"
	"gopkg.ilharper.com/koi/app/util"
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/koicmd"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/koishell"
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
	systray.Run(tray.onReady, tray.onExit)

	return nil
}

func (tray *TrayDaemon) onReady() {
	l := do.MustInvoke[*logger.Logger](tray.i)
	p := do.MustInvoke[*message.Printer](tray.i)

	l.Debug(p.Sprintf("Tray ready."))

	// systray.SetTitle("Koishi")
	if runtime.GOOS != "windows" {
		systray.SetTooltip("Koishi")
	}
	systray.SetTemplateIcon(icon.Koishi, icon.Koishi)

	mStarting := systray.AddMenuItem(p.Sprintf("Starting..."), "")
	mStarting.Disable()
	tray.chanReg = append(tray.chanReg, mStarting.ClickedCh)

	tray.addItemsAfter()

	if _, err := tray.manager.Ensure(true); err != nil {
		l.Error(err)
	}

	tray.rebuild()
}

func (tray *TrayDaemon) onExit() {
	l := do.MustInvoke[*logger.Logger](tray.i)
	wg := do.MustInvoke[*sync.WaitGroup](tray.i)

	_ = tray.i.Shutdown()
	l.Close()
	wg.Wait()

	os.Exit(0)
}

func (tray *TrayDaemon) rebuild() {
	var err error

	l := do.MustInvoke[*logger.Logger](tray.i)
	p := do.MustInvoke[*message.Printer](tray.i)

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
		l.Error(p.Sprintf("failed to get result, response is nil"))

		return
	}

	if response.Type != proto.TypeResponseResult {
		l.Error(p.Sprintf("failed to parse result %#+v: response type '%s' is not '%s': %v", response, response.Type, proto.TypeResponseResult, err))

		return
	}

	err = mapstructure.Decode(response.Data, &result)
	if err != nil {
		l.Error(p.Sprintf("failed to parse response %#+v: %v", response, err))

		return
	}

	if result.Code != 0 {
		s, ok := result.Data.(string)
		if !ok {
			l.Error(p.Sprintf("result data %#+v is not string", result.Data))

			return
		}
		l.Error(s)

		return
	}

	var resultPs koicmd.ResultPs
	err = mapstructure.Decode(result.Data, &resultPs)
	if err != nil {
		l.Error(p.Sprintf("failed to parse result %#+v: %v", result, err))

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
		mOpen := m.AddSubMenuItem(p.Sprintf("Open"), p.Sprintf("Open"))
		mOpen.SetTemplateIcon(icon.Open, icon.Open)
		mStart := m.AddSubMenuItem(p.Sprintf("Start"), p.Sprintf("Start"))
		mStart.SetTemplateIcon(icon.Start, icon.Start)
		mRestart := m.AddSubMenuItem(p.Sprintf("Restart"), p.Sprintf("Restart"))
		mRestart.SetTemplateIcon(icon.Restart, icon.Restart)
		mStop := m.AddSubMenuItem(p.Sprintf("Stop"), p.Sprintf("Stop"))
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

				l.Debug(p.Sprintf("Opening instance %s", name))

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
						l.Error(p.Sprintf("failed to get result, response is nil"))

						break
					}
					if response.Type == proto.TypeResponseResult {
						err = mapstructure.Decode(response.Data, &result)
						if err != nil {
							l.Error(p.Sprintf("failed to parse response %#+v: %v", response, err))

							break
						}

						break
					}
					// Ignore other type of responses
				}

				if result.Code != 0 {
					s, ok := result.Data.(string)
					if !ok {
						l.Error(p.Sprintf("result data %#+v is not string", result.Data))

						continue
					}

					l.Error(s)

					continue
				}

				<-time.After(refreshWaitDuration)
				l.Debug(p.Sprintf("Rebuilding tray"))
				tray.rebuild()
			}
		}(instance.Name)

		go func(name string) {
			for {
				_, ok := <-mStart.ClickedCh
				if !ok {
					break
				}

				l.Debug(p.Sprintf("Starting instance %s", name))

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
						l.Error(p.Sprintf("failed to get result, response is nil"))

						break
					}
					if response.Type == proto.TypeResponseResult {
						err = mapstructure.Decode(response.Data, &result)
						if err != nil {
							l.Error(p.Sprintf("failed to parse response %#+v: %v", response, err))

							break
						}

						break
					}
					// Ignore other type of responses
				}

				if result.Code != 0 {
					s, ok := result.Data.(string)
					if !ok {
						l.Error(p.Sprintf("result data %#+v is not string", result.Data))

						continue
					}

					l.Error(s)

					continue
				}

				<-time.After(refreshWaitDuration)
				l.Debug(p.Sprintf("Rebuilding tray"))
				tray.rebuild()
			}
		}(instance.Name)

		go func(name string) {
			for {
				_, ok := <-mStop.ClickedCh
				if !ok {
					break
				}

				l.Debug(p.Sprintf("Stopping instance %s", name))

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
						l.Error(p.Sprintf("failed to get result, response is nil"))

						break
					}
					if response.Type == proto.TypeResponseResult {
						err = mapstructure.Decode(response.Data, &result)
						if err != nil {
							l.Error(p.Sprintf("failed to parse response %#+v: %v", response, err))

							break
						}

						break
					}
					// Ignore other type of responses
				}

				if result.Code != 0 {
					s, ok := result.Data.(string)
					if !ok {
						l.Error(p.Sprintf("result data %#+v is not string", result.Data))

						continue
					}

					l.Error(s)

					continue
				}

				<-time.After(refreshWaitDuration)
				l.Debug(p.Sprintf("Rebuilding tray"))
				tray.rebuild()
			}
		}(instance.Name)

		go func(name string) {
			for {
				_, ok := <-mRestart.ClickedCh
				if !ok {
					break
				}

				l.Debug(p.Sprintf("Restarting instance %s", name))

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
						l.Error(p.Sprintf("failed to get result, response is nil"))

						break
					}
					if response.Type == proto.TypeResponseResult {
						err = mapstructure.Decode(response.Data, &result)
						if err != nil {
							l.Error(p.Sprintf("failed to parse response %#+v: %v", response, err))

							break
						}

						break
					}
					// Ignore other type of responses
				}

				if result.Code != 0 {
					s, ok := result.Data.(string)
					if !ok {
						l.Error(p.Sprintf("result data %#+v is not string", result.Data))

						continue
					}

					l.Error(s)

					continue
				}

				<-time.After(refreshWaitDuration)
				l.Debug(p.Sprintf("Rebuilding tray"))
				tray.rebuild()
			}
		}(instance.Name)
	}

	systray.AddSeparator()
	tray.addItemsAfter()
}

func (tray *TrayDaemon) addItemsBefore() {
	p := do.MustInvoke[*message.Printer](tray.i)

	mTitle := systray.AddMenuItem(p.Sprintf("Koishi Desktop"), "")
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
	p := do.MustInvoke[*message.Printer](tray.i)
	shell := do.MustInvoke[*koishell.KoiShell](tray.i)

	mAdvanced := systray.AddMenuItem(p.Sprintf("Advanced"), "")
	mRefresh := mAdvanced.AddSubMenuItem(p.Sprintf("Refresh"), "")
	mRefresh.SetTemplateIcon(icon.Restart, icon.Restart)
	mStartDaemon := mAdvanced.AddSubMenuItem(p.Sprintf("Start Daemon"), "")
	mStartDaemon.SetTemplateIcon(icon.Start, icon.Start)
	mStopDaemon := mAdvanced.AddSubMenuItem(p.Sprintf("Stop Daemon"), "")
	mStopDaemon.SetTemplateIcon(icon.Stop, icon.Stop)
	mKillDaemon := mAdvanced.AddSubMenuItem(p.Sprintf("Kill Daemon"), "")
	mKillDaemon.SetTemplateIcon(icon.Kill, icon.Kill)
	mExit := mAdvanced.AddSubMenuItem(p.Sprintf("Stop and Exit"), "")
	mExit.SetTemplateIcon(icon.Exit, icon.Exit)
	mAbout := systray.AddMenuItem("About", "")
	mQuit := systray.AddMenuItem(p.Sprintf("Hide"), "")
	mQuit.SetTemplateIcon(icon.Hide, icon.Hide)

	tray.chanReg = append(tray.chanReg, mAdvanced.ClickedCh)
	tray.chanReg = append(tray.chanReg, mRefresh.ClickedCh)
	tray.chanReg = append(tray.chanReg, mStartDaemon.ClickedCh)
	tray.chanReg = append(tray.chanReg, mStopDaemon.ClickedCh)
	tray.chanReg = append(tray.chanReg, mKillDaemon.ClickedCh)
	tray.chanReg = append(tray.chanReg, mExit.ClickedCh)
	tray.chanReg = append(tray.chanReg, mAbout.ClickedCh)
	tray.chanReg = append(tray.chanReg, mQuit.ClickedCh)

	go func() {
		for {
			_, ok := <-mRefresh.ClickedCh
			if !ok {
				break
			}
			l.Debug(p.Sprintf("Rebuilding tray"))
			tray.rebuild()
		}
	}()

	go func() {
		for {
			_, ok := <-mStartDaemon.ClickedCh
			if !ok {
				break
			}
			l.Debug(p.Sprintf("Starting daemon"))
			err := tray.manager.Start(true)
			if err != nil {
				l.Error(err)

				continue
			}
			l.Debug(p.Sprintf("Rebuilding tray"))
			tray.rebuild()
		}
	}()

	go func() {
		for {
			_, ok := <-mStopDaemon.ClickedCh
			if !ok {
				break
			}
			l.Debug(p.Sprintf("Stopping daemon"))
			tray.manager.Stop()
			l.Debug(p.Sprintf("Rebuilding tray"))
			tray.rebuild()
		}
	}()

	go func() {
		for {
			_, ok := <-mKillDaemon.ClickedCh
			if !ok {
				break
			}
			l.Debug(p.Sprintf("Killing daemon"))
			tray.manager.Kill()
			l.Debug(p.Sprintf("Rebuilding tray"))
			tray.rebuild()
		}
	}()

	go func() {
		for {
			_, ok := <-mExit.ClickedCh
			if !ok {
				break
			}
			l.Debug(p.Sprintf("Stopping daemon"))
			tray.manager.Stop()
			l.Debug(p.Sprintf("Exiting systray"))
			systray.Quit()
		}
	}()

	go func() {
		for {
			_, ok := <-mAbout.ClickedCh
			if !ok {
				break
			}
			l.Debug("Showing about dialog")
			shell.About()
		}
	}()

	go func() {
		_, ok := <-mQuit.ClickedCh
		if ok {
			l.Debug(p.Sprintf("Exiting systray"))
			systray.Quit()
		}
	}()
}
