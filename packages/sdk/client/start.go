package client

import (
	"fmt"

	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/core/util/net"
	"gopkg.ilharper.com/x/rpl"
)

func Start(
	conn *Options,
	instances []string,
) (<-chan *proto.Response, <-chan *rpl.Log, error) {
	var err error

	ws, err := Connect(conn)
	if err != nil {
		return nil, nil, err
	}

	request := proto.NewCommandRequest(
		"start",
		map[string]any{
			"instances": instances,
		},
	)

	err = net.JSON.Send(ws, request)
	if err != nil {
		return nil, nil, fmt.Errorf("websocket send error: %w", err)
	}

	wsRespC := make(chan *proto.Response)

	go func() {
		for {
			var resp proto.Response
			rErr := net.JSON.Receive(ws, &resp)
			if rErr != nil {
				wsRespC <- nil

				break
			}
			wsRespC <- &resp
		}
	}()

	respC, logC := logger.FilterLog(wsRespC)

	return respC, logC, nil
}
