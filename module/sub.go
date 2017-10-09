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

		v, e := NewModule(res.Client).LoadProt(prot)
		if v != nil && e == nil {
			res.Write(res.Pack(*prot, v))
		}

		ps.AddSub(prot, res)

	}
}

// 取消订阅
func (s *Sub) Remove(ps *Pubsub, req *Request, res *Response) {
	if prot, err := s.parse(req.Body); err == nil {
		ps.RemoveSub(prot, res)
	}
}
