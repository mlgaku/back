package service

import (
	"github.com/gorilla/websocket"
	"github.com/mlgaku/back/types"
)

type (
	rp struct {
		res *Response
		pro *types.Prot
	}
	Pubsub struct {
		list map[string]map[*websocket.Conn]*rp
	}
)

func (p *Pubsub) getID(pro *types.Prot) string {
	return pro.Mod + "_" + pro.Act
}

// 发布
func (p *Pubsub) Publish(pro *types.Prot) {
	id := p.getID(pro)
	if p.list[id] == nil {
		return
	}

	val := map[string][]byte{}
	for _, v := range p.list[id] {

		if _, ok := val[v.pro.Body]; !ok {
			w, e := NewModule(v.res.Client).LoadProt(v.pro)
			if w == nil || e != nil {
				val[v.pro.Body] = nil
				continue
			}
			val[v.pro.Body] = v.res.Pack(*v.pro, w)
		}

		if w := val[v.pro.Body]; w != nil {
			v.res.Write(w)
		}

	}
}

// 添加订阅
func (p *Pubsub) AddSub(pro *types.Prot, res *Response) {
	id := p.getID(pro)
	if p.list[id] == nil {
		p.list[id] = make(map[*websocket.Conn]*rp)
	}
	p.list[id][res.Client.Connection] = &rp{pro: pro, res: res}
}

// 取消订阅
func (p *Pubsub) RemoveSub(pro *types.Prot, res *Response) {
	id := p.getID(pro)
	if p.list[id] != nil {
		delete(p.list[id], res.Client.Connection)
	}
}

// 获得 NewPubsub 实例
func NewPubsub() *Pubsub {
	return &Pubsub{
		list: make(map[string]map[*websocket.Conn]*rp),
	}
}
