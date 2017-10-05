package types

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type Request struct {
	Body       []byte          // 字节内容
	BodyStr    string          // 字符串内容
	Http       *http.Request   // 原始HTTP请求
	Connection *websocket.Conn // 客户连接
}
