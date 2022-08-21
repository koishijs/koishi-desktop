package proc

import (
	"path/filepath"

	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/koiconfig"
)

func NewNodeProc(
	i *do.Injector,
	ch uint16,
	command []string,
	cwd string,
) *KoiProc {
	cfg := do.MustInvoke[*koiconfig.Config](i)

	return NewKoiProc(
		i,
		ch,
		cfg.Computed.DirNodeExe,
		"node",
		command,
		cwd,
	)
}

func NewYarnProc(
	i *do.Injector,
	ch uint16,
	command []string,
	cwd string,
) *KoiProc {
	cfg := do.MustInvoke[*koiconfig.Config](i)

	return NewNodeProc(
		i,
		ch,
		append([]string{filepath.Join(cfg.Computed.DirNodeExe, "yarn.cjs")}, command...),
		cwd,
	)
}
