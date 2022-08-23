//nolint:wrapcheck
package client

import (
	"fmt"
	"net"

	"golang.org/x/net/websocket"
	"gopkg.ilharper.com/koi/core/god"
)

type Options struct {
	Host string
	Port string
}

// Connect tries to connect to Koi god daemon
// and returns a bare [websocket.Conn].
func Connect(options *Options) (*websocket.Conn, error) {
	return websocket.Dial(
		fmt.Sprintf("ws://%s%s", net.JoinHostPort(options.Host, options.Port), god.DaemonEndpoint),
		"",
		fmt.Sprintf("http://%s/", net.JoinHostPort(options.Host, options.Port)),
	)
}
