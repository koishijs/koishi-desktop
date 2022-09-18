package daemonserv

import (
	"github.com/samber/do"
	"golang.org/x/net/websocket"
	"gopkg.ilharper.com/koi/core/god/task"
	"gopkg.ilharper.com/koi/core/koicmd"
)

type DaemonService struct {
	// The [god.Task] registry.
	tasks task.TaskRegistry

	// The [websocket.Handler].
	//
	// Functions are pointers so just store value.
	Handler websocket.Handler
}

func NewDaemonService(i *do.Injector) (*DaemonService, error) {
	do.Provide(i, koicmd.NewKoiCmdRegistry)

	service := &DaemonService{}
	service.Handler = buildHandle(i, service)

	return service, nil
}
