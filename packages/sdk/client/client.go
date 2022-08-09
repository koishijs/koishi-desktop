package client

import (
	"fmt"
	"golang.org/x/net/websocket"
)

type Options struct {
	Host     string
	Port     string
	Endpoint string
}

// Connect tries to connect to Koi god daemon
// and returns a bare [websocket.Conn].
func Connect(options *Options) (client *websocket.Conn, err error) {
	return websocket.Dial(
		fmt.Sprintf("ws://%s:%s%s", options.Host, options.Port, options.Endpoint),
		"",
		fmt.Sprintf("http://%s:%s/", options.Host, options.Port),
	)
}
