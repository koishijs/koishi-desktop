package client

import (
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/util/net"
)

func Stop(conn *Options) (err error) {
	ws, err := Connect(conn)
	if err != nil {
		return
	}

	request := proto.NewRequest("stop", nil)

	err = net.JSON.Send(ws, request)
	if err != nil {
		return
	}

	return
}
