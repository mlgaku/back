package service

import (
	"github.com/gorilla/websocket"
	"github.com/mlgaku/back/types"
)

type pubsub struct {
	list map[string]map[*websocket.Conn]*types.Response
}

// 发布
func (p *pubsub) publish(id string) {
	if p.list[id] == nil {
		return
	}

	for _, v := range p.list[id] {
		v.Write([]byte("hello world."))
	}
}

// 添加订阅
func (p *pubsub) addSub(id string, res *types.Response) {
	if p.list[id] == nil {
		p.list[id] = make(map[*websocket.Conn]*types.Response)
	}
	p.list[id][res.Client.Connection] = res
}

// 取消订阅
func (p *pubsub) removeSub(id string, res *types.Response) {
	if p.list[id] != nil {
		delete(p.list[id], res.Client.Connection)
	}
}

// 创建替身
func (p *pubsub) pseudo() *types.Pubsub {
	return &types.Pubsub{
		Publish: func(id string) {
			p.publish(id)
		},
		AddSub: func(id string, res *types.Response) {
			p.addSub(id, res)
		},
		RemoveSub: func(id string, res *types.Response) {
			p.removeSub(id, res)
		},
	}
}

// 获得 newPubsub 实例
func newPubsub() *pubsub {
	return &pubsub{
		list: make(map[string]map[*websocket.Conn]*types.Response),
	}
}
