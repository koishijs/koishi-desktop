package daemonproc

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/samber/do"
	"golang.org/x/text/message"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/core/proc"
	"gopkg.ilharper.com/koi/core/ui/webview"
	"gopkg.ilharper.com/koi/core/util/instance"
	"gopkg.ilharper.com/koi/core/util/strutil"
	"gopkg.ilharper.com/x/rpl"
)

const deltaCh uint16 = 1000

var ErrAlreadyStarted = errors.New("instance already started")

type DProc struct {
	koiProc *proc.KoiProc
	Listen  string

	wvProcessRegLock sync.Mutex
	wvProcessReg     []*exec.Cmd
}

type DProcMeta struct {
	Pid    int
	Listen string
}

type DaemonProcess struct {
	// The mutex lock.
	//
	// There's no need to use [sync.RWMutex]
	// because almost all ops are write.
	mutex sync.Mutex
	wg    sync.WaitGroup

	i *do.Injector

	reg     [256]*DProc
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
	p := do.MustInvoke[*message.Printer](daemonProc.i)
	cfg, err := do.Invoke[*koiconfig.Config](daemonProc.i)
	if err != nil {
		return err
	}

	l.Info(p.Sprintf("Starting these instances:\n%s", strings.Join(cfg.Data.Start, ", ")))

	daemonProc.mutex.Lock()
	defer daemonProc.mutex.Unlock()

	for _, name := range cfg.Data.Start {
		exists, existsErr := instance.IsInstanceExists(daemonProc.i, name)
		if existsErr != nil {
			l.Warn(existsErr)

			continue
		}
		if !exists {
			l.Warn(p.Sprintf("Instance %s doesn't exist. Skipped.", name))

			continue
		}

		err = daemonProc.startIntl(name)
		if err != nil {
			l.Warn(p.Sprintf("Failed to start %s: %v", name, err))
		}
	}

	return nil
}

func (daemonProc *DaemonProcess) Start(name string) error {
	exists, existsErr := instance.IsInstanceExists(daemonProc.i, name)
	if existsErr != nil {
		return fmt.Errorf("check instance %s status failed: %w", name, existsErr)
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
	p := do.MustInvoke[*message.Printer](daemonProc.i)
	cfg := do.MustInvoke[*koiconfig.Config](daemonProc.i)

	index := daemonProc.getIndex(name)

	dp := daemonProc.reg[index]
	if dp != nil {
		return ErrAlreadyStarted
	}

	dp = &DProc{}
	daemonProc.reg[index] = dp

	daemonProc.wg.Add(1)
	go func() {
		for {
			dp.koiProc = proc.NewYarnProc(
				daemonProc.i,
				deltaCh+uint16(index),
				[]string{"start"},
				filepath.Join(cfg.Computed.DirInstance, name),
			)

			dp.koiProc.Register(do.MustInvoke[*rpl.Receiver](daemonProc.i))
			dp.koiProc.HookOutput = func(msg string) {
				go func() {
					if strings.Contains(msg, " server listening at ") {
						listen := msg[strings.Index(msg, "http"):]                     //nolint:gocritic
						listen = listen[:strings.Index(listen, strutil.ColorStartCtr)] //nolint:gocritic
						l.Debug(p.Sprintf("Parsed %s.", listen))
						dp.Listen = listen
						dp.wvProcessRegLock.Lock()
						if len(dp.wvProcessReg) == 0 {
							cmd := webview.Open(daemonProc.i, name, listen)
							if cmd != nil {
								dp.wvProcessReg = append(dp.wvProcessReg, cmd)
							}
						}
						dp.wvProcessRegLock.Unlock()
					}
				}()
			}

			err := dp.koiProc.Run()

			if exitErr, ok := err.(*exec.ExitError); ok {
				if exitErr.ExitCode() == 52 {
					l.Info(p.Sprintf("Instance %s exited with code restart (52). Restarting.", name))
				} else {
					l.Warn(p.Sprintf("Instance %s exited with: %v.", name, exitErr))
				}
			} else if err != nil {
				l.Warn(p.Sprintf("Instance %s exited: %v.", name, err))

				break
			} else {
				l.Info(p.Sprintf("Instance %s exited.", name))

				break
			}
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
		return fmt.Errorf("check instance %s status failed: %w", name, existsErr)
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
	l := do.MustInvoke[*logger.Logger](daemonProc.i)
	p := do.MustInvoke[*message.Printer](daemonProc.i)

	dp := daemonProc.reg[daemonProc.nameReg[name]]
	if err := dp.koiProc.Stop(); err != nil {
		l.Debug(p.Sprintf("failed to gracefully stop process %d: %v. Trying kill", dp.koiProc.Pid(), err))

		return dp.koiProc.Kill() //nolint:wrapcheck
	}

	return nil
}

func (daemonProc *DaemonProcess) Shutdown() error {
	daemonProc.mutex.Lock()

	for _, dp := range daemonProc.reg {
		if dp != nil {
			err := dp.koiProc.Stop()
			if err != nil {
				_ = dp.koiProc.Kill()
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

// GetMeta find and return meta info of instance.
//
// Returns nil if instance is not running.
func (daemonProc *DaemonProcess) GetMeta(name string) *DProcMeta {
	daemonProc.mutex.Lock()
	defer daemonProc.mutex.Unlock()

	dp := daemonProc.reg[daemonProc.getIndex(name)]
	if dp == nil {
		return nil
	}

	return &DProcMeta{
		Pid:    dp.koiProc.Pid(),
		Listen: dp.Listen,
	}
}

func (daemonProc *DaemonProcess) GetDProc(name string) *DProc {
	daemonProc.mutex.Lock()
	defer daemonProc.mutex.Unlock()

	return daemonProc.reg[daemonProc.getIndex(name)]
}

func (dp *DProc) AddWebViewProcess(cmd *exec.Cmd) {
	if cmd == nil {
		return
	}

	dp.wvProcessRegLock.Lock()
	defer dp.wvProcessRegLock.Unlock()

	dp.wvProcessReg = append(dp.wvProcessReg, cmd)
}
