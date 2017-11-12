package module

import (
	"github.com/mlgaku/back/db"
	"github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
)

type Notice struct {
	db db.Notice

	service.Di
}

// 获取通知列表
func (n *Notice) List() Value {
	return &Succ{
		Data: n.db.FindByMaster(n.Ses().Get("user").(*db.User).Id),
	}
}

// 移除通知
func (n *Notice) Remove() Value {
	notice := n.db.Find(db.NewNotice(n.Req().Body, "b").Id)

	if notice.Master != n.Ses().Get("user").(*db.User).Id {
		return &Fail{Msg: "你不能移除别人的通知"}
	}

	n.db.ChangeReadById(notice.Id, true)

	n.Ps().Publish(&Prot{Mod: "notice", Act: "list"})
	return &Succ{}
}
