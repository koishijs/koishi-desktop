package client

import (
	"fmt"
	"golang.org/x/net/websocket"
)

type KoiClient struct {
	conn *websocket.Conn
}

func Connect(
	host string,
	port string,
	endpoint string,
) (client *KoiClient, err error) {
	ws, err := websocket.Dial(
		fmt.Sprintf("ws://%s:%s%s", host, port, endpoint),
		"",
		fmt.Sprintf("http://%s:%s/", host, port),
	)
	if err != nil {
		return
	}

	client = &KoiClient{
		conn: ws,
	}
	return
}

func (client *KoiClient) Close() error {
	return client.conn.Close()
}
