package service

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 写超时时间
	writeWait = 10 * time.Second

	// 读超时时间
	pongWait = 60 * time.Second

	// ping 周期
	pingPeriod = (pongWait * 9) / 10

	// 最大消息大小(2M)
	maxMessageSize = 1024 * 1024 * 2
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// debug
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type client struct {
	send       chan []byte     // 待发送数据
	http       *http.Request   // 原始HTTP请求
	server     *server         // 客户所属 server
	connection *websocket.Conn // 客户 ws 连接
}

// 读事件
func (c *client) readPump() {
	defer func() {
		c.server.unregister <- c
		c.connection.Close()
	}()

	c.connection.SetReadLimit(maxMessageSize)
	c.connection.SetReadDeadline(time.Now().Add(pongWait))
	c.connection.SetPongHandler(func(string) error { c.connection.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, msg, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.server.broadcast <- &message{c, bytes.TrimSpace(msg)}
	}
}

// 写事件
func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.connection.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.connection.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				c.connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(msg)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// 获得 client 实例
func newClient(ser *server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	cli := &client{
		send:       make(chan []byte, 256),
		http:       r,
		server:     ser,
		connection: conn,
	}
	ser.register <- cli

	go cli.writePump()
	go cli.readPump()
}
