package types

import "github.com/gorilla/websocket"

type Response struct {
	Write      func(Value)     // 写内容
	Connection *websocket.Conn // 客户连接
}
