package module

import (
	"encoding/json"
	"github.com/mlgaku/back/db"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
)

type Notice struct {
	Db db.Notice
}

// 获取通知列表
func (n *Notice) List(bd *Database, req *Request) Value {
	notice, _ := db.NewNotice(req.Body)

	dat, err := n.Db.FindByMaster(bd, notice.Master)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: dat}
}

// 移除通知
func (n *Notice) Remove(ps *Pubsub, bd *Database, req *Request) Value {
	notice, _ := db.NewNotice(req.Body)

	if err := n.Db.ChangeReadById(bd, notice.Id, true); err != nil {
		return &Fail{Msg: err.Error()}
	}

	ps.Publish(&Prot{Mod: "notice", Act: "list"})
	return &Succ{}
}
