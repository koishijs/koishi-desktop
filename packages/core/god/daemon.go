package god

import (
	"github.com/samber/do"
	"golang.org/x/net/websocket"
	"gopkg.ilharper.com/koi/core/koicmd"
)

type Daemon struct {
	// The [god.Task] registry.
	tasks taskRegistry

	// The [websocket.Handler].
	//
	// Functions are pointers so just store value.
	Handler websocket.Handler
}

func NewDaemon(i *do.Injector) *Daemon {
	do.Provide(i, koicmd.NewKoiCmdRegistry)

	daemon := &Daemon{}
	daemon.Handler = buildHandle(i, daemon)
	return daemon
}
