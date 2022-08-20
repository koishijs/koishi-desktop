package client

import (
	"fmt"

	"golang.org/x/net/websocket"
	"gopkg.ilharper.com/koi/core/god"
)

type Options struct {
	Host string
	Port string
}

// Connect tries to connect to Koi god daemon
// and returns a bare [websocket.Conn].
func Connect(options *Options) (client *websocket.Conn, err error) {
	return websocket.Dial(
		fmt.Sprintf("ws://%s:%s%s", options.Host, options.Port, god.DaemonEndpoint),
		"",
		fmt.Sprintf("http://%s:%s/", options.Host, options.Port),
	)
}
