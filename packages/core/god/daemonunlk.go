package god

import (
	"os"
	"path/filepath"

	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/logger"
)

type daemonUnlocker struct {
	l      *logger.Logger
	config *koiconfig.Config
}

func newDaemonUnlocker(i *do.Injector) (*daemonUnlocker, error) {
	cfg, err := do.Invoke[*koiconfig.Config](i)
	if err != nil {
		return nil, err
	}

	return &daemonUnlocker{
		l:      do.MustInvoke[*logger.Logger](i),
		config: cfg,
	}, nil
}

func (unlocker *daemonUnlocker) Shutdown() error {
	err := os.Remove(filepath.Join(unlocker.config.Computed.DirLock, "daemon.lock"))
	if err != nil {
		unlocker.l.Errorf("failed to delete daemon lock: %s", err)
	}
	// Do not short other do.Shutdownable
	return nil
}
