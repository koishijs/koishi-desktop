package logger

import (
	"golang.org/x/net/websocket"
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/util/net"
	"gopkg.ilharper.com/x/rpl"
)

type WsSender struct {
	c chan rpl.Log
}

func NewWsSender(localL *rpl.Logger, ws *websocket.Conn) *WsSender {
	var channel chan rpl.Log
	go func(l *rpl.Logger, w *websocket.Conn, c chan rpl.Log) {
		for {
			log := <-c
			err := net.JSON.Send(w, proto.NewLog(log))
			if err != nil {
				l.Errorf("failed to send log %v: %s", log, err)
			}
		}
	}(localL, ws, channel)
	return &WsSender{c: channel}
}

func (sender *WsSender) Writer() chan<- rpl.Log {
	return sender.c
}
