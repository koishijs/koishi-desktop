package daemonunlk

import (
	"os"
	"path/filepath"

	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/koiconfig"
)

type DaemonUnlocker struct {
	config *koiconfig.Config
}

func NewDaemonUnlocker(i *do.Injector) (*DaemonUnlocker, error) {
	cfg, err := do.Invoke[*koiconfig.Config](i)
	if err != nil {
		return nil, err
	}

	return &DaemonUnlocker{
		config: cfg,
	}, nil
}

func (unlocker *DaemonUnlocker) Shutdown() error {
	_ = os.Remove(filepath.Join(unlocker.config.Computed.DirLock, "daemon.lock"))

	// Do not short other do.Shutdownable
	return nil
}
