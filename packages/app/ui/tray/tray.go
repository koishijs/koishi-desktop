package tray

import (
	"errors"
	"fmt"
	util2 "gopkg.ilharper.com/koi/core/util"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"fyne.io/systray"
	"github.com/goccy/go-json"
	"github.com/mitchellh/mapstructure"
	"github.com/samber/do"
	"github.com/shirou/gopsutil/v3/process"
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
	"gopkg.ilharper.com/x/browser"
)

var refreshWaitDuration = 2 * time.Second

type TrayLock struct {
	Pid int `json:"pid" mapstructure:"pid"`
}

type TrayDaemon struct { //nolint:golint
	i       *do.Injector
	chanReg util2.ChannelRegistry[struct{}]
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
	var err error

	l := do.MustInvoke[*logger.Logger](tray.i)
	p := do.MustInvoke[*message.Printer](tray.i)
	cfg := do.MustInvoke[*koiconfig.Config](tray.i)
	shell := do.MustInvoke[*koishell.KoiShell](tray.i)

	// Tray mutex
	trayLockPath := filepath.Join(cfg.Computed.DirLock, "tray.lock")
	_, err = os.Stat(trayLockPath)
	if err != nil && (!(errors.Is(err, fs.ErrNotExist))) {
		return fmt.Errorf("failed to stat %s: %w", trayLockPath, err)
	}
	if err == nil {
		// tray.lock exists
		pid, aliveErr := checkTrayAlive(trayLockPath)
		if aliveErr == nil {
			shell.AlreadyRunning()
			return fmt.Errorf("tray running, PID=%d", pid)
		}

		_ = os.Remove(trayLockPath)
	}

	do.Provide(tray.i, NewTrayUnlocker)

	// tray.lock does not exist. Writing
	l.Debug(p.Sprintf("Writing tray.lock..."))
	lock, err := os.OpenFile(
		trayLockPath,
		os.O_WRONLY|os.O_CREATE|os.O_EXCL, // Must create new file and write only
		0o444,                             // -r--r--r--
	)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", trayLockPath, err)
	}

	trayLock := &TrayLock{
		Pid: os.Getpid(),
	}
	trayLockJSON, err := json.Marshal(trayLock)
	if err != nil {
		return fmt.Errorf("failed to generate tray lock data: %w", err)
	}
	_, err = lock.Write(trayLockJSON)
	if err != nil {
		return fmt.Errorf("failed to write tray lock data: %w", err)
	}
	err = lock.Close()
	if err != nil {
		return fmt.Errorf("failed to close tray lock: %w", err)
	}

	systray.Run(tray.onReady, tray.onExit)

	return nil
}

func (tray *TrayDaemon) onReady() {
	l := do.MustInvoke[*logger.Logger](tray.i)
	p := do.MustInvoke[*message.Printer](tray.i)

	l.Debug(p.Sprintf("Tray ready."))

	// systray.SetTitle("Cordis")
	if runtime.GOOS != "windows" {
		systray.SetTooltip("Cordis")
	}
	systray.SetTemplateIcon(icon.Koishi, icon.Koishi)

	mStarting := systray.AddMenuItem(p.Sprintf("Starting..."), "")
	mStarting.Disable()
	tray.chanReg.Register(mStarting.ClickedCh)

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

		tray.chanReg.Register(mOpen.ClickedCh)
		tray.chanReg.Register(mStart.ClickedCh)
		tray.chanReg.Register(mRestart.ClickedCh)
		tray.chanReg.Register(mStop.ClickedCh)

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

				respC, logC, err := client.Open(conn, []string{name})
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

				respC, logC, err := client.Start(conn, []string{name})
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

				respC, logC, err := client.Stop(conn, []string{name})
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

				respC, logC, err := client.Restart(conn, []string{name})
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

	mTitle := systray.AddMenuItem(p.Sprintf("Cordis Desktop"), "")
	mTitle.Disable()
	mTitle.SetTemplateIcon(icon.Koishi, icon.Koishi)
	version := "v" + util.AppVersion
	mVersion := systray.AddMenuItem(version, "")
	mVersion.Disable()

	tray.chanReg.Register(mTitle.ClickedCh)
	tray.chanReg.Register(mVersion.ClickedCh)
}

