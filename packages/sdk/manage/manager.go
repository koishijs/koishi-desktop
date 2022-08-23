package manage

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/goccy/go-json"
	"github.com/shirou/gopsutil/v3/process"
	"gopkg.ilharper.com/koi/core/god"
	"gopkg.ilharper.com/koi/sdk/client"
)

type KoiManager struct {
	exe  string
	lock string
}

func NewKoiManager(exe string, dirLock string) *KoiManager {
	return &KoiManager{
		exe:  exe,
		lock: filepath.Join(dirLock, "daemon.lock"),
	}
}

// Ensure god daemon available and can be used.
//
// Ensure handles all situations
// and you should generally use this method
// to get [client.Options] of god daemon.
func (manager *KoiManager) Ensure() (*client.Options, error) {
	var err error

	conn, availErr := manager.Available()
	if availErr == nil {
		return conn, nil
	}

	manager.Stop()
	err = manager.Start()
	if err != nil {
		return nil, err
	}

	return manager.Available()
}

// Available detects whether god daemon can be used.
//
// Available won't refresh daemon.
// If you want get [client.Options] of
// a available god daemon,
// use Ensure instead.
func (manager *KoiManager) Available() (*client.Options, error) {
	var err error

	conn, err := manager.Conn()
	if err != nil {
		return nil, err
	}

	err = client.Ping(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to ping god daemon: %w", err)
	}

	return conn, nil
}

// Start god daemon.
func (manager *KoiManager) Start() error {
	var cmd *exec.Cmd
	// THE B3ST S0lUt!0N
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/C", "start", "/b", manager.exe, "run", "daemon") //nolint:gosec
	} else {
		cmd = exec.Command("sh", "-c", fmt.Sprintf("%s run daemon &", manager.exe)) //nolint:gosec
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run daemon bootstrap shell: %w", err)
	}

	<-time.After(4 * time.Second)

	return nil
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
	defer func() {
		<-time.After(2 * time.Second)
	}()

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

	_ = manager.Kill()
}

// Kill the running daemon if exists.
//
// Do not use this method directly,
// use Stop instead.
func (manager *KoiManager) Kill() uint16 {
	killed := manager.tryKillProcesses()
	manager.tryDeleteLock()

	return killed
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

func (manager *KoiManager) tryKillProcesses() uint16 {
	killed, _ := manager.killProcesses()

	return killed
}

func (manager *KoiManager) killProcesses() (uint16, error) {
	var err error

	processes, err := manager.processes()
	if err != nil {
		return 0, err
	}

	var killed uint16 = 0
	for _, p := range processes {
		if p.Kill() == nil {
			killed++
		}
	}

	return killed, nil
}

func (manager *KoiManager) processes() ([]*process.Process, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes running: %w", err)
	}

	pss := make([]*process.Process, 0, len(processes))

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

	return pss, nil
}
