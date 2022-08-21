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
		return err
	}

	request := proto.NewRequest("stop", nil)

	err = net.JSON.Send(ws, request)
	if err != nil {
		return err
	}

	var resp proto.Response
	err = net.JSON.Receive(ws, &resp)
	if err != nil {
		return err
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
		return errors.New(result.Data.(string))
	}

	return nil
}
