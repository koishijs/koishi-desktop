package manage

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/shirou/gopsutil/v3/process"
	"gopkg.ilharper.com/koi/core/god"
	"gopkg.ilharper.com/koi/core/util"
	"gopkg.ilharper.com/koi/sdk/client"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type KoiManager struct {
	exe  string
	lock string
}

func NewKoiManager(exe string, dirLock string) (manager *KoiManager) {
	manager = &KoiManager{
		exe:  exe,
		lock: filepath.Join(dirLock, "daemon.lock"),
	}
	return
}

// Ensure god daemon available and can be used.
//
// Ensure handles all situations
// and you should generally use this method
// to get [client.Options] of god daemon.
func (manager *KoiManager) Ensure() (conn *client.Options, err error) {
	conn, availErr := manager.Available()
	if availErr == nil {
		return
	}

	manager.Stop()
	err = manager.Start()
	if err != nil {
		return
	}

	<-time.After(util.TimeWait)
	conn, err = manager.Available()
	return
}

// Available detects whether god daemon can be used.
//
// Available won't refresh daemon.
// If you want get [client.Options] of
// a available god daemon,
// use Ensure instead.
func (manager *KoiManager) Available() (conn *client.Options, err error) {
	conn, err = manager.Conn()
	if err != nil {
		return
	}

	err = client.Ping(conn)
	if err != nil {
		return
	}

	return
}

// Start god daemon.
func (manager *KoiManager) Start() (err error) {
	cmd := exec.Cmd{
		Path: manager.exe,
		Args: []string{"run", "daemon"},
		Dir:  filepath.Dir(manager.exe),
	}
	err = cmd.Start()
	if err != nil {
		return
	}
	err = cmd.Process.Release()
	if err != nil {
		return
	}
	return
}

// Stop the running god daemon if exists.
//
// Stop promise that god daemon exits,
// and you should generally use this method
// to stop a god daemon.
//
// Stop will request a graceful stop first,
// wait a while ([util.TimeWait]),
// and if process still exists Stop will call Kill
// to ensure daemon dead.
func (manager *KoiManager) Stop() {
	conn, connErr := manager.Conn()
	done := false
	if connErr == nil {
		// Successfully get connection options.
		// Try gracefully shutdown.
		stopErr := client.Stop(conn)
		done = stopErr == nil
	}
	if done {
		manager.tryDeleteLock()
		return
	}

	manager.Kill()
	return
}

// Kill the running daemon if exists.
//
// Do not use this method directly,
// use Stop instead.
func (manager *KoiManager) Kill() {
	_ = manager.killProcesses()
	manager.tryDeleteLock()
}

func (manager *KoiManager) tryDeleteLock() {
	_ = os.Remove(manager.lock)
}

// Conn directly get [client.Options] from daemon.lock file.
//
// Do not use this method directly,
// use Ensure or Available instead.
func (manager *KoiManager) Conn() (*client.Options, error) {
	lock, err := manager.Lock()
	if err != nil {
		return nil, err
	}

	return &client.Options{
		Host: lock.Host,
		Port: lock.Port,
	}, nil
}

// Lock directly get [god.DaemonLock] from daemon.lock file.
//
// Do not use this method directly,
// use Ensure or Available instead.
func (manager *KoiManager) Lock() (*god.DaemonLock, error) {
	lockFile, err := os.ReadFile(manager.lock)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", manager.lock, err)
	}

	var lock god.DaemonLock
	err = json.Unmarshal(lockFile, &lock)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", manager.lock, err)
	}

	return &lock, nil
}

func (manager *KoiManager) killProcesses() (err error) {
	processes, err := manager.processes()
	if err != nil {
		return err
	}

	for _, p := range processes {
		_ = p.Kill()
	}

	return
}

func (manager *KoiManager) processes() (pss []*process.Process, err error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	selfPid := int32(os.Getpid())
	for _, p := range processes {
		if p.Pid == selfPid {
			continue
		}
		e, eErr := p.Exe()
		if eErr != nil {
			continue
		}
		if filepath.Clean(e) != manager.exe {
			continue
		}
		pss = append(pss, p)
	}

	return
}
