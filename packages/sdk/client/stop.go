package client

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/util/net"
)

func Stop(conn *Options) error {
	var err error

	ws, err := Connect(conn)
	if err != nil {
		return fmt.Errorf("failed to stop daemon: %w", err)
	}

	request := proto.NewRequest("stop", nil)

	err = net.JSON.Send(ws, request)
	if err != nil {
		return fmt.Errorf("failed to stop daemon: %w", err)
	}

	var resp proto.Response
	err = net.JSON.Receive(ws, &resp)
	if err != nil {
		return fmt.Errorf("failed to stop daemon: %w", err)
	}

	if resp.Type != proto.TypeResponseResult {
		return errors.New("failed to stop daemon: response type is not 'result'")
	}
	var result proto.Result
	err = mapstructure.Decode(resp.Data, &result)
	if err != nil {
		return fmt.Errorf("failed to stop daemon: failed to parse result %#+v: %w", resp.Data, err)
	}
	if result.Code != 0 {
		s, ok := result.Data.(string)
		if !ok {
			return fmt.Errorf("result data %#+v is not string", result.Data)
		}

		return errors.New(s)
	}

	return nil
}
