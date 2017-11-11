package module

import (
	"encoding/json"
	"github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
)

type Sub struct {
	service.Di
}

func (*Sub) parse(body []byte) (*Prot, error) {
	p := &Prot{}
	return p, json.Unmarshal(body, p)
}

// 添加订阅
func (s *Sub) Add() {

	if prot, err := s.parse(s.Req().Body); err == nil {

		res := s.Res()

		v, e := service.NewModule(res.Client).LoadProt(prot)
		if e != nil {
			res.Write(res.Pack(*prot, &Fail{Msg: e.Error()}))
			return
		}

		if v != nil {
			res.Write(res.Pack(*prot, v))
		}
		s.Ps().AddSub(prot, res)

	}
}

// 取消订阅
func (s *Sub) Remove() {
	if prot, err := s.parse(s.Req().Body); err == nil {
		s.Ps().RemoveSub(prot, s.Res())
	}
}
