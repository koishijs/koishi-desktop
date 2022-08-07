package god

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/samber/do"
	"golang.org/x/net/websocket"
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/koicmd"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/core/util/di"
	"gopkg.ilharper.com/koi/core/util/net"
)

// Handle request.
// Here's already a new goroutine started by [websocket.Handler].
func buildHandle(i *do.Injector, daemon *Daemon) func(ws *websocket.Conn) {
	l := do.MustInvoke[*logger.Logger](i)

	return func(ws *websocket.Conn) {
		var err error

		var request proto.Request
		err = net.JSON.Receive(ws, &request)
		if err != nil {
			l.Error(fmt.Errorf("failed to parse JSON request: %w", err))
			return
		}

		switch request.Type {
		case proto.TypeRequestCommand:
			var commandRequest proto.CommandRequest
			err = mapstructure.Decode(request.Data, &commandRequest)
			if err != nil {
				l.Error(fmt.Errorf("failed to parse command: %w", err))
			}
			err = handleCommand(i, daemon, ws, &commandRequest)
			if err != nil {
				l.Error(err)
			}
			return
		default:
			l.Errorf("unknown request type: %s", request.Type)
			return
		}
	}
}

func handleCommand(
	i *do.Injector,
	daemon *Daemon,
	ws *websocket.Conn,
	command *proto.CommandRequest,
) error {
	localL := do.MustInvoke[*logger.Logger](i)

	// Create scoped injector
	scopedI := di.Scope(i)

	// Acquire Task
	daemon.tasks.Acquire(scopedI)
	task := do.MustInvoke[*Task](scopedI)

	// Build remote procedure Logger
	// Then override Logger
	do.Override(i, logger.BuildNewLogger(uint16(task.Id)))

	// Build Response channel
	ch := make(chan *proto.Response)
	defer close(ch)
	do.ProvideNamedValue(i, koicmd.ServiceKoiCmdResponseChan, ch)

	// Build RPL Response Sender
	do.Provide(i, logger.NewResponseSender)
	l := do.MustInvoke[*logger.Logger](i)
	defer l.Close()
	// Register Senders
	l.Register(do.MustInvoke[*logger.ConsoleTarget](i))
	l.Register(do.MustInvoke[*logger.ResponseSender](i))

	// Get command registry
	reg := do.MustInvoke[*koicmd.Registry](i)
	// Get command
	kCmd, ok := (*reg)[command.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", command.Name)
	}

	// Start sending response
	go func(
		localL1 *logger.Logger,
		ws1 *websocket.Conn,
		ch1 <-chan *proto.Response,
	) {
		for {
			resp := <-ch1
			if resp == nil {
				break
			}

			err := net.JSON.Send(ws1, resp)
			if err != nil {
				localL1.Error(fmt.Errorf("failed to send response: %w", err))
			}
		}
	}(localL, ws, ch)

	// Invoke command
	response := kCmd(scopedI)
	if response != nil {
		ch <- response
	}
	return nil
}
