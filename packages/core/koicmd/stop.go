package koicmd

import (
	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/god/daemonproc"
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/koierr"
	"gopkg.ilharper.com/koi/core/logger"
)

func koiCmdStop(i *do.Injector) *proto.Response {
	var err error

	l := do.MustInvoke[*logger.Logger](i)
	command := do.MustInvoke[*proto.CommandRequest](i)
	daemonProc := do.MustInvoke[*daemonproc.DaemonProcess](i)

	l.Debug("Trigger KoiCmd stop")

	// Parse command
	instances, ok := command.Flags["instances"].([]any)
	if !ok {
		return proto.NewErrorResult(koierr.ErrBadRequest)
	}

	for _, instanceAny := range instances {
		instance := instanceAny.(string)

		l.Infof("Stopping instance %s...", instance)

		err = daemonProc.Stop(instance)
		if err != nil {
			return proto.NewErrorResult(koierr.NewErrInternalError(err))
		}
	}

	return proto.NewSuccessResult(nil)
}
