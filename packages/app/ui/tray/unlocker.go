package tray

import (
	"os"
	"path/filepath"

	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/koiconfig"
)

type TrayUnlocker struct {
	config *koiconfig.Config
}

func NewTrayUnlocker(i *do.Injector) (*TrayUnlocker, error) {
	cfg, err := do.Invoke[*koiconfig.Config](i)
	if err != nil {
		return nil, err
	}

	return &TrayUnlocker{
		config: cfg,
	}, nil
}

func (unlocker *TrayUnlocker) Shutdown() error {
	_ = os.Remove(filepath.Join(unlocker.config.Computed.DirLock, "tray.lock"))

	// Do not short other do.Shutdownable
	return nil
}
