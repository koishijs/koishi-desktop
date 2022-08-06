package koiconn

import (
	"fmt"
	"golang.org/x/net/websocket"
)

type Option struct {
	Host     string
	Port     string
	Endpoint string
}

type KoiConn struct {
	Conn   *websocket.Conn
	option *Option
}

func Connect(option *Option) (conn *KoiConn, err error) {
	ws, err := websocket.Dial(
		fmt.Sprintf("ws://%s:%s%s", option.Host, option.Port, option.Endpoint),
		"",
		fmt.Sprintf("http://%s:%s/", option.Host, option.Port),
	)
	if err != nil {
		return
	}

	conn = &KoiConn{
		Conn:   ws,
		option: option,
	}
	return
}

func (conn *KoiConn) Close() error {
	return conn.Conn.Close()
}
