package koicmd

import (
	"fmt"

	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/god/daemonproc"
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/koierr"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/core/ui/webview"
)

func koiCmdOpen(i *do.Injector) *proto.Response {
	l := do.MustInvoke[*logger.Logger](i)
	command := do.MustInvoke[*proto.CommandRequest](i)
	daemonProc := do.MustInvoke[*daemonproc.DaemonProcess](i)

	l.Debug("Trigger KoiCmd open")

	// Parse command
	instances, ok := command.Flags["instances"].([]any)
	if !ok {
		return proto.NewErrorResult(koierr.ErrBadRequest)
	}

	for _, instanceAny := range instances {
		instance := instanceAny.(string)

		l.Infof("Opening instance %s...", instance)

		meta := daemonProc.GetMeta(instance)
		if meta == nil {
			return proto.NewErrorResult(koierr.NewErrInternalError(fmt.Errorf("cannot get meta of instance %s", instance)))
		}

		webview.Open(i, instance, meta.Listen)
	}

	return proto.NewSuccessResult(nil)
}
