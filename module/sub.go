package module

import (
	"encoding/json"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
)

type Sub struct{}

func (*Sub) parse(body []byte) (*Prot, error) {
	p := &Prot{}
	return p, json.Unmarshal(body, p)
}

// 添加订阅
func (s *Sub) Add(ps *Pubsub, req *Request, res *Response) {
	if prot, err := s.parse(req.Body); err == nil {
		ps.AddSub(prot, res)
	}
}

// 取消订阅
func (s *Sub) Remove(ps *Pubsub, req *Request, res *Response) {
	if prot, err := s.parse(req.Body); err == nil {
		ps.RemoveSub(prot, res)
	}
}

func (s *Sub) Pub(ps *Pubsub) {
	ps.Publish("node_list")
}
