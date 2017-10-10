package service

import "github.com/gorilla/websocket"

type Session struct {
	conn *websocket.Conn
	data map[*websocket.Conn]map[string]interface{}
}

// 获取
func (s *Session) Get(key string) interface{} {
	return APP.Session.data[s.conn][key]
}

// 设置
func (s *Session) Set(key string, val interface{}) {
	APP.Session.data[s.conn][key] = val
}

// 删除
func (s *Session) Remove(key string) {
	delete(APP.Session.data[s.conn], key)
}

// 销毁
func (*Session) Destroy(conn *websocket.Conn) {
	delete(APP.Session.data, conn)
}

// 获得 Session 实例
func NewSession(conn *websocket.Conn) *Session {
	if conn != nil && APP.Session.data[conn] == nil {
		APP.Session.data[conn] = make(map[string]interface{})
	}

	return &Session{
		conn: conn,
		data: make(map[*websocket.Conn]map[string]interface{}),
	}
}
