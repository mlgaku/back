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

type Client struct {
	Send       chan []byte     // 待发送数据
	Http       *http.Request   // 原始HTTP请求
	Server     *Server         // 客户所属 server
	Connection *websocket.Conn // 客户 ws 连接
}

// 读事件
func (c *Client) readPump() {
	defer func() {
		c.Server.unregister <- c
		c.Connection.Close()
	}()

	c.Connection.SetReadLimit(maxMessageSize)
	c.Connection.SetReadDeadline(time.Now().Add(pongWait))
	c.Connection.SetPongHandler(func(string) error { c.Connection.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, msg, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.Server.broadcast <- &message{c, bytes.TrimSpace(msg)}
	}
}

// 写事件
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Connection.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Connection.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(msg)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// 获得 Client 实例
func NewClient(ser *Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	cli := &Client{
		Send:       make(chan []byte, 256),
		Http:       r,
		Server:     ser,
		Connection: conn,
	}
	ser.register <- cli

	go cli.writePump()
	go cli.readPump()
}
