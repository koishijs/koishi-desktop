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

func NewWsSender(l *rpl.Logger, ws *websocket.Conn) *WsSender {
	var c chan rpl.Log
	go func(l1 *rpl.Logger, w *websocket.Conn, c1 chan rpl.Log) {
		for {
			log := <-c1
			go func(log1 rpl.Log) {
				err := net.JSON.Send(w, proto.NewLog(log1))
				if err != nil {
					l1.Errorf("failed to send log %v: %s", log1, err)
				}
			}(log)
		}
	}(l, ws, c)
	return &WsSender{c: c}
}

func (sender *WsSender) Writer() chan<- rpl.Log {
	return sender.c
}
