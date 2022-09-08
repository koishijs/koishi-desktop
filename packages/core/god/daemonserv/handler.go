package daemonserv

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"github.com/samber/do"
	"golang.org/x/net/websocket"
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/god/task"
	"gopkg.ilharper.com/koi/core/koicmd"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/core/util/net"
)

// Handle request.
// Here's already a new goroutine started by [websocket.Handler].
func buildHandle(i *do.Injector, daemon *daemonService) func(ws *websocket.Conn) {
	l := do.MustInvoke[*logger.Logger](i)

	return func(ws *websocket.Conn) {
		var err error

		remoteAddr := ws.Request().RemoteAddr
		l.Debugf("Client connected at %s", remoteAddr)

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
			l.Infof("Stopping god daemon as request of %s...", remoteAddr)
			err = net.JSON.Send(ws, proto.NewSuccessResult(nil))
			if err != nil {
				l.Error(fmt.Errorf("failed to send stop response: %w", err))
			}
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
	daemon *daemonService,
	ws *websocket.Conn,
	command *proto.CommandRequest,
) error {
	localL := do.MustInvoke[*logger.Logger](i)

	// Create scoped injector
	scopedI := i.Scope()

	// Acquire Task
	daemon.tasks.Acquire(scopedI)
	t := do.MustInvoke[*task.Task](scopedI)
	defer func() {
		localL.Debugf("Releasing task %d", t.ID)
		daemon.tasks.Release(scopedI)
	}()

	localL.Debugf("Acquired task %d", t.ID)

	// Build remote procedure Logger
	// Then override Logger
	do.Override(scopedI, logger.BuildNewLogger(uint16(t.ID)))

	// Build Response channel
	ch := make(chan *proto.Response)
	do.ProvideValue(scopedI, ch)

	// Provide command
	do.ProvideValue(scopedI, command)

	// Build RPL Response Sender
	do.Provide(scopedI, logger.NewResponseSender)
	l := do.MustInvoke[*logger.Logger](scopedI)
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
	if response := kCmd(scopedI); response != nil {
		ch <- response
	}

	// Close ResponseSender here.
	// DO NOT call l.Close(), as it will close
	// KoiFileTarget the same time, which is used by localL.
	// Also, this must invoke synchronously before ch closed.
	do.MustInvoke[*logger.ResponseSender](scopedI).Close()
	ch <- nil

	// Wait the final send finish
	<-send

	return nil
}