func (tray *TrayDaemon) addItemsAfter() {
	l := do.MustInvoke[*logger.Logger](tray.i)
	p := do.MustInvoke[*message.Printer](tray.i)
	cfg := do.MustInvoke[*koiconfig.Config](tray.i)
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
	mOpenDataFolder := mAdvanced.AddSubMenuItem(p.Sprintf("Open Data Folder"), "")
	mOpenTerminal := mAdvanced.AddSubMenuItem(p.Sprintf("Open Terminal"), "")
	mExit := mAdvanced.AddSubMenuItem(p.Sprintf("Stop and Exit"), "")
	mExit.SetTemplateIcon(icon.Exit, icon.Exit)
	mAbout := systray.AddMenuItem(p.Sprintf("About"), "")
	mQuit := systray.AddMenuItem(p.Sprintf("Hide"), "")
	mQuit.SetTemplateIcon(icon.Hide, icon.Hide)

	tray.chanReg.Register(mAdvanced.ClickedCh)
	tray.chanReg.Register(mRefresh.ClickedCh)
	tray.chanReg.Register(mStartDaemon.ClickedCh)
	tray.chanReg.Register(mStopDaemon.ClickedCh)
	tray.chanReg.Register(mKillDaemon.ClickedCh)
	tray.chanReg.Register(mOpenDataFolder.ClickedCh)
	tray.chanReg.Register(mOpenTerminal.ClickedCh)
	tray.chanReg.Register(mExit.ClickedCh)
	tray.chanReg.Register(mAbout.ClickedCh)
	tray.chanReg.Register(mQuit.ClickedCh)

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
			_, ok := <-mOpenDataFolder.ClickedCh
			if !ok {
				break
			}
			l.Debug(p.Sprintf("Opening data folder"))
			err := browser.OpenURL(cfg.Computed.DirData)
			if err != nil {
				l.Error(p.Sprintf("Failed to open data folder: %v", err))
			}
		}
	}()

	go func() {
		for {
			_, ok := <-mOpenTerminal.ClickedCh
			if !ok {
				break
			}
			l.Debug(p.Sprintf("Opening terminal"))

			var err error

			exe, err := os.Executable()
			if err != nil {
				l.Error(p.Sprintf("Failed to get executable: %v", err))
				continue
			}
			dirExe := filepath.Dir(exe)
			println(dirExe)

			var cmd *exec.Cmd
			if runtime.GOOS == "windows" {
				cmd = exec.Command(
					os.Getenv("COMSPEC"),
					"/C",
					"START",
					"Cordis Desktop Terminal",
					"/D",
					dirExe,
					"cmd", // Use "cmd" here
					"/K",
					"echo Cordis Desktop Terminal - You can start running koi command here.",
				)
			} else {
				cmd = exec.Command(
					"open",
					"-a",
					"Terminal",
					dirExe,
				)
			}

			err = cmd.Run()
			if err != nil {
				l.Error(p.Sprintf("Failed to open terminal: %v", err))
			}
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
			l.Debug(p.Sprintf("Showing about dialog"))
			go shell.About()
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

func checkTrayAlive(lockPath string) (int32, error) {
	var err error

	lockFile, err := os.ReadFile(lockPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read %s: %w", lockPath, err)
	}

	var lock TrayLock
	err = json.Unmarshal(lockFile, &lock)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s: %w", lockPath, err)
	}

	pid := int32(lock.Pid)
	proc, err := process.NewProcess(pid)
	if err != nil {
		return 0, fmt.Errorf("failed to get process %d: %w", pid, err)
	}

	isRunning, err := proc.IsRunning()
	if err != nil {
		return 0, fmt.Errorf("failed to get process %d state: %w", pid, err)
	}

	if !isRunning {
		return 0, fmt.Errorf("process %d is not running", pid)
	}

	return pid, nil
}
