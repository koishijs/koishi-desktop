package koicmd

import (
	"path/filepath"

	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/koierr"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/core/proc"
)

func koiCmdYarn(i *do.Injector) *proto.Response {
	var err error

	l := do.MustInvoke[*logger.Logger](i)
	command := do.MustInvoke[*proto.CommandRequest](i)
	cfg := do.MustInvoke[*koiconfig.Config](i)

	l.Debug("Trigger KoiCmd yarn")

	// Parse command
	instance, ok := command.Flags["instance"].(string)
	if !ok {
		return proto.NewErrorResult(koierr.ErrBadRequest)
	}
	argsAny, ok := command.Flags["args"].([]any)
	if !ok {
		return proto.NewErrorResult(koierr.ErrBadRequest)
	}

	args := make([]string, 0, len(argsAny))
	for _, argAny := range argsAny {
		args = append(args, argAny.(string))
	}

	koiProc := proc.NewYarnProc(
		i,
		2000,
		args,
		filepath.Join(cfg.Computed.DirInstance, instance),
	)

	koiProc.Register(do.MustInvoke[*logger.KoiFileTarget](i))
	koiProc.Register(do.MustInvoke[*logger.ResponseSender](i))

	l.Infof("Running command:\n%#+v\nOn instance: %s", args, instance)

	err = koiProc.Run()
	if err != nil {
		return proto.NewErrorResult(koierr.NewErrInternalError(err))
	}

	return proto.NewSuccessResult(nil)
}
