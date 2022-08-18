package koicmd

import (
	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/koierr"
)

func koiCmdImport(i *do.Injector) *proto.Response {
	return proto.NewErrorResult(koierr.ErrNotImplemented)
}
