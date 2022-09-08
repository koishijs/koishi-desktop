package koicmd

import (
	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/god/daemonproc"
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/koierr"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/core/util/instance"
)

type ResultPsInstance struct {
	Name    string `json:"name" mapstructure:"name"`
	Running bool   `json:"running" mapstructure:"running"`
	Pid     int    `json:"pid" mapstructure:"pid"`
}

type ResultPs struct {
	Instances []*ResultPsInstance `json:"instances" mapstructure:"instances"`
}

func koiCmdPs(i *do.Injector) *proto.Response {
	var err error

	l := do.MustInvoke[*logger.Logger](i)
	command := do.MustInvoke[*proto.CommandRequest](i)
	daemonProc := do.MustInvoke[*daemonproc.DaemonProcess](i)

	l.Debug("Trigger KoiCmd ps")

	// Parse command
	all, ok := command.Flags["all"].(bool)
	if !ok {
		return proto.NewErrorResult(koierr.ErrBadRequest)
	}

	instanceNames, err := instance.Instances(i)
	if err != nil {
		return proto.NewErrorResult(koierr.NewErrInternalError(err))
	}

	result := &ResultPs{}

	for _, name := range instanceNames {
		pid := daemonProc.GetPid(name)
		running := pid != 0
		if (!running) && (!all) {
			continue
		}
		result.Instances = append(result.Instances, &ResultPsInstance{
			Name:    name,
			Running: running,
			Pid:     pid,
		})
	}

	return proto.NewSuccessResult(result)
}
