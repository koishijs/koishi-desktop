package god

import (
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/goccy/go-json"
	"github.com/samber/do"
	"github.com/shirou/gopsutil/v3/process"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/logger"
)

// Daemon the world.
//
// This will serve as the main goroutine
// and run during the whole lifecycle.
func Daemon(i *do.Injector) error {
	var err error

	cfg, err := do.Invoke[*koiconfig.Config](i)
	if err != nil {
		return err
	}

	// Daemon mutex
	daemonLockPath := filepath.Join(cfg.Computed.DirLock, "daemon.lock")
	_, err = os.Stat(daemonLockPath)
	if err != nil && (!(errors.Is(err, fs.ErrNotExist))) {
		return fmt.Errorf("failed to stat %s: %w", daemonLockPath, err)
	}
	if err == nil {
		// daemon.lock exists
		pid, aliveErr := checkDaemonAlive(daemonLockPath)
		if aliveErr == nil {
			return fmt.Errorf("god daemon running, PID=%d\nCannot start another god daemon when there's already one.\nIf that daemon crashes, use 'koi daemon' to fix it.\nIf you just want to restart daemon, use 'koi daemon restart'", pid)
		}

		_ = os.Remove(daemonLockPath)
	}

	// Register daemonProcess synchronously,
	do.Provide(i, newDaemonProcess)
	// And start it in a new goroutine as early as possible.
	// This ensures Koishi starts quickly first.
	err = do.MustInvoke[*daemonProcess](i).init()
	if err != nil {
		return err
	}

	l := do.MustInvoke[*logger.Logger](i)

	// Provide daemonUnlocker.
	// It will try to remove daemon.lock when shutdown.
	do.Provide(i, newDaemonUnlocker)

	// Construct TCP listener
	listener, err := net.Listen("tcp4", "localhost:")
	if err != nil {
		return fmt.Errorf("failed to start daemon: %w", err)
	}
	addr := listener.Addr().String()
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("failed to parse addr %s: %w", addr, err)
	}

	// daemon.lock does not exist. Writing
	l.Debug("Writing daemon.lock...")
	lock, err := os.OpenFile(
		daemonLockPath,
		os.O_WRONLY|os.O_CREATE|os.O_EXCL, // Must create new file and write only
		0o444,                             // -r--r--r--
	)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", daemonLockPath, err)
	}

	daemonLock := &DaemonLock{
		Pid:  os.Getpid(),
		Host: host,
		Port: port,
	}
	daemonLockJSON, err := json.Marshal(daemonLock)
	if err != nil {
		return fmt.Errorf("failed to generate daemon lock data: %w", err)
	}
	_, err = lock.Write(daemonLockJSON)
	if err != nil {
		return fmt.Errorf("failed to write daemon lock data: %w", err)
	}
	err = lock.Close()
	if err != nil {
		return fmt.Errorf("failed to close daemon lock: %w", err)
	}

	// Construct daemonService
	service := newDaemonService(i)

	mux := http.NewServeMux()
	mux.Handle(DaemonEndpoint, service.Handler)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,

		ReadTimeout:       3 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      3 * time.Second,
	}
	do.ProvideValue(i, server)
	l.Debug("Serving daemon...")
	err = server.Serve(listener)
	if !(errors.Is(err, http.ErrServerClosed)) {
		return fmt.Errorf("daemon closed: %w", err)
	}

	return nil
}

func checkDaemonAlive(lockPath string) (int32, error) {
	var err error

	lockFile, err := os.ReadFile(lockPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read %s: %w", lockPath, err)
	}

	var lock DaemonLock
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
