package god

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/samber/do"
	"golang.org/x/net/websocket"
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/koicmd"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/core/util/di"
	"gopkg.ilharper.com/koi/core/util/net"
	"net/http"
)

// Handle request.
// Here's already a new goroutine started by [websocket.Handler].
func buildHandle(i *do.Injector, daemon *Daemon) func(ws *websocket.Conn) {
	l := do.MustInvoke[*logger.Logger](i)

	return func(ws *websocket.Conn) {
		var err error

		l.Debugf("Client connected at %s", ws.RemoteAddr())

		defer func(ws *websocket.Conn) {
			closeErr := ws.Close()
			if closeErr != nil {
				l.Error(fmt.Errorf("failed to close ws connection: %w", closeErr))
			}
		}(ws)

		var request proto.Request
		err = net.JSON.Receive(ws, &request)
		if err != nil {
			l.Error(fmt.Errorf("failed to parse JSON request: %w", err))
			return
		}
		l.Debugf("Parsed request type: %s", request.Type)

		switch request.Type {
		case "ping":
			err = net.JSON.Send(ws, proto.NewResponse("pong", nil))
			if err != nil {
				l.Error(fmt.Errorf("failed to send 'pong': %w", err))
			}
			l.Debug("Send pong back")
			return
		case "stop":
			err = do.MustInvoke[*http.Server](i).Shutdown(context.Background())
			if err != nil {
				l.Error(fmt.Errorf("failed to close http server: %w", err))
			}
			return
		case proto.TypeRequestCommand:
			var commandRequest proto.CommandRequest
			err = mapstructure.Decode(request.Data, &commandRequest)
			if err != nil {
				l.Error(fmt.Errorf("failed to parse command: %w", err))
				return
			}
			l.Debugf("Parsed command: %s", commandRequest.Name)

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

	localL.Debugf("Acquired task %d", task.Id)

	// Build remote procedure Logger
	// Then override Logger
	do.Override(scopedI, logger.BuildNewLogger(uint16(task.Id)))

	// Build Response channel
	ch := make(chan *proto.Response)
	do.ProvideNamedValue(scopedI, koicmd.ServiceKoiCmdResponseChan, ch)

	// Build RPL Response Sender
	do.Provide(scopedI, logger.NewResponseSender)
	l := do.MustInvoke[*logger.Logger](scopedI)
	defer l.Close()
	// Register Senders
	l.Register(do.MustInvoke[*logger.KoiFileTarget](scopedI))
	l.Register(do.MustInvoke[*logger.ResponseSender](scopedI))

	// Get command registry
	// Use i here as registry is global provided
	reg := do.MustInvoke[*koicmd.Registry](i)
	// Get command
	kCmd, ok := (*reg)[command.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", command.Name)
	}

	send := make(chan bool)

	// Start sending response
	go func() {
		for {
			resp := <-ch
			if resp == nil {
				close(send)
				break
			}

			err := net.JSON.Send(ws, resp)
			if err != nil {
				localL.Error(fmt.Errorf("failed to send response: %w", err))
			}
		}
	}()

	// Invoke command
	response := kCmd(scopedI)
	if response != nil {
		ch <- response
	}
	close(ch)

	// Wait the final send finish
	<-send

	return nil
}
