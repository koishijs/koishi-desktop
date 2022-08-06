package net

import (
	"github.com/goccy/go-json"
	"golang.org/x/net/websocket"
)

func jsonMarshal(v any) (msg []byte, payloadType byte, err error) {
	msg, err = json.Marshal(v)
	return msg, 1 /* TextFrame */, err
}

func jsonUnmarshal(msg []byte, payloadType byte, v any) (err error) {
	return json.Unmarshal(msg, v)
}

var JSON = websocket.Codec{
	Marshal:   jsonMarshal,
	Unmarshal: jsonUnmarshal,
}
