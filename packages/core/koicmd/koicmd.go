package koicmd

import (
	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/god/proto"
)

const (
	ServiceKoiCmdResponseChan = "gopkg.ilharper.com/koi/core/koicmd.KoiCmdResponseChan"
)

type KoiCmd func(i *do.Injector) *proto.Response

type Registry map[string]KoiCmd

func NewKoiCmdRegistry(i *do.Injector) (*Registry, error) {
	return &Registry{}, nil
}
