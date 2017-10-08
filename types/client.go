package types

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Send       *chan []byte    // 待发送数据
	Connection *websocket.Conn // 客户连接
}
