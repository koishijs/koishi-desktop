package daemonproc

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/core/proc"
	"gopkg.ilharper.com/koi/core/util/instance"
	"gopkg.ilharper.com/koi/core/util/strutil"
	"gopkg.ilharper.com/x/browser"
)

const deltaCh uint16 = 1000

var ErrAlreadyStarted = errors.New("instance already started")

type DaemonProcess struct {
	// The mutex lock.
	//
	// There's no need to use [sync.RWMutex]
	// because almost all ops are write.
	mutex sync.Mutex
	wg    sync.WaitGroup

	i *do.Injector

	reg     [256]*proc.KoiProc
	nameReg map[string]uint8
}

func NewDaemonProcess(i *do.Injector) (*DaemonProcess, error) {
	return &DaemonProcess{
		i:       i,
		nameReg: make(map[string]uint8),
	}, nil
}

func (daemonProc *DaemonProcess) Init() error {
	var err error

	l := do.MustInvoke[*logger.Logger](daemonProc.i)
	cfg, err := do.Invoke[*koiconfig.Config](daemonProc.i)
	if err != nil {
		return err
	}

	l.Infof("Starting these instances:\n%s", strings.Join(cfg.Data.Start, ", "))

	daemonProc.mutex.Lock()
	defer daemonProc.mutex.Unlock()

	for _, name := range cfg.Data.Start {
		exists, existsErr := instance.IsInstanceExists(daemonProc.i, name)
		if existsErr != nil {
			l.Warn(existsErr)

			continue
		}
		if !exists {
			l.Warnf("Instance %s doesn't exist. Skipped.", name)

			continue
		}

		err = daemonProc.startIntl(name)
		if err != nil {
			l.Warnf("Failed to start %s: %s", name, err.Error())
		}
	}

	return nil
}

func (daemonProc *DaemonProcess) Start(name string) error {
	exists, existsErr := instance.IsInstanceExists(daemonProc.i, name)
	if existsErr != nil {
		return existsErr
	}
	if !exists {
		return fmt.Errorf("instance %s dows not exist", name)
	}

	daemonProc.mutex.Lock()
	defer daemonProc.mutex.Unlock()

	return daemonProc.startIntl(name)
}

// Must ensure lock before calling this method.
func (daemonProc *DaemonProcess) startIntl(name string) error {
	l := do.MustInvoke[*logger.Logger](daemonProc.i)
	cfg := do.MustInvoke[*koiconfig.Config](daemonProc.i)

	index := daemonProc.getIndex(name)

	koiProc := daemonProc.reg[index]
	if koiProc != nil {
		return ErrAlreadyStarted
	}

	koiProc = proc.NewYarnProc(
		daemonProc.i,
		deltaCh+uint16(index),
		[]string{"koishi", "start"},
		filepath.Join(cfg.Computed.DirInstance, name),
	)
	daemonProc.reg[index] = koiProc

	koiProc.Register(do.MustInvoke[*logger.KoiFileTarget](daemonProc.i))

	if cfg.Data.Open {
		koiProc.HookOutput = func(msg string) {
			go func() {
				if strings.Contains(msg, " server listening at ") {
					s := msg[strings.Index(msg, "http"):]           //nolint:gocritic
					s = s[:strings.Index(s, strutil.ColorStartCtr)] //nolint:gocritic
					l.Debugf("Parsed %s. Try opening browser.", s)
					err := browser.OpenURL(s)
					if err != nil {
						l.Warnf("cannot open browser: %s", err.Error())
					}
				}
			}()
		}
	}

	daemonProc.wg.Add(1)
	go func() {
		err := koiProc.Run()
		if err == nil {
			l.Infof("Instance %s exited.", name)
		} else {
			l.Warnf("Instance %s exited with: %s", name, err.Error())
		}

		defer daemonProc.wg.Done()
		daemonProc.mutex.Lock()
		defer daemonProc.mutex.Unlock()

		daemonProc.reg[index] = nil
	}()

	return nil
}

func (daemonProc *DaemonProcess) Stop(name string) error {
	exists, existsErr := instance.IsInstanceExists(daemonProc.i, name)
	if existsErr != nil {
		return existsErr
	}
	if !exists {
		return fmt.Errorf("instance %s dows not exist", name)
	}

	daemonProc.mutex.Lock()
	defer daemonProc.mutex.Unlock()

	return daemonProc.stopIntl(name)
}

// Must ensure lock before calling this method.
func (daemonProc *DaemonProcess) stopIntl(name string) error {
	return daemonProc.reg[daemonProc.nameReg[name]].Stop() //nolint:wrapcheck
}

func (daemonProc *DaemonProcess) Shutdown() error {
	daemonProc.mutex.Lock()

	for _, koiProc := range daemonProc.reg {
		if koiProc != nil {
			err := koiProc.Stop()
			if err != nil {
				_ = koiProc.Kill()
			}
		}
	}

	daemonProc.mutex.Unlock()
	daemonProc.wg.Wait()

	// Do not short other do.Shutdownable
	return nil
}

// getIndex finds the reg index of instance name.
//
// Must ensure lock before calling this method.
func (daemonProc *DaemonProcess) getIndex(name string) uint8 {
	var index uint8 = 0
	for n, i := range daemonProc.nameReg {
		if name == n {
			return i
		}
		index++
	}
	daemonProc.nameReg[name] = index

	return index
}

// GetPid find and return PID of instance.
//
// Returns 0 if instance is not running.
func (daemonProc *DaemonProcess) GetPid(name string) int {
	daemonProc.mutex.Lock()
	defer daemonProc.mutex.Unlock()

	koiProc := daemonProc.reg[daemonProc.getIndex(name)]
	if koiProc == nil {
		return 0
	}
	return koiProc.Pid()
}
