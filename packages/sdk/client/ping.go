package client

import (
	"fmt"

	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/util/net"
)

func Ping(conn *Options) (err error) {
	ws, err := Connect(conn)
	if err != nil {
		return
	}

	request := proto.NewRequest("ping", nil)

	err = net.JSON.Send(ws, request)
	if err != nil {
		return
	}

	var resp proto.Response
	err = net.JSON.Receive(ws, &resp)
	if err != nil {
		return
	}
	if resp.Type != "pong" {
		return fmt.Errorf("pingpong failed: response not 'pong'")
	}

	return
}
