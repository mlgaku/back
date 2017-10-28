package module

import (
	"github.com/mlgaku/back/db"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
)

type Notice struct {
	Db db.Notice
}

// 获取通知列表
func (n *Notice) List(bd *Database, ses *Session) Value {
	dat, err := n.Db.FindByMaster(
		bd, ses.Get("user").(*db.User).Id,
	)

	if err != nil {
		return &Fail{Msg: err.Error()}
	}
	return &Succ{Data: dat}
}

// 移除通知
func (n *Notice) Remove(ps *Pubsub, bd *Database, req *Request, ses *Session) Value {
	notice, _ := db.NewNotice(req.Body)

	if err := n.Db.Find(bd, notice.Id, notice); err != nil {
		return &Fail{Msg: err.Error()}
	}

	if notice.Master != ses.Get("user").(*db.User).Id {
		return &Fail{Msg: "你不能移除别人的通知"}
	}

	if err := n.Db.ChangeReadById(bd, notice.Id, true); err != nil {
		return &Fail{Msg: err.Error()}
	}

	ps.Publish(&Prot{Mod: "notice", Act: "list"})
	return &Succ{}
}
