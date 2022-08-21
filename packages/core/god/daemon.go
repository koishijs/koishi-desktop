package god

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/goccy/go-json"
	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/logger"
)

// Daemon the world.
//
// This will serve as the main goroutine
// and run during the whole lifecycle.
func Daemon(i *do.Injector) error {
	var err error

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

	cfg, err := do.Invoke[*koiconfig.Config](i)
	if err != nil {
		return err
	}

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

	l.Debug("Writing daemon.lock...")
	lock, err := os.OpenFile(
		filepath.Join(cfg.Computed.DirLock, "daemon.lock"),
		os.O_WRONLY|os.O_CREATE|os.O_EXCL, // Must create new file and write only
		0o444,                             // -r--r--r--
	)

	daemonLock := &DaemonLock{
		Pid:  os.Getpid(),
		Host: host,
		Port: port,
	}
	daemonLockJson, err := json.Marshal(daemonLock)
	if err != nil {
		return fmt.Errorf("failed to generate daemon lock data: %w", err)
	}
	_, err = lock.Write(daemonLockJson)
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

	server := &http.Server{Addr: addr, Handler: mux}
	do.ProvideValue(i, server)
	l.Debug("Serving daemon...")
	err = server.Serve(listener)
	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	} else {
		err = fmt.Errorf("daemon closed: %w", err)
	}

	return nil
}
