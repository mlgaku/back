package module

import (
	"encoding/json"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
)

type Sub struct {
	Mod string `json:"mod"`
	Act string `json:"act"`
}

func (*Sub) parse(body []byte) (string, error) {
	sub := &Sub{}
	return sub.Mod + "_" + sub.Act, json.Unmarshal(body, sub)
}

// 添加订阅
func (s *Sub) Add(ps *Pubsub, req *Request, res *Response) Value {
	id, _ := s.parse(req.Body)
	ps.AddSub(id, res)
	return id
}

// 取消订阅
func (s *Sub) Remove(ps *Pubsub, req *Request, res *Response) Value {
	id, _ := s.parse(req.Body)
	ps.RemoveSub(id, res)
	return id
}

func (s *Sub) Pub(ps *Pubsub, req *Request) Value {
	id, _ := s.parse(req.Body)
	ps.Publish(id)
	return id
}
