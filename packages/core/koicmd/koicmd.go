package koicmd

import (
	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/god/proto"
)

type KoiCmd func(i *do.Injector) *proto.Response

type Registry map[string]KoiCmd

func NewKoiCmdRegistry(i *do.Injector) (*Registry, error) {
	return &Registry{
		"import": koiCmdImport,
	}, nil
}
