package service

import (
	"github.com/gorilla/websocket"
)

type Pubsub struct {
	list map[string]map[*websocket.Conn]*Response
}

// 发布
func (p *Pubsub) Publish(id string) {
	if p.list[id] == nil {
		return
	}

	for _, v := range p.list[id] {
		v.Write([]byte("hello world."))
	}
}

// 添加订阅
func (p *Pubsub) AddSub(id string, res *Response) {
	if p.list[id] == nil {
		p.list[id] = make(map[*websocket.Conn]*Response)
	}
	p.list[id][res.Client.Connection] = res
}

// 取消订阅
func (p *Pubsub) RemoveSub(id string, res *Response) {
	if p.list[id] != nil {
		delete(p.list[id], res.Client.Connection)
	}
}

// 获得 NewPubsub 实例
func NewPubsub() *Pubsub {
	return &Pubsub{
		list: make(map[string]map[*websocket.Conn]*Response),
	}
}
